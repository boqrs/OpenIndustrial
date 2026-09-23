package iot

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// Repository defines persistence operations required by IoT.
//
// IoT does not have a separate persistence model. Runtime state is
// stored on the existing Device entity.
type Repository interface {
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
	) (
		*model.DeviceCommand,
		error,
	)

	UpdateCommand(
		ctx context.Context,
		cmd *model.DeviceCommand,
	) error

	ListDeviceCommands(
		ctx context.Context,
		deviceID uint,
	) (
		[]*model.DeviceCommand,
		error,
	)
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
	) (
		*CommandResponse,
		error,
	)

	GetCommand(
		ctx context.Context,
		id uint,
	) (
		*CommandResponse,
		error,
	)

	ListDeviceCommands(
		ctx context.Context,
		deviceID uint,
	) (
		[]*CommandResponse,
		error,
	)

	AcknowledgeCommand(
		ctx context.Context,
		req *CommandAckRequest,
	) error
}
