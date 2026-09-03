package executors

import (
	"context"
	"errors"
	"fmt"
)


// --- Supporting Abstractions ---

// SerialNumberGenerator defines the interface for a service that can generate
// unique serial numbers based on product and other parameters.
type SerialNumberGenerator interface {
	Generate(
		ctx context.Context,
		productID uint,
		parameters map[string]any,
	) (string, error)
}

// SNGenerateExecutor is an OperationExecutor that generates a serial number.
type SNGenerateExecutor struct {
	generator SerialNumberGenerator
}

// NewSNGenerateExecutor creates a new SNGenerateExecutor.
func NewSNGenerateExecutor(generator SerialNumberGenerator) *SNGenerateExecutor {
	return &SNGenerateExecutor{
		generator: generator,
	}
}

// Type returns the executor's unique type code.
func (e *SNGenerateExecutor) Type() string {
	return OperationTypeSNGenerate
}

// Validate checks the input for the SN generation operation.
func (e *SNGenerateExecutor) Validate(ctx context.Context, input *OperationInput) error {
	if input == nil {
		return errors.New("operation input is nil")
	}
	if input.ProductID == 0 {
		return errors.New("product ID is required for serial number generation")
	}
	if e.generator == nil {
		return errors.New("serial number generator is not configured")
	}
	return nil
}

// Execute generates the serial number by calling the injected generator service.
func (e *SNGenerateExecutor) Execute(ctx context.Context, input *OperationInput) (*OperationOutput, error) {
	sn, err := e.generator.Generate(ctx, input.ProductID, input.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}
	return &OperationOutput{
		Result: map[string]any{
			"serial_number": sn,
		},
	}, nil
}