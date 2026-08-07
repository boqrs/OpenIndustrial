package eventbus

// Handler defines the function signature for an event handler.
// It takes a DomainEvent and returns an error if processing fails.
type Handler func(event DomainEvent) error

// Bus defines the interface for the event bus.
// It provides methods to publish events and subscribe to them.
type Bus interface {
	// Publish sends an event to all registered handlers for its type.
	Publish(event DomainEvent) error

	// Subscribe registers a handler for a specific event type.
	Subscribe(eventType string, handler Handler)
}