package mqtt

import "crypto/tls"

type Config struct {
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
