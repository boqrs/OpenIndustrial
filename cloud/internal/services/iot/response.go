package iot

import (
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"time"
)

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

type CommandResponse struct {
	ID uint `json:"id"`

	DeviceID uint `json:"device_id"`

	Command string `json:"command"`

	Payload string `json:"payload"`

	Status string `json:"status"`

	MessageID string `json:"message_id"`

	ErrorMessage string `json:"error_message"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	SentAt *time.Time `json:"sent_at"`

	AcknowledgedAt *time.Time `json:"acknowledged_at"`
}

func commandToResponse(
	cmd *model.DeviceCommand,
) *CommandResponse {

	if cmd == nil {

		return nil

	}

	return &CommandResponse{

		ID: cmd.ID,

		DeviceID: cmd.DeviceID,

		Command: cmd.Command,

		Payload: cmd.Payload,

		Status: string(cmd.Status),

		MessageID: cmd.MessageID,

		ErrorMessage: cmd.ErrorMessage,

		CreatedAt: cmd.CreatedAt,

		UpdatedAt: cmd.UpdatedAt,

		SentAt: cmd.SentAt,

		AcknowledgedAt: cmd.AcknowledgedAt,
	}

}
