package bom

import (
	"context"
	"errors"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	mSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/material"
	pSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrBOMNotFound = errors.New("bom not found")
	ErrInvalidBOM = errors.New("invalid bom")
	ErrInvalidBOMState = errors.New("invalid bom state")
	ErrBOMMustHaveItems = errors.New("bom must contain at least one item")
	ErrInvalidBOMItem = errors.New("invalid bom item")
	ErrDuplicateBOMVersion = errors.New("bom version already exists")
	ErrMaterialNotFound = errors.New("material not found")
	ErrProductNotFound = errors.New("product not found")
)

// service implements the bom.Service interface.
type service struct {
	repo         Repository
	materialRepo mSrv.Repository
	productSvc   pSrv.Service
	uow          postgres.UnitOfWork
}

// NewService creates a new BOM service.
func NewService(
	repo Repository,
	materialRepo mSrv.Repository,
	productSvc pSrv.Service,
	uow postgres.UnitOfWork,
) Service {
	return &service{
		repo:         repo,
		materialRepo: materialRepo,
		productSvc:   productSvc,
		uow:          uow,
	}
}

func (s *service) Create(ctx context.Context, tenantID uuid.UUID, req *CreateRequest) (*Response, error) {
	if err := s.validateCreateRequest(ctx, tenantID, req); err != nil {
		return nil, err
	}

	bom := &model.BOM{
		TenantID:    tenantID,
		ProductID:   req.ProductID,
		BOMNo:       strings.TrimSpace(req.BOMNo),
		Version:     req.Version,
		Status:      model.BOMStatusDraft, // Always starts as draft
		Description: req.Description,
	}

	items := make([]*model.BOMItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = &model.BOMItem{
			TenantID:      tenantID,
			MaterialID:    itemReq.MaterialID,
			Quantity:      itemReq.Quantity,
			Unit:          itemReq.Unit,
			Sequence:      itemReq.Sequence,
			OperationCode: itemReq.OperationCode,
			IsOptional:    itemReq.IsOptional,
			Description:   itemReq.Description,
		}
	}


	err := s.uow.Execute(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(ctx, bom); err != nil {
			return err
		}
		if len(items) > 0 {
			for _, item := range items {
				item.BOMID = bom.ID
			}
			if err := s.repo.CreateItems(ctx, items); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.buildResponse(bom, items), nil
}

func (s *service) GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*Response, error) {
	bom, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBOMNotFound
		}
		return nil, err
	}

	items, err := s.repo.GetItems(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	return s.buildResponse(bom, items), nil
}

func (s *service) List(ctx context.Context, tenantID, productID uuid.UUID, offset, limit int) ([]*Response, int64, error) {
	boms, total, err := s.repo.List(ctx, tenantID, productID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*Response, len(boms))
	for i, bom := range boms {
		// Note: This creates an N+1 query problem. For production, consider a batch-fetch for items.
		items, err := s.repo.GetItems(ctx, tenantID, bom.ID)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = s.buildResponse(bom, items)
	}

	return responses, total, nil
}

func (s *service) Update(ctx context.Context, tenantID uuid.UUID, id uint, req *UpdateRequest) (*Response, error) {
	bom, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBOMNotFound
		}
		return nil, err
	}

	if bom.Status != model.BOMStatusDraft {
		return nil, ErrInvalidBOMState
	}

	if err := s.validateItems(ctx, tenantID, req.Items); err != nil {
		return nil, err
	}

	items := make([]*model.BOMItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = &model.BOMItem{
			TenantID:      tenantID,
			BOMID:         id,
			MaterialID:    itemReq.MaterialID,
			Quantity:      itemReq.Quantity,
			Unit:          itemReq.Unit,
			Sequence:      itemReq.Sequence,
			OperationCode: itemReq.OperationCode,
			IsOptional:    itemReq.IsOptional,
			Description:   itemReq.Description,
		}
	}

	bom.Description = req.Description

	err = s.uow.Execute(ctx, func(txCtx context.Context) error {
		if err := s.repo.Update(ctx, bom); err != nil {
			return err
		}
		if err := s.repo.DeleteItems(ctx, tenantID, id); err != nil {
			return err
		}
		if len(items) > 0 {
			if err := s.repo.CreateItems(ctx, items); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.buildResponse(bom, items), nil
}

func (s *service) Release(ctx context.Context, tenantID uuid.UUID, id uint) (*Response, error) {
	bom, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBOMNotFound
		}
		return nil, err
	}

	if bom.Status != model.BOMStatusDraft {
		return nil, ErrInvalidBOMState
	}

	items, err := s.repo.GetItems(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, ErrBOMMustHaveItems
	}

	// Re-validate all materials on release
	for _, item := range items {
		if _, err := s.materialRepo.GetByID(ctx, tenantID, item.MaterialID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrMaterialNotFound
			}
			return nil, err
		}
	}

	bom.Status = model.BOMStatusReleased
	if err := s.repo.Update(ctx, bom); err != nil {
		return nil, err
	}

	return s.buildResponse(bom, items), nil
}

func (s *service) Obsolete(ctx context.Context, tenantID uuid.UUID, id uint) error {
	bom, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBOMNotFound
		}
		return err
	}

	if bom.Status != model.BOMStatusReleased {
		return ErrInvalidBOMState
	}

	bom.Status = model.BOMStatusObsolete
	return s.repo.Update(ctx, bom)
}

// --- Validation Helpers ---

func (s *service) validateCreateRequest(ctx context.Context, tenantID uuid.UUID, req *CreateRequest) error {
	// Validate Product
	if _, err := s.productSvc.GetProductModel(ctx, req.ProductID); err != nil {
		// Assuming gorm.ErrRecordNotFound or a similar error is returned for not found
		return ErrProductNotFound
	}

	// Validate BOM Number and Version uniqueness
	bomNo := strings.TrimSpace(req.BOMNo)
	if bomNo == "" || req.Version <= 0 {
		return ErrInvalidBOM
	}
	existing, err := s.repo.GetByNoVersion(ctx, tenantID, bomNo, req.Version)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return ErrDuplicateBOMVersion
	}

	// Validate Items
	return s.validateItems(ctx, tenantID, req.Items)
}

func (s *service) validateItems(ctx context.Context, tenantID uuid.UUID, items []ItemRequest) error {
	if len(items) == 0 {
		return nil // It's valid to create a BOM with no items initially
	}

	materialIDs := make(map[uint]struct{})
	for _, item := range items {
		if item.Quantity <= 0 || item.Unit == "" {
			return ErrInvalidBOMItem
		}
		if _, exists := materialIDs[item.MaterialID]; exists {
			// Optional: could check for duplicate material IDs if that's a business rule
		}
		materialIDs[item.MaterialID] = struct{}{}

		// Check if material exists
		if _, err := s.materialRepo.GetByID(ctx, tenantID, item.MaterialID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrMaterialNotFound
			}
			return err
		}
	}
	return nil
}

// --- Response Builder ---

func (s *service) buildResponse(bom *model.BOM, items []*model.BOMItem) *Response {
	itemResponses := make([]ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = ItemResponse{
			ID:            item.ID,
			MaterialID:    item.MaterialID,
			Quantity:      item.Quantity,
			Unit:          item.Unit,
			Sequence:      item.Sequence,
			OperationCode: item.OperationCode,
			IsOptional:    item.IsOptional,
			Description:   item.Description,
		}
	}

	return &Response{
		ID:          bom.ID,
		TenantID:    bom.TenantID,
		ProductID:   bom.ProductID,
		BOMNo:       bom.BOMNo,
		Version:     bom.Version,
		Status:      bom.Status,
		Description: bom.Description,
		Items:       itemResponses,
		CreatedAt:   bom.CreatedAt,
		UpdatedAt:   bom.UpdatedAt,
	}
}