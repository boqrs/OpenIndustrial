package telemetry

import (
	"encoding/json"
	"time"
)

// DeviceState provides a real-time snapshot of a device's status and latest metric values.
// This is the "hot" data, stored in a fast-access store like Redis.
type DeviceState struct {
	DeviceID string
	Status   string // e.g., "online", "offline"
	LastSeen time.Time
	Values   json.RawMessage // A JSON object of { "metric_code": value }
}