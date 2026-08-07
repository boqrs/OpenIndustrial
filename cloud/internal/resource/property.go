package resource

import "encoding/json"

// Property represents a dynamic key-value attribute for a Resource instance.
// This allows storing type-specific data without altering the core Resource schema.
type Property struct {
	// ResourceID links the property to its parent resource.
	ResourceID string `json:"resource_id"`
	// Key is the name of the attribute, e.g., "ip_address", "axis_count".
	Key string `json:"key"`
	// Value is the data associated with the key, stored flexibly as a JSON object.
	Value json.RawMessage `json:"value"`
}