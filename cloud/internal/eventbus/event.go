package eventbus

import "time"

// DomainEvent represents a significant change or fact that occurred within a specific domain.
// It is the fundamental communication block in an event-driven architecture.
type DomainEvent struct {
	// ID is the unique identifier for this specific event instance.
	ID string

	// Type is a string that uniquely identifies the kind of event, e.g., "workorder.completed".
	// It follows the convention: "domain.aggregate.action".
	Type string

	// AggregateID is the identifier of the root entity that this event pertains to.
	// For example, the WorkOrder ID or the Product SN.
	AggregateID string

	// AggregateType is the type of the root entity, e.g., "WorkOrder", "ProductInstance".
	AggregateType string

	// Data contains the payload of the event, serialized as JSON.
	// It provides the specific details of what happened.
	Data []byte

	// CreatedAt is the timestamp when the event was created.
	CreatedAt time.Time
}