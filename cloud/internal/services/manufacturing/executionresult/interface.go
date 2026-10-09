package executionresult

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context, entity *model.ExecutionResult) error
	CreateTx(ctx context.Context, entity *model.ExecutionResult) error
	GetByID(ctx context.Context, tenantID uint, id uint) (*model.ExecutionResult, error)
	GetByExecutionID(ctx context.Context, tenantID uint, executionID uint) (*model.ExecutionResult, error)
	Update(ctx context.Context, entity *model.ExecutionResult) error
	UpdateTx(ctx context.Context, entity *model.ExecutionResult) error
	GetByIDForUpdateTx(ctx context.Context, tenantID uint, id uint) (*model.ExecutionResult, error)
}

type Service interface {
	CreateResult(ctx context.Context, req *CreateResultRequest) (*Response, error)
	GetResult(ctx context.Context, id uint) (*Response, error)
	GetResultByExecutionID(ctx context.Context, executionID uint) (*Response, error)
	CancelResult(ctx context.Context, id uint) error
}
