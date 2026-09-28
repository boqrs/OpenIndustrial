package aws

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type client struct {
	client mqtt.Client
}

func NewClient(
	config Config,
) (*client, error) {
	if config.Endpoint == "" {
		return nil, errors.New(
			"aws iot endpoint is required",
		)
	}

	if config.ClientID == "" {
		return nil, errors.New(
			"aws iot client id is required",
		)
	}

	if config.CertificateFile == "" {
		return nil, errors.New(
			"aws iot certificate file is required",
		)
	}

	if config.PrivateKeyFile == "" {
		return nil, errors.New(
			"aws iot private key file is required",
		)
	}

	certificate, err := tls.LoadX509KeyPair(
		config.CertificateFile,
		config.PrivateKeyFile,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"load aws iot certificate: %w",
			err,
		)
	}

	rootCAs, err := loadRootCAs(
		config.RootCAFile,
	)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: config.Endpoint,
		Certificates: []tls.Certificate{
			certificate,
		},
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
	}

	options := mqtt.NewClientOptions()

	options.AddBroker(
		fmt.Sprintf(
			"ssl://%s:8883",
			config.Endpoint,
		),
	)

	options.SetClientID(
		config.ClientID,
	)

	options.SetTLSConfig(
		tlsConfig,
	)

	options.SetAutoReconnect(true)

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

	options.SetCleanSession(
		config.CleanSession,
	)

	return &client{
		client: mqtt.NewClient(options),
	}, nil
}

func (c *client) Connect(
	ctx context.Context,
) error {
	if c == nil || c.client == nil {
		return errors.New(
			"aws iot mqtt client is not configured",
		)
	}

	if c.client.IsConnected() {
		return nil
	}

	token := c.client.Connect()

	if !token.Wait() {
		return errors.New(
			"aws iot mqtt connect timeout",
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

	if c.client.IsConnected() {
		c.client.Disconnect(250)
	}

	return nil
}

func (c *client) IsConnected() bool {
	return c != nil &&
		c.client != nil &&
		c.client.IsConnected()
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
			"aws iot mqtt client is not configured",
		)
	}

	if !c.client.IsConnected() {
		return "", errors.New(
			"aws iot mqtt client is not connected",
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
			"aws iot mqtt publish timeout",
		)
	}

	if err := token.Error(); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"aws-%d",
		time.Now().UTC().UnixNano(),
	), nil
}

func (c *client) Subscribe(
	ctx context.Context,
	topic string,
	qos byte,
	handler func(
		context.Context,
		string,
		[]byte,
	),
) error {
	if c == nil || c.client == nil {
		return errors.New(
			"aws iot mqtt client is not configured",
		)
	}

	if !c.client.IsConnected() {
		return errors.New(
			"aws iot mqtt client is not connected",
		)
	}

	if handler == nil {
		return errors.New(
			"aws iot mqtt handler is nil",
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
			"aws iot mqtt subscribe timeout",
		)
	}

	if err := token.Error(); err != nil {
		return err
	}

	return nil
}

func loadRootCAs(
	rootCAFile string,
) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		pool = x509.NewCertPool()
	}

	if rootCAFile == "" {
		return pool, nil
	}

	data, err := os.ReadFile(rootCAFile)
	if err != nil {
		return nil, fmt.Errorf(
			"read aws iot root ca: %w",
			err,
		)
	}

	if ok := pool.AppendCertsFromPEM(data); !ok {
		return nil, errors.New(
			"failed to append aws iot root ca",
		)
	}

	return pool, nil
}
