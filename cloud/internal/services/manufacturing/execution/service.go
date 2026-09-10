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
	repository       Repository
	workOrderSvc     workorder.Service
	routingSvc       routing.Service
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
		repository:       repository,
		workOrderSvc:     workOrderSvc,
		routingSvc:       routingService,
		executorRegistry: executorRegistry,
	}
}

func (s *serviceImpl) createExecution(
	ctx context.Context,
	req *CreateExecutionRequest,
	useTx bool,
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

	if req.ResourceID == 0 {
		return nil, fmt.Errorf("execution resource ID is required")
	}

	wo, err := s.workOrderSvc.GetByID(
		ctx,
		tenantID,
		req.WorkOrderID,
	)
	if err != nil {
		return nil, err
	}

	if wo == nil {
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
		return nil, err
	}

	if rt == nil {
		return nil, ErrRoutingNotFound
	}

	if rt.Status != model.RoutingStatusActive {
		return nil, ErrRoutingNotActive
	}

	if rt.ProductID != wo.ProductID {
		return nil, ErrRoutingProductMismatch
	}

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

	entity := &model.ProductionExecution{
		ResourceID:     req.ResourceID,
		TenantID:       tenantID,
		WorkOrderID:    wo.ID,
		ProductID:      wo.ProductID,
		RoutingID:      wo.RoutingID,
		RoutingVersion: rt.Version,
		Status:         model.ProductionExecutionStatusPending,
	}

	operations := make(
		[]*model.ExecutionOperation,
		0,
		len(routingOperations),
	)

	for _, op := range routingOperations {
		if op == nil {
			continue
		}

		parameters := append(
			[]byte(nil),
			op.Parameters...,
		)

		workstationID := op.WorkStationID

		operations = append(
			operations,
			&model.ExecutionOperation{
				RoutingOperationID: &op.ID,
				Sequence:           op.Sequence,
				Code:               op.Code,
				Name:               op.Name,
				Description:        op.Description,
				WorkstationID:      &workstationID,
				Parameters:         parameters,
				Status:             model.ExecutionOperationStatusPending,
			},
		)
	}

	if len(operations) == 0 {
		return nil, ErrRoutingHasNoOperations
	}

	if useTx {
		if err := s.repository.CreateExecutionTx(
			ctx,
			entity,
			operations,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to create execution: %w",
				err,
			)
		}
	} else {
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
	}

	return toExecutionResponse(entity), nil
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
	return s.createExecution(ctx, req, false)
}

func (s *serviceImpl) CreateExecutionTx(
	ctx context.Context,
	req *CreateExecutionRequest,
) (*ExecutionResponse, error) {
	return s.createExecution(ctx, req, true)
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

func (s *serviceImpl) CancelExecution(
	ctx context.Context,
	id uint,
) error {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}

	return withExecutionTransaction(ctx, func(txCtx context.Context) error {
		exec, err := s.repository.GetExecutionByIDForUpdateTx(
			txCtx,
			tenantID,
			id,
		)
		if err != nil {
			return err
		}

		if exec == nil {
			return ErrExecutionNotFound
		}

		if exec.Status != model.ProductionExecutionStatusPending &&
			exec.Status != model.ProductionExecutionStatusInProgress {
			return ErrInvalidExecutionState
		}

		exec.Status = model.ProductionExecutionStatusCancelled

		if err := s.repository.UpdateExecution(
			txCtx,
			exec,
		); err != nil {
			return fmt.Errorf(
				"failed to cancel execution: %w",
				err,
			)
		}

		return nil
	})
}

// StartOperation starts an execution operation.
//
// An operation can only start when:
//   - the execution is in progress;
//   - the operation is pending;
//   - the immediately preceding operation is completed or skipped;
//   - the corresponding executor exists;
//   - the executor accepts the operation parameters.
//
// The operation row is loaded with FOR UPDATE. The caller is expected to
// invoke this method inside a transaction.
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
	// 1. Lock and validate Execution
	// -------------------------------------------------------------------------

	exec, err := s.repository.GetExecutionByIDForUpdateTx(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return err
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

	// -------------------------------------------------------------------------
	// 2. Lock and load Operation
	// -------------------------------------------------------------------------

	op, err := s.repository.GetOperationForUpdateTx(
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

	operations, err := s.repository.ListOperations(
		ctx,
		executionID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to list execution operations for sequence validation: %w",
			err,
		)
	}

	sort.Slice(
		operations,
		func(i, j int) bool {
			return operations[i].Sequence < operations[j].Sequence
		},
	)

	currentIndex := -1

	for i, operation := range operations {
		if operation == nil {
			continue
		}

		if operation.ID == op.ID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return fmt.Errorf(
			"consistency error: current operation ID %d not found in execution %d",
			op.ID,
			executionID,
		)
	}

	if currentIndex > 0 {
		previousOperation := operations[currentIndex-1]

		if previousOperation == nil {
			return fmt.Errorf(
				"consistency error: previous operation is nil",
			)
		}

		if previousOperation.Status != model.ExecutionOperationStatusCompleted &&
			previousOperation.Status != model.ExecutionOperationStatusSkipped {
			return ErrPriorOperationIncomplete
		}
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
		if err := json.Unmarshal(
			op.Parameters,
			&parameters,
		); err != nil {
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
		Parameters:           parameters,
	}

	// -------------------------------------------------------------------------
	// 7. Validate with Executor
	// -------------------------------------------------------------------------

	if err := executor.Validate(
		ctx,
		input,
	); err != nil {
		return fmt.Errorf(
			"operation %d validation failed: %w",
			op.ID,
			err,
		)
	}

	// -------------------------------------------------------------------------
	// 8. Start Operation
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
// When all operations belonging to the execution are completed or skipped,
// the execution itself is automatically marked as completed.
//
// The caller is expected to invoke this method inside a transaction.
func (s *serviceImpl) CompleteOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
	result map[string]any,
) error {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}

	// -------------------------------------------------------------------------
	// 1. Lock and validate Execution
	// -------------------------------------------------------------------------

	exec, err := s.repository.GetExecutionByIDForUpdateTx(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return err
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

	// -------------------------------------------------------------------------
	// 2. Lock and load Operation
	// -------------------------------------------------------------------------

	op, err := s.repository.GetOperationForUpdateTx(
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

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	// -------------------------------------------------------------------------
	// 4. Persist Operation result
	// -------------------------------------------------------------------------

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf(
			"marshal operation result: %w",
			err,
		)
	}

	now := time.Now()

	op.Status = model.ExecutionOperationStatusCompleted
	op.CompletedAt = &now
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
		exec,
	); err != nil {
		return fmt.Errorf(
			"failed to complete execution: %w",
			err,
		)
	}

	return nil
}

// FailOperation marks an execution operation as failed and fails the whole
// execution.
//
// The caller is expected to invoke this method inside a transaction.
func (s *serviceImpl) FailOperation(
	ctx context.Context,
	executionID uint,
	operationID uint,
	result map[string]any,
) error {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return fmt.Errorf("tenant ID not found in context")
	}

	// -------------------------------------------------------------------------
	// 1. Lock and validate Execution
	// -------------------------------------------------------------------------

	exec, err := s.repository.GetExecutionByIDForUpdateTx(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return err
	}

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

	// -------------------------------------------------------------------------
	// 2. Lock and load Operation
	// -------------------------------------------------------------------------

	op, err := s.repository.GetOperationForUpdateTx(
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

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	// -------------------------------------------------------------------------
	// 4. Persist Operation failure result
	// -------------------------------------------------------------------------

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf(
			"marshal operation result: %w",
			err,
		)
	}

	now := time.Now()

	op.Status = model.ExecutionOperationStatusFailed
	op.CompletedAt = &now
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

	// -------------------------------------------------------------------------
	// 5. Fail Execution
	// -------------------------------------------------------------------------

	exec.Status = model.ProductionExecutionStatusFailed
	exec.CompletedAt = &now

	if err := s.repository.UpdateExecution(
		ctx,
		exec,
	); err != nil {
		return fmt.Errorf(
			"failed to fail execution: %w",
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
	exec *model.ProductionExecution,
) error {

	if exec == nil {
		return ErrExecutionNotFound
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return nil
	}

	operations, err := s.repository.ListOperations(
		ctx,
		exec.ID,
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
		ResourceID:     entity.ResourceID,
		TenantID:       entity.TenantID,
		WorkOrderID:    entity.WorkOrderID,
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
		ID:          entity.ID,
		ExecutionID: entity.ExecutionID,
		Sequence:    entity.Sequence,
		Status:      entity.Status,
		StartedAt:   entity.StartedAt,
		CompletedAt: entity.CompletedAt,
		Code:        entity.Code,
		Name:        entity.Name,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	if entity.WorkstationID != nil {
		response.WorkstationID = *entity.WorkstationID
	}

	if entity.RoutingOperationID != nil {
		response.RoutingOperationID = *entity.RoutingOperationID
	}

	return response
}
