package iot

import "context"

// Consumer represents a message consumer from a broker like MQTT.
// It uses a handler to process the messages.
type Consumer struct {
	handler *MQTTHandler
}

// NewConsumer is a constructor for the Consumer.
// It takes a handler that will be used to process messages.
func NewConsumer(handler *MQTTHandler) *Consumer {
	return &Consumer{
		handler: handler,
	}
}

// Run starts the consumer loop.
// In a real implementation, this would connect to a message broker.
func (c *Consumer) Run(ctx context.Context) error {
	// This is a placeholder. A real implementation would listen for messages
	// and call c.handler.Handle(topic, payload) for each message.
	<-ctx.Done()
	return nil
}