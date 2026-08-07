package iot

import (
	"context"
	"encoding/json"
	"fmt"
)

// MQTTHandler is responsible for handling raw MQTT messages and passing them to the service layer.
type MQTTHandler struct {
	service Service
}

// NewMQTTHandler creates a new MQTTHandler.
func NewMQTTHandler(service Service) *MQTTHandler {
	return &MQTTHandler{
		service: service,
	}
}

// Handle parses an MQTT message and delegates to the IoT service.
func (h *MQTTHandler) Handle(topic string, payload []byte) error {
	// Example topic: gateway/{gatewayId}/device/{deviceId}/event
	if _, err := ParseGatewayEventTopic(topic); err != nil {
		// For now, we just log and ignore unhandled topics.
		// In a real system, you might want metrics or different error handling.
		fmt.Printf("unhandled topic: %s, error: %v\n", topic, err)
		return nil
	}

	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal iot.Message: %w", err)
	}

	// The service layer will convert the message to a domain event and publish it.
	return h.service.HandleMessage(context.Background(), msg)
}