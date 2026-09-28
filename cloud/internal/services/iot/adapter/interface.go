package adapter

import "context"

// MessageHandler handles a Device -> Cloud message.
type MessageHandler func(
	ctx context.Context,
	topic string,
	payload []byte,
)

// DeviceMessageAdapter is the infrastructure abstraction used by IoT.
//
// IoT business services do not care whether messages are delivered through:
//
//   - a self-hosted MQTT broker
//   - AWS IoT Core
//
// Both implementations must satisfy this interface.
type DeviceMessageAdapter interface {
	Connect(ctx context.Context) error

	Disconnect(ctx context.Context) error

	PublishCommand(
		ctx context.Context,
		resourceID uint,
		commandID uint,
		command string,
		payload string,
	) (messageID string, err error)

	SubscribeDeviceStatus(
		ctx context.Context,
		resourceID uint,
		handler MessageHandler,
	) error

	IsConnected() bool
}
