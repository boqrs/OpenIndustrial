package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/adapter"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/protocol"
)

const (
	commandQoS byte = 1
	statusQoS  byte = 1
)

type Adapter struct {
	client *client
}

func NewAdapter(
	config Config,
) (*Adapter, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		client: client,
	}, nil
}

func (a *Adapter) Connect(
	ctx context.Context,
) error {
	if a == nil || a.client == nil {
		return errors.New(
			"mqtt adapter is not configured",
		)
	}

	return a.client.Connect(ctx)
}

func (a *Adapter) Disconnect(
	ctx context.Context,
) error {
	if a == nil || a.client == nil {
		return nil
	}

	return a.client.Disconnect(ctx)
}

func (a *Adapter) IsConnected() bool {
	if a == nil || a.client == nil {
		return false
	}

	return a.client.IsConnected()
}

func (a *Adapter) PublishCommand(
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

	message, err := json.Marshal(
		struct {
			CommandID uint   `json:"command_id"`
			Command   string `json:"command"`
			Payload   string `json:"payload"`
		}{
			CommandID: commandID,
			Command:   command,
			Payload:   payload,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"marshal mqtt command: %w",
			err,
		)
	}

	return a.client.Publish(
		ctx,
		protocol.DeviceCommandTopic(resourceID),
		message,
		commandQoS,
		false,
	)
}

func (a *Adapter) SubscribeDeviceStatus(
	ctx context.Context,
	resourceID uint,
	handler adapter.MessageHandler,
) error {
	if resourceID == 0 {
		return protocol.ErrInvalidResourceID
	}

	if handler == nil {
		return errors.New(
			"mqtt status handler is nil",
		)
	}

	return a.client.Subscribe(
		ctx,
		protocol.DeviceStatusTopic(resourceID),
		statusQoS,
		handler,
	)
}

var _ adapter.DeviceMessageAdapter = (*Adapter)(nil)
