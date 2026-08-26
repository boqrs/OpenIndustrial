package workorder

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotFound = errors.New("work order not found")
	ErrWorkOrderExists = errors.New("work order already exists")
	ErrInvalidWorkOrder = errors.New("invalid work order")
	ErrWorkOrderImmutable = errors.New("work order cannot be modified in current state")
	ErrInvalidWorkOrderState = errors.New("invalid work order state transition")
	ErrWorkOrderQuantityExceeded = errors.New("work order quantity exceeds production plan quantity")
	ErrProductionPlanNotFound = errors.New("production plan not found")
	ErrProductionPlanMismatch = errors.New("work order does not match production plan")
)

type serviceImpl struct {
	repository Repository
	psrv       planning.Service
}

func NewService(repository Repository, productionPlanService planning.Service) Service {
	return &serviceImpl{
		repository: repository,
		psrv:       productionPlanService,
	}
}

func (s *serviceImpl) CreateWorkOrder(ctx context.Context, req *CreateWorkOrderRequest) (*WorkOrderResponse, error) {
	tenantID := tenantIDFromContext(ctx)

	if req == nil {
		return nil, ErrInvalidWorkOrder
	}

	code := strings.TrimSpace(req.Code)

	if code == "" {
		return nil, ErrInvalidWorkOrder
	}

	if tenantID == uuid.Nil ||
		req.ProductionPlanID == 0 ||
		req.ProductID == 0 ||
		req.RoutingID == 0 {
		return nil, ErrInvalidWorkOrder
	}

	if req.PlannedQuantity <= 0 {
		return nil, ErrInvalidWorkOrder
	}

	// DueDate is optional, but if provided, it must be valid.
	if req.DueDate != nil && req.DueDate.IsZero() {
		return nil, ErrInvalidWorkOrder
	}

	existing, err := s.repository.GetByCode(
		ctx,
		tenantID,
		code,
	)

	if err == nil && existing != nil {
		return nil, ErrWorkOrderExists
	}

	if err != nil &&
		!errors.Is(err, ErrWorkOrderNotFound) {
		return nil, fmt.Errorf(
			"check work order: %w",
			err,
		)
	}

	plan, err := s.psrv.GetProductionPlanByID(
		ctx,
		req.ProductionPlanID,
	)
	if err != nil {
		return nil, err
	}

	if plan == nil ||
		plan.ID == 0 {
		return nil, ErrProductionPlanNotFound
	}

	if plan.TenantID != tenantID {
		return nil, ErrProductionPlanNotFound
	}

	if plan.ProductID != req.ProductID {
		return nil, ErrProductionPlanMismatch
	}

	totalPlanned, err := s.repository.SumPlannedQuantityByPlanID(
		ctx,
		tenantID,
		req.ProductionPlanID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"calculate production plan quantity: %w",
			err,
		)
	}

	if totalPlanned+req.PlannedQuantity >
		plan.PlannedQuantity {
		return nil, ErrWorkOrderQuantityExceeded
	}

	entity := &model.WorkOrder{
		ResourceUUID:     uuid.New(), // Assign a new resource UUID
		TenantID:         tenantID,
		Code:             code,
		ProductionPlanID: req.ProductionPlanID,
		ProductID:        req.ProductID,
		RoutingID:        req.RoutingID,
		PlannedQuantity:  req.PlannedQuantity,
		DueDate:          req.DueDate,
		Status:           model.WorkOrderStatusDraft,
		Priority:         req.Priority,
		//Description:      strings.TrimSpace(req.Description),
	}

	if err := s.repository.Create(
		ctx,
		entity,
	); err != nil {
		return nil, fmt.Errorf(
			"create work order: %w",
			err,
		)
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) GetWorkOrder(ctx context.Context, id uint) (*WorkOrderResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) ListWorkOrders(ctx context.Context,status *model.WorkOrderStatus,productionPlanID *uint) ([]*WorkOrderResponse, error) {
	tenantID := tenantIDFromContext(ctx)

	entities, err := s.repository.List(
		ctx,
		tenantID,
		status,
		productionPlanID,
	)
	if err != nil {
		return nil, err
	}

	result := make(
		[]*WorkOrderResponse,
		0,
		len(entities),
	)

	for _, entity := range entities {
		result = append(
			result,
			toResponse(entity),
		)
	}

	return result, nil
}

func (s *serviceImpl) UpdateWorkOrder(ctx context.Context, id uint, req *UpdateWorkOrderRequest) (*WorkOrderResponse, error) {
	if req == nil {
		return nil, ErrInvalidWorkOrder
	}

	tenantID := tenantIDFromContext(ctx)

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return nil, err
	}

	if entity.Status != model.WorkOrderStatusDraft {
		return nil, ErrWorkOrderImmutable
	}

	originalPlannedQuantity := entity.PlannedQuantity

	if req.PlannedQuantity != nil {
		if *req.PlannedQuantity <= 0 {
			return nil, ErrInvalidWorkOrder
		}
		entity.PlannedQuantity = *req.PlannedQuantity
	}

	if req.DueDate != nil {
		entity.DueDate = req.DueDate
	}

	if req.Priority != nil {
		entity.Priority = *req.Priority
	}

	// if req.Description != nil {
	// 	entity.Description =
	// 		strings.TrimSpace(*req.Description)
	// }

	plan, err := s.psrv.GetProductionPlanByID(
		ctx,
		entity.ProductionPlanID,
	)
	if err != nil {
		return nil, err
	}

	totalPlanned, err :=
		s.repository.SumPlannedQuantityByPlanID(
			ctx,
			tenantID,
			entity.ProductionPlanID,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"calculate production plan quantity: %w",
			err,
		)
	}

	// The repository sum includes the current entity's old value.
	// We subtract the old value and add the new one to check the new total.
	totalPlanned = totalPlanned - originalPlannedQuantity + entity.PlannedQuantity

	if totalPlanned > plan.PlannedQuantity {
		return nil, ErrWorkOrderQuantityExceeded
	}

	if err := s.repository.Update(
		ctx,
		entity,
	); err != nil {
		return nil, fmt.Errorf(
			"update work order: %w",
			err,
		)
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) ReleaseWorkOrder(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return err
	}

	if entity.Status !=
		model.WorkOrderStatusDraft {
		return ErrInvalidWorkOrderState
	}

	entity.Status =
		model.WorkOrderStatusReleased

	// In a real system, this might trigger other events,
	// like allocating materials.
	return s.repository.Update(
		ctx,
		entity,
	)
}

func (s *serviceImpl) StartWorkOrder(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return err
	}

	if entity.Status !=
		model.WorkOrderStatusReleased {
		return ErrInvalidWorkOrderState
	}

	now := time.Now()
	entity.Status = model.WorkOrderStatusInProgress
	entity.StartedAt = &now

	return s.repository.Update(
		ctx,
		entity,
	)
}

func (s *serviceImpl) CancelWorkOrder(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return err
	}

	switch entity.Status {
	case model.WorkOrderStatusDraft,
		model.WorkOrderStatusReleased:

		entity.Status =
			model.WorkOrderStatusCancelled

		return s.repository.Update(
			ctx,
			entity,
		)

	default:
		return ErrInvalidWorkOrderState
	}
}

func toResponse(entity *model.WorkOrder) *WorkOrderResponse {
	if entity == nil {
		return nil
	}
	return &WorkOrderResponse{
		ID:               entity.ID,
		ResourceUUID:     entity.ResourceUUID,
		TenantID:         entity.TenantID,
		Code:             entity.Code,
		ProductionPlanID: entity.ProductionPlanID,
		ProductID:        entity.ProductID,
		RoutingID:        entity.RoutingID,
		PlannedQuantity:  entity.PlannedQuantity,
		DueDate:          entity.DueDate,
		Status:           entity.Status,
		Priority:         entity.Priority,
		//Description:      entity.Description,
		StartedAt:        entity.StartedAt,
		CompletedAt:      entity.CompletedAt,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	value := ctx.Value("tenant_id")

	if id, ok := value.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}