package execution

import (
	"context"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// --- Repository Interface ---
type Repository interface {
	// Execution methods
	CreateExecution(ctx context.Context, execution *model.ProductionExecution, operations []*model.ExecutionOperation) error
	GetExecutionByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.ProductionExecution, error)
	ListExecutions(ctx context.Context, tenantID uuid.UUID, workOrderID *uint, deviceID *uint, status *model.ProductionExecutionStatus) ([]*model.ProductionExecution, error)
	UpdateExecution(ctx context.Context, execution *model.ProductionExecution) error
	CountExecutions(ctx context.Context, tenantID uuid.UUID, workOrderID uint) (int64, error)

	// Operation methods
	GetOperation(ctx context.Context, executionID, operationID uint) (*model.ExecutionOperation, error)
	ListOperations(ctx context.Context,executionID uint) ([]*model.ExecutionOperation, error)
	UpdateOperation(ctx context.Context, operation *model.ExecutionOperation) error
	GetCurrentOperation(ctx context.Context, tenantID uuid.UUID, executionID uint) (*model.ExecutionOperation, error)
}

// --- Service Interface ---
type Service interface {
	// Execution methods
	CreateExecution(ctx context.Context, req *CreateExecutionRequest) (*ExecutionResponse, error)
	GetExecution(ctx context.Context, id uint) (*ExecutionResponse, error)
	ListExecutions(ctx context.Context, workOrderID *uint, deviceID *uint, status *model.ProductionExecutionStatus) ([]*ExecutionResponse, error)
	StartExecution(ctx context.Context, id uint) error
	CancelExecution(ctx context.Context, id uint) error

	// Operation methods
	StartOperation(ctx context.Context, executionID uint, operationID uint) error
	CompleteOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error
	FailOperation(ctx context.Context, executionID uint, operationID uint, result map[string]any) error
	ListOperations(ctx context.Context, executionID uint) ([]*OperationResponse, error)
}