package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/protocol"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/provider"
)

const commandQoS byte = 1

const statusQoS byte = 1

// Provider implements the IoT MQTT infrastructure interfaces
// using AWS IoT Core.
type Provider struct {
	client provider.MQTTClient
}

// NewProvider creates an AWS IoT provider.
func NewProvider(
	client provider.MQTTClient,
) *Provider {
	return &Provider{
		client: client,
	}
}

// Connect connects to AWS IoT Core.
func (p *Provider) Connect(
	ctx context.Context,
) error {
	if p == nil ||
		p.client == nil {

		return errors.New(
			"aws iot mqtt client is not configured",
		)
	}

	return p.client.Connect(
		ctx,
	)
}

// Disconnect disconnects from AWS IoT Core.
func (p *Provider) Disconnect(
	ctx context.Context,
) error {
	if p == nil ||
		p.client == nil {

		return nil
	}

	return p.client.Disconnect(
		ctx,
	)
}

// IsConnected returns whether the MQTT connection is active.
func (p *Provider) IsConnected() bool {
	if p == nil ||
		p.client == nil {

		return false
	}

	return p.client.IsConnected()
}

// PublishCommand publishes a command to:
//
// devices/{resourceID}/command
func (p *Provider) PublishCommand(
	ctx context.Context,
	resourceID uint,
	commandID uint,
	command string,
	payload string,
) (string, error) {
	if resourceID == 0 {
		return "", protocol.ErrInvalidResourceID
	}

	if commandID == 0 {
		return "", errors.New(
			"invalid command id",
		)
	}

	if command == "" {
		return "", errors.New(
			"command is required",
		)
	}

	if p == nil ||
		p.client == nil {

		return "", errors.New(
			"aws iot mqtt client is not configured",
		)
	}

	message, err := json.Marshal(
		struct {
			CommandID uint   `json:"command_id"`
			Command   string `json:"command"`
			Payload   string `json:"payload"`
		}{
			CommandID: commandID,

			Command: command,

			Payload: payload,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"marshal mqtt command: %w",
			err,
		)
	}

	return p.client.Publish(
		ctx,

		protocol.DeviceCommandTopic(
			resourceID,
		),

		message,

		commandQoS,

		false,
	)
}

// SubscribeDeviceStatus subscribes Cloud to:
//
// devices/{resourceID}/status
//
// This is intentionally NOT SubscribeDeviceCommand.
//
// Cloud receives Device status.
// Cloud publishes Device commands.
func (p *Provider) SubscribeDeviceStatus(
	ctx context.Context,
	resourceID uint,
	handler provider.MessageHandler,
) error {
	if resourceID == 0 {
		return protocol.ErrInvalidResourceID
	}

	if handler == nil {
		return errors.New(
			"mqtt status handler is nil",
		)
	}

	if p == nil ||
		p.client == nil {

		return errors.New(
			"aws iot mqtt client is not configured",
		)
	}

	return p.client.Subscribe(
		ctx,

		protocol.DeviceStatusTopic(
			resourceID,
		),

		statusQoS,

		handler,
	)
}
