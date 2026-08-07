package device

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

// DeviceTelemetryReceivedEvent is the event payload for when new telemetry data is received.
type DeviceTelemetryReceivedEvent struct {
	DeviceID  string
	Timestamp time.Time
	Points    map[string]interface{}
}

// ToDomainEvent converts the specific event to a generic event.Event for the bus.
func (dtre DeviceTelemetryReceivedEvent) ToDomainEvent() event.Event {
	return event.Event{
		Name: DeviceTelemetryReceivedEventName,
		Data: dtre,
	}
}