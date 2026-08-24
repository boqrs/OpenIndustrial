package workorder

import (
	"context"
	"errors"
	"fmt"
	"strings"

	//"time"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

var (
	ErrWorkOrderNotFound = errors.New("work order not found")

	ErrWorkOrderExists = errors.New(
		"work order already exists",
	)

	ErrInvalidWorkOrder = errors.New(
		"invalid work order",
	)

	ErrWorkOrderImmutable = errors.New(
		"work order cannot be modified in current state",
	)

	ErrInvalidWorkOrderState = errors.New(
		"invalid work order state transition",
	)

	ErrWorkOrderQuantityExceeded = errors.New(
		"work order quantity exceeds production plan quantity",
	)

	ErrProductionPlanNotFound = errors.New(
		"production plan not found",
	)

	ErrProductionPlanMismatch = errors.New(
		"work order does not match production plan",
	)
)

type serviceImpl struct {
	repository Repository

	/*
		ProductionPlanService is intentionally represented
		through a small interface.

		We don't import the Product or Factory implementation.

		That keeps WorkOrder from becoming coupled to
		the internals of other domains.
	*/
	psrv planning.Service
}

type ProductionPlanInfo struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	ProductID       uuid.UUID
	FactoryID       uuid.UUID
	PlannedQuantity int
	Status          string
}

func NewService(repository Repository,productionPlanService planning.Service) Service {
	return &serviceImpl{
		repository:           repository,
		psrv: productionPlanService,
	}
}

func (s *serviceImpl) CreateWorkOrder(ctx context.Context,req *CreateWorkOrderRequest) (*WorkOrderResponse, error) {
	tenantID := tenantIDFromContext(ctx)

	if req == nil {
		return nil, ErrInvalidWorkOrder
	}

	orderNo := strings.TrimSpace(req.OrderNo)

	if orderNo == "" {
		return nil, ErrInvalidWorkOrder
	}

	if tenantID == uuid.Nil ||
		req.ProductionPlanID == uuid.Nil ||
		req.ProductID == uuid.Nil ||
		req.FactoryID == uuid.Nil {
		return nil, ErrInvalidWorkOrder
	}

	if req.PlannedQuantity <= 0 {
		return nil, ErrInvalidWorkOrder
	}

	if req.PlannedStartAt.IsZero() ||
		req.PlannedEndAt.IsZero() ||
		!req.PlannedEndAt.After(req.PlannedStartAt) {
		return nil, ErrInvalidWorkOrder
	}

	existing, err := s.repository.GetByOrderNo(
		ctx,
		tenantID,
		orderNo,
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
		plan.ID == uuid.Nil {
		return nil, ErrProductionPlanNotFound
	}

	if plan.TenantID != tenantID {
		return nil, ErrProductionPlanNotFound
	}

	if plan.ProductID != req.ProductID ||
		plan.FactoryID != req.FactoryID {
		return nil, ErrProductionPlanMismatch
	}

	totalPlanned, err := s.repository.SumPlannedQuantity(
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
		TenantID: tenantID,

		OrderNo: orderNo,

		ProductionPlanID: req.ProductionPlanID,

		ProductID: req.ProductID,

		FactoryID: req.FactoryID,

		PlannedQuantity: req.PlannedQuantity,

		CompletedQuantity: 0,

		PlannedStartAt: req.PlannedStartAt,

		PlannedEndAt: req.PlannedEndAt,

		Status: model.WorkOrderStatusDraft,

		Priority: req.Priority,

		Description: strings.TrimSpace(
			req.Description,
		),
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

func (s *serviceImpl) GetWorkOrder(ctx context.Context,id uuid.UUID) (*WorkOrderResponse, error) {
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

func (s *serviceImpl) ListWorkOrders(
	ctx context.Context,
	status *model.WorkOrderStatus,
	productionPlanID *uuid.UUID,
) ([]*WorkOrderResponse, error) {
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

func (s *serviceImpl) UpdateWorkOrder(ctx context.Context,id uuid.UUID,req *UpdateWorkOrderRequest) (*WorkOrderResponse, error) {
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

	if req.PlannedQuantity != nil {
		if *req.PlannedQuantity <= 0 {
			return nil, ErrInvalidWorkOrder
		}

		entity.PlannedQuantity = *req.PlannedQuantity
	}

	if req.PlannedStartAt != nil {
		entity.PlannedStartAt =
			*req.PlannedStartAt
	}

	if req.PlannedEndAt != nil {
		entity.PlannedEndAt =
			*req.PlannedEndAt
	}

	if !entity.PlannedEndAt.After(
		entity.PlannedStartAt,
	) {
		return nil, ErrInvalidWorkOrder
	}

	if req.Priority != nil {
		entity.Priority = *req.Priority
	}

	if req.Description != nil {
		entity.Description =
			strings.TrimSpace(*req.Description)
	}

	if entity.PlannedQuantity <
		entity.CompletedQuantity {
		return nil, ErrInvalidWorkOrder
	}

	/*
		Quantity must still fit inside the Production Plan.

		We calculate all other work orders first,
		then add the new quantity.
	*/
	plan, err := s.psrv.GetProductionPlanByID(
		ctx,
		entity.ProductionPlanID,
	)
	if err != nil {
		return nil, err
	}

	totalPlanned, err :=
		s.repository.SumPlannedQuantity(
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

	/*
		The repository sum includes the current entity.
		Remove it before checking the new value.
	*/
	totalPlanned -= entity.PlannedQuantity

	if totalPlanned+entity.PlannedQuantity >
		plan.PlannedQuantity {
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

func (s *serviceImpl) ReleaseWorkOrder(ctx context.Context,id uuid.UUID) error {
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

	return s.repository.Update(
		ctx,
		entity,
	)
}

func (s *serviceImpl) StartWorkOrder(ctx context.Context,id uuid.UUID) error {
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

	entity.Status =
		model.WorkOrderStatusInProgress

	return s.repository.Update(
		ctx,
		entity,
	)
}

func (s *serviceImpl) CancelWorkOrder(ctx context.Context,id uuid.UUID) error {
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
	return &WorkOrderResponse{
		ID: entity.ID,

		TenantID: entity.TenantID,

		OrderNo: entity.OrderNo,

		ProductionPlanID:
			entity.ProductionPlanID,

		ProductID: entity.ProductID,

		FactoryID: entity.FactoryID,

		PlannedQuantity:
			entity.PlannedQuantity,

		CompletedQuantity:
			entity.CompletedQuantity,

		RemainingQuantity:
			entity.PlannedQuantity -
				entity.CompletedQuantity,

		PlannedStartAt:
			entity.PlannedStartAt,

		PlannedEndAt:
			entity.PlannedEndAt,

		Status: entity.Status,

		Priority: entity.Priority,

		Description: entity.Description,

		CreatedAt: entity.CreatedAt,

		UpdatedAt: entity.UpdatedAt,
	}
}

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	value := ctx.Value("tenant_id")

	if id, ok := value.(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}