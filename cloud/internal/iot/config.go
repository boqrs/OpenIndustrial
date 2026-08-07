package iot

type MQTTConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
	Topics   []string
}