package mqtt

import "crypto/tls"

type Config struct {
	BrokerURL string `mapstructure:"broker_url" json:"broker_url" yaml:":"broker_url"`

	ClientID string `mapstructure:"client_id" json:"client_id" yaml:":"client_id"`

	Username string `mapstructure:"username" json:"username" yaml:":"username"`

	Password string `mapstructure:"password" json:"password" yaml:":"password"`

	KeepAliveSeconds int `mapstructure:"keep_alive_seconds" json:"keep_alive_seconds" yaml:":"keep_alive_seconds"`

	ConnectTimeoutSeconds int `mapstructure:"connect_timeout_seconds" json:"connect_timeout_seconds" yaml:":"connect_timeout_seconds"`

	MaxReconnectSeconds int `mapstructure:"max_reconnect_seconds" json:"max_reconnect_seconds" yaml:":"max_reconnect_seconds"`

	CleanSession bool `mapstructure:"clean_session" json:"clean_session" yaml:":"clean_session"`

	TLSConfig *tls.Config `mapstructure:"tls_config" json:"tls_config" yaml:":"tls_config"`
}
