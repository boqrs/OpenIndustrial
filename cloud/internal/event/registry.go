package event

import (
	"sync"
)

// Registry stores event handlers.
type Registry struct {
	mu       sync.RWMutex
	handlers map[Type][]Handler
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[Type][]Handler),
	}
}

// Register registers a handler for a specific event type.
func (r *Registry) Register(
	t Type,
	handler Handler,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[t] = append(
		r.handlers[t],
		handler,
	)
}

// Get retrieves all handlers for a specific event type.
func (r *Registry) Get(
	t Type,
) []Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent race conditions on the slice
	return append(
		[]Handler{},
		r.handlers[t]...,
	)
}