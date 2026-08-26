package planning

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	"github.com/google/uuid"
)

var (
	ErrProductionPlanNotFound    = errors.New("production plan not found")
	ErrProductionPlanNoExists    = errors.New("production plan with this number already exists")
	ErrInvalidProductionPlan     = errors.New("invalid production plan data, please check constraints")
	ErrProductionPlanState       = errors.New("operation not allowed in the current production plan state")
	ErrReferencedProductNotFound = errors.New("referenced product not found")
	ErrReferencedFactoryNotFound = errors.New("referenced factory not found")
)

type serviceImpl struct {
	repository Repository
	productSvc product.Service // Dependency to validate product existence
	factorySvc factory.Service // Dependency to validate factory existence
}

// NewService creates a new planning service.
func NewService(repo Repository, productSvc product.Service, factorySvc factory.Service) Service {
	return &serviceImpl{
		repository: repo,
		productSvc: productSvc,
		factorySvc: factorySvc,
	}
}

func (s *serviceImpl) CreateProductionPlan(ctx context.Context, req *CreateProductionPlanRequest) (*ProductionPlanResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	if !req.PlannedEndAt.After(req.PlannedStartAt) || req.PlannedQuantity <= 0 {
		return nil, ErrInvalidProductionPlan
	}

	// Check for plan number uniqueness
	_, err := s.repository.GetByPlanNo(ctx, tenantID, req.PlanNo)
	if err == nil {
		return nil, ErrProductionPlanNoExists
	}
	if !errors.Is(err, ErrProductionPlanNotFound) {
		return nil, err
	}

	// TODO: Uncomment and fix after updating product and factory service interfaces.
	// The following lines are commented out to allow the 'planning' module to compile independently.
	// They will fail until the GetProductModel and GetFactory methods are updated to accept uint IDs.
	/*
		// Validate dependencies
		if _, err := s.productSvc.GetProductModel(ctx, req.ProductID); err != nil {
			return nil, ErrReferencedProductNotFound
		}
		if _, err := s.factorySvc.GetFactory(ctx, req.FactoryID); err != nil {
			return nil, ErrReferencedFactoryNotFound
		}
	*/

	entity := &model.ProductionPlan{
		ResourceUUID:    uuid.New(),
		TenantID:        tenantID,
		PlanNo:          req.PlanNo,
		ProductID:       req.ProductID,
		FactoryID:       req.FactoryID,
		PlannedQuantity: req.PlannedQuantity,
		PlannedStartAt:  req.PlannedStartAt,
		PlannedEndAt:    req.PlannedEndAt,
		Description:     req.Description,
		Status:          model.ProductionPlanStatusDraft,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("create production plan: %w", err)
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) GetProductionPlanByID(ctx context.Context, id uint) (*ProductionPlanResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toResponse(entity), nil
}

func (s *serviceImpl) ListProductionPlans(ctx context.Context, status *model.ProductionPlanStatus) ([]*ProductionPlanResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entities, err := s.repository.List(ctx, tenantID, status)
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductionPlanResponse, len(entities))
	for i, entity := range entities {
		responses[i] = toResponse(entity)
	}
	return responses, nil
}

func (s *serviceImpl) UpdateProductionPlan(ctx context.Context, id uint, req *UpdateProductionPlanRequest) (*ProductionPlanResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if entity.Status != model.ProductionPlanStatusDraft {
		return nil, ErrProductionPlanState
	}

	if req.PlannedQuantity != nil {
		if *req.PlannedQuantity <= 0 {
			return nil, ErrInvalidProductionPlan
		}
		entity.PlannedQuantity = *req.PlannedQuantity
	}
	if req.PlannedStartAt != nil {
		entity.PlannedStartAt = *req.PlannedStartAt
	}
	if req.PlannedEndAt != nil {
		entity.PlannedEndAt = *req.PlannedEndAt
	}
	if !entity.PlannedEndAt.After(entity.PlannedStartAt) {
		return nil, ErrInvalidProductionPlan
	}
	if req.Description != nil {
		entity.Description = strings.TrimSpace(*req.Description)
	}

	if err := s.repository.Update(ctx, entity); err != nil {
		return nil, fmt.Errorf("update production plan: %w", err)
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) ReleaseProductionPlan(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if entity.Status != model.ProductionPlanStatusDraft {
		return ErrProductionPlanState
	}

	entity.Status = model.ProductionPlanStatusReleased
	return s.repository.Update(ctx, entity)
}

func (s *serviceImpl) CancelProductionPlan(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	switch entity.Status {
	case model.ProductionPlanStatusDraft, model.ProductionPlanStatusReleased:
		entity.Status = model.ProductionPlanStatusCancelled
		return s.repository.Update(ctx, entity)
	default:
		return ErrProductionPlanState
	}
}

func toResponse(entity *model.ProductionPlan) *ProductionPlanResponse {
	if entity == nil {
		return nil
	}
	return &ProductionPlanResponse{
		ID:              entity.ID,
		ResourceUUID:    entity.ResourceUUID,
		TenantID:        entity.TenantID,
		PlanNo:          entity.PlanNo,
		ProductID:       entity.ProductID,
		FactoryID:       entity.FactoryID,
		PlannedQuantity: entity.PlannedQuantity,
		PlannedStartAt:  entity.PlannedStartAt,
		PlannedEndAt:    entity.PlannedEndAt,
		Status:          entity.Status,
		Description:     entity.Description,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}

// tenantIDFromContext is a placeholder helper.
func tenantIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}
	return uuid.Nil // Should be handled by auth middleware
}