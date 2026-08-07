package eventbus

import "sync"

// Registry holds the mapping between event types and their handlers.
type Registry struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string][]Handler),
	}
}

// Add registers a handler for a given event type.
func (r *Registry) Add(eventType string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

// GetHandlers returns all handlers for a given event type.
func (r *Registry) GetHandlers(eventType string) []Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.handlers[eventType]
}