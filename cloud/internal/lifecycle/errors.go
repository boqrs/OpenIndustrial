package lifecycle

import "errors"

var (
	// ErrInvalidStateTransition is returned when a state transition is not allowed.
	ErrInvalidStateTransition = errors.New("invalid state transition")

	// ErrLifecycleDefinitionNotFound is returned when a lifecycle definition cannot be found.
	ErrLifecycleDefinitionNotFound = errors.New("lifecycle definition not found")

	// ErrLifecycleInstanceNotFound is returned when a lifecycle instance cannot be found.
	ErrLifecycleInstanceNotFound = errors.New("lifecycle instance not found")

	// ErrInvalidInitialState is returned when the provided initial state is not valid for the definition.
	ErrInvalidInitialState = errors.New("invalid initial state for this lifecycle definition")
)