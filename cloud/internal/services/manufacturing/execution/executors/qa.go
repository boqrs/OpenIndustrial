package executors

import (
	"context"
	"errors"
	"fmt"
)

// FunctionTestExecutor represents a production functional test.
//
// The actual test may be performed by an external tester, PLC,
// workstation application, or other production equipment.
//
// The executor itself only defines and validates the standard
// result structure. The concrete production site integration can
// provide the actual test result later.
type FunctionTestExecutor struct {
}

func NewFunctionTestExecutor() *FunctionTestExecutor {
	return &FunctionTestExecutor{}
}

func (e *FunctionTestExecutor) Type() string {
	return OperationTypeFunctionTest
}

func (e *FunctionTestExecutor) Validate(
	ctx context.Context,
	input *OperationInput,
) error {
	if input == nil {
		return errors.New("operation input is nil")
	}

	if input.ExecutionID == 0 {
		return errors.New("execution ID is required")
	}

	if input.ExecutionOperationID == 0 {
		return errors.New("execution operation ID is required")
	}

	if input.ProductID == 0 {
		return errors.New("product ID is required")
	}

	return nil
}

func (e *FunctionTestExecutor) Execute(
	ctx context.Context,
	input *OperationInput,
) (*OperationOutput, error) {
	if err := e.Validate(ctx, input); err != nil {
		return nil, err
	}

	// FUNCTION_TEST is normally completed by an external production
	// system. Execute therefore does not invent a test result.
	//
	// The actual result should be supplied later through
	// CompleteOperation().
	return nil, fmt.Errorf(
		"function test requires external execution result",
	)
}
