package model

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================
// Device lifecycle status
//
// This status represents the business lifecycle of a device.
//
// It does NOT represent network connection state.
// ============================================================

type DeviceStatus string

const (

	// Device created from manufacturing result.
	DeviceStatusCreated DeviceStatus = "created"

	// Device activated by customer.
	DeviceStatusActivated DeviceStatus = "activated"

	// Device permanently disabled.
	DeviceStatusDisabled DeviceStatus = "disabled"
)

func (s DeviceStatus) String() string {
	return string(s)
}

func (s DeviceStatus) IsValid() bool {

	switch s {

	case DeviceStatusCreated,
		DeviceStatusActivated,
		DeviceStatusDisabled:

		return true

	default:

		return false
	}
}

// ============================================================
// IoT connection state
//
// Managed by MQTT/AWS IoT Core connection events.
//
// This is runtime state.
// ============================================================

type ConnectionStatus string

const (
	ConnectionStatusDisconnected ConnectionStatus = "disconnected"

	ConnectionStatusConnected ConnectionStatus = "connected"
)

func (s ConnectionStatus) String() string {

	return string(s)

}

func (s ConnectionStatus) IsValid() bool {

	switch s {

	case ConnectionStatusDisconnected,
		ConnectionStatusConnected:

		return true

	default:

		return false
	}
}

// ============================================================
// Device
//
// Device is created ONLY during manufacturing.
//
// IoT activation never creates Device.
// ============================================================

type Device struct {
	ID uint `gorm:"primaryKey"`

	// --------------------------------------------------------
	// Resource identity
	// --------------------------------------------------------

	ResourceID uint `gorm:"not null;index"`

	// Static product definition

	ProductID uint `gorm:"not null;index"`

	// --------------------------------------------------------
	// Manufacturing provenance
	// --------------------------------------------------------

	WorkOrderID uint `gorm:"not null;index"`

	ExecutionID uint `gorm:"not null;index"`

	ExecutionResultID uint `gorm:"not null;index"`

	// --------------------------------------------------------
	// Physical identity
	// --------------------------------------------------------

	SerialNumber string `gorm:"size:255;not null;uniqueIndex"`

	HardwareID string `gorm:"size:255;index"`

	// --------------------------------------------------------
	// Device lifecycle
	// --------------------------------------------------------

	Status DeviceStatus `gorm:"size:50;not null;index"`

	// --------------------------------------------------------
	// Customer activation
	// --------------------------------------------------------

	CustomerUUID *uuid.UUID `gorm:"type:uuid;index"`

	ActivatedAt *time.Time

	// --------------------------------------------------------
	// IoT runtime
	// --------------------------------------------------------

	ConnectionStatus ConnectionStatus `gorm:"size:50;not null;default:'disconnected'"`

	// MQTT client identifier.
	//
	// Usually mapped from certificate identity.

	ClientID string `gorm:"size:255;index"`

	LastOnlineAt *time.Time

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (*Device) TableName() string {

	return "devices"

}
