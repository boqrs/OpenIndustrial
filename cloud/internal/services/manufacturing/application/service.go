package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/google/uuid"
)

var (
	ErrWorkOrderNotReleasable    = errors.New("work order is not in a releasable state for execution")
	ErrRoutingNotActive          = errors.New("routing is not active")
	ErrRoutingHasNoOperations    = errors.New("routing has no operations")
	ErrWorkOrderQuantityExceeded = errors.New("the planned quantity for the work order has already been met or exceeded")
	ErrExecutionResultNotFound   = errors.New("execution result not found")
	ErrExecutionNotCompleted     = errors.New("execution is not completed")
	ErrExecutionAlreadyFailed    = errors.New("execution has failed")
	ErrExecutionResultInvalid    = errors.New("execution result is invalid")
)

type service struct {
	uow                 postgres.UnitOfWork
	resources           resource.Service
	workOrders          workorder.Repository
	routings            routing.Repository
	executions          execution.Service
	executionResults    Repository
	executionRepository execution.Repository
	devices             device.Service
}

// NewService creates a new manufacturing application service.
func NewService(uow postgres.UnitOfWork, workOrders workorder.Repository, routings routing.Repository, executions execution.Service, executionResults Repository, executionRepository execution.Repository, devices device.Service, resources resource.Service) Service {
	return &service{
		uow:                 uow,
		resources:           resources,
		workOrders:          workOrders,
		routings:            routings,
		executions:          executions,
		executionResults:    executionResults,
		executionRepository: executionRepository,
		devices:             devices,
	}
}

func (s *service) CreateProductionExecution(
	ctx context.Context,
	workOrderID uint,
) (*execution.ExecutionResponse, error) {

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID not found in context")
	}

	if workOrderID == 0 {
		return nil, execution.ErrWorkOrderNotFound
	}

	returnResult := (*execution.ExecutionResponse)(nil)

	err := s.uow.Execute(ctx, func(txCtx context.Context) error {

		workOrder, err := s.workOrders.GetByIDForUpdateTx(
			txCtx,
			tenantID,
			workOrderID,
		)
		if err != nil {
			return fmt.Errorf(
				"get work order: %w",
				err,
			)
		}

		if workOrder == nil {
			return execution.ErrWorkOrderNotFound
		}

		if workOrder.Status != model.WorkOrderStatusReleased &&
			workOrder.Status != model.WorkOrderStatusInProgress {
			return execution.ErrWorkOrderNotExecutable
		}

		existingCount, err := s.executionRepository.CountExecutions(
			txCtx,
			tenantID,
			workOrder.ID,
		)
		if err != nil {
			return fmt.Errorf(
				"count production executions: %w",
				err,
			)
		}

		if existingCount >= workOrder.PlannedQuantity {
			return ErrWorkOrderQuantityExceeded
		}

		// Execution is a Resource-backed entity.
		resourceEntity, err := s.resources.CreateResourceTx(
			txCtx,
			&resource.CreateResource{
				Type:     "PRODUCTION_EXECUTION",
				Name:     fmt.Sprintf("Execution-%d", workOrder.ID),
				Status:   "pending",
				TenantID: tenantID,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"create execution resource: %w",
				err,
			)
		}

		if resourceEntity == nil || resourceEntity.ID == 0 {
			return errors.New(
				"create execution resource returned invalid resource",
			)
		}

		result, err := s.executions.CreateExecutionTx(
			txCtx,
			&execution.CreateExecutionRequest{
				WorkOrderID: workOrder.ID,
				ResourceID:  resourceEntity.ID,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"create execution: %w",
				err,
			)
		}

		returnResult = result

		return nil
	},
	)

	if err != nil {
		return nil, err
	}

	return returnResult, nil
}

type deviceIdentity struct {
	SerialNumber string
	HardwareID   string
}

func extractDeviceIdentity(
	operations []*model.ExecutionOperation,
) (*deviceIdentity, error) {
	if len(operations) == 0 {
		return nil, fmt.Errorf(
			"%w: execution has no operations",
			ErrExecutionResultInvalid,
		)
	}

	identity := &deviceIdentity{}

	for _, operation := range operations {
		if operation == nil || len(operation.Result) == 0 {
			continue
		}

		var result map[string]any

		if err := json.Unmarshal(operation.Result, &result); err != nil {
			return nil, fmt.Errorf(
				"%w: invalid result of operation %d: %v",
				ErrExecutionResultInvalid,
				operation.ID,
				err,
			)
		}

		if value, ok := stringValue(result["serial_number"]); ok {
			if identity.SerialNumber != "" &&
				identity.SerialNumber != value {
				return nil, fmt.Errorf(
					"%w: conflicting serial_number",
					ErrExecutionResultInvalid,
				)
			}

			identity.SerialNumber = value
		}

		if value, ok := stringValue(result["hardware_id"]); ok {
			if identity.HardwareID != "" &&
				identity.HardwareID != value {
				return nil, fmt.Errorf(
					"%w: conflicting hardware_id",
					ErrExecutionResultInvalid,
				)
			}

			identity.HardwareID = value
		}
	}

	if identity.SerialNumber == "" {
		return nil, fmt.Errorf(
			"%w: no serial_number found in execution operations",
			ErrExecutionResultInvalid,
		)
	}

	return identity, nil
}

// ConfirmExecutionResult confirms the final production result.
// // // This is the manufacturing completion boundary:
//
//	// // ExecutionResult // ↓
//
// // validate Execution // ↓
//
//	// finalize qualified product identities // ↓
//	// update WorkOrder.CompletedQuantity // ↓
//	// confirm ExecutionResult //
//
// // All operations are executed inside one UnitOfWork transaction.
func (s *service) ConfirmExecutionResult(
	ctx context.Context,
	executionResultID uint,
) error {
	tenantID := pkg.TenantIDFromContext(ctx)

	if tenantID == uuid.Nil {
		return errors.New("tenant ID not found in context")
	}

	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		// 1. Lock ExecutionResult.
		result, err := s.executionResults.GetByIDForUpdateTx(
			txCtx,
			tenantID,
			executionResultID,
		)
		if err != nil {
			return fmt.Errorf(
				"get execution result: %w",
				err,
			)
		}

		if result == nil {
			return ErrExecutionResultNotFound
		}

		// Idempotent.
		if result.Status ==
			model.ExecutionResultStatusConfirmed {
			return nil
		}

		if result.Status !=
			model.ExecutionResultStatusDraft {
			return fmt.Errorf(
				"%w: status=%s",
				ErrExecutionResultInvalid,
				result.Status,
			)
		}

		// 2. Quantity validation.
		if result.ProducedQuantity < 0 ||
			result.QualifiedQuantity < 0 ||
			result.RejectedQuantity < 0 {
			return fmt.Errorf(
				"%w: negative quantity",
				ErrExecutionResultInvalid,
			)
		}

		if result.QualifiedQuantity+
			result.RejectedQuantity !=
			result.ProducedQuantity {
			return fmt.Errorf(
				"%w: produced=%d qualified=%d rejected=%d",
				ErrExecutionResultInvalid,
				result.ProducedQuantity,
				result.QualifiedQuantity,
				result.RejectedQuantity,
			)
		}

		// Current manufacturing model:
		// one Execution represents one physical product.
		if result.ProducedQuantity != 1 {
			return fmt.Errorf(
				"%w: execution result produced quantity must be 1, got %d",
				ErrExecutionResultInvalid,
				result.ProducedQuantity,
			)
		}

		// 3. Lock Execution.
		exec, err :=
			s.executionRepository.GetExecutionByIDForUpdateTx(
				txCtx,
				tenantID,
				result.ExecutionID,
			)
		if err != nil {
			return fmt.Errorf(
				"get execution for update: %w",
				err,
			)
		}

		if exec == nil {
			return execution.ErrExecutionNotFound
		}

		if exec.WorkOrderID != result.WorkOrderID {
			return fmt.Errorf(
				"%w: execution work order mismatch",
				ErrExecutionResultInvalid,
			)
		}

		if exec.Status ==
			model.ProductionExecutionStatusFailed {
			return ErrExecutionAlreadyFailed
		}

		if exec.Status !=
			model.ProductionExecutionStatusCompleted {
			return ErrExecutionNotCompleted
		}

		// 4. Lock WorkOrder.
		workOrder, err :=
			s.workOrders.GetByIDForUpdateTx(
				txCtx,
				tenantID,
				result.WorkOrderID,
			)
		if err != nil {
			return fmt.Errorf(
				"get work order for update: %w",
				err,
			)
		}

		if workOrder == nil {
			return execution.ErrWorkOrderNotFound
		}

		if workOrder.Status ==
			model.WorkOrderStatusCancelled {
			return fmt.Errorf(
				"%w: work order is cancelled",
				ErrExecutionResultInvalid,
			)
		}

		if workOrder.Status != model.WorkOrderStatusReleased && workOrder.Status != model.WorkOrderStatusInProgress {
			return execution.ErrWorkOrderNotExecutable
		}
		// existingCount, err := s.executionRepository.CountExecutions(txCtx, tenantID, workOrder.ID)
		// if err != nil {
		// 	return fmt.Errorf("count production executions: %w", err)
		// }
		// if existingCount >= workOrder.PlannedQuantity {
		// 	return ErrWorkOrderQuantityExceeded
		// }

		// A confirmed result must never make the WorkOrder
		// exceed its planned production quantity.
		if workOrder.CompletedQuantity+result.QualifiedQuantity > workOrder.PlannedQuantity {
			return fmt.Errorf("%w: completed=%d qualified=%d planned=%d", ErrExecutionResultInvalid, workOrder.CompletedQuantity, result.QualifiedQuantity, workOrder.PlannedQuantity)
		}

		// 5. Qualified product → Device.
		if result.QualifiedQuantity == 1 {
			operations, err :=
				s.executionRepository.ListOperations(
					txCtx,
					exec.ID,
				)
			if err != nil {
				return fmt.Errorf(
					"list execution operations: %w",
					err,
				)
			}

			identity, err := extractDeviceIdentity(operations)

			if err != nil {
				return err
			}

			req := &device.CreateDeviceFromExecutionResultRequest{
				ProductID:         exec.ProductID,
				WorkOrderID:       exec.WorkOrderID,
				ExecutionID:       exec.ID,
				ExecutionResultID: result.ID,
				SerialNumber:      identity.SerialNumber,
				HardwareID:        identity.HardwareID,
			}

			if _, err :=
				s.devices.CreateFromExecutionResultTx(
					txCtx,
					req,
				); err != nil {
				return fmt.Errorf(
					"create device: %w",
					err,
				)
			}
		}

		// 6. Update WorkOrder.
		workOrder.CompletedQuantity += result.ProducedQuantity
		if workOrder.CompletedQuantity >= workOrder.PlannedQuantity {
			now := time.Now()

			workOrder.Status = model.WorkOrderStatusCompleted
			workOrder.CompletedAt = &now
		}

		if err := s.workOrders.UpdateTx(
			txCtx,
			workOrder,
		); err != nil {
			return fmt.Errorf(
				"update work order: %w",
				err,
			)
		}

		// 7. Confirm result.
		now := time.Now()

		result.Status =
			model.ExecutionResultStatusConfirmed

		result.ConfirmedAt = &now

		if err := s.executionResults.UpdateTx(
			txCtx,
			result,
		); err != nil {
			return fmt.Errorf(
				"confirm execution result: %w",
				err,
			)
		}

		return nil
	})
}

func stringValue(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case nil:
		return "", false
	default:
		return "", false
	}
}

// StartProductionExecution starts a production execution.
//
// The Execution row is locked first, followed by the associated WorkOrder.
// This prevents concurrent requests from starting the same execution twice
// and keeps the lock order consistent with other manufacturing workflows.
func (s *service) StartProductionExecution(
	ctx context.Context,
	executionID uint,
) error {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return errors.New("tenant ID not found in context")
	}

	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		// -----------------------------------------------------------------
		// 1. Lock Execution
		// -----------------------------------------------------------------

		exec, err := s.executionRepository.GetExecutionByIDForUpdateTx(
			txCtx,
			tenantID,
			executionID,
		)
		if err != nil {
			return fmt.Errorf("get execution for update: %w", err)
		}

		if exec == nil {
			return execution.ErrExecutionNotFound
		}

		if exec.Status != model.ProductionExecutionStatusPending {
			return execution.ErrInvalidExecutionState
		}

		// -----------------------------------------------------------------
		// 2. Lock WorkOrder
		// -----------------------------------------------------------------

		workOrder, err := s.workOrders.GetByIDForUpdateTx(
			txCtx,
			tenantID,
			exec.WorkOrderID,
		)
		if err != nil {
			return fmt.Errorf("get work order for update: %w", err)
		}

		if workOrder == nil {
			return execution.ErrWorkOrderNotFound
		}

		if workOrder.Status != model.WorkOrderStatusReleased &&
			workOrder.Status != model.WorkOrderStatusInProgress {
			return execution.ErrWorkOrderNotExecutable
		}

		// -----------------------------------------------------------------
		// 3. Start WorkOrder if necessary
		// -----------------------------------------------------------------

		now := time.Now()

		if workOrder.Status == model.WorkOrderStatusReleased {
			workOrder.Status = model.WorkOrderStatusInProgress
			workOrder.StartedAt = &now

			if err := s.workOrders.UpdateTx(
				txCtx,
				workOrder,
			); err != nil {
				return fmt.Errorf(
					"start work order: %w",
					err,
				)
			}
		}

		// -----------------------------------------------------------------
		// 4. Start Execution
		// -----------------------------------------------------------------

		exec.Status = model.ProductionExecutionStatusInProgress
		exec.StartedAt = &now

		if err := s.executionRepository.UpdateExecution(
			txCtx,
			exec,
		); err != nil {
			return fmt.Errorf(
				"start execution: %w",
				err,
			)
		}

		return nil
	})
}

// StartProductionOperation starts one execution operation.
//
// The operation state transition is executed inside a transaction so the
// operation row can be locked with FOR UPDATE and the state transition is
// atomic.
func (s *service) StartProductionOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
) error {
	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		return s.executions.StartOperation(
			txCtx,
			executionID,
			operationID,
		)
	})
}

// CompleteProductionOperation completes one execution operation.
//
// Completing an operation may also complete the whole ProductionExecution.
// Both state changes must happen in the same transaction.
func (s *service) CompleteProductionOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
	result map[string]any,
) error {
	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		return s.executions.CompleteOperation(
			txCtx,
			executionID,
			operationID,
			result,
		)
	})
}

// FailProductionOperation marks one execution operation as failed.
//
// Operation failure also fails the whole ProductionExecution. Both updates
// must happen atomically.
func (s *service) FailProductionOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
	result map[string]any,
) error {
	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		return s.executions.FailOperation(
			txCtx,
			executionID,
			operationID,
			result,
		)
	})
}

func (s *service) CancelProductionExecution(
	ctx context.Context,
	executionID uint,
) error {

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return errors.New("tenant ID not found in context")
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {

			exec, err := s.executionRepository.GetExecutionByIDForUpdateTx(
				txCtx,
				tenantID,
				executionID,
			)
			if err != nil {
				return fmt.Errorf(
					"get execution for update: %w",
					err,
				)
			}

			if exec == nil {
				return execution.ErrExecutionNotFound
			}

			if exec.Status != model.ProductionExecutionStatusPending &&
				exec.Status != model.ProductionExecutionStatusInProgress {
				return execution.ErrInvalidExecutionState
			}

			exec.Status = model.ProductionExecutionStatusCancelled

			if err := s.executionRepository.UpdateExecution(
				txCtx,
				exec,
			); err != nil {
				return fmt.Errorf(
					"cancel execution: %w",
					err,
				)
			}

			return nil
		},
	)
}
