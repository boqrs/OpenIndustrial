package executionresult	

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"

)

var (
	ErrResultNotFound       = errors.New("execution result not found")
	ErrResultAlreadyExists  = errors.New("execution result already exists")
	ErrInvalidResultState   = errors.New("invalid execution result state")
	ErrInvalidQuantity      = errors.New("invalid execution result quantity")
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) CreateResult(
	ctx context.Context,
	req *CreateResultRequest,
) (*Response, error) {
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

	existing, err := s.repository.GetByExecutionID(
		ctx,
		tenantID,
		req.ExecutionID,
	)

	if err == nil && existing != nil {
		return nil, ErrResultAlreadyExists
	}

	entity := &model.ExecutionResult{
		TenantID:          tenantID,
		WorkOrderID:       req.WorkOrderID,
		ExecutionID:       req.ExecutionID,
		ProducedQuantity:  req.ProducedQuantity,
		QualifiedQuantity: req.QualifiedQuantity,
		RejectedQuantity:  req.RejectedQuantity,
		Status:            model.ExecutionResultStatusDraft,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("create execution result: %w", err)
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

func (s *service) ConfirmResult(
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

	// Confirmation is idempotent.
	if entity.Status == model.ExecutionResultStatusConfirmed {
		return nil
	}

	if entity.Status != model.ExecutionResultStatusDraft {
		return fmt.Errorf(
			"%w: status=%s",
			ErrInvalidResultState,
			entity.Status,
		)
	}

	if err := validateQuantities(
		entity.ProducedQuantity,
		entity.QualifiedQuantity,
		entity.RejectedQuantity,
	); err != nil {
		return err
	}

	now := time.Now()

	entity.Status = model.ExecutionResultStatusConfirmed
	entity.ConfirmedAt = &now

	if err := s.repository.Update(ctx, entity); err != nil {
		return fmt.Errorf("confirm execution result: %w", err)
	}

	return nil
}

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
