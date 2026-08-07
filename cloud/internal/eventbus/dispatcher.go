package eventbus

import (
	"sync"
)

// Dispatcher is responsible for dispatching events to registered handlers.
// It holds the registry of handlers and manages the dispatching logic.
type Dispatcher struct {
	registry *Registry
	mu       sync.RWMutex
}

// NewDispatcher creates a new Dispatcher.
func NewDispatcher(registry *Registry) *Dispatcher {
	return &Dispatcher{
		registry: registry,
	}
}

// Dispatch finds all handlers for the given event and invokes them.
// For simplicity, this is a synchronous implementation.
// A production system might use goroutines for concurrent handling.
func (d *Dispatcher) Dispatch(event DomainEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers := d.registry.GetHandlers(event.Type)
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			// In a real system, you'd need a strategy for handling errors.
			// e.g., logging, retries, dead-letter queue.
			return err
		}
	}
	return nil
}