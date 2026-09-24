package aws

import (
	"os"
	"strconv"
	"strings"
)

// Config contains AWS IoT Core MQTT configuration.
type Config struct {
	Endpoint string

	Region string

	ClientID string

	CertificateFile string

	PrivateKeyFile string

	RootCAFile string

	KeepAliveSeconds int

	ConnectTimeoutSeconds int

	MaxReconnectSeconds int

	CleanSession bool
}

// LoadConfig loads AWS IoT configuration from environment variables.
//
// Required:
//
//	AWS_IOT_ENDPOINT
//	AWS_REGION
//	AWS_IOT_CLIENT_ID
//	AWS_IOT_CERTIFICATE_FILE
//	AWS_IOT_PRIVATE_KEY_FILE
//
// Optional:
//
//	AWS_IOT_ROOT_CA_FILE
//	AWS_IOT_KEEP_ALIVE
//	AWS_IOT_CONNECT_TIMEOUT
//	AWS_IOT_MAX_RECONNECT
//	AWS_IOT_CLEAN_SESSION
func LoadConfig() Config {
	return Config{
		Endpoint: os.Getenv(
			"AWS_IOT_ENDPOINT",
		),

		Region: os.Getenv(
			"AWS_REGION",
		),

		ClientID: os.Getenv(
			"AWS_IOT_CLIENT_ID",
		),

		CertificateFile: os.Getenv(
			"AWS_IOT_CERTIFICATE_FILE",
		),

		PrivateKeyFile: os.Getenv(
			"AWS_IOT_PRIVATE_KEY_FILE",
		),

		RootCAFile: os.Getenv(
			"AWS_IOT_ROOT_CA_FILE",
		),

		KeepAliveSeconds: envInt(
			"AWS_IOT_KEEP_ALIVE",
			60,
		),

		ConnectTimeoutSeconds: envInt(
			"AWS_IOT_CONNECT_TIMEOUT",
			30,
		),

		MaxReconnectSeconds: envInt(
			"AWS_IOT_MAX_RECONNECT",
			30,
		),

		CleanSession: envBool(
			"AWS_IOT_CLEAN_SESSION",
			false,
		),
	}
}

func envInt(
	name string,
	defaultValue int,
) int {
	value := strings.TrimSpace(
		os.Getenv(name),
	)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(
		value,
	)
	if err != nil {
		return defaultValue
	}

	return result
}

func envBool(
	name string,
	defaultValue bool,
) bool {
	value := strings.TrimSpace(
		os.Getenv(name),
	)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(
		value,
	)
	if err != nil {
		return defaultValue
	}

	return result
}
