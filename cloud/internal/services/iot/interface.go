package iot

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

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
