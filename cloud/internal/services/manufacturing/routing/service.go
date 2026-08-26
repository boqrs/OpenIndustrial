package routing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	//"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// --- Errors ---
var (
	ErrRoutingNotFound      = errors.New("routing not found")
	ErrRoutingExists        = errors.New("routing with this name and version already exists for the product")
	ErrOperationNotFound    = errors.New("operation not found")
	ErrOperationExists      = errors.New("operation with this code already exists in the routing")
	ErrOperationSeqExists   = errors.New("operation with this sequence already exists in the routing")
	ErrInvalidRoutingState  = errors.New("operation not allowed in current routing state")
	ErrRoutingInUse         = errors.New("routing is in use and cannot be modified or deactivated")
	ErrDefaultRoutingExists = errors.New("a default routing for this product already exists")
)

// --- Service Implementation ---

type serviceImpl struct {
	repository Repository
	// productSvc product.Service // TODO: Inject when product service is ready
}

func NewService(repository Repository) Service {
	return &serviceImpl{
		repository: repository,
	}
}

// --- Routing Methods ---

func (s *serviceImpl) CreateRouting(ctx context.Context, req *CreateRoutingRequest) (*RoutingResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	name := strings.TrimSpace(req.Name)
	if name == "" || req.ProductID == 0 {
		return nil, fmt.Errorf("product ID and routing name are required")
	}

	// TODO: Validate product existence
	// _, err := s.productSvc.GetProduct(ctx, req.ProductID)
	// if err != nil {
	// 	return nil, err
	// }

	// For simplicity, versioning is handled manually for now.
	// A more robust system might auto-increment versions.
	const version = 1
	existing, err := s.repository.GetRoutingByNameAndVersion(ctx, tenantID, req.ProductID, name, version)
	if err == nil && existing != nil {
		return nil, ErrRoutingExists
	}

	if req.IsDefault {
		if err := s.repository.DeactivateOtherRoutings(ctx, tenantID, req.ProductID, 0); err != nil {
			return nil, fmt.Errorf("failed to handle default routing logic: %w", err)
		}
	}

	entity := &model.Routing{
		ResourceUUID: uuid.New(),
		//TenantID:     tenantID,
		ProductID:    req.ProductID,
		Name:         name,
		Version:      version,
		Description:  req.Description,
		Status:       model.RoutingStatusInactive,
		//IsDefault:    req.IsDefault,
	}

	if err := s.repository.CreateRouting(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to create routing: %w", err)
	}

	return toRoutingResponse(entity), nil
}

func (s *serviceImpl) GetRouting(ctx context.Context, id uint) (*RoutingResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toRoutingResponse(entity), nil
}

func (s *serviceImpl) ListRoutings(ctx context.Context, productID *uint, status *model.RoutingStatus) ([]*RoutingResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entities, err := s.repository.ListRoutings(ctx, tenantID, productID, status)
	if err != nil {
		return nil, err
	}
	responses := make([]*RoutingResponse, len(entities))
	for i, entity := range entities {
		responses[i] = toRoutingResponse(entity)
	}
	return responses, nil
}

func (s *serviceImpl) UpdateRouting(ctx context.Context, id uint, req *UpdateRoutingRequest) (*RoutingResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if entity.Status == model.RoutingStatusActive {
		return nil, ErrInvalidRoutingState
	}

	if req.Name != nil {
		entity.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.IsDefault != nil {
		if *req.IsDefault {
			if err := s.repository.DeactivateOtherRoutings(ctx, tenantID, entity.ProductID, entity.ID); err != nil {
				return nil, fmt.Errorf("failed to handle default routing logic: %w", err)
			}
		}
		//entity.IsDefault = *req.IsDefault
	}

	if err := s.repository.UpdateRouting(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to update routing: %w", err)
	}

	return toRoutingResponse(entity), nil
}

func (s *serviceImpl) ActivateRouting(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if entity.Status == model.RoutingStatusActive {
		return nil // Already active
	}

	// Ensure there's at least one operation before activating
	count, err := s.repository.CountOperations(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to count operations: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("cannot activate a routing with no operations")
	}

	entity.Status = model.RoutingStatusActive
	return s.repository.UpdateRouting(ctx, entity)
}

func (s *serviceImpl) DeactivateRouting(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if entity.Status == model.RoutingStatusInactive {
		return nil // Already inactive
	}

	// TODO: Check if the routing is used by any non-completed work orders.
	// This check is crucial to prevent data inconsistency.

	entity.Status = model.RoutingStatusInactive
	return s.repository.UpdateRouting(ctx, entity)
}

// --- Operation Methods ---

func (s *serviceImpl) AddOperation(ctx context.Context, routingID uint, req *CreateOperationRequest) (*OperationResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	code := strings.TrimSpace(req.Code)
	if code == "" || req.Name == "" || req.Sequence <= 0 {
		return nil, fmt.Errorf("operation code, name, and a positive sequence are required")
	}

	routing, err := s.repository.GetRoutingByID(ctx, tenantID, routingID)
	if err != nil {
		return nil, err
	}
	if routing.Status == model.RoutingStatusActive {
		return nil, ErrInvalidRoutingState
	}

	// TODO: Validate WorkCenter existence

	if op, _ := s.repository.GetOperationByCode(ctx, tenantID, routingID, code); op != nil {
		return nil, ErrOperationExists
	}
	if op, _ := s.repository.GetOperationBySequence(ctx, tenantID, routingID, req.Sequence); op != nil {
		return nil, ErrOperationSeqExists
	}

	entity := &model.RoutingOperation{
		//ID:             uuid.New(),
		//TenantID:       tenantID,
		RoutingID:      routingID,
		Code:           code,
		Name:           req.Name,
		Description:    req.Description,
		//WorkCenterID:   req.WorkCenterID,
		Sequence:       req.Sequence,
		//SetupTime:      req.SetupTime,
		//ProcessingTime: req.ProcessingTime,
	}

	if err := s.repository.CreateOperation(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to create operation: %w", err)
	}

	return toOperationResponse(entity), nil
}

func (s *serviceImpl) ListOperations(ctx context.Context, routingID uint) ([]*OperationResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entities, err := s.repository.ListOperations(ctx, tenantID, routingID)
	if err != nil {
		return nil, err
	}
	responses := make([]*OperationResponse, len(entities))
	for i, entity := range entities {
		responses[i] = toOperationResponse(entity)
	}
	return responses, nil
}

func (s *serviceImpl) UpdateOperation(ctx context.Context, routingID uint, operationID uuid.UUID, req *UpdateOperationRequest) (*OperationResponse, error) {
	tenantID := tenantIDFromContext(ctx)

	routing, err := s.repository.GetRoutingByID(ctx, tenantID, routingID)
	if err != nil {
		return nil, err
	}
	if routing.Status == model.RoutingStatusActive {
		return nil, ErrInvalidRoutingState
	}

	entity, err := s.repository.GetOperation(ctx, tenantID, routingID, operationID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	// if req.WorkCenterID != nil {
	// 	entity.WorkCenterID = *req.WorkCenterID
	// }
	if req.Sequence != nil {
		if *req.Sequence <= 0 {
			return nil, fmt.Errorf("sequence must be positive")
		}
		if op, _ := s.repository.GetOperationBySequence(ctx, tenantID, routingID, *req.Sequence); op != nil && op.ID != entity.ID {
			return nil, ErrOperationSeqExists
		}
		entity.Sequence = *req.Sequence
	}
	// if req.SetupTime != nil {
	// 	entity.SetupTime = *req.SetupTime
	// }
	// if req.ProcessingTime != nil {
	// 	entity.ProcessingTime = *req.ProcessingTime
	// }

	if err := s.repository.UpdateOperation(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to update operation: %w", err)
	}

	return toOperationResponse(entity), nil
}

func (s *serviceImpl) DeleteOperation(ctx context.Context, routingID uint, operationID uuid.UUID) error {
	tenantID := tenantIDFromContext(ctx)

	routing, err := s.repository.GetRoutingByID(ctx, tenantID, routingID)
	if err != nil {
		return err
	}
	if routing.Status == model.RoutingStatusActive {
		return ErrInvalidRoutingState
	}

	return s.repository.DeleteOperation(ctx, tenantID, routingID, operationID)
}

// --- Mappers ---

func toRoutingResponse(entity *model.Routing) *RoutingResponse {
	if entity == nil {
		return nil
	}
	return &RoutingResponse{
		ID:           entity.ID,
		ResourceUUID: entity.ResourceUUID,
		//TenantID:     entity.TenantID,
		ProductID:    entity.ProductID,
		Name:         entity.Name,
		Version:      entity.Version,
		Description:  entity.Description,
		Status:       entity.Status,
		//IsDefault:    entity.IsDefault,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func toOperationResponse(entity *model.RoutingOperation) *OperationResponse {
	if entity == nil {
		return nil
	}
	return &OperationResponse{
		ID:             entity.ID,
		RoutingID:      entity.RoutingID,
		Code:           entity.Code,
		Name:           entity.Name,
		Description:    entity.Description,
		//WorkCenterID:   entity.WorkCenterID,
		Sequence:       entity.Sequence,
		//SetupTime:      entity.SetupTime,
		//ProcessingTime: entity.ProcessingTime,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

// --- Helpers ---

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}