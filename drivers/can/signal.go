package can

import "time"

// Signal defines how to decode a physical value from a CAN frame's data payload.
type Signal struct {
	ID        string    `json:"id"`
	FrameID   uint32    `json:"frameId"`
	StartBit  uint      `json:"startBit"`
	Length    uint      `json:"length"`
	ByteOrder ByteOrder `json:"byteOrder"`
	DataType  DataType  `json:"dataType"`
	Scale     float64   `json:"scale"`
	Offset    float64   `json:"offset"`
}

// SignalValue represents a decoded signal with its physical value.
// This is the unified data structure that will be used by the cache and publisher.
type SignalValue struct {
	ID        string
	Value     any
	Timestamp time.Time
}