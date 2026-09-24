package iot

// MQTTAction defines an MQTT operation.
type MQTTAction string

const (
	MQTTActionPublish MQTTAction = "publish"

	MQTTActionSubscribe MQTTAction = "subscribe"
)

// AuthenticateMQTTRequest contains the certificate identity
// presented by the MQTT client.
//
// The certificate itself is validated by Security.
type AuthenticateMQTTRequest struct {
	CertificateFingerprint string `json:"certificate_fingerprint"`
}

// AuthorizeMQTTRequest contains an MQTT authorization request.
type AuthorizeMQTTRequest struct {
	ResourceID uint `json:"resource_id"`

	Action MQTTAction `json:"action"`

	Topic string `json:"topic"`
}

// CreateCommandRequest creates a command for a Device.
type CreateCommandRequest struct {
	DeviceID uint `json:"device_id"`

	Command string `json:"command"`

	Payload string `json:"payload"`
}

// CommandAckRequest represents a command acknowledgement received
// from a Device.
type CommandAckRequest struct {
	ResourceID uint `json:"-"`

	CommandID uint `json:"command_id"`

	Success bool `json:"success"`

	Error string `json:"error"`
}
