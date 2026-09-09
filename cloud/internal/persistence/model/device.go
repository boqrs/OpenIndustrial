package model

import (
	"time"
)

// DeviceStatus represents the runtime status of a device, distinct from its resource lifecycle status.
type DeviceStatus string

const (
	// DeviceStatusCreated means the device has been registered in the system but has never connected.
	DeviceStatusCreated DeviceStatus = "created"
	// DeviceStatusOnline means the device is currently connected and communicating.
	DeviceStatusOnline DeviceStatus = "online"
	// DeviceStatusOffline means the device was previously online but is now disconnected.
	DeviceStatusOffline DeviceStatus = "offline"
	// DeviceStatusFault means the device has reported an error state.
	DeviceStatusFault DeviceStatus = "fault"
	// DeviceStatusMaintenance means the device is temporarily out of service for maintenance.
	DeviceStatusMaintenance DeviceStatus = "maintenance"
)

// Device represents a physical device instance in the real world.
// It is an instantiation of a static ProductModel.
type Device struct {
	ID         uint `gorm:"primaryKey"`
	ResourceID uint `gorm:"not null;index"`
	ProductID  uint `gorm:"not null;index"`

	// Manufacturing provenance
	WorkOrderID       uint `gorm:"not null;index"`
	ExecutionID       uint `gorm:"not null;index"`
	ExecutionResultID uint `gorm:"not null;index"`

	// Identity
	SerialNumber string `gorm:"size:255;not null;index"`
	HardwareID   string `gorm:"size:255;index"`

	// Runtime state
	Status DeviceStatus `gorm:"size:50;not null"`

	ActivatedAt  *time.Time
	LastOnlineAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
