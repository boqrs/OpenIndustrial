package iot

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// Repository defines the persistence operations required by IoT.
//
// IoT does not own Device persistence. Runtime state is stored on
// the existing Device entity.
type Repository interface {
	GetDeviceByID(
		ctx context.Context,
		deviceID uint,
	) (*model.Device, error)

	GetDeviceByResourceID(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	SetOnline(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	SetOffline(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	UpdateLastOnline(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	CreateCommand(
		ctx context.Context,
		cmd *model.DeviceCommand,
	) error

	GetCommand(
		ctx context.Context,
		id uint,
	) (*model.DeviceCommand, error)

	UpdateCommand(
		ctx context.Context,
		cmd *model.DeviceCommand,
	) error

	ListDeviceCommands(
		ctx context.Context,
		deviceID uint,
	) ([]*model.DeviceCommand, error)
}

// MQTTCommandPublisher is the infrastructure abstraction used by
// IoT to publish commands to devices.
type MQTTCommandPublisher interface {
	PublishCommand(
		ctx context.Context,
		resourceID uint,
		commandID uint,
		command string,
		payload string,
	) (messageID string, err error)
}

// MQTTStatusSubscriber is the infrastructure abstraction used by
// IoT to receive Device status messages.
type MQTTStatusSubscriber interface {
	SubscribeDeviceStatus(
		ctx context.Context,
		resourceID uint,
		handler func(
			ctx context.Context,
			topic string,
			payload []byte,
		),
	) error
}

// Service defines the IoT business layer.
type Service interface {
	AuthenticateMQTT(
		ctx context.Context,
		req AuthenticateMQTTRequest,
	) (*MQTTAuthenticationResponse, error)

	AuthorizeMQTT(
		ctx context.Context,
		req AuthorizeMQTTRequest,
	) error

	GetDeviceTopics(
		ctx context.Context,
		resourceID uint,
	) (*DeviceTopicsResponse, error)

	DeviceOnline(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	DeviceOffline(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	DeviceHeartbeat(
		ctx context.Context,
		resourceID uint,
	) (*model.Device, error)

	CreateCommand(
		ctx context.Context,
		req *CreateCommandRequest,
	) (*CommandResponse, error)

	SendCommand(
		ctx context.Context,
		commandID uint,
	) (*CommandResponse, error)

	GetCommand(
		ctx context.Context,
		id uint,
	) (*CommandResponse, error)

	ListDeviceCommands(
		ctx context.Context,
		deviceID uint,
	) ([]*CommandResponse, error)

	AcknowledgeCommand(
		ctx context.Context,
		req *CommandAckRequest,
	) error

	HandleStatusMessage(
		ctx context.Context,
		resourceID uint,
		payload []byte,
	) error
}
