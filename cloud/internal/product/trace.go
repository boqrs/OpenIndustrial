package product

import (
	"encoding/json"
	"time"
)

// ProductionEvent is a record of something that happened to a product instance at a specific station.
// This is the core of the product timeline.
type ProductionEvent struct {
	ID                string
	ProductInstanceID string
	StationID         string // The resource ID of the station where the event occurred.
	EventType         string // e.g., "assembly.finished", "test.started", "test.passed"
	OperatorID        string // The user ID of the person who triggered the event.
	Data              json.RawMessage
	CreatedAt         time.Time
}

// QualityRecord stores detailed data from a quality check or test.
type QualityRecord struct {
	ID                string
	ProductInstanceID string
	StationID         string
	Result            string // e.g., "pass", "fail"
	Data              json.RawMessage
	CreatedAt         time.Time
}