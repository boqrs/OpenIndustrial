package iot

import "time"

// MQTTAuthenticationResponse represents the identity resolved from // an authenticated device certificate.
type MQTTAuthenticationResponse struct {
	Authenticated bool   `json:"authenticated"`
	ResourceID    uint   `json:"resource_id"`
	CertificateID string `json:"certificate_id"`
}

// DeviceTopicsResponse contains all MQTT topics belonging to a // Resource.
type DeviceTopicsResponse struct {
	ResourceID     uint   `json:"resource_id"`
	TelemetryTopic string `json:"telemetry_topic"`
	StatusTopic    string `json:"status_topic"`
	EventTopic     string `json:"event_topic"`
	CommandTopic   string `json:"command_topic"`
}

// DeviceRuntime represents the IoT runtime state persisted on Device.
type DeviceRuntime struct {
	DeviceID     uint
	ResourceID   uint
	Status       string
	LastOnlineAt *time.Time
}

// DeviceRuntimeResponse exposes Device runtime state through the // IoT service.
type DeviceRuntimeResponse struct {
	DeviceID     uint       `json:"device_id"`
	ResourceID   uint       `json:"resource_id"`
	Status       string     `json:"status"`
	LastOnlineAt *time.Time `json:"last_online_at,omitempty"`
}
