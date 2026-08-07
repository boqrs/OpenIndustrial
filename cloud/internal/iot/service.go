package iot

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
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
	eventName, err := mapMessageTypeToEventName(msg.Type)
	if err != nil {
		return err
	}

	e := event.Event{
		Name:        eventName,
		Aggregate:   "Device", // All IoT events are aggregated by Device for now
		AggregateID: msg.DeviceID,
		Source:      "gateway",
		Data:        msg,
		CreatedAt:   time.Now().UTC(),
	}

	// For gateway-level events, the AggregateID should be the GatewayID.
	if msg.Type == MessageGatewayOnline || msg.Type == MessageGatewayOffline {
		e.Aggregate = "Gateway"
		e.AggregateID = msg.GatewayID
	}

	return s.eventBus.Publish(e)
}

func mapMessageTypeToEventName(msgType string) (string, error) {
	switch msgType {
	case MessageGatewayOnline:
		return "iot.gateway.online", nil
	case MessageGatewayOffline:
		return "iot.gateway.offline", nil
	case MessagePointChanged:
		return "iot.device.point.changed", nil
	case MessageDeviceAlarm:
		return "iot.device.alarm", nil
	default:
		return "", fmt.Errorf("unknown message type: %s", msgType)
	}
}