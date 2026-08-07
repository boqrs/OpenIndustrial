package device

import "time"

// DeviceState represents the real-time state of a single data point for a device.
// This model is intended for use in a key-value store like Redis.
type DeviceState struct {
	DeviceID  string      `json:"deviceId"`
	PointID   string      `json:"pointId"`
	Value     interface{} `json:"value"`
	Quality   uint8       `json:"quality"`
	Timestamp time.Time   `json:"timestamp"`
}