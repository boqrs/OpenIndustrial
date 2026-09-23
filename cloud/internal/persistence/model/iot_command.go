package model

import (
	"time"

	"github.com/google/uuid"
)

type CommandStatus string

const (
	CommandStatusCreated CommandStatus = "created"

	CommandStatusSent CommandStatus = "sent"

	CommandStatusAcknowledged CommandStatus = "acknowledged"

	CommandStatusFailed CommandStatus = "failed"

	CommandStatusTimeout CommandStatus = "timeout"
)

type DeviceCommand struct {
	ID uint `gorm:"primaryKey"`

	// target device

	DeviceID uint `gorm:"not null;index"`

	// operator from JWT user_id

	OperatorID uuid.UUID `gorm:"type:uuid;not null;index"`
	// business command

	Command string `gorm:"size:100;not null"`

	// JSON payload

	Payload string `gorm:"type:jsonb"`

	Status CommandStatus `gorm:"size:50;not null;index"`

	// provider message id

	MessageID string `gorm:"size:255;index"`

	SentAt *time.Time

	AcknowledgedAt *time.Time

	ErrorMessage string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (DeviceCommand) TableName() string {

	return "device_commands"

}
