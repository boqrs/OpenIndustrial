package executionresult

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
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
	uow        postgres.UnitOfWork
	repository Repository
	executions execution.Service
}

func NewService(
	uow postgres.UnitOfWork,
	repository Repository,
	executions execution.Service,
) Service {
	return &service{
		uow:        uow,
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

	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil, errors.New("tenant id not found")
	}

	var result *Response

	err := s.uow.Execute(ctx, func(txCtx context.Context) error {
		// 1. 锁定执行记录。
		// 同一 Execution 的并发结果创建请求会在这里串行化。
		exec, err := s.executions.GetExecutionForUpdateTx(
			txCtx,
			req.ExecutionID,
		)
		if err != nil {
			if errors.Is(err, execution.ErrExecutionNotFound) {
				return execution.ErrExecutionNotFound
			}
			return fmt.Errorf("get execution for update: %w", err)
		}
		if exec == nil || exec.TenantID != tenantID {
			return execution.ErrExecutionNotFound
		}

		// 2. 执行必须已经完成。
		if exec.Status != model.ProductionExecutionStatusCompleted {
			return fmt.Errorf(
				"%w: execution status=%s",
				ErrInvalidResultState,
				exec.Status,
			)
		}

		// 3. 在持有 Execution 行锁的情况下检查重复结果。
		existing, err := s.repository.GetByExecutionID(
			txCtx,
			tenantID,
			req.ExecutionID,
		)
		if err == nil && existing != nil {
			return ErrResultAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check existing execution result: %w",
				err,
			)
		}

		// 4. 从已校验的 Execution 派生关联字段。
		entity := &model.ExecutionResult{
			TenantID:          tenantID,
			WorkOrderID:       exec.WorkOrderID,
			ExecutionID:       exec.ID,
			ProducedQuantity:  req.ProducedQuantity,
			QualifiedQuantity: req.QualifiedQuantity,
			RejectedQuantity:  req.RejectedQuantity,
			Status:            model.ExecutionResultStatusDraft,
		}

		if err := s.repository.CreateTx(txCtx, entity); err != nil {
			return fmt.Errorf("create execution result: %w", err)
		}

		result = toResponse(entity)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetResult(
	ctx context.Context,
	id uint,
) (*Response, error) {
	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil, errors.New("tenant id not found")
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
	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil, errors.New("tenant id not found")
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

func (s *service) CancelResult(
	ctx context.Context,
	id uint,
) error {
	if id == 0 {
		return ErrResultNotFound
	}

	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return errors.New("tenant id not found")
	}

	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		entity, err := s.repository.GetByIDForUpdateTx(
			txCtx,
			tenantID,
			id,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrResultNotFound
			}
			return fmt.Errorf("get execution result for cancellation: %w", err)
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

		if err := s.repository.UpdateTx(txCtx, entity); err != nil {
			return fmt.Errorf("cancel execution result: %w", err)
		}

		return nil
	})
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
