package eventbus

import (
	"sync"
)

// Event represents a system event.
type Event struct {
	Topic string      // The topic of the event (e.g., "device.created", "point.updated")
	Data  interface{} // The actual data associated with the event
}

// EventHandler is a function that handles a specific event.
type EventHandler func(event Event)

// EventBus defines the interface for publishing and subscribing to events.
type EventBus interface {
	// Publish sends an event to all subscribed handlers for the given topic.
	Publish(event Event)
	// Subscribe registers an EventHandler for a specific topic.
	// Returns a function that can be called to unsubscribe the handler.
	Subscribe(topic string, handler EventHandler) func()
}

// NewEventBus creates a new in-memory EventBus.
func NewEventBus() EventBus {
	return &inMemoryEventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// inMemoryEventBus is a simple, concurrency-safe in-memory implementation of the EventBus interface.
type inMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventHandler
}

// Publish sends an event to all subscribed handlers for the given topic.
func (eb *inMemoryEventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := eb.subscribers[event.Topic]
	eb.mu.RUnlock()

	// Execute handlers in separate goroutines to avoid blocking the publisher
	for _, handler := range handlers {
		go handler(event)
	}
}

// Subscribe registers an EventHandler for a specific topic.
// Returns a function that can be called to unsubscribe the handler.
func (eb *inMemoryEventBus) Subscribe(topic string, handler EventHandler) func() {
	eb.mu.Lock()
	eb.subscribers[topic] = append(eb.subscribers[topic], handler)
	eb.mu.Unlock()

	// Return an unsubscribe function
	return func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()

		if handlers, ok := eb.subscribers[topic]; ok {
			for i, h := range handlers {
				if &h == &handler { // Compare handler addresses to find the specific one
					eb.subscribers[topic] = append(handlers[:i], handlers[i+1:]...)
					break
				}
			}
		}
	}
}