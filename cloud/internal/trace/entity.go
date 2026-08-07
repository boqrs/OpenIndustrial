package trace

import (
	"encoding/json"
	"time"
)

// TraceRecord is a single, immutable entry in a product's lifecycle history.
// It captures a significant event, providing the "who, what, where, when, and why"
// for complete traceability.
type TraceRecord struct {
	ID        string `json:"id"`
	// ProductID links the record to the specific product instance (e.g., by serial number).
	ProductID string `json:"product_id"`
	// EventType describes the nature of the event (e.g., "assembled", "tested", "shipped").
	EventType string `json:"event_type"`
	// ResourceID identifies the resource involved (e.g., the station, machine, or production line).
	ResourceID string `json:"resource_id,omitempty"`
	// OperatorID identifies the person or system that performed the action.
	OperatorID string `json:"operator_id,omitempty"`
	// Data contains context-specific information about the event, such as test results,
	// material batch numbers, or configuration parameters.
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}