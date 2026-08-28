package execution

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/routing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/google/uuid"
)

// --- Errors ---
var (
	ErrExecutionNotFound        = errors.New("execution not found")
	ErrOperationNotFound        = errors.New("execution operation not found")
	ErrInvalidExecutionState    = errors.New("operation not allowed in current execution state")
	ErrWorkOrderNotFound        = errors.New("associated work order not found")
	ErrWorkOrderMismatch        = errors.New("work order details do not match")
	ErrInvalidOperationState    = errors.New("invalid operation state transition")
	ErrPriorOperationIncomplete = errors.New("prior operation is not yet complete")
)

// --- Service Implementation ---

type serviceImpl struct {
	repository   Repository
	workOrderSvc workorder.Service
	routingSvc   routing.Service
}

func NewService(repository Repository,workOrderSvc workorder.Service,routingService routing.Service) Service {
	return &serviceImpl{
		repository:   repository,
		workOrderSvc: workOrderSvc,
		routingSvc:   routingService,
	}
}

// --- Execution Methods ---

func (s *serviceImpl) CreateExecution(ctx context.Context,req *CreateExecutionRequest) (*ExecutionResponse, error) {

	tenantID := tenantIDFromContext(ctx)

	if req == nil {
		return nil, ErrInvalidExecutionState
	}

	if req.WorkOrderID == 0 || req.Quantity <= 0 {
		return nil, fmt.Errorf(
			"work order ID and a positive quantity are required",
		)
	}

	// --------------------------------------------------
	// 1. Load WorkOrder
	// --------------------------------------------------

	wo, err := s.workOrderSvc.GetByID(
		ctx,
		tenantID,
		req.WorkOrderID,
	)
	if err != nil {
		return nil, ErrWorkOrderNotFound
	}

	// --------------------------------------------------
	// 2. Validate WorkOrder
	// --------------------------------------------------

	if wo.Status != model.WorkOrderStatusReleased &&
		wo.Status != model.WorkOrderStatusInProgress {

		return nil, fmt.Errorf(
			"work order is not available for execution",
		)
	}

	if req.Quantity > int64(wo.PlannedQuantity) {
		return nil, fmt.Errorf(
			"execution quantity exceeds work order planned quantity",
		)
	}

	// --------------------------------------------------
	// 3. Load Routing
	// --------------------------------------------------

	routing, err := s.routingSvc.GetRouting(
		ctx,
		wo.RoutingID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get routing: %w",
			err,
		)
	}

	if routing.Status != "active" {
		return nil, fmt.Errorf(
			"work order routing is not active",
		)
	}

	// --------------------------------------------------
	// 4. Validate Product
	// --------------------------------------------------

	if routing.ProductID != wo.ProductID {
		return nil, ErrWorkOrderMismatch
	}

	// --------------------------------------------------
	// 5. Load Routing Operations
	// --------------------------------------------------

	routingOperations, err := s.routingSvc.ListOperations(
		ctx,
		wo.RoutingID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get routing operations: %w",
			err,
		)
	}

	if len(routingOperations) == 0 {
		return nil, fmt.Errorf(
			"cannot create execution for routing without operations",
		)
	}

	// --------------------------------------------------
	// 6. Create Execution
	// --------------------------------------------------

	entity := &model.ProductionExecution{
		ResourceUUID:   uuid.New(),
		TenantID:       tenantID,
		WorkOrderID:    wo.ID,
		ProductID:      wo.ProductID,
		RoutingID:      wo.RoutingID,
		RoutingVersion: routing.Version,
		DeviceID:       req.DeviceID,
		Status:         model.ProductionExecutionStatusPending,
	}

	// --------------------------------------------------
	// 7. Build Execution Operations
	// --------------------------------------------------

	operations := make(
		[]*model.ExecutionOperation,
		0,
		len(routingOperations),
	)

	for _, op := range routingOperations {
		operations = append(
			operations,
			&model.ExecutionOperation{
			//	ID:                 uuid.New(), // Assign a new UUID for the snapshot
			//	TenantID:           tenantID,
				ExecutionID:        entity.ID,
				RoutingOperationID: &op.ID,
				Sequence:           op.Sequence,
				Code:               op.Code,
				Name:               op.Name,
				Description:        op.Description,
				//:      op.WorkstationID,
				Status:             model.ExecutionOperationStatusPending,
			},
		)
	}

	// --------------------------------------------------
	// 8. Persist
	// --------------------------------------------------

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

func (s *serviceImpl) StartExecution(ctx context.Context,id uint) error {

	tenantID := tenantIDFromContext(ctx)

	exec, err := s.repository.GetExecutionByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return err
	}

	if exec.Status != model.ProductionExecutionStatusPending {
		return ErrInvalidExecutionState
	}

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

		return fmt.Errorf(
			"work order is not available for execution",
		)
	}

	// If WorkOrder is just 'Released', start it.
	if wo.Status == model.WorkOrderStatusReleased {
		if err := s.workOrderSvc.Start(
			ctx,
			tenantID,
			exec.WorkOrderID,
		); err != nil {
			return fmt.Errorf("failed to start associated work order: %w", err)
		}
	}

	now := time.Now()

	exec.Status = model.ProductionExecutionStatusInProgress
	exec.StartedAt = &now

	return s.repository.UpdateExecution(
		ctx,
		exec,
	)
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

func (s *serviceImpl) StartOperation(ctx context.Context, executionID uint, operationID uint) error {
	tenantID := tenantIDFromContext(ctx)

	exec, err := s.repository.GetExecutionByID(ctx, tenantID, executionID)
	if err != nil {
		return err
	}

	if exec.Status != model.ProductionExecutionStatusInProgress {
		return ErrInvalidExecutionState
	}

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
				if prevOp.Status != model.ExecutionOperationStatusCompleted && prevOp.Status != model.ExecutionOperationStatusSkipped {
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

func (s *serviceImpl) CompleteOperation(ctx context.Context,executionID uint,operationID uint,result map[string]any) error {

	tenantID := tenantIDFromContext(ctx)

	op, err := s.repository.GetOperation(
		ctx,
		tenantID,
		executionID,
		operationID,
	)
	if err != nil {
		return err
	}

	if op.Status != model.ExecutionOperationStatusInProgress {
		return ErrInvalidOperationState
	}

	now := time.Now()

	op.Status = model.ExecutionOperationStatusCompleted
	op.CompletedAt = &now

	if err := s.repository.UpdateOperation(ctx, op); err != nil {
		return err
	}

	return s.tryCompleteExecution(
		ctx,
		executionID,
	)
}

func (s *serviceImpl) FailOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error {
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
		return err
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
		return err
	}

	if len(operations) == 0 {
		return nil
	}

	for _, op := range operations {
		if op.Status != model.ExecutionOperationStatusCompleted &&
			op.Status != model.ExecutionOperationStatusSkipped {

			return nil
		}
	}

	now := time.Now()

	exec.Status = model.ProductionExecutionStatusCompleted
	exec.CompletedAt = &now

	return s.repository.UpdateExecution(
		ctx,
		exec,
	)
}

// --- Mappers ---

func toExecutionResponse(entity *model.ProductionExecution) *ExecutionResponse {
	if entity == nil {
		return nil
	}
	return &ExecutionResponse{
		ID:             entity.ID,
		ResourceUUID:   entity.ResourceUUID,
		TenantID:       entity.TenantID,
		WorkOrderID:    entity.WorkOrderID,
		DeviceID:       entity.DeviceID,
		Status:         entity.Status,
		StartedAt:      entity.StartedAt,
		CompletedAt:    entity.CompletedAt,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		ProductID:      entity.ProductID,
		RoutingID:      entity.RoutingID,
		RoutingVersion: entity.RoutingVersion,
	}
}

func toOperationResponse(entity *model.ExecutionOperation) *OperationResponse {
	if entity == nil {
		return nil
	}
	return &OperationResponse{
		ID:                 entity.ID,
		ExecutionID:        entity.ExecutionID,
		Sequence:           entity.Sequence,
		Status:             entity.Status,
		StartedAt:          entity.StartedAt,
		CompletedAt:        entity.CompletedAt,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		Code:               entity.Code,
		Name:               entity.Name,
		Description:        entity.Description,
		WorkstationID:      *entity.WorkstationID,
		RoutingOperationID: *entity.RoutingOperationID,
	}
}

// --- Helpers ---

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}