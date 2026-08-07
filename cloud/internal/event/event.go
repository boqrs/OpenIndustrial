package event

import (
	"time"
)

// Event represents a domain event.
//
// All modules communicate through events.
//
// Example:
//
// production.finished
// device.alarm
// user.created
//
type Event struct {
	// Unique event id
	ID string

	// Event name
	//
	// example:
	//
	// production.finished
	//
	Name string

	// Aggregate type
	//
	// Product
	// WorkOrder
	// Device
	//
	Aggregate string

	// Aggregate identifier
	AggregateID string

	// Event source service
	//
	// mes
	// gateway
	// product
	//
	Source string

	// Event payload
	Data any

	CreatedAt time.Time
}