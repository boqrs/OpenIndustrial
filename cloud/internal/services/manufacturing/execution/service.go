package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution/executors"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/google/uuid"
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
)

// -----------------------------------------------------------------------------
// Service implementation
// -----------------------------------------------------------------------------

type serviceImpl struct {
	repository   Repository
	workOrderSvc workorder.Service
	routingSvc   routing.Service
	executorRegistry *executors.OperationExecutorRegistry
}

// NewService creates the production execution service.
func NewService(
	repository Repository,
	workOrderSvc workorder.Service,
	routingService routing.Service,
	executorRegistry *executors.OperationExecutorRegistry,
) Service {

	return &serviceImpl{
		repository:   repository,
		workOrderSvc: workOrderSvc,
		routingSvc:   routingService,
		executorRegistry: executorRegistry,
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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, fmt.Errorf("tenant ID not found in context")
	}
	if req == nil {
		return nil, ErrInvalidExecutionState
	}
	if req.WorkOrderID == 0 {
		return nil, ErrWorkOrderNotFound
	}

	wo, err := s.workOrderSvc.GetByID(
		ctx,
		tenantID,
		req.WorkOrderID,
	)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}


	if wo.Status != model.WorkOrderStatusReleased &&
		wo.Status != model.WorkOrderStatusInProgress {

		return nil, ErrWorkOrderNotExecutable
	}

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
		TenantID:       tenantID,
		WorkOrderID:    wo.ID,
		ProductID:       wo.ProductID,
		RoutingID:      wo.RoutingID,
		RoutingVersion: rt.Version,
		DeviceID:       req.DeviceID,
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
		//TODO: routing需要增加参数
	   // parameters := append([]byte(nil), op.Parameters...)

		operations = append(
			operations,
			&model.ExecutionOperation{
				ExecutionID:        entity.ID,
				RoutingOperationID: &op.ID,
				Sequence:           op.Sequence,
				Code:               op.Code,
				Name:               op.Name,
				Description:        op.Description,
				//WorkstationID:      op.WorkstationID, // 快照 WorkstationID
				//Parameters:         op.Parameters,   // 快照 Parameters
				///Parameters:         parameters,
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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, fmt.Errorf("tenant ID not found in context")
	}

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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, fmt.Errorf("tenant ID not found in context")
	}
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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}
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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}
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
// An operation cannot start until its immediately preceding operation has
// completed or been skipped.
func (s *serviceImpl) StartOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
) error {

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}
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

	// Fetch all sibling operations to determine the correct execution order.
	operations, err := s.repository.ListOperations(ctx, executionID)
	if err != nil {
		return fmt.Errorf("failed to list execution operations for sequence validation: %w", err)
	}

	// Sort operations by sequence to establish the definitive order.
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Sequence < operations[j].Sequence
	})

	// Find the index of the current operation in the sorted list.
	currentIndex := -1
	for i, operation := range operations {
		if operation.ID == op.ID {
			currentIndex = i
			break
		}
	}

	// If the operation is not the first one in the sequence, check its predecessor.
	if currentIndex > 0 {
		previousOperation := operations[currentIndex-1]

		// A required preceding operation must be completed or skipped.
		if previousOperation.Status != model.ExecutionOperationStatusCompleted &&
			previousOperation.Status != model.ExecutionOperationStatusSkipped {
			return fmt.Errorf(
				"previous operation %d is not completed or skipped",
				previousOperation.ID,
			)
		}
	} else if currentIndex == -1 {
		// This should not happen if the operation was loaded correctly before.
		return fmt.Errorf("consistency error: current operation ID %d not found in its own execution %d", op.ID, executionID)
	}

		// -------------------------------------------------------------------------
	// 5. Resolve Executor
	// -------------------------------------------------------------------------

	executor, ok := s.executorRegistry.Get(op.Code)
	if !ok {
		return fmt.Errorf(
			"executor not found for operation code: %s",
			op.Code,
		)
	}

	// -------------------------------------------------------------------------
	// 6. Build OperationInput
	// -------------------------------------------------------------------------
	
	var parameters map[string]any
	if len(op.Parameters) > 0 {
		if err := json.Unmarshal(op.Parameters, &parameters); err != nil {
			return fmt.Errorf(
				"invalid operation parameters for operation %d: %w",
				op.ID,
				err,
			)
		}
	}

	input := &executors.OperationInput{
		ExecutionID:          exec.ID,
		ExecutionOperationID: op.ID,
		WorkOrderID:          exec.WorkOrderID,
		ProductID:            exec.ProductID,
		DeviceID:             exec.DeviceID,
		Parameters:           parameters,
	}

	// -------------------------------------------------------------------------
	// 7. Validate Operation with Executor
	// -------------------------------------------------------------------------

	if err := executor.Validate(ctx, input); err != nil {
		return fmt.Errorf(
			"operation %d validation failed: %w",
			op.ID,
			err,
		)
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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return  fmt.Errorf("tenant ID not found in context")
	}
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
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal operation result: %w", err)
	}
	op.Result = resultJSON


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

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return  fmt.Errorf("tenant ID not found in context")
	}
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

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal operation result: %w", err)
	}
	op.Result = resultJSON

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

	//tenantID := tenantIDFromContext(ctx)

	entities, err := s.repository.ListOperations(
		ctx,
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


func (s *serviceImpl) tryCompleteExecution(
	ctx context.Context,
	executionID uint,
) error {

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}
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

