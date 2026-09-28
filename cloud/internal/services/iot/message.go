package iot

import (
	"context"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/protocol"
)

// MQTTMessageHandler handles MQTT messages received from devices.
type MQTTMessageHandler struct {
	service Service
}

// NewMQTTMessageHandler creates a MQTT message handler.
func NewMQTTMessageHandler(service Service) *MQTTMessageHandler {
	return &MQTTMessageHandler{
		service: service,
	}
}

// HandleStatusMessage handles a raw MQTT status message.
//
// The topic must be:
//
//	devices/{resourceID}/status
func (s *service) HandleStatusMessage(
	ctx context.Context,
	resourceID uint,
	payload []byte,
) error {
	if resourceID == 0 {
		return ErrInvalidResourceID
	}

	if len(payload) == 0 {
		return errors.New("empty MQTT status payload")
	}

	message, err := protocol.ParseStatusMessage(payload)
	if err != nil {
		return err
	}

	switch message.Type {
	case protocol.StatusMessageHeartbeat:
		_, err := s.DeviceHeartbeat(ctx, resourceID)
		return err

	case protocol.StatusMessageOnline:
		_, err := s.DeviceOnline(ctx, resourceID)
		return err

	case protocol.StatusMessageOffline:
		_, err := s.DeviceOffline(ctx, resourceID)
		return err

	case protocol.StatusMessageCommandAck:
		success := false

		if message.Success != nil {
			success = *message.Success
		}

		return s.AcknowledgeCommand(
			ctx,
			&CommandAckRequest{
				ResourceID: resourceID,
				CommandID:  message.CommandID,
				Success:    success,
				Error:      message.Error,
			},
		)

	case protocol.StatusMessageEvent:
		// Event bus integration will be connected here later.
		return nil

	default:
		return fmt.Errorf(
			"unsupported status message type: %q",
			message.Type,
		)
	}
}

// HandleStatus is a convenience wrapper for handling a device status
// payload without checking the MQTT topic.
func (h *MQTTMessageHandler) HandleStatus(
	ctx context.Context,
	resourceID uint,
	payload []byte,
) error {
	if h == nil || h.service == nil {
		return errors.New("iot message handler is not configured")
	}

	return h.service.HandleStatusMessage(ctx, resourceID, payload)
}

// GetDeviceTopics returns the MQTT topics belonging to a device.
//
// A device has exactly two MQTT business topics:
//
//	devices/{resourceID}/status
//	devices/{resourceID}/command
func (s *service) GetDeviceTopics(
	ctx context.Context,
	resourceID uint,
) (*DeviceTopicsResponse, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New("iot repository is not configured")
	}

	if _, err := s.repo.GetDeviceByResourceID(ctx, resourceID); err != nil {
		return nil, err
	}

	topics, err := protocol.GetDeviceTopics(resourceID)
	if err != nil {
		return nil, err
	}

	return &DeviceTopicsResponse{
		ResourceID:   topics.ResourceID,
		StatusTopic:  topics.StatusTopic,
		CommandTopic: topics.CommandTopic,
	}, nil
}

// DeviceOnline marks a device as online.
func (s *service) DeviceOnline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New("iot repository is not configured")
	}

	return s.repo.SetOnline(ctx, resourceID)
}

// DeviceOffline marks a device as offline.
func (s *service) DeviceOffline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New("iot repository is not configured")
	}

	return s.repo.SetOffline(ctx, resourceID)
}

// DeviceHeartbeat updates the device's last-online timestamp.
func (s *service) DeviceHeartbeat(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	if resourceID == 0 {
		return nil, ErrInvalidResourceID
	}

	if s.repo == nil {
		return nil, errors.New("iot repository is not configured")
	}

	return s.repo.UpdateLastOnline(ctx, resourceID)
}
