package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttClient struct {
	client mqtt.Client

	config MQTTConfig
}

// NewMQTTClient creates a generic MQTT client.
//
// This package does not know anything about AWS IoT Core.
func NewMQTTClient(
	config MQTTConfig,
) (MQTTClient, error) {
	if config.BrokerURL == "" {
		return nil, errors.New(
			"mqtt broker url is required",
		)
	}

	if config.ClientID == "" {
		return nil, errors.New(
			"mqtt client id is required",
		)
	}

	opts := mqtt.NewClientOptions()

	opts.AddBroker(
		config.BrokerURL,
	)

	opts.SetClientID(
		config.ClientID,
	)

	opts.SetCleanSession(
		config.CleanSession,
	)

	if config.Username != "" {
		opts.SetUsername(
			config.Username,
		)
	}

	if config.Password != "" {
		opts.SetPassword(
			config.Password,
		)
	}

	if config.KeepAliveSeconds > 0 {
		opts.SetKeepAlive(
			time.Duration(
				config.KeepAliveSeconds,
			) * time.Second,
		)
	}

	opts.SetAutoReconnect(true)

	if config.MaxReconnectSeconds > 0 {
		opts.SetMaxReconnectInterval(
			time.Duration(
				config.MaxReconnectSeconds,
			) * time.Second,
		)
	}

	if config.ConnectTimeoutSeconds > 0 {
		opts.SetConnectTimeout(
			time.Duration(
				config.ConnectTimeoutSeconds,
			) * time.Second,
		)
	}

	if config.TLSConfig != nil {
		opts.SetTLSConfig(
			config.TLSConfig,
		)
	}

	client := mqtt.NewClient(
		opts,
	)

	return &mqttClient{
		client: client,
		config: config,
	}, nil
}

func (c *mqttClient) Connect(
	ctx context.Context,
) error {
	if c == nil ||
		c.client == nil {

		return errors.New(
			"mqtt client is not configured",
		)
	}

	if c.client.IsConnected() {
		return nil
	}

	token := c.client.Connect()

	if token.Wait() &&
		token.Error() != nil {

		return token.Error()
	}

	select {

	case <-ctx.Done():

		return ctx.Err()

	default:

	}

	if !c.client.IsConnected() {
		return errors.New(
			"mqtt client failed to connect",
		)
	}

	return nil
}

func (c *mqttClient) Disconnect(
	ctx context.Context,
) error {
	if c == nil ||
		c.client == nil {

		return nil
	}

	if !c.client.IsConnected() {
		return nil
	}

	select {

	case <-ctx.Done():

		return ctx.Err()

	default:

	}

	c.client.Disconnect(
		250,
	)

	return nil
}

func (c *mqttClient) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
	qos byte,
	retained bool,
) (string, error) {
	if c == nil ||
		c.client == nil {

		return "", errors.New(
			"mqtt client is not configured",
		)
	}

	if !c.client.IsConnected() {
		return "", errors.New(
			"mqtt client is not connected",
		)
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	token := c.client.Publish(
		topic,
		qos,
		retained,
		payload,
	)

	if !token.Wait() {
		return "", errors.New(
			"mqtt publish timeout",
		)
	}

	if err := token.Error(); err != nil {
		return "", err
	}

	// Paho does not expose the MQTT packet identifier through
	// mqtt.Token. The IoT domain only needs a stable identifier
	// for the persisted command delivery record.
	messageID := fmt.Sprintf(
		"%s-%d",
		c.config.ClientID,
		time.Now().UTC().UnixNano(),
	)

	return messageID, nil
}

func (c *mqttClient) Subscribe(
	ctx context.Context,
	topic string,
	qos byte,
	handler MessageHandler,
) error {
	if c == nil ||
		c.client == nil {

		return errors.New(
			"mqtt client is not configured",
		)
	}

	if !c.client.IsConnected() {
		return errors.New(
			"mqtt client is not connected",
		)
	}

	if handler == nil {
		return errors.New(
			"mqtt message handler is nil",
		)
	}

	token := c.client.Subscribe(
		topic,
		qos,
		func(
			client mqtt.Client,
			message mqtt.Message,
		) {
			handler(
				context.Background(),
				message.Topic(),
				message.Payload(),
			)
		},
	)

	if !token.Wait() {
		return errors.New(
			"mqtt subscribe timeout",
		)
	}

	if err := token.Error(); err != nil {
		return err
	}

	select {

	case <-ctx.Done():

		return ctx.Err()

	default:

	}

	return nil
}

func (c *mqttClient) IsConnected() bool {
	if c == nil ||
		c.client == nil {

		return false
	}

	return c.client.IsConnected()
}
