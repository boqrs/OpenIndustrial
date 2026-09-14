package application

import (
	"context"
	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
)

type Repository interface {
	CreateTx(
		ctx context.Context,
		entity *model.ExecutionResult,
	) error

	Create(
		ctx context.Context,
		entity *model.ExecutionResult,
	) error

	GetByIDForUpdateTx(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint,
	) (*model.ExecutionResult, error)

	GetByID(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint,
	) (*model.ExecutionResult, error)

	GetByExecutionID(
		ctx context.Context,
		tenantID uuid.UUID,
		executionID uint,
	) (*model.ExecutionResult, error)

	UpdateTx(
		ctx context.Context,
		entity *model.ExecutionResult,
	) error

	Update(
		ctx context.Context,
		entity *model.ExecutionResult,
	) error
}

type Service interface {
	CreateProductionExecution(ctx context.Context, workOrderID uint) (*execution.ExecutionResponse, error)
	ConfirmExecutionResult(ctx context.Context, executionResultID uint) error
	StartProductionExecution(ctx context.Context, executionID uint) error
	StartProductionOperation(ctx context.Context, executionID uint, operationID uint) error
	FailProductionOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error
	CompleteProductionOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error
	CancelProductionExecution(ctx context.Context, executionID uint) error
}
