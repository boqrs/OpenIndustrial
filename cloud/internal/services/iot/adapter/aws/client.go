package aws

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/provider"
)

// NewClient creates a generic MQTT client configured for AWS IoT Core.
//
// AWS IoT Core MQTT endpoint:
//
//	ssl://<endpoint>:8883
//
// Authentication:
//
//	X.509 client certificate + private key.
func NewClient(
	config Config,
) (provider.MQTTClient, error) {
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

		RootCAs: rootCAs,

		InsecureSkipVerify: false,
	}

	mqttConfig := provider.MQTTConfig{
		BrokerURL: fmt.Sprintf(
			"ssl://%s:8883",
			config.Endpoint,
		),

		ClientID: config.ClientID,

		KeepAliveSeconds: config.KeepAliveSeconds,

		ConnectTimeoutSeconds: config.ConnectTimeoutSeconds,

		MaxReconnectSeconds: config.MaxReconnectSeconds,

		CleanSession: config.CleanSession,

		TLSConfig: tlsConfig,
	}

	return provider.NewMQTTClient(
		mqttConfig,
	)
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

	data, err := os.ReadFile(
		rootCAFile,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"read aws iot root ca: %w",
			err,
		)
	}

	if ok := pool.AppendCertsFromPEM(
		data,
	); !ok {
		return nil, errors.New(
			"failed to append aws iot root ca",
		)
	}

	return pool, nil
}
