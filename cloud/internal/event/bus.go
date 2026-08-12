package event

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// HandlerFunc defines the signature for a function that can handle a specific event.
// It receives the event envelope and is expected to perform the business logic.
// Returning an error will signal that the event processing failed.
type HandlerFunc func(ctx context.Context, event *Envelope) error

// Bus is the interface for an event bus that allows subscribing to event types
// and dispatching events to the appropriate handlers.
type Bus interface {
	// Subscribe registers a handler for a given event type.
	Subscribe(eventType string, handler HandlerFunc)

	// Dispatch finds all handlers for the event's type and executes them.
	// If any handler returns an error, dispatching stops and the error is returned.
	Dispatch(ctx context.Context, event *Envelope) error
}

// InMemoryBus is a simple, thread-safe, in-memory implementation of the Bus interface.
type InMemoryBus struct {
	handlers map[string][]HandlerFunc
	mu       sync.RWMutex
}

// NewInMemoryBus creates a new instance of InMemoryBus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]HandlerFunc),
	}
}

// Subscribe registers a handler for a specific event type.
func (b *InMemoryBus) Subscribe(eventType string, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	log.Printf("INFO: Handler registered for event type '%s'", eventType)
}

// Dispatch finds and executes all handlers registered for the given event's type.
// Execution is sequential. If a handler fails, the process stops and returns the error.
func (b *InMemoryBus) Dispatch(ctx context.Context, event *Envelope) error {
	b.mu.RLock()
	handlers, ok := b.handlers[event.Type]
	b.mu.RUnlock()

	if !ok {
		// It's not necessarily an error if no one is listening, but it's good to know.
		log.Printf("WARN: No handlers registered for event type '%s'", event.Type)
		return nil
	}

	for _, handler := range handlers {
		// In a real-world scenario, you might want to add more observability here,
		// like metrics on handler execution time.
		if err := handler(ctx, event); err != nil {
			// The caller (e.g., the stream subscriber) can use this error to decide
			// not to ACK the message, allowing for a retry.
			return fmt.Errorf("handler for event type '%s' failed: %w", event.Type, err)
		}
	}

	return nil
}