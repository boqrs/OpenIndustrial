package execution

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
)

// -----------------------------------------------------------------------------
// Errors
// -----------------------------------------------------------------------------

var (
	ErrExecutionNotFound = errors.New(
		"execution not found",
	)

	ErrOperationNotFound = errors.New(
		"execution operation not found",
	)

	ErrInvalidExecutionState = errors.New(
		"operation not allowed in current execution state",
	)

	ErrInvalidOperationState = errors.New(
		"invalid operation state transition",
	)

	ErrWorkOrderNotFound = errors.New(
		"associated work order not found",
	)

	ErrWorkOrderNotExecutable = errors.New(
		"work order is not available for execution",
	)

	ErrWorkOrderMismatch = errors.New(
		"work order details do not match",
	)

	ErrRoutingNotFound = errors.New(
		"associated routing not found",
	)

	ErrRoutingNotActive = errors.New(
		"routing is not active",
	)

	ErrRoutingProductMismatch = errors.New(
		"routing product does not match work order product",
	)

	ErrRoutingHasNoOperations = errors.New(
		"routing has no operations",
	)

	ErrPriorOperationIncomplete = errors.New(
		"prior operation is not yet complete",
	)

	// ErrExecutionQuantityExceeded = errors.New(
	// 	"execution quantity exceeds remaining work order quantity",
	// )

	// ErrInvalidExecutionQuantity = errors.New(
	// 	"execution quantity must be greater than zero",
	// )
)

// -----------------------------------------------------------------------------
// Service implementation
// -----------------------------------------------------------------------------

type serviceImpl struct {
	repository   Repository
	workOrderSvc workorder.Service
	routingSvc   routing.Service
}

// NewService creates the production execution service.
func NewService(
	repository Repository,
	workOrderSvc workorder.Service,
	routingService routing.Service,
) Service {
	return &serviceImpl{
		repository:   repository,
		workOrderSvc: workOrderSvc,
		routingSvc:   routingService,
	}
}

// -----------------------------------------------------------------------------
// Execution
// -----------------------------------------------------------------------------

// CreateExecution creates one actual production execution from a WorkOrder.
//
// The caller only provides the WorkOrder. Product and Routing are derived from
// the WorkOrder to prevent the execution from becoming detached from the
// production task.
func (s *serviceImpl) CreateExecution(
	ctx context.Context,
	req *CreateExecutionRequest,
) (*ExecutionResponse, error) {

	tenantID := tenantIDFromContext(ctx)

	if req == nil {
		return nil, ErrInvalidExecutionState
	}

	if req.WorkOrderID == 0 {
		return nil, ErrWorkOrderNotFound
	}

	// if req.Quantity <= 0 {
	// 	return nil, ErrInvalidExecutionQuantity
	// }

	// -------------------------------------------------------------------------
	// 1. Load WorkOrder
	// -------------------------------------------------------------------------

	wo, err := s.workOrderSvc.GetByID(
		ctx,
		tenantID,
		req.WorkOrderID,
	)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	// -------------------------------------------------------------------------
	// 2. Validate WorkOrder state
	// -------------------------------------------------------------------------

	if wo.Status != model.WorkOrderStatusReleased &&
		wo.Status != model.WorkOrderStatusInProgress {

		return nil, ErrWorkOrderNotExecutable
	}

	// -------------------------------------------------------------------------
	// 3. Validate remaining quantity
	//
	// A WorkOrder may have multiple executions:
	//
	//     WO = 100
	//       Execution A = 30
	//       Execution B = 40
	//       remaining  = 30
	//
	// The next execution cannot request more than 30.
	// -------------------------------------------------------------------------

	// if err := s.validateExecutionQuantity(
	// 	ctx,
	// 	tenantID,
	// 	wo,
	// 	//req.Quantity,
	// ); err != nil {
	// 	return nil, err
	// }

	// -------------------------------------------------------------------------
	// 4. Load Routing
	//
	// Routing is derived from WorkOrder, not supplied by the caller.
	// -------------------------------------------------------------------------

	rt, err := s.routingSvc.GetRouting(
		ctx,
		wo.RoutingID,
	)
	if err != nil {
		return nil, ErrRoutingNotFound
	}

	// -------------------------------------------------------------------------
	// 5. Validate Routing
	// -------------------------------------------------------------------------

	if rt.Status != model.RoutingStatusActive {
		return nil, ErrRoutingNotActive
	}

	if rt.ProductID != wo.ProductID {
		return nil, ErrRoutingProductMismatch
	}

	// -------------------------------------------------------------------------
	// 6. Load Routing Operations
	// -------------------------------------------------------------------------

	routingOperations, err := s.routingSvc.ListOperations(
		ctx,
		wo.RoutingID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list routing operations: %w",
			err,
		)
	}

	if len(routingOperations) == 0 {
		return nil, ErrRoutingHasNoOperations
	}

	// -------------------------------------------------------------------------
	// 7. Create ProductionExecution
	// -------------------------------------------------------------------------

	entity := &model.ProductionExecution{
		//ResourceUUID:   uuid.New(),
		TenantID:       tenantID,
		WorkOrderID:    wo.ID,
		ProductID:       wo.ProductID,
		RoutingID:      wo.RoutingID,
		RoutingVersion: rt.Version,
		DeviceID:       req.DeviceID,
		//Quantity:       req.Quantity,
		Status:         model.ProductionExecutionStatusPending,
	}

	// -------------------------------------------------------------------------
	// 8. Create ExecutionOperation snapshots
	//
	// The execution operation is a snapshot of the routing operation at the
	// time execution is created.
	//
	// This is important because Routing is a definition while Execution is
	// historical production data.
	// -------------------------------------------------------------------------

	operations := make(
		[]*model.ExecutionOperation,
		0,
		len(routingOperations),
	)

	for _, op := range routingOperations {

		if op == nil {
			continue
		}

		operations = append(
			operations,
			&model.ExecutionOperation{
				ExecutionID:        entity.ID,
				RoutingOperationID: &op.ID,
				Sequence:           op.Sequence,
				Code:               op.Code,
				Name:               op.Name,
				Description:        op.Description,
				Status:             model.ExecutionOperationStatusPending,
			},
		)
	}

	if len(operations) == 0 {
		return nil, ErrRoutingHasNoOperations
	}

	// -------------------------------------------------------------------------
	// 9. Persist
	//
	// Repository owns persistence details. We deliberately do not introduce
	// transaction handling here.
	// -------------------------------------------------------------------------

	if err := s.repository.CreateExecution(
		ctx,
		entity,
		operations,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create execution: %w",
			err,
		)
	}

	return toExecutionResponse(entity), nil
}

// GetExecution returns one execution.
func (s *serviceImpl) GetExecution(
	ctx context.Context,
	id uint,
) (*ExecutionResponse, error) {

	tenantID := tenantIDFromContext(ctx)

	entity, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, ErrExecutionNotFound
	}

	return toExecutionResponse(entity), nil
}

// ListExecutions returns executions matching the supplied filters.
func (s *serviceImpl) ListExecutions(
	ctx context.Context,
	workOrderID *uint,
	deviceID *uint,
	status *model.ProductionExecutionStatus,
) ([]*ExecutionResponse, error) {

	tenantID := tenantIDFromContext(ctx)

	entities, err := s.repository.ListExecutions(
		ctx,
		tenantID,
		workOrderID,
		deviceID,
		status,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*ExecutionResponse,
		0,
		len(entities),
	)

	for _, entity := range entities {
		if entity == nil {
			continue
		}

		responses = append(
			responses,
			toExecutionResponse(entity),
		)
	}

	return responses, nil
}

// StartExecution starts an execution.
//
// If the associated WorkOrder is still Released, starting the first execution
// also starts the WorkOrder.
func (s *serviceImpl) StartExecution(
	ctx context.Context,
	id uint,
) error {

	tenantID := tenantIDFromContext(ctx)

	// -------------------------------------------------------------------------
	// 1. Load execution
	// -------------------------------------------------------------------------

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return ErrExecutionNotFound
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	// -------------------------------------------------------------------------
	// 2. Validate execution state
	// -------------------------------------------------------------------------

	if exec.Status != model.ProductionExecutionStatusPending {
		return ErrInvalidExecutionState
	}

	// -------------------------------------------------------------------------
	// 3. Load WorkOrder
	// -------------------------------------------------------------------------

	wo, err := s.workOrderSvc.GetByID(
		ctx,
		tenantID,
		exec.WorkOrderID,
	)
	if err != nil {
		return ErrWorkOrderNotFound
	}

	if wo.Status != model.WorkOrderStatusReleased &&
		wo.Status != model.WorkOrderStatusInProgress {

		return ErrWorkOrderNotExecutable
	}

	// -------------------------------------------------------------------------
	// 4. Revalidate Routing
	//
	// Execution can sit in Pending state for some time. Therefore we must not
	// assume that the Routing was still active when the execution was created.
	// -------------------------------------------------------------------------

	rt, err := s.routingSvc.GetRouting(
		ctx,
		exec.RoutingID,
	)
	if err != nil {
		return ErrRoutingNotFound
	}

	if rt.Status != model.RoutingStatusActive {
		return ErrRoutingNotActive
	}

	if rt.ProductID != exec.ProductID {
		return ErrRoutingProductMismatch
	}

	// -------------------------------------------------------------------------
	// 5. Start WorkOrder if necessary
	// -------------------------------------------------------------------------

	if wo.Status == model.WorkOrderStatusReleased {

		if err := s.workOrderSvc.Start(
			ctx,
			tenantID,
			exec.WorkOrderID,
		); err != nil {
			return fmt.Errorf(
				"failed to start associated work order: %w",
				err,
			)
		}
	}

	// -------------------------------------------------------------------------
	// 6. Start execution
	// -------------------------------------------------------------------------

	now := time.Now()

	exec.Status = model.ProductionExecutionStatusInProgress
	exec.StartedAt = &now

	if err := s.repository.UpdateExecution(
		ctx,
		exec,
	); err != nil {
		return fmt.Errorf(
			"failed to start execution: %w",
			err,
		)
	}

	return nil
}

// CancelExecution cancels an execution.
func (s *serviceImpl) CancelExecution(
	ctx context.Context,
	id uint,
) error {

	tenantID := tenantIDFromContext(ctx)

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return ErrExecutionNotFound
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status == model.ProductionExecutionStatusCompleted ||
		exec.Status == model.ProductionExecutionStatusCancelled {

		return ErrInvalidExecutionState
	}

	exec.Status = model.ProductionExecutionStatusCancelled

	if err := s.repository.UpdateExecution(
		ctx,
		exec,
	); err != nil {
		return fmt.Errorf(
			"failed to cancel execution: %w",
			err,
		)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Operation
// -----------------------------------------------------------------------------

// StartOperation starts one execution operation.
//
// Operations must be executed sequentially:
//
//     10 → 20 → 30 → ...
//
// An operation cannot start until its immediately preceding operation has
// completed or been skipped.
func (s *serviceImpl) StartOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
) error {

	tenantID := tenantIDFromContext(ctx)

	// -------------------------------------------------------------------------
	// 1. Validate Execution
	// -------------------------------------------------------------------------

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return ErrExecutionNotFound
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

	// -------------------------------------------------------------------------
	// 2. Load Operation
	// -------------------------------------------------------------------------

	op, err := s.repository.GetOperation(
		ctx,
		tenantID,
		executionID,
		operationID,
	)
	if err != nil {
		return ErrOperationNotFound
	}

	if op == nil {
		return ErrOperationNotFound
	}

	// -------------------------------------------------------------------------
	// 3. Validate Operation state
	// -------------------------------------------------------------------------

	if op.Status != model.ExecutionOperationStatusPending {
		return ErrInvalidOperationState
	}

	// -------------------------------------------------------------------------
	// 4. Validate previous operation
	// -------------------------------------------------------------------------

	if op.Sequence > 1 {

		operations, err := s.repository.ListOperations(
			ctx,
			tenantID,
			executionID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to list execution operations: %w",
				err,
			)
		}

		previousFound := false

		for _, previous := range operations {

			if previous == nil {
				continue
			}

			if previous.Sequence != op.Sequence-1 {
				continue
			}

			previousFound = true

			if previous.Status != model.ExecutionOperationStatusCompleted &&
				previous.Status != model.ExecutionOperationStatusSkipped {

				return ErrPriorOperationIncomplete
			}

			break
		}

		if !previousFound {
			return ErrPriorOperationIncomplete
		}
	}

	// -------------------------------------------------------------------------
	// 5. Start Operation
	// -------------------------------------------------------------------------

	now := time.Now()

	op.Status = model.ExecutionOperationStatusInProgress
	op.StartedAt = &now

	if err := s.repository.UpdateOperation(
		ctx,
		op,
	); err != nil {
		return fmt.Errorf(
			"failed to start execution operation: %w",
			err,
		)
	}

	return nil
}

// CompleteOperation completes an execution operation.
//
// Once the last operation is completed, the Execution itself is automatically
// completed.
func (s *serviceImpl) CompleteOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
	result map[string]any,
) error {

	tenantID := tenantIDFromContext(ctx)

	// -------------------------------------------------------------------------
	// 1. Validate Execution
	// -------------------------------------------------------------------------

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return ErrExecutionNotFound
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

	// -------------------------------------------------------------------------
	// 2. Load Operation
	// -------------------------------------------------------------------------

	op, err := s.repository.GetOperation(
		ctx,
		tenantID,
		executionID,
		operationID,
	)
	if err != nil {
		return ErrOperationNotFound
	}

	if op == nil {
		return ErrOperationNotFound
	}

	// -------------------------------------------------------------------------
	// 3. Validate state
	// -------------------------------------------------------------------------

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	// -------------------------------------------------------------------------
	// 4. Complete Operation
	// -------------------------------------------------------------------------

	now := time.Now()

	op.Status = model.ExecutionOperationStatusCompleted
	op.CompletedAt = &now

	// The current model/repository does not currently expose a guaranteed
	// Result persistence field in this service contract. Keep the result at
	// the application boundary until that capability is intentionally added.
	_ = result

	if err := s.repository.UpdateOperation(
		ctx,
		op,
	); err != nil {
		return fmt.Errorf(
			"failed to complete execution operation: %w",
			err,
		)
	}

	// -------------------------------------------------------------------------
	// 5. Try to complete Execution
	// -------------------------------------------------------------------------

	if err := s.tryCompleteExecution(
		ctx,
		executionID,
	); err != nil {
		return fmt.Errorf(
			"failed to complete execution: %w",
			err,
		)
	}

	return nil
}

// FailOperation marks an execution operation as failed.
func (s *serviceImpl) FailOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
	result map[string]any,
) error {

	tenantID := tenantIDFromContext(ctx)

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return ErrExecutionNotFound
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

	op, err := s.repository.GetOperation(
		ctx,
		tenantID,
		executionID,
		operationID,
	)
	if err != nil {
		return ErrOperationNotFound
	}

	if op == nil {
		return ErrOperationNotFound
	}

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	now := time.Now()

	op.Status = model.ExecutionOperationStatusFailed
	op.CompletedAt = &now

	_ = result

	if err := s.repository.UpdateOperation(
		ctx,
		op,
	); err != nil {
		return fmt.Errorf(
			"failed to fail execution operation: %w",
			err,
		)
	}

	return nil
}

// ListOperations lists operations belonging to an execution.
func (s *serviceImpl) ListOperations(
	ctx context.Context,
	executionID uint,
) ([]*OperationResponse, error) {

	tenantID := tenantIDFromContext(ctx)

	entities, err := s.repository.ListOperations(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*OperationResponse,
		0,
		len(entities),
	)

	for _, entity := range entities {

		if entity == nil {
			continue
		}

		responses = append(
			responses,
			toOperationResponse(entity),
		)
	}

	return responses, nil
}

// -----------------------------------------------------------------------------
// Internal business logic
// -----------------------------------------------------------------------------

// validateExecutionQuantity verifies that the requested execution quantity
// does not exceed the remaining quantity of the WorkOrder.
//
// Only completed and in-progress executions consume the WorkOrder quantity.
// Cancelled and pending executions do not count as produced quantity.
// func (s *serviceImpl) validateExecutionQuantity(
// 	ctx context.Context,
// 	tenantID uuid.UUID,
// 	wo *workorder.Response,
// 	requestedQuantity int64,
// ) error {

// 	if requestedQuantity <= 0 {
// 		return ErrInvalidExecutionQuantity
// 	}

// 	executions, err := s.repository.ListExecutions(
// 		ctx,
// 		tenantID,
// 		&wo.ID,
// 		nil,
// 		nil,
// 	)
// 	if err != nil {
// 		return fmt.Errorf(
// 			"failed to list existing executions: %w",
// 			err,
// 		)
// 	}

// 	var consumedQuantity int64

// 	for _, execution := range executions {

// 		if execution == nil {
// 			continue
// 		}

// 		switch execution.Status {

// 		case model.ProductionExecutionStatusInProgress,
// 			model.ProductionExecutionStatusCompleted:

// 			//consumedQuantity += execution.Quantity
// 		}
// 	}

// 	remainingQuantity := int64(wo.PlannedQuantity) - consumedQuantity

// 	if requestedQuantity > remainingQuantity {
// 		return ErrExecutionQuantityExceeded
// 	}

// 	return nil
// }

// tryCompleteExecution checks whether every operation belonging to an
// execution has reached a terminal successful state.
//
// Terminal successful states:
//
//     Completed
//     Skipped
//
// Failed is deliberately not considered successful.
func (s *serviceImpl) tryCompleteExecution(
	ctx context.Context,
	executionID uint,
) error {

	tenantID := tenantIDFromContext(ctx)

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return ErrExecutionNotFound
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return nil
	}

	operations, err := s.repository.ListOperations(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to list execution operations: %w",
			err,
		)
	}

	if len(operations) == 0 {
		return ErrRoutingHasNoOperations
	}

	for _, op := range operations {

		if op == nil {
			continue
		}

		if op.Status != model.ExecutionOperationStatusCompleted &&
			op.Status != model.ExecutionOperationStatusSkipped {

			return nil
		}
	}

	// -------------------------------------------------------------------------
	// All operations completed.
	// -------------------------------------------------------------------------

	now := time.Now()

	exec.Status = model.ProductionExecutionStatusCompleted
	exec.CompletedAt = &now

	if err := s.repository.UpdateExecution(
		ctx,
		exec,
	); err != nil {
		return fmt.Errorf(
			"failed to complete execution: %w",
			err,
		)
	}

	// -------------------------------------------------------------------------
	// Execution completion may cause WorkOrder completion.
	// -------------------------------------------------------------------------

	if err := s.tryCompleteWorkOrder(
		ctx,
		exec.WorkOrderID,
	); err != nil {
		return fmt.Errorf(
			"failed to evaluate work order completion: %w",
			err,
		)
	}

	return nil
}

// tryCompleteWorkOrder checks whether the total completed execution quantity
// has reached the WorkOrder planned quantity.
//
// Example:
//
//     WorkOrder = 100
//
//     Execution A = 30 completed
//     Execution B = 40 completed
//     Execution C = 30 completed
//
//     => WorkOrder Completed
func (s *serviceImpl) tryCompleteWorkOrder(
	ctx context.Context,
	workOrderID uint,
) error {

	tenantID := tenantIDFromContext(ctx)

	wo, err := s.workOrderSvc.GetByID(
		ctx,
		tenantID,
		workOrderID,
	)
	if err != nil {
		return ErrWorkOrderNotFound
	}

	if wo.Status == model.WorkOrderStatusCompleted ||
		wo.Status == model.WorkOrderStatusCancelled {

		return nil
	}

	executions, err := s.repository.ListExecutions(
		ctx,
		tenantID,
		&workOrderID,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to list work order executions: %w",
			err,
		)
	}

	var completedQuantity int64

	for _, execution := range executions {

		if execution == nil {
			continue
		}

		if execution.Status != model.ProductionExecutionStatusCompleted {
			continue
		}

		//completedQuantity += execution.Quantity
	}

	if completedQuantity < int64(wo.PlannedQuantity) {
		return nil
	}

	// WorkOrder service owns the WorkOrder state transition.
	if err := s.workOrderSvc.Complete(
		ctx,
		tenantID,
		workOrderID,
	); err != nil {
		return fmt.Errorf(
			"failed to complete work order: %w",
			err,
		)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Mappers
// -----------------------------------------------------------------------------

func toExecutionResponse(
	entity *model.ProductionExecution,
) *ExecutionResponse {

	if entity == nil {
		return nil
	}

	return &ExecutionResponse{
		ID:             entity.ID,
		ResourceID:   entity.ResourceID,
		TenantID:       entity.TenantID,
		WorkOrderID:    entity.WorkOrderID,
		DeviceID:       entity.DeviceID,
		//Quantity:       entity.Quantity,
		Status:         entity.Status,
		StartedAt:      entity.StartedAt,
		CompletedAt:    entity.CompletedAt,
		ProductID:      entity.ProductID,
		RoutingID:      entity.RoutingID,
		RoutingVersion: entity.RoutingVersion,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

func toOperationResponse(
	entity *model.ExecutionOperation,
) *OperationResponse {

	if entity == nil {
		return nil
	}

	response := &OperationResponse{
		ID:                 entity.ID,
		ExecutionID:        entity.ExecutionID,
		Sequence:           entity.Sequence,
		Status:             entity.Status,
		StartedAt:          entity.StartedAt,
		CompletedAt:        entity.CompletedAt,
		Code:               entity.Code,
		Name:               entity.Name,
		Description:        entity.Description,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
	}

	if entity.WorkstationID != nil {
		response.WorkstationID = *entity.WorkstationID
	}

	if entity.RoutingOperationID != nil {
		response.RoutingOperationID = *entity.RoutingOperationID
	}

	return response
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func tenantIDFromContext(ctx context.Context) uuid.UUID {

	if ctx == nil {
		return uuid.Nil
	}

	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}