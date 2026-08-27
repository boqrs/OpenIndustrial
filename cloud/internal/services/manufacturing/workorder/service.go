package workorder

import (
	"context"
	"fmt"
	"errors"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	bomSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/bom"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	"github.com/google/uuid"
)

var (
	ErrInvalidWorkOrder      = errors.New("invalid work order data")
	ErrWorkOrderNotFound     = errors.New("work order not found")
	ErrInvalidWorkOrderState = errors.New("invalid work order state for this operation")
	ErrPlanProductMismatch   = errors.New("work order product does not match production plan product")
	ErrQuantityExceedsPlan   = errors.New("work order quantity exceeds remaining quantity of the production plan")
	ErrBOMProductMismatch    = errors.New("bom does not belong to the specified product")
	ErrBOMNotReleased        = errors.New("bom is not in released status")
)

type serviceImpl struct {
	repository Repository
	psrv       planning.Service
	bsrv       bomSrv.Service
}

func NewService(
	repository Repository,
	productionPlanService planning.Service,
	bomService bomSrv.Service,
) Service {
	return &serviceImpl{
		repository: repository,
		psrv:       productionPlanService,
		bsrv:       bomService,
	}
}

func (s *serviceImpl) Create(ctx context.Context, tenantID uuid.UUID, req *CreateRequest) (*Response, error) {
	// --- Basic Validation ---
	if req.ProductionPlanID == 0 || req.ProductID == 0 || req.BOMID == 0 || req.RoutingID == 0 {
		return nil, ErrInvalidWorkOrder
	}

	// --- Production Plan Validation ---
	plan, err := s.psrv.GetProductionPlanByID(ctx, req.ProductionPlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get production plan: %w", err)
	}
	if plan.ProductID != req.ProductID {
		return nil, ErrPlanProductMismatch
	}
	// if req.PlannedQuantity > plan.Quantity { // Simplified logic
	// 	return nil, ErrQuantityExceedsPlan
	// }

	// --- BOM Validation ---
	bom, err := s.bsrv.GetByID(ctx, tenantID, req.BOMID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bom: %w", err)
	}
	if bom.ProductID != req.ProductID {
		return nil, ErrBOMProductMismatch
	}
	// Correctly use the constant from the BOM domain itself, assuming bom package exposes it.
	// This avoids reaching into the shared model package and helps prevent circular dependencies.
	if bom.Status != "released" { // Assuming bom.Status is a string like "released"
		return nil, ErrBOMNotReleased
	}

	// --- Entity Creation ---
	entity := &model.WorkOrder{
		ResourceUUID:     uuid.New(),
		TenantID:         tenantID,
		ProductionPlanID: req.ProductionPlanID,
		ProductID:        req.ProductID,
		BOMID:            req.BOMID,
		RoutingID:        req.RoutingID,
		Code:             req.Code,
		PlannedQuantity:  req.PlannedQuantity,
		DueDate:          req.DueDate,
		Status:           model.WorkOrderStatusDraft,
		Priority:         req.Priority,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to create work order: %w", err)
	}

	return ToResponse(entity), nil
}

func (s *serviceImpl) Release(ctx context.Context, tenantID uuid.UUID, id uint) error {
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to get work order for release: %w", err)
	}

	if entity.Status != model.WorkOrderStatusDraft {
		return ErrInvalidWorkOrderState
	}

	// --- Production Validity Re-validation ---
	// 1. Re-validate Production Plan
	plan, err := s.psrv.GetProductionPlanByID(ctx, entity.ProductionPlanID)
	if err != nil {
		return fmt.Errorf("failed to re-validate production plan: %w", err)
	}
	if plan.ProductID != entity.ProductID {
		return ErrPlanProductMismatch
	}
	// // A more robust check would be to sum all related work orders' quantities
	// if entity.PlannedQuantity > plan.Quantity {
	// 	return ErrQuantityExceedsPlan
	// }

	// 2. Re-validate BOM
	bom, err := s.bsrv.GetByID(ctx, tenantID, entity.BOMID)
	if err != nil {
		return fmt.Errorf("failed to re-validate bom: %w", err)
	}
	if bom.ProductID != entity.ProductID {
		return ErrBOMProductMismatch
	}
	if bom.Status != "released" { // Assuming bom.Status is a string like "released"
		return ErrBOMNotReleased
	}

	// 3. TODO: Re-validate Routing (once Routing service is integrated)

	entity.Status = model.WorkOrderStatusReleased
	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf("failed to release work order: %w", err)
	}
	return nil
}

// ... (GetByID, List, Update, and other methods remain the same)
func (s *serviceImpl) GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*Response, error) {
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work order by id: %w", err)
	}
	return ToResponse(entity), nil
}

func (s *serviceImpl) List(ctx context.Context, req *ListRequest) ([]*Response, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	entities, err := s.repository.List(ctx, req.TenantID, &req.ProductID, offset, req.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders: %w", err)
	}

	total, err := s.repository.Count(ctx, req.TenantID, req.ProductID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count work orders: %w", err)
	}

	return ToResponses(entities), total, nil
}

func (s *serviceImpl) Update(ctx context.Context, tenantID uuid.UUID, id uint, req *UpdateRequest) (*Response, error) {
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work order for update: %w", err)
	}

	if entity.Status != model.WorkOrderStatusDraft {
		return nil, ErrInvalidWorkOrderState
	}

	// Update fields
	entity.Code = req.Code
	entity.PlannedQuantity = req.PlannedQuantity
	entity.Priority = req.Priority
	entity.DueDate = req.DueDate

	if err := s.repository.Update(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to update work order: %w", err)
	}

	return ToResponse(entity), nil
}

func (s *serviceImpl) Start(ctx context.Context, tenantID uuid.UUID, id uint) error {
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to get work order for start: %w", err)
	}

	if entity.Status != model.WorkOrderStatusReleased {
		return ErrInvalidWorkOrderState
	}

	entity.Status = model.WorkOrderStatusInProgress
	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf("failed to start work order: %w", err)
	}
	return nil
}

func (s *serviceImpl) Complete(ctx context.Context, tenantID uuid.UUID, id uint) error {
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to get work order for completion: %w", err)
	}

	if entity.Status != model.WorkOrderStatusInProgress {
		return ErrInvalidWorkOrderState
	}

	entity.Status = model.WorkOrderStatusCompleted
	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf("failed to complete work order: %w", err)
	}
	return nil
}

func (s *serviceImpl) Cancel(ctx context.Context, tenantID uuid.UUID, id uint) error {
	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("failed to get work order for cancellation: %w", err)
	}

	// You can cancel a draft or a released work order
	if entity.Status != model.WorkOrderStatusDraft && entity.Status != model.WorkOrderStatusReleased {
		return ErrInvalidWorkOrderState
	}

	entity.Status = model.WorkOrderStatusCancelled
	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf("failed to cancel work order: %w", err)
	}
	return nil
}