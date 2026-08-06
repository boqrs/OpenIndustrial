package lifecycle

import (
	"context"
)

// Lifecycle defines the interface for components that have a managed lifecycle.
// Components implementing this interface can be started and stopped by the runtime.
type Lifecycle interface {
	// Start initializes and starts the component.
	// It should block until the component is running or an error occurs.
	// The provided context can be used for cancellation signals during startup.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the component.
	// It should block until the component is stopped.
	// The provided context can be used for cancellation signals during shutdown.
	Stop(ctx context.Context) error
}