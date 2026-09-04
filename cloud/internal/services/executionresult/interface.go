package executionresult

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, entity *model.ExecutionResult) error
	GetByID(ctx context.Context,tenantID uuid.UUID,id uint,) (*model.ExecutionResult, error)
	GetByExecutionID(ctx context.Context,tenantID uuid.UUID,executionID uint) (*model.ExecutionResult, error)
	Update(ctx context.Context, entity *model.ExecutionResult) error
}

type Service interface {
	Create(
		ctx context.Context,
		tenantID uuid.UUID,
		req *CreateRequest,
	) (*Response, error)

	GetByID(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint,
	) (*Response, error)

	Confirm(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint,
	) error

	Cancel(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint,
	) error
}