package aws

// Config contains AWS IoT Core MQTT configuration.
type Config struct {
	Endpoint string `mapstructure:"endpoint"json:"endpoint"yaml:":"endpoint"`

	Region string `mapstructure:"region" json:"region" yaml:":"region"`

	ClientID string `mapstructure:"client_id" json:"client_id" yaml:":"client_id"`

	CertificateFile string `mapstructure:"certificate_file" json:"certificate_file" yaml:":"certificate_file"`

	PrivateKeyFile string `mapstructure:"private_key_file" json:"private_key_file" yaml:":"private_key_file"`

	RootCAFile string `mapstructure:"root_ca_file" json:"root_ca_file" yaml:":"root_ca_file"`

	KeepAliveSeconds int `mapstructure:"keep_alive_seconds" json:"keep_alive_seconds" yaml:":"keep_alive_seconds"`

	ConnectTimeoutSeconds int `mapstructure:"connect_timeout_seconds" json:"connect_timeout_seconds" yaml:":"connect_timeout_seconds"`

	MaxReconnectSeconds int `mapstructure:"max_reconnect_seconds" json:"max_reconnect_seconds" yaml:":"max_reconnect_seconds"`

	CleanSession bool `mapstructure:"clean_session" json:"clean_session" yaml:":"clean_session"`
}
