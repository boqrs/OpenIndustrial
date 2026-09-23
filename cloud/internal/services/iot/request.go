package iot

// MQTTAction defines an MQTT operation.
type MQTTAction string

const (
	MQTTActionPublish   MQTTAction = "publish"
	MQTTActionSubscribe MQTTAction = "subscribe"
)

// AuthenticateMQTTRequest contains the certificate identity presented
// by the MQTT client.
//
// The certificate itself is validated by the Security service.
// IoT only uses the canonical certificate fingerprint to resolve
// ResourceID.
type AuthenticateMQTTRequest struct {
	CertificateFingerprint string `json:"certificate_fingerprint"`
}

// AuthorizeMQTTRequest contains an MQTT authorization request.
type AuthorizeMQTTRequest struct {
	ResourceID uint       `json:"resource_id"`
	Action     MQTTAction `json:"action"`
	Topic      string     `json:"topic"`
}

type CreateCommandRequest struct {
	DeviceID uint `json:"device_id"`

	Command string `json:"command"`

	Payload string `json:"payload"`
}

type CommandAckRequest struct {
	CommandID uint `json:"command_id"`

	Success bool `json:"success"`

	Error string `json:"error"`
}
