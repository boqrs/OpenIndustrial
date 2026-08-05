package profinet

import "time"

// Sample represents a single data point collected from a device.
// This is the unified data model for all drivers.
type Sample struct {
	PointID   string
	Value     interface{}
	Timestamp time.Time
	Quality   Quality
	Source    string // e.g., "profinet"
}

// Quality defines the quality of a data point.
type Quality string

const (
	QualityGood         Quality = "good"
	QualityBad          Quality = "bad"
	QualityUncertain    Quality = "uncertain"
	QualityDisconnected Quality = "disconnected"
)

// DeviceInfo represents a discovered PROFINET device on the network.
type DeviceInfo struct {
	ID          string
	MAC         string
	IP          string
	StationName string
	Vendor      string
	Product     string
}

// RecordRequest represents a request to read an acyclic record from a device.
type RecordRequest struct {
	DeviceID string
	Slot     uint16
	SubSlot  uint16
	Index    uint16
}

// RecordResponse represents the response for an acyclic record read.
type RecordResponse struct {
	DeviceID string
	Data     []byte
}