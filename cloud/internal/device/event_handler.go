package device

import (
	"context"
	"fmt"
	"reflect"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/pkg/event"
)

// EventHandler handles events related to the device domain.
type EventHandler struct {
	// Dependencies like a repository or service can be added here.
}

// NewEventHandler creates a new event handler for the device domain.
func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// RegisterSubscriptions registers all event handlers for this domain on the event bus.
func (h *EventHandler) RegisterSubscriptions(bus event.Bus) {
	// The event type string is derived from the event struct type.
	bus.Subscribe(reflect.TypeOf(&DeviceRegisteredEvent{}).String(), h.handleDeviceRegistered)
	bus.Subscribe(reflect.TypeOf(&DeviceTelemetryReceivedEvent{}).String(), h.handleDeviceTelemetry)
}

// handleDeviceRegistered is the specific handler for when a device is registered.
func (h *EventHandler) handleDeviceRegistered(ctx context.Context, evt *DeviceRegisteredEvent) error {
	fmt.Printf("Handling DeviceRegisteredEvent: DeviceID=%s, Name=%s\n", evt.DeviceID, evt.Name)
	// Here you could, for example, provision the device in a cloud IoT service.
	return nil
}

// handleDeviceTelemetry is the specific handler for when telemetry is received.
func (h *EventHandler) handleDeviceTelemetry(ctx context.Context, evt *DeviceTelemetryReceivedEvent) error {
	fmt.Printf("Handling DeviceTelemetryReceivedEvent: DeviceID=%s, Payload=%v\n", evt.DeviceID, evt.Payload)
	// Here you could, for example, save the telemetry data to a time-series database.
	return nil
}