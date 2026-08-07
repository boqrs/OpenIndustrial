package event

import "context"

// Bus defines the interface for a generic, decoupled event bus.
// It is a shared infrastructure component, completely unaware of any specific business domain.
type Bus interface {
	// Publish sends an event. The event is passed as an interface{},
	// allowing any type of event to be transmitted.
	Publish(ctx context.Context, event interface{}) error

	// Subscribe registers a handler for a specific event type.
	// The eventType is a string representation of the event's type,
	// e.g., "*device.DeviceRegisteredEvent".
	// The handler must be a function, e.g., func(ctx context.Context, e *device.DeviceRegisteredEvent) error.
	Subscribe(eventType string, handler interface{}) error
}