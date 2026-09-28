package mqtt

import (
	"context"
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type client struct {
	client mqtt.Client
	config Config
}

func NewClient(
	config Config,
) (*client, error) {
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

	options := mqtt.NewClientOptions()

	options.AddBroker(
		config.BrokerURL,
	)

	options.SetClientID(
		config.ClientID,
	)

	options.SetCleanSession(
		config.CleanSession,
	)

	if config.Username != "" {
		options.SetUsername(
			config.Username,
		)
	}

	if config.Password != "" {
		options.SetPassword(
			config.Password,
		)
	}

	if config.KeepAliveSeconds > 0 {
		options.SetKeepAlive(
			time.Duration(
				config.KeepAliveSeconds,
			) * time.Second,
		)
	}

	if config.ConnectTimeoutSeconds > 0 {
		options.SetConnectTimeout(
			time.Duration(
				config.ConnectTimeoutSeconds,
			) * time.Second,
		)
	}

	if config.MaxReconnectSeconds > 0 {
		options.SetMaxReconnectInterval(
			time.Duration(
				config.MaxReconnectSeconds,
			) * time.Second,
		)
	}

	options.SetAutoReconnect(true)

	if config.TLSConfig != nil {
		options.SetTLSConfig(
			config.TLSConfig,
		)
	}

	mqttClient := mqtt.NewClient(options)

	return &client{
		client: mqttClient,
		config: config,
	}, nil
}

func (c *client) Connect(
	ctx context.Context,
) error {
	if c == nil || c.client == nil {
		return errors.New(
			"mqtt client is not configured",
		)
	}

	if c.client.IsConnected() {
		return nil
	}

	token := c.client.Connect()

	if !token.Wait() {
		return errors.New(
			"mqtt connect timeout",
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

	if !c.client.IsConnected() {
		return errors.New(
			"mqtt client failed to connect",
		)
	}

	return nil
}

func (c *client) Disconnect(
	ctx context.Context,
) error {
	if c == nil || c.client == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if !c.client.IsConnected() {
		return nil
	}

	c.client.Disconnect(250)

	return nil
}

func (c *client) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
	qos byte,
	retained bool,
) (string, error) {
	if c == nil || c.client == nil {
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

	return fmt.Sprintf(
		"%s-%d",
		c.config.ClientID,
		time.Now().UTC().UnixNano(),
	), nil
}

func (c *client) Subscribe(
	ctx context.Context,
	topic string,
	qos byte,
	handler func(
		ctx context.Context,
		topic string,
		payload []byte,
	),
) error {
	if c == nil || c.client == nil {
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
			_ mqtt.Client,
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

func (c *client) IsConnected() bool {
	if c == nil || c.client == nil {
		return false
	}

	return c.client.IsConnected()
}
