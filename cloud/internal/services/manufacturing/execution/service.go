package execution

import (
	"context"
	//"encoding/json"
	"errors"
	"fmt"
	//"strings"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// --- Errors ---
var (
	ErrExecutionNotFound    = errors.New("execution not found")
	ErrOperationNotFound    = errors.New("execution operation not found")
	ErrInvalidExecutionState = errors.New("operation not allowed in current execution state")
	ErrWorkOrderNotFound    = errors.New("associated work order not found")
	ErrWorkOrderMismatch    = errors.New("work order details do not match")
	ErrInvalidOperationState = errors.New("invalid operation state transition")
	ErrPriorOperationIncomplete = errors.New("prior operation is not yet complete")
)

// --- Service Implementation ---

type serviceImpl struct {
	repository    Repository
	workOrderSvc workorder.Service
	// routingSvc routing.Service // Needed to get operations for a work order's routing
}

func NewService(repository Repository, workOrderSvc workorder.Service) Service {
	return &serviceImpl{
		repository:    repository,
		workOrderSvc: workOrderSvc,
	}
}

// --- Execution Methods ---

func (s *serviceImpl) CreateExecution(ctx context.Context, req *CreateExecutionRequest) (*ExecutionResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	if req.WorkOrderID == 0 || req.Quantity <= 0 {
		return nil, fmt.Errorf("work order ID and a positive quantity are required")
	}

	// --- Temporarily commented out to allow independent compilation ---
	/*
	wo, err := s.workOrderSvc.GetWorkOrder(ctx, req.WorkOrderID)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	// Validate work order status, etc.
	if wo.Status != model.WorkOrderStatusReleased && wo.Status != model.WorkOrderStatusInProgress {
		return nil, fmt.Errorf("work order is not in a state that allows execution")
	}

	// TODO: Fetch operations from the work order's routing
	// routingOps, err := s.routingSvc.ListOperations(ctx, wo.RoutingID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get routing operations: %w", err)
	// }
	// if len(routingOps) == 0 {
	// 	return nil, fmt.Errorf("cannot create execution for a routing with no operations")
	// }
	*/

	entity := &model.ProductionExecution{
		ResourceUUID: uuid.New(),
		TenantID:     tenantID,
		WorkOrderID:  req.WorkOrderID,
		DeviceID:     req.DeviceID,
		//Quantity:     req.Quantity,
		Status:       model.ProductionExecutionStatusPending,
	}

	// --- This part needs the routing operations ---
	// var execOps []*model.ExecutionOperation
	// for _, op := range routingOps {
	// 	execOps = append(execOps, &model.ExecutionOperation{
	// 		ID:           uuid.New(),
	// 		TenantID:     tenantID,
	// 		ExecutionID:  entity.ID, // This will be set by the DB
	// 		Sequence:     op.Sequence,
	// 		WorkCenterID: op.WorkCenterID,
	// 		Status:       model.ExecutionOperationStatusPending,
	// 	})
	// }

	// For now, create without operations to allow compilation
	if err := s.repository.CreateExecution(ctx, entity, nil); err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	return toExecutionResponse(entity), nil
}

func (s *serviceImpl) GetExecution(ctx context.Context, id uint) (*ExecutionResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entity, err := s.repository.GetExecutionByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toExecutionResponse(entity), nil
}

func (s *serviceImpl) ListExecutions(ctx context.Context, workOrderID *uint, deviceID *uint, status *model.ProductionExecutionStatus) ([]*ExecutionResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entities, err := s.repository.ListExecutions(ctx, tenantID, workOrderID, deviceID, status)
	if err != nil {
		return nil, err
	}
	responses := make([]*ExecutionResponse, len(entities))
	for i, entity := range entities {
		responses[i] = toExecutionResponse(entity)
	}
	return responses, nil
}

func (s *serviceImpl) StartExecution(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)
	exec, err := s.repository.GetExecutionByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if exec.Status != model.ProductionExecutionStatusPending {
		return ErrInvalidExecutionState
	}

	now := time.Now()
	exec.Status = model.ProductionExecutionStatusInProgress
	exec.StartedAt = &now

	return s.repository.UpdateExecution(ctx, exec)
}

func (s *serviceImpl) CancelExecution(ctx context.Context, id uint) error {
	tenantID := tenantIDFromContext(ctx)
	exec, err := s.repository.GetExecutionByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if exec.Status == model.ProductionExecutionStatusCompleted || exec.Status == model.ProductionExecutionStatusCancelled {
		return ErrInvalidExecutionState
	}

	exec.Status = model.ProductionExecutionStatusCancelled
	return s.repository.UpdateExecution(ctx, exec)
}

// --- Operation Methods ---

func (s *serviceImpl) StartOperation(ctx context.Context, executionID uint, operationID uuid.UUID) error {
	tenantID := tenantIDFromContext(ctx)
	op, err := s.repository.GetOperation(ctx, tenantID, executionID, operationID)
	if err != nil {
		return err
	}

	if op.Status != model.ExecutionOperationStatusPending {
		return ErrInvalidOperationState
	}

	// Check if this is the first operation or if the previous one is complete
	if op.Sequence > 1 {
		ops, err := s.repository.ListOperations(ctx, tenantID, executionID)
		if err != nil {
			return err
		}
		for _, prevOp := range ops {
			if prevOp.Sequence == op.Sequence-1 {
				if prevOp.Status != model.ExecutionOperationStatusCompleted {
					return ErrPriorOperationIncomplete
				}
				break
			}
		}
	}

	now := time.Now()
	op.Status = model.ExecutionOperationStatusInProgress
	op.StartedAt = &now

	return s.repository.UpdateOperation(ctx, op)
}

func (s *serviceImpl) CompleteOperation(ctx context.Context, executionID uint, operationID uuid.UUID, result map[string]any) error {
	tenantID := tenantIDFromContext(ctx)
	op, err := s.repository.GetOperation(ctx, tenantID, executionID, operationID)
	if err != nil {
		return err
	}

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	now := time.Now()
	op.Status = model.ExecutionOperationStatusCompleted
	op.CompletedAt = &now
	//op.Result = result

	return s.repository.UpdateOperation(ctx, op)
}

func (s *serviceImpl) FailOperation(ctx context.Context, executionID uint, operationID uuid.UUID, result map[string]any) error {
	tenantID := tenantIDFromContext(ctx)
	op, err := s.repository.GetOperation(ctx, tenantID, executionID, operationID)
	if err != nil {
		return err
	}

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	now := time.Now()
	op.Status = model.ExecutionOperationStatusFailed
	op.CompletedAt = &now
	//op.Result = result

	return s.repository.UpdateOperation(ctx, op)
}

func (s *serviceImpl) ListOperations(ctx context.Context, executionID uint) ([]*OperationResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	entities, err := s.repository.ListOperations(ctx, tenantID, executionID)
	if err != nil {
		return nil, err
	}
	responses := make([]*OperationResponse, len(entities))
	for i, entity := range entities {
		responses[i] = toOperationResponse(entity)
	}
	return responses, nil
}

// --- Mappers ---

func toExecutionResponse(entity *model.ProductionExecution) *ExecutionResponse {
	if entity == nil {
		return nil
	}
	return &ExecutionResponse{
		ID:           entity.ID,
		ResourceUUID: entity.ResourceUUID,
		TenantID:     entity.TenantID,
		WorkOrderID:  entity.WorkOrderID,
		DeviceID:     entity.DeviceID,
		//Quantity:     entity.Quantity,
		Status:       entity.Status,
		StartedAt:    entity.StartedAt,
		CompletedAt:  entity.CompletedAt,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func toOperationResponse(entity *model.ExecutionOperation) *OperationResponse {
	if entity == nil {
		return nil
	}
	return &OperationResponse{
		ID:           entity.ID,
		ExecutionID:  entity.ExecutionID,
		Sequence:     entity.Sequence,
		//WorkCenterID: entity.WorkCenterID,
		Status:       entity.Status,
		//Result:       entity.Result,
		StartedAt:    entity.StartedAt,
		CompletedAt:  entity.CompletedAt,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

// --- Helpers ---

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}