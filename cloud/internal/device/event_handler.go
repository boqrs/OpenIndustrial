package device

import (
	"context"
	"fmt"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

const (
	DeviceTelemetryReceivedEventName = "device.telemetry.received"
)

// EventHandler handles device-related events.
type EventHandler struct {
	service Service
}

// NewEventHandler creates a new device event handler.
func NewEventHandler(service Service) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

// Register subscribes to relevant device events.
func (h *EventHandler) Register(bus event.Bus) {
	bus.Subscribe(DeviceTelemetryReceivedEventName, h)
}

// Handle processes incoming events.
func (h *EventHandler) Handle(e event.Event) error {
	switch e.Name {
	case DeviceTelemetryReceivedEventName:
		return h.handleDeviceTelemetry(e)
	}
	return nil
}

func (h *EventHandler) handleDeviceTelemetry(e event.Event) error {
	payload, ok := e.Data.(DeviceTelemetryReceivedEvent)
	if !ok {
		err := fmt.Errorf("invalid payload type for event %s", DeviceTelemetryReceivedEventName)
		fmt.Println(err)
		return err
	}

	err := h.service.RecordTelemetry(context.Background(), payload.DeviceID, payload.Timestamp, payload.Points)
	if err != nil {
		fmt.Printf("Error recording telemetry: %v\n", err)
		return err
	}

	return nil
}