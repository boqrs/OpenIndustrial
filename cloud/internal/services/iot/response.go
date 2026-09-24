package iot

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// MQTTAuthenticationResponse represents the Device identity resolved
// from an authenticated certificate.
type MQTTAuthenticationResponse struct {
	Authenticated bool `json:"authenticated"`

	ResourceID uint `json:"resource_id"`

	CertificateID string `json:"certificate_id"`
}

// DeviceTopicsResponse contains the MQTT topics belonging to a Device.
type DeviceTopicsResponse struct {
	ResourceID uint `json:"resource_id"`

	StatusTopic string `json:"status_topic"`

	CommandTopic string `json:"command_topic"`
}

// DeviceRuntime represents runtime information of a Device.
type DeviceRuntime struct {
	DeviceID uint

	ResourceID uint

	Status string

	LastOnlineAt *time.Time
}

// DeviceRuntimeResponse exposes Device runtime information.
type DeviceRuntimeResponse struct {
	DeviceID uint `json:"device_id"`

	ResourceID uint `json:"resource_id"`

	Status string `json:"status"`

	LastOnlineAt *time.Time `json:"last_online_at,omitempty"`
}

// CommandResponse represents a Device command.
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
