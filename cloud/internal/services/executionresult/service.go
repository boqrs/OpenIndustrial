package executionresult

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
)

var (
	ErrResultNotFound      = errors.New("execution result not found")
	ErrResultAlreadyExists = errors.New("execution result already exists")
	ErrInvalidResultState  = errors.New("invalid execution result state")
	ErrInvalidQuantity     = errors.New("invalid execution result quantity")
)

type service struct {
	repository Repository
	executions  execution.Service

}

func NewService(repository Repository, executions execution.Service,) Service {
	return &service{
		repository: repository,
		executions: executions,
	}
}

func (s *service) CreateResult(
    ctx context.Context,
    req *CreateResultRequest,
) (*Response, error) {
    if req == nil {
        return nil, errors.New("create execution result request is required")
    }

    if req.ExecutionID == 0 {
        return nil, errors.New("execution id is required")
    }

    if err := validateQuantities(
        req.ProducedQuantity,
        req.QualifiedQuantity,
        req.RejectedQuantity,
    ); err != nil {
        return nil, err
    }

    tenantID := pkg.TenantIDFromContext(ctx)
    if tenantID == uuid.Nil {
        return nil, errors.New("tenant id not found in context")
    }

    // Execution is the source of WorkOrderID.
    exec, err := s.executions.GetExecution(ctx, req.ExecutionID)
    if err != nil {
        return nil, fmt.Errorf("get execution: %w", err)
    }

    if exec == nil {
        return nil, errors.New("execution not found")
    }

    // A production result can only be created after the execution
    // has finished all routing operations.
    if exec.Status != model.ProductionExecutionStatusCompleted {
        return nil, fmt.Errorf(
            "%w: execution status=%s",
            ErrInvalidResultState,
            exec.Status,
        )
    }

    existing, err := s.repository.GetByExecutionID(
        ctx,
        tenantID,
        req.ExecutionID,
    )

    if err == nil && existing != nil {
        return nil, ErrResultAlreadyExists
    }

    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf(
            "check existing execution result: %w",
            err,
        )
    }

    entity := &model.ExecutionResult{
        TenantID:          tenantID,
        WorkOrderID:       exec.WorkOrderID,
        ExecutionID:       exec.ID,
        ProducedQuantity:  req.ProducedQuantity,
        QualifiedQuantity: req.QualifiedQuantity,
        RejectedQuantity:  req.RejectedQuantity,
        Status:            model.ExecutionResultStatusDraft,
    }

    if err := s.repository.Create(ctx, entity); err != nil {
        return nil, fmt.Errorf(
            "create execution result: %w",
            err,
        )
    }

    return toResponse(entity), nil
}

func (s *service) GetResult(
	ctx context.Context,
	id uint,
) (*Response, error) {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant id not found in context")
	}

	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, fmt.Errorf("get execution result: %w", err)
	}

	if entity == nil {
		return nil, ErrResultNotFound
	}

	return toResponse(entity), nil
}

func (s *service) GetResultByExecutionID(
	ctx context.Context,
	executionID uint,
) (*Response, error) {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant id not found in context")
	}

	entity, err := s.repository.GetByExecutionID(
		ctx,
		tenantID,
		executionID,
	)
	if err != nil {
		return nil, fmt.Errorf("get execution result: %w", err)
	}

	if entity == nil {
		return nil, ErrResultNotFound
	}

	return toResponse(entity), nil
}

// func (s *service) ConfirmResult(
// 	ctx context.Context,
// 	id uint,
// ) error {
// 	tenantID := pkg.TenantIDFromContext(ctx)
// 	if tenantID == uuid.Nil {
// 		return errors.New("tenant id not found in context")
// 	}

// 	entity, err := s.repository.GetByID(ctx, tenantID, id)
// 	if err != nil {
// 		return fmt.Errorf("get execution result: %w", err)
// 	}

// 	if entity == nil {
// 		return ErrResultNotFound
// 	}

// 	// Confirmation is idempotent.
// 	if entity.Status == model.ExecutionResultStatusConfirmed {
// 		return nil
// 	}

// 	if entity.Status != model.ExecutionResultStatusDraft {
// 		return fmt.Errorf(
// 			"%w: status=%s",
// 			ErrInvalidResultState,
// 			entity.Status,
// 		)
// 	}

// 	if err := validateQuantities(
// 		entity.ProducedQuantity,
// 		entity.QualifiedQuantity,
// 		entity.RejectedQuantity,
// 	); err != nil {
// 		return err
// 	}

// 	now := time.Now()

// 	entity.Status = model.ExecutionResultStatusConfirmed
// 	entity.ConfirmedAt = &now

// 	if err := s.repository.Update(ctx, entity); err != nil {
// 		return fmt.Errorf("confirm execution result: %w", err)
// 	}

// 	return nil
// }

func (s *service) CancelResult(
	ctx context.Context,
	id uint,
) error {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return errors.New("tenant id not found in context")
	}

	entity, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		return fmt.Errorf("get execution result: %w", err)
	}

	if entity == nil {
		return ErrResultNotFound
	}

	if entity.Status == model.ExecutionResultStatusCancelled {
		return nil
	}

	if entity.Status != model.ExecutionResultStatusDraft {
		return fmt.Errorf(
			"%w: status=%s",
			ErrInvalidResultState,
			entity.Status,
		)
	}

	entity.Status = model.ExecutionResultStatusCancelled

	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf("cancel execution result: %w", err)
	}

	return nil
}

func validateQuantities(
	produced int64,
	qualified int64,
	rejected int64,
) error {
	if produced < 0 || qualified < 0 || rejected < 0 {
		return ErrInvalidQuantity
	}

	if qualified+rejected != produced {
		return fmt.Errorf(
			"%w: produced=%d qualified=%d rejected=%d",
			ErrInvalidQuantity,
			produced,
			qualified,
			rejected,
		)
	}

	return nil
}
