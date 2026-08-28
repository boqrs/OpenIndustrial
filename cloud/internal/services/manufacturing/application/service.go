package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotReleasable      = errors.New("work order is not in a releasable state for execution")
	ErrRoutingNotActive            = errors.New("routing is not active")
	ErrRoutingHasNoOperations      = errors.New("routing has no operations")
	ErrWorkOrderQuantityExceeded   = errors.New("the planned quantity for the work order has already been met or exceeded")
)


type service struct {
	uow        postgres.UnitOfWork
	workOrders workorder.Repository
	routings   routing.Repository
	executions execution.Repository
}

// NewService creates a new manufacturing application service.
func NewService(uow postgres.UnitOfWork, workOrders workorder.Repository, routings routing.Repository, executions execution.Repository) Service {
	return &service{
		uow:        uow,
		workOrders: workOrders,
		routings:   routings,
		executions: executions,
	}
}

// CreateProductionExecution is a transactional use case that creates a new production execution from a work order.
func (s *service) CreateProductionExecution(ctx context.Context, workOrderID uint, deviceID *uint, quantity int64) (*execution.ExecutionResponse, error) {
	tenantID := tenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID not found in context")
	}

	var result *execution.ExecutionResponse
	var err error

	err = s.uow.Execute(ctx, func(txCtx context.Context) error {
		// 1. Get the work order. In a real scenario, you'd lock the row for update.
		// The GetByID method needs to be implemented in the workorder repository.
		workOrder, err := s.workOrders.GetByID(txCtx, tenantID, workOrderID)
		if err != nil {
			return fmt.Errorf("failed to get work order: %w", err)
		}

		// 2. Validate the work order's state.
		if err := validateWorkOrderForExecution(workOrder); err != nil {
			return err
		}

		// 3. Check if the planned quantity has already been met.
		totalExecuted, err := s.executions.CountExecutions(txCtx, tenantID, workOrderID)
		if err != nil {
			return fmt.Errorf("failed to count existing executions: %w", err)
		}
		if totalExecuted >= workOrder.PlannedQuantity {
			return ErrWorkOrderQuantityExceeded
		}

		// 4. Get the associated routing.
		routingEntity, err := s.routings.GetRoutingByID(txCtx, tenantID, workOrder.RoutingID)
		if err != nil {
			return fmt.Errorf("failed to get routing: %w", err)
		}
		if err := validateRoutingForExecution(routingEntity); err != nil {
			return err
		}

		// 5. Get the operations from the routing.
		routingOperations, err := s.routings.ListOperations(txCtx, tenantID, routingEntity.ID)
		if err != nil {
			return fmt.Errorf("failed to list routing operations: %w", err)
		}
		if len(routingOperations) == 0 {
			return ErrRoutingHasNoOperations
		}

		// 6. Build the new execution entity.
		entity := &model.ProductionExecution{
			//ResourceUUID: uuid.New(),
			TenantID:     tenantID,
			WorkOrderID:  workOrder.ID,
			DeviceID:     deviceID,
			//Quantity:     quantity,
			Status:       model.ProductionExecutionStatusPending,
		}

		// 7. Build the associated execution operation entities.
		operations := buildExecutionOperations(tenantID, routingOperations)

		// 8. Persist the new execution and its operations.
		if err := s.executions.CreateExecution(txCtx, entity, operations); err != nil {
			return fmt.Errorf("failed to create execution in repository: %w", err)
		}

		// 9. Prepare the response DTO.
		result = toExecutionResponse(entity, operations)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// --- Helper Functions ---

func validateWorkOrderForExecution(wo *model.WorkOrder) error {
	if wo.Status != model.WorkOrderStatusReleased {
		return fmt.Errorf("%w: current status is '%s'", ErrWorkOrderNotReleasable, wo.Status)
	}
	return nil
}

func validateRoutingForExecution(r *model.Routing) error {
	if r.Status != model.RoutingStatusActive {
		return fmt.Errorf("%w: current status is '%s'", ErrRoutingNotActive, r.Status)
	}
	return nil
}

func buildExecutionOperations(tenantID uuid.UUID, routingOps []*model.RoutingOperation) []*model.ExecutionOperation {
	if len(routingOps) == 0 {
		return nil
	}

	execOps := make([]*model.ExecutionOperation, len(routingOps))
	for i, op := range routingOps {
		execOps[i] = &model.ExecutionOperation{
			//ID:           uuid.New(),
			//TenantID:     tenantID,
			// ExecutionID is set by the database and linked via the transaction.
			Sequence:     op.Sequence,
			//WorkCenterID: op.WorkCenterID,
			Status:       model.ExecutionOperationStatusPending,
		}
	}
	return execOps
}

func toExecutionResponse(exec *model.ProductionExecution, ops []*model.ExecutionOperation) *execution.ExecutionResponse {
	// This mapper now aligns with the execution.ExecutionResponse DTO
	return &execution.ExecutionResponse{
		ID:           exec.ID,
		ResourceID: exec.ResourceID,
		TenantID:     exec.TenantID,
		WorkOrderID:  exec.WorkOrderID,
		DeviceID:     exec.DeviceID,
		//Quantity:     exec.Quantity,
		Status:       exec.Status,
		StartedAt:    exec.StartedAt,
		CompletedAt:  exec.CompletedAt,
		CreatedAt:    exec.CreatedAt,
		UpdatedAt:    exec.UpdatedAt,
	}
}

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}