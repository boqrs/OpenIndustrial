package iot

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/device"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/pkg/event"
	"github.com/google/uuid"
)

// Service defines the business logic for handling IoT messages.
type Service interface {
	HandleMessage(ctx context.Context, msg Message) error
}

type service struct {
	eventBus event.Bus
}

// NewService creates a new IoT service.
func NewService(bus event.Bus) Service {
	return &service{
		eventBus: bus,
	}
}

// HandleMessage converts an IoT message into a domain event and publishes it.
func (s *service) HandleMessage(ctx context.Context, msg Message) error {
	// For now, we only handle telemetry data.
	// In the future, we can handle other message types here.
	if msg.Type != MessagePointChanged {
		return fmt.Errorf("unhandled message type: %s", msg.Type)
	}

	deviceID, err := uuid.Parse(msg.DeviceID)
	if err != nil {
		return fmt.Errorf("invalid device id: %w", err)
	}

	payload := map[string]interface{}{
		"pointId": msg.PointID,
		"value":   msg.Value,
		"quality": msg.Quality,
	}

	evt := &device.DeviceTelemetryReceivedEvent{
		DeviceID:  deviceID,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	}

	return s.eventBus.Publish(ctx, evt)
}