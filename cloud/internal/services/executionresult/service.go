package executionresult

import (
	"context"
	"fmt"
	"time"
	"errors"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
	"github.com/google/uuid"
)

var (
	ErrInvalidRequest = errors.New(
		"invalid execution result request",
	)

	ErrExecutionResultNotFound = errors.New(
		"execution result not found",
	)

	ErrExecutionResultAlreadyExists = errors.New(
		"execution result already exists",
	)

	ErrExecutionResultAlreadyConfirmed = errors.New(
		"execution result already confirmed",
	)

	ErrExecutionNotCompleted = errors.New(
		"execution is not completed",
	)

	ErrInvalidQuantity = errors.New(
		"invalid production quantity",
	)

	ErrCompletedQuantityExceeded = errors.New(
		"completed quantity exceeds work order planned quantity",
	)

	ErrExecutionHasNoOperations = errors.New(
		"execution has no operations",
	)

	ErrQualifiedItemsMismatch = errors.New(
		"qualified quantity does not match qualified production items",
	)

	ErrDeviceFinalizationFailed = errors.New(
		"failed to finalize device",
	)

	ErrInvalidExecutionResultState = errors.New(
		"invalid execution result state",
	)
)

type serviceImpl struct {
	repository      Repository
	executionService execution.Service
}

func NewService(
	repository Repository,
	executionService execution.Service,
) Service {
	return &serviceImpl{
		repository:       repository,
		executionService: executionService,
	}
}

func (s *serviceImpl) Create(
	ctx context.Context,
	tenantID uuid.UUID,
	req *CreateRequest,
) (*Response, error) {

	if req == nil || req.ExecutionID == 0 {
		return nil, ErrInvalidRequest
	}

	if err := validateQuantity(req); err != nil {
		return nil, err
	}

	existing, err := s.repository.GetByExecutionID(
		ctx,
		tenantID,
		req.ExecutionID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get execution result: %w",
			err,
		)
	}

	if existing != nil {
		return nil, ErrExecutionResultAlreadyExists
	}

	exec, err := s.executionService.GetExecution(
		ctx,
		req.ExecutionID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get execution: %w",
			err,
		)
	}

	if exec == nil {
		return nil, ErrExecutionNotCompleted
	}

	if exec.Status != model.ProductionExecutionStatusCompleted {
		return nil, ErrExecutionNotCompleted
	}

	entity := &model.ExecutionResult{
		TenantID:          tenantID,
		ExecutionID:       req.ExecutionID,
		WorkOrderID:       exec.WorkOrderID,
		ProducedQuantity:  req.ProducedQuantity,
		QualifiedQuantity: req.QualifiedQuantity,
		RejectedQuantity:  req.RejectedQuantity,
		Status:            model.ExecutionResultStatusDraft,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf(
			"failed to create execution result: %w",
			err,
		)
	}

	return ToResponse(entity), nil
}

func (s *serviceImpl) GetByID(
	ctx context.Context,
	tenantID uuid.UUID,
	id uint,
) (*Response, error) {

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get execution result: %w",
			err,
		)
	}

	if entity == nil {
		return nil, ErrExecutionResultNotFound
	}

	return ToResponse(entity), nil
}

func (s *serviceImpl) Cancel(
	ctx context.Context,
	tenantID uuid.UUID,
	id uint,
) error {

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get execution result: %w",
			err,
		)
	}

	if entity == nil {
		return ErrExecutionResultNotFound
	}

	if entity.Status != model.ExecutionResultStatusDraft {
		return ErrInvalidRequest
	}

	entity.Status = model.ExecutionResultStatusCancelled

	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf(
			"failed to cancel execution result: %w",
			err,
		)
	}

	return nil
}

func (s *serviceImpl) Confirm(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint)  error {
    result, err := s.repository.GetByID(ctx, tenantID, id)
    if err != nil {
        return fmt.Errorf("failed to get execution result: %w", err)
    }

    if result == nil {
        return ErrExecutionResultNotFound
    }

    // 幂等
    if result.Status == model.ExecutionResultStatusConfirmed {
        return nil
    }

    if result.Status != model.ExecutionResultStatusDraft {
        return ErrInvalidExecutionResultState
    }

    execution, err := s.executionService.GetExecution(ctx, result.ExecutionID)
    if err != nil {
        return fmt.Errorf("failed to get execution: %w", err)
    }

    if execution.Status != model.ProductionExecutionStatusCompleted {
        return ErrExecutionNotCompleted
    }

    // if err := validateQuantity(
    //     result.ProducedQuantity,
    //     result.QualifiedQuantity,
    //     result.RejectedQuantity,
    // ); err != nil {
    //     return err
    // }

    // 这里开始处理 Qualified Items
    //
    // ExecutionOperation.Result
    //        ↓
    //   extract items
    //        ↓
    // validate identity
    //        ↓
    // Device.FinalizeFromExecution
    //
    // 然后：
    //
    // WorkOrder.CompletedQuantity += QualifiedQuantity
    //
    // 最后：
    // result.Status = Confirmed

    return nil
}

func validateQuantity(req *CreateRequest) error {
	if req.ProducedQuantity < 0 ||
		req.QualifiedQuantity < 0 ||
		req.RejectedQuantity < 0 {
		return ErrInvalidQuantity
	}

	if req.QualifiedQuantity+req.RejectedQuantity >
		req.ProducedQuantity {
		return ErrInvalidQuantity
	}

	return nil
}

func now() time.Time {
	return time.Now()
}

