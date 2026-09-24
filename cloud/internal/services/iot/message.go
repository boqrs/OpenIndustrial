package iot

import (
	"context"
	"encoding/json"
)

// MQTTMessageHandler converts broker messages into IoT domain actions.
//
// MQTT transport is intentionally separated from IoT business logic.
// The broker adapter should call this handler after topic routing.
type MQTTMessageHandler struct {
	service Service
}

func NewMQTTMessageHandler(service Service) *MQTTMessageHandler {
	return &MQTTMessageHandler{
		service: service,
	}
}

// HandleHeartbeat handles:
//
// device/{resource_id}/heartbeat
//
// Heartbeat only updates runtime information. It does not create devices.
func (h *MQTTMessageHandler) HandleHeartbeat(
	ctx context.Context,
	resourceID uint,
	payload []byte,
) error {
	if _, err := h.service.DeviceHeartbeat(ctx, resourceID); err != nil {
		return err
	}

	return nil
}

// HandleOnline handles explicit connection events.
func (h *MQTTMessageHandler) HandleOnline(
	ctx context.Context,
	resourceID uint,
) error {
	_, err := h.service.DeviceOnline(ctx, resourceID)
	return err
}

// HandleOffline handles broker disconnect events.
func (h *MQTTMessageHandler) HandleOffline(
	ctx context.Context,
	resourceID uint,
) error {
	_, err := h.service.DeviceOffline(ctx, resourceID)
	return err
}

// CommandAckPayload is the wire representation sent by devices.
type CommandAckPayload struct {
	CommandID uint   `json:"command_id"`
	Success   bool   `json:"success"`
	Error     string `json:"error"`
}

// HandleCommandAck handles:
//
// device/{resource_id}/command/ack
func (h *MQTTMessageHandler) HandleCommandAck(
	ctx context.Context,
	payload []byte,
) error {
	var ack CommandAckPayload

	if err := json.Unmarshal(payload, &ack); err != nil {
		return err
	}

	return h.service.AcknowledgeCommand(
		ctx,
		&CommandAckRequest{
			CommandID: ack.CommandID,
			Success:   ack.Success,
			Error:     ack.Error,
		},
	)
}
