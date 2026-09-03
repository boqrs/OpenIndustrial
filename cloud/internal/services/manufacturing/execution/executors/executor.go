package executors

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

// OperationExecutor defines the interface for a pluggable business logic
// component that can be invoked by the execution engine.
type OperationExecutor interface {
	// Type returns the unique code that identifies this executor. This code
	// corresponds to the RoutingOperation.Code.
	Type() string

	// Validate checks if the input parameters are valid for this executor.
	Validate(ctx context.Context, input *OperationInput) error

	// Execute performs the business logic of the operation.
	Execute(ctx context.Context, input *OperationInput) (*OperationOutput, error)
}

// --- Executor Registry ---

// OperationExecutorRegistry holds a map of registered OperationExecutor
// instances, keyed by their type code.
type OperationExecutorRegistry struct {
	executors map[string]OperationExecutor
}

// NewOperationExecutorRegistry creates a new, empty executor registry.
func NewOperationExecutorRegistry() *OperationExecutorRegistry {
	return &OperationExecutorRegistry{
		executors: make(map[string]OperationExecutor),
	}
}

// Register adds an executor to the registry.
func (r *OperationExecutorRegistry) Register(executor OperationExecutor) {
	if executor == nil {
		return
	}
	if r.executors == nil {
		r.executors = make(map[string]OperationExecutor)
	}
	r.executors[executor.Type()] = executor
}

// Get retrieves an executor from the registry by its type code.
func (r *OperationExecutorRegistry) Get(operationType string) (OperationExecutor, bool) {
	if r == nil || r.executors == nil {
		return nil, false
	}
	executor, ok := r.executors[operationType]
	return executor, ok
}


// --- Registry Builder ---

// BuildOperationExecutorRegistry constructs the registry and registers all
// standard, first-party executors.
func BuildOperationExecutorRegistry(
	snGenerator SerialNumberGenerator,
	ca security.CertificateAuthority,
) *OperationExecutorRegistry {

	registry := NewOperationExecutorRegistry()

	registry.Register(
		NewSNGenerateExecutor(snGenerator),
	)

	registry.Register(
		NewCertificateIssueExecutor(ca),
	)

	return registry
}