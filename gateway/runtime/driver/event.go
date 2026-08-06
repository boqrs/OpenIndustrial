package driver

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/model/point"
)

// EventType defines the type of driver event.
type EventType uint8

const (
	// EventDataChanged indicates a point value changed.
	EventDataChanged EventType = iota

	// EventDeviceOnline indicates device connected.
	EventDeviceOnline

	// EventDeviceOffline indicates device disconnected.
	EventDeviceOffline

	// EventAlarm indicates an alarm event.
	EventAlarm
)

// Event is the unified data event produced by drivers.
//
// All protocol drivers should convert their native data
// into this structure before sending to runtime.
//
type Event struct {
	// Type of event.
	Type EventType

	// DriverID identifies the driver instance.
	//
	// Example:
	// modbus-meter-01
	//
	DriverID string

	// AssetID identifies the logical asset.
	//
	// Example:
	// battery-station-01
	//
	AssetID string

	// DeviceID identifies the physical device.
	//
	// Example:
	// meter-001
	//
	DeviceID string

	// PointID identifies the measurement point.
	//
	// Example:
	// voltage
	//
	PointID string

	// Value is the actual value.
	Value any

	// Quality indicates data quality.
	Quality point.Quality

	// Timestamp is the source timestamp.
	Timestamp time.Time

	// Metadata stores extra protocol information.
	//
	// Example:
	//
	// {
	//   "protocol":"modbus",
	//   "address":"40001"
	// }
	//
	Metadata map[string]any
}

func NewDataEvent(
	driverID string,
	deviceID string,
	pointID string,
	value any,
	quality point.Quality,
	timestamp time.Time,
) Event {
	return Event{
		Type:      EventDataChanged,
		DriverID:  driverID,
		DeviceID:  deviceID,
		PointID:   pointID,
		Value:     value,
		Quality:   quality,
		Timestamp: timestamp,
		Metadata:  make(map[string]any),
	}
}

// WithMetadata adds extra information.
func (e Event) WithMetadata(
	key string,
	value any,
) Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}
	e.Metadata[key] = value
	return e
}