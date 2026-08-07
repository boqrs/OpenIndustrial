package resource

import (
	"encoding/json"
	"time"
)

// Twin holds the last known state of a resource, acting as its digital counterpart.
// It is the core of the real-time digital twin functionality.
type Twin struct {
	// ResourceID is the unique identifier for the resource this twin represents.
	ResourceID string `json:"resource_id"`
	// State is a JSON object representing the last reported state of the resource.
	// e.g., {"status": "running", "speed": 3000, "alarm": false}
	State json.RawMessage `json:"state"`
	// Version is an atomically incrementing number to prevent race conditions.
	Version int64 `json:"version"`
	// UpdatedAt is the timestamp of the last state update.
	UpdatedAt time.Time `json:"updated_at"`
}