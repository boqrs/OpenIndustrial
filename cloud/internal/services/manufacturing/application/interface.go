package application

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
)

type Service interface {
	CreateProductionExecution(ctx context.Context, workOrderID uint) (*execution.ExecutionResponse, error)
	ConfirmExecutionResult(ctx context.Context, executionResultID uint) error
	StartProductionExecution(ctx context.Context, executionID uint) error
	StartProductionOperation(ctx context.Context, executionID uint, operationID uint) error
	FailProductionOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error
	CompleteProductionOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error
	CancelProductionExecution(ctx context.Context, executionID uint) error
}
