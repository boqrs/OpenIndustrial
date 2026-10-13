
package routing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
)

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

type serviceImpl struct {
	repository  Repository
	resourceSvc resource.Service
}

func NewService(repository Repository, resourceSvc resource.Service) Service {
	return &serviceImpl{
		repository:  repository,
		resourceSvc: resourceSvc,
	}
}

func (s *serviceImpl) CreateRouting(
	ctx context.Context,
	req *CreateRoutingRequest,
) (*RoutingResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" || req.ProductID == 0 {
		return nil, errors.New("product ID and routing name are required")
	}

	// ProductID is a product_models.id, not a resources.id.
	// Validate ownership through product_models.resource_id.
	belongs, err := s.repository.ProductBelongsToTenant(ctx, tenantID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("validate product ownership: %w", err)
	}
	if !belongs {
		return nil, errors.New("product not found in current tenant")
	}

	const version = 1
	existing, err := s.repository.GetRoutingByNameAndVersion(
		ctx, tenantID, req.ProductID, name, version,
	)
	if err != nil {
		return nil, fmt.Errorf("check existing routing: %w", err)
	}
	if existing != nil {
		return nil, ErrRoutingExists
	}

	code := "RT-" + strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	resourceCode := code

	resourceEntity, err := s.resourceSvc.CreateResource(ctx, &resource.CreateResource{
		Type:     "ROUTING",
		Name:     name,
		Code:     &resourceCode,
		Status:   model.StatusPending,
		TenantID: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("create routing resource: %w", err)
	}

	entity := &model.Routing{
		ResourceID:  resourceEntity.ID,
		ProductID:   req.ProductID,
		Code:        code,
		Name:        name,
		Version:     version,
		Description: strings.TrimSpace(req.Description),
		Status:      model.RoutingStatusDraft,
	}

	if err := s.repository.CreateRouting(ctx, entity); err != nil {
		if cleanupErr := s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.ID); cleanupErr != nil {
			return nil, fmt.Errorf(
				"create routing: %w; cleanup resource %d also failed: %v",
				err, resourceEntity.ID, cleanupErr,
			)
		}
		return nil, fmt.Errorf("create routing: %w", err)
	}

	return toRoutingResponse(entity), nil
}

func (s *serviceImpl) GetRouting(ctx context.Context, id uint) (*RoutingResponse, error) {
	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toRoutingResponse(entity), nil
}

func (s *serviceImpl) ListRoutings(
	ctx context.Context,
	productID *uint,
	status *model.RoutingStatus,
) ([]*RoutingResponse, error) {
	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	if productID != nil {
		belongs, err := s.repository.ProductBelongsToTenant(ctx, tenantID, *productID)
		if err != nil {
			return nil, fmt.Errorf("validate product ownership: %w", err)
		}
		if !belongs {
			return nil, errors.New("product not found in current tenant")
		}
	}

	entities, err := s.repository.ListRoutings(ctx, tenantID, productID, status)
	if err != nil {
		return nil, err
	}

	responses := make([]*RoutingResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, toRoutingResponse(entity))
	}
	return responses, nil
}

func (s *serviceImpl) UpdateRouting(
	ctx context.Context,
	id uint,
	req *UpdateRoutingRequest,
) (*RoutingResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if entity.Status == model.RoutingStatusActive {
		return nil, ErrInvalidRoutingState
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("routing name cannot be empty")
		}
		entity.Name = name
	}
	if req.Description != nil {
		entity.Description = strings.TrimSpace(*req.Description)
	}

	if err := s.repository.UpdateRouting(ctx, entity); err != nil {
		return nil, fmt.Errorf("update routing: %w", err)
	}
	return toRoutingResponse(entity), nil
}

func (s *serviceImpl) ActivateRouting(ctx context.Context, id uint) error {
	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return errors.New("tenant id not found")
	}

	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if entity.Status == model.RoutingStatusActive {
		return nil
	}
	if entity.Status == model.RoutingStatusArchived {
		return ErrInvalidRoutingState
	}

	count, err := s.repository.CountOperations(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("count routing operations: %w", err)
	}
	if count == 0 {
		return errors.New("cannot activate a routing with no operations")
	}

	entity.Status = model.RoutingStatusActive
	return s.repository.UpdateRouting(ctx, entity)
}

func (s *serviceImpl) DeactivateRouting(ctx context.Context, id uint) error {
	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return errors.New("tenant id not found")
	}

	entity, err := s.repository.GetRoutingByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if entity.Status == model.RoutingStatusInactive {
		return nil
	}

	// TODO: Before deactivation, check whether unfinished work orders use this routing.
	entity.Status = model.RoutingStatusInactive
	return s.repository.UpdateRouting(ctx, entity)
}

func (s *serviceImpl) AddOperation(
	ctx context.Context,
	routingID uint,
	req *CreateOperationRequest,
) (*OperationResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" || req.Sequence <= 0 {
		return nil, errors.New("operation code, name, and a positive sequence are required")
	}
	if req.StandardDurationSeconds < 0 {
		return nil, errors.New("standard duration cannot be negative")
	}
	if len(req.Parameters) > 0 && !json.Valid(req.Parameters) {
		return nil, errors.New("parameters must contain valid JSON")
	}

	routingEntity, err := s.repository.GetRoutingByID(ctx, tenantID, routingID)
	if err != nil {
		return nil, err
	}
	if routingEntity.Status == model.RoutingStatusActive ||
		routingEntity.Status == model.RoutingStatusArchived {
		return nil, ErrInvalidRoutingState
	}

	if req.WorkstationID != nil {
		workstation, err := s.resourceSvc.GetResourceByID(ctx, tenantID, *req.WorkstationID)
		if err != nil {
			return nil, fmt.Errorf("workstation not found in current tenant: %w", err)
		}
		if workstation.ResourceType != string(resource.ResourceTypeWorkStation) {
			return nil, errors.New("selected resource is not a workstation")
		}
	}

	existingCode, err := s.repository.GetOperationByCode(ctx, tenantID, routingID, code)
	if err != nil {
		return nil, err
	}
	if existingCode != nil {
		return nil, ErrOperationExists
	}

	existingSequence, err := s.repository.GetOperationBySequence(ctx, tenantID, routingID, req.Sequence)
	if err != nil {
		return nil, err
	}
	if existingSequence != nil {
		return nil, ErrOperationSeqExists
	}

	required := true
	if req.Required != nil {
		required = *req.Required
	}

	entity := &model.RoutingOperation{
		RoutingID:               routingID,
		Code:                    code,
		Name:                    name,
		Description:             strings.TrimSpace(req.Description),
		WorkstationID:           req.WorkstationID,
		Sequence:                req.Sequence,
		StandardDurationSeconds: req.StandardDurationSeconds,
		Required:                required,
		Parameters:              append([]byte(nil), req.Parameters...),
	}

	if err := s.repository.CreateOperation(ctx, entity); err != nil {
		return nil, fmt.Errorf("create routing operation: %w", err)
	}
	return toOperationResponse(entity), nil
}

func (s *serviceImpl) ListOperations(
	ctx context.Context,
	routingID uint,
) ([]*OperationResponse, error) {
	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	// Verify the routing belongs to the current tenant before listing operations.
	if _, err := s.repository.GetRoutingByID(ctx, tenantID, routingID); err != nil {
		return nil, err
	}

	entities, err := s.repository.ListOperations(ctx, tenantID, routingID)
	if err != nil {
		return nil, err
	}

	responses := make([]*OperationResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, toOperationResponse(entity))
	}
	return responses, nil
}

func (s *serviceImpl) UpdateOperation(
	ctx context.Context,
	routingID, operationID uint,
	req *UpdateOperationRequest,
) (*OperationResponse, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return nil, errors.New("tenant id not found")
	}

	routingEntity, err := s.repository.GetRoutingByID(ctx, tenantID, routingID)
	if err != nil {
		return nil, err
	}
	if routingEntity.Status == model.RoutingStatusActive ||
		routingEntity.Status == model.RoutingStatusArchived {
		return nil, ErrInvalidRoutingState
	}

	entity, err := s.repository.GetOperation(ctx, tenantID, routingID, operationID)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, errors.New("operation code cannot be empty")
		}
		existing, err := s.repository.GetOperationByCode(ctx, tenantID, routingID, code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != entity.ID {
			return nil, ErrOperationExists
		}
		entity.Code = code
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("operation name cannot be empty")
		}
		entity.Name = name
	}
	if req.Description != nil {
		entity.Description = strings.TrimSpace(*req.Description)
	}
	if req.WorkstationID != nil {
		workstation, err := s.resourceSvc.GetResourceByID(ctx, tenantID, *req.WorkstationID)
		if err != nil {
			return nil, fmt.Errorf("workstation not found in current tenant: %w", err)
		}
		if workstation.ResourceType != string(resource.ResourceTypeWorkStation) {
			return nil, errors.New("selected resource is not a workstation")
		}
		entity.WorkstationID = req.WorkstationID
	}
	if req.Sequence != nil {
		if *req.Sequence <= 0 {
			return nil, errors.New("sequence must be positive")
		}
		existing, err := s.repository.GetOperationBySequence(ctx, tenantID, routingID, *req.Sequence)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != entity.ID {
			return nil, ErrOperationSeqExists
		}
		entity.Sequence = *req.Sequence
	}
	if req.StandardDurationSeconds != nil {
		if *req.StandardDurationSeconds < 0 {
			return nil, errors.New("standard duration cannot be negative")
		}
		entity.StandardDurationSeconds = *req.StandardDurationSeconds
	}
	if req.Required != nil {
		entity.Required = *req.Required
	}
	if len(req.Parameters) > 0 {
		if !json.Valid(req.Parameters) {
			return nil, errors.New("parameters must contain valid JSON")
		}
		entity.Parameters = append([]byte(nil), req.Parameters...)
	}

	if err := s.repository.UpdateOperation(ctx, entity); err != nil {
		return nil, fmt.Errorf("update routing operation: %w", err)
	}
	return toOperationResponse(entity), nil
}

func (s *serviceImpl) DeleteOperation(
	ctx context.Context,
	routingID, operationID uint,
) error {
	tenantID, ok := pkg.TenantIDFromContext(ctx)
	if !ok {
		return errors.New("tenant id not found")
	}

	routingEntity, err := s.repository.GetRoutingByID(ctx, tenantID, routingID)
	if err != nil {
		return err
	}
	if routingEntity.Status == model.RoutingStatusActive ||
		routingEntity.Status == model.RoutingStatusArchived {
		return ErrInvalidRoutingState
	}

	return s.repository.DeleteOperation(ctx, tenantID, routingID, operationID)
}

func toRoutingResponse(entity *model.Routing) *RoutingResponse {
	if entity == nil {
		return nil
	}
	return &RoutingResponse{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		ProductID:   entity.ProductID,
		Code:        entity.Code,
		Name:        entity.Name,
		Version:     entity.Version,
		Description: entity.Description,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func toOperationResponse(entity *model.RoutingOperation) *OperationResponse {
	if entity == nil {
		return nil
	}

	var parameters json.RawMessage
	if len(entity.Parameters) > 0 && json.Valid(entity.Parameters) {
		parameters = append(json.RawMessage(nil), entity.Parameters...)
	}

	return &OperationResponse{
		ID:                      entity.ID,
		RoutingID:               entity.RoutingID,
		Code:                    entity.Code,
		Name:                    entity.Name,
		Description:             entity.Description,
		WorkstationID:           entity.WorkstationID,
		Sequence:                entity.Sequence,
		StandardDurationSeconds: entity.StandardDurationSeconds,
		Required:                entity.Required,
		Parameters:              parameters,
		CreatedAt:               entity.CreatedAt,
		UpdatedAt:               entity.UpdatedAt,
	}
}