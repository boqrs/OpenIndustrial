package application

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/execution"
)

type Service interface {
	CreateProductionExecution(ctx context.Context, workOrderID uint, deviceID *uint) (*execution.ExecutionResponse, error)
	ConfirmExecutionResult(ctx context.Context, executionResultID uint) error
}
