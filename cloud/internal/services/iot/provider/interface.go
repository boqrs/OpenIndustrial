package provider

import (
	"context"
	"crypto/tls"
)

// MessageHandler handles an MQTT message.
//
// The provider is responsible only for receiving the MQTT message.
// Business interpretation belongs to IoT.
type MessageHandler func(
	ctx context.Context,
	topic string,
	payload []byte,
)

// MQTTClient is the generic MQTT transport abstraction.
type MQTTClient interface {
	Connect(
		ctx context.Context,
	) error

	Disconnect(
		ctx context.Context,
	) error

	Publish(
		ctx context.Context,
		topic string,
		payload []byte,
		qos byte,
		retained bool,
	) (messageID string, err error)

	Subscribe(
		ctx context.Context,
		topic string,
		qos byte,
		handler MessageHandler,
	) error

	IsConnected() bool
}

// MQTTConfig contains provider-independent MQTT configuration.
type MQTTConfig struct {
	BrokerURL string

	ClientID string

	Username string

	Password string

	KeepAliveSeconds int

	ConnectTimeoutSeconds int

	MaxReconnectSeconds int

	CleanSession bool

	TLSConfig *tls.Config
}

// MQTTCommandPublisher publishes IoT commands.
type MQTTCommandPublisher interface {
	PublishCommand(
		ctx context.Context,
		resourceID uint,
		commandID uint,
		command string,
		payload string,
	) (messageID string, err error)
}

// MQTTStatusSubscriber subscribes to Device status messages.
type MQTTStatusSubscriber interface {
	SubscribeDeviceStatus(
		ctx context.Context,
		resourceID uint,
		handler MessageHandler,
	) error
}
