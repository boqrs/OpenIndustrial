package device

import (
	"time"

	"github.com/google/uuid"
)

// DeviceRegisteredEvent is published when a new device is successfully registered.
// This is the single source of truth for this event's definition.
type DeviceRegisteredEvent struct {
	DeviceID  uuid.UUID
	OrgID     uuid.UUID
	Name      string
	Timestamp time.Time
}

// DeviceTelemetryReceivedEvent is published when new telemetry data arrives from a device.
// This is the single source of truth for this event's definition.
type DeviceTelemetryReceivedEvent struct {
	DeviceID  uuid.UUID
	Payload   map[string]interface{}
	Timestamp time.Time
}