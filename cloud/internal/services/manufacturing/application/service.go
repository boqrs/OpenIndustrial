package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/executionresult"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/google/uuid"
)

var ( 
	ErrWorkOrderNotReleasable = errors.New("work order is not in a releasable state for execution") 
	ErrRoutingNotActive = errors.New("routing is not active") 
	ErrRoutingHasNoOperations = errors.New("routing has no operations") 
	ErrWorkOrderQuantityExceeded = errors.New("the planned quantity for the work order has already been met or exceeded")
	ErrExecutionResultNotFound = errors.New("execution result not found")
	ErrExecutionNotCompleted = errors.New("execution is not completed") 
	ErrExecutionAlreadyFailed = errors.New("execution has failed") 
	ErrExecutionResultInvalid = errors.New("execution result is invalid") 
)


type service struct {
	uow        postgres.UnitOfWork
	workOrders workorder.Repository
	routings   routing.Repository
	executions execution.Repository
	executionResults executionresult.Repository
	devices device.Service
}

// NewService creates a new manufacturing application service.
func NewService(uow postgres.UnitOfWork, workOrders workorder.Repository, routings routing.Repository, executions execution.Repository, executionResults executionresult.Repository, devices device.Service) Service {
	return &service{
		uow:        uow,
		workOrders: workOrders,
		routings:   routings,
		executions: executions,
		executionResults: executionResults,
		devices: devices,
	}
}

// CreateProductionExecution is a transactional use case that creates a new production execution from a work order.
func (s *service) CreateProductionExecution(ctx context.Context, workOrderID uint, deviceID *uint) (*execution.ExecutionResponse, error) {
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
		operations := buildExecutionOperations(routingOperations)

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

func buildExecutionOperations(
    routingOps []*model.RoutingOperation,
) []*model.ExecutionOperation {
    operations := make([]*model.ExecutionOperation, 0, len(routingOps))

    for _, routingOp := range routingOps {
        operations = append(operations, &model.ExecutionOperation{
            RoutingOperationID: &routingOp.ID,
            Code:               routingOp.Code,
            Name:               routingOp.Name,
            Description:        routingOp.Description,
            Sequence:           routingOp.Sequence,
            WorkstationID:      routingOp.WorkstationID,
            Parameters:         routingOp.Parameters,
            Status:             model.ExecutionOperationStatusPending,
        })
    }

    return operations
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

// ConfirmExecutionResult confirms the final production result. 
// // // This is the manufacturing completion boundary:
//  // // ExecutionResult // ↓ 
// // validate Execution // ↓
//  // finalize qualified product identities // ↓
//  // update WorkOrder.CompletedQuantity // ↓
//  // confirm ExecutionResult // 
// // All operations are executed inside one UnitOfWork transaction.
 func (s *service) ConfirmExecutionResult( ctx context.Context, executionResultID uint, ) error {
	 tenantID := pkg.TenantIDFromContext(ctx) 
	 if tenantID == uuid.Nil {
		 return errors.New("tenant ID not found in context")
	} 
	return s.uow.Execute(ctx, func(txCtx context.Context) error { 
		// ------------------------------------------------------------ 
		// // 1. Load ExecutionResult // ------------------------------------------------------------ 
		result, err := s.executionResults.GetByID( txCtx, tenantID, executionResultID) 
		if err != nil { 
			return fmt.Errorf("get execution result: %w", err) 
		} 

		if result == nil { 
			return ErrExecutionResultNotFound 
		}
		 // Idempotent confirmation. 
		 if result.Status == model.ExecutionResultStatusConfirmed { 
				return nil 
			} 
			if result.Status != model.ExecutionResultStatusDraft {
				 return fmt.Errorf( "%w: status=%s", ErrExecutionResultInvalid, result.Status, ) 
			}
				  // ------------------------------------------------------------ // 2. Validate quantities // 
				  // ------------------------------------------------------------ 
			if result.ProducedQuantity < 0 || result.QualifiedQuantity < 0 || result.RejectedQuantity < 0 { 
				return fmt.Errorf( "%w: negative quantity", ErrExecutionResultInvalid, ) 
			} 

			if result.QualifiedQuantity+result.RejectedQuantity != result.ProducedQuantity { 
				return fmt.Errorf( "%w: produced=%d qualified=%d rejected=%d", ErrExecutionResultInvalid, result.ProducedQuantity, result.QualifiedQuantity, result.RejectedQuantity, ) 
			} // ------------------------------------------------------------ // 3. Load Execution // ------------------------------------------------------------ 

			exec, err := s.executions.GetExecutionByID( txCtx, tenantID, result.ExecutionID) 
			if err != nil { 
				return fmt.Errorf("get execution: %w", err) 
			} 
			if exec == nil {
				 return errors.New("execution not found")
			} 

			if exec.Status != model.ProductionExecutionStatusCompleted {
				 if exec.Status == model.ProductionExecutionStatusFailed {
					 return ErrExecutionAlreadyFailed 
				 } 
					return ErrExecutionNotCompleted 
				} 
				if exec.WorkOrderID != result.WorkOrderID {
					 return fmt.Errorf( "%w: execution work order mismatch", ErrExecutionResultInvalid, ) 
				} // ------------------------------------------------------------ // 4. Load WorkOrder // ------------------------------------------------------------ 
				workOrder, err := s.workOrders.GetByID( txCtx, tenantID, result.WorkOrderID, )
				if err != nil { 
					return fmt.Errorf("get work order: %w", err) 
				}
				
				if workOrder == nil { 
					return errors.New("work order not found") 
				}

				if workOrder.Status == model.WorkOrderStatusCancelled {
					 return errors.New("cannot confirm result for cancelled work order") 
				} // ------------------------------------------------------------ // 5. Validate cumulative quantity // ------------------------------------------------------------ 
				newCompletedQuantity := workOrder.CompletedQuantity + result.QualifiedQuantity 
				if newCompletedQuantity > workOrder.PlannedQuantity {
					 return ErrWorkOrderQuantityExceeded 
				} 
					   // ------------------------------------------------------------ // 6. Finalize qualified product identities // ------------------------------------------------------------ // // IMPORTANT: // // We intentionally do NOT create devices based only on // QualifiedQuantity. // // A quantity tells us "97 products passed". // It does not tell us which 97 physical products they are. // // The actual item-level identity comes from ExecutionOperation // results. // // This extraction will be implemented as the next step: // // ExecutionOperation.Result // ↓ // item_key // ↓ // serial_number // ↓ // Device // // For now we explicitly reject a non-zero qualified quantity // without item-level identity data instead of silently creating // incorrect Devices. //
				if result.QualifiedQuantity > 0 { 
					operations, err := s.executions.ListOperations( txCtx, exec.ID, )
					if err != nil { 
						return fmt.Errorf( "list execution operations: %w", err, ) 
					} 
					if err := validateQualifiedItems( operations, result.QualifiedQuantity); err != nil {
						 return err 
					}
					
					items := extractQualifiedItems(operations) 
					for _, item := range items {
						req := &device.CreateDeviceFromExecutionResultRequest{
							 ProductID: exec.ProductID, 
							 WorkOrderID: exec.WorkOrderID, 
							 ExecutionID: exec.ID, 
							 ExecutionResultID: result.ID, 
							 SerialNumber: item.SerialNumber, 
							 HardwareID: item.HardwareID,
							} 
						if _, err := s.devices.CreateFromExecutionResult( txCtx, req, ); err != nil {
							 return fmt.Errorf( "create device for item %s: %w", item.ItemKey, err) 
						} 
					} 
				} 
							 // ------------------------------------------------------------ // 7. Update WorkOrder // ------------------------------------------------------------
				workOrder.CompletedQuantity = newCompletedQuantity
				if newCompletedQuantity >= workOrder.PlannedQuantity {
					 now := time.Now()
					 workOrder.Status = model.WorkOrderStatusCompleted 
					 workOrder.CompletedAt = &now }
					  if err := s.workOrders.Update( txCtx, workOrder, ); err != nil { 
						return fmt.Errorf( "update work order: %w", err, )
					  } 
// ------------------------------------------------------------ // 8. Confirm ExecutionResult // ------------------------------------------------------------ 
					now := time.Now() 
					result.Status = model.ExecutionResultStatusConfirmed 
					result.ConfirmedAt = &now 
					if err := s.executionResults.Update( txCtx, result, ); err != nil { 
							return fmt.Errorf( "confirm execution result: %w", err, ) 
						} 
							return nil }) 
}
 type qualifiedItem struct { 
	ItemKey string 
	SerialNumber string
	HardwareID string 
} 
func validateQualifiedItems( operations []*model.ExecutionOperation, qualifiedQuantity int64, ) error { 
	if qualifiedQuantity == 0 { 
		return nil 
	}
	if len(operations) == 0 {
		 return fmt.Errorf( "%w: execution has no operations", ErrExecutionResultInvalid) 
	} 
	items := extractQualifiedItems(operations)
	if int64(len(items)) != qualifiedQuantity {
		 return fmt.Errorf( "%w: qualified quantity=%d but identity items=%d", ErrExecutionResultInvalid, qualifiedQuantity, len(items))
	} 
	return nil 
} 
						
func extractQualifiedItems( operations []*model.ExecutionOperation, ) []qualifiedItem { // TODO: // // Parse the standardized: // // { // "items": [ // { // "item_key": "000001", // "data": { // "serial_number": "SN000001", // "hardware_id": "HW000001" // } // } // ] // } // // from ExecutionOperation.Result. // // We deliberately leave this parser isolated from the // Executor implementations. 
							
						return nil
}





