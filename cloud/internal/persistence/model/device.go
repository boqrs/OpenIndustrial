package model

import (
	"time"

	"github.com/google/uuid"
)

// DeviceStatus represents the runtime status of a device.
//
// DeviceStatus is an IoT/runtime state. It does not represent
// warehouse or shipment lifecycle.
type DeviceStatus string

const (
	// DeviceStatusCreated means the device has been created by
	// the manufacturing process but has never connected to IoT.
	DeviceStatusCreated DeviceStatus = "created"

	// DeviceStatusOnline means the device is currently connected
	// and communicating with the IoT infrastructure.
	DeviceStatusOnline DeviceStatus = "online"

	// DeviceStatusOffline means the device is not currently connected.
	DeviceStatusOffline DeviceStatus = "offline"

	// DeviceStatusFault means the device has reported a fault state.
	DeviceStatusFault DeviceStatus = "fault"

	// DeviceStatusMaintenance means the device is temporarily
	// unavailable because it is under maintenance.
	DeviceStatusMaintenance DeviceStatus = "maintenance"
)

func (s DeviceStatus) String() string {
	return string(s)
}

func (s DeviceStatus) IsValid() bool {
	switch s {
	case DeviceStatusCreated,
		DeviceStatusOnline,
		DeviceStatusOffline,
		DeviceStatusFault,
		DeviceStatusMaintenance:
		return true

	default:
		return false
	}
}

// Device represents a physical device instance in the real world.
//
// A Device is created only from the manufacturing process.
// IoT provisioning or activation must never create a Device.
//
// Device identity is anchored by:
//
//	Device
//	    └── Resource
//	          ├── ResourceIdentity
//	          └── ResourceCertificate
//
// Warehouse and shipment lifecycle is managed by WMS.
// IoT runtime state is managed through this model.
type Device struct {
	ID         uint `gorm:"primaryKey"`
	ResourceID uint `gorm:"not null;index"`
	ProductID  uint `gorm:"not null;index"`

	// Manufacturing provenance
	WorkOrderID       uint `gorm:"not null;index"`
	ExecutionID       uint `gorm:"not null;index"`
	ExecutionResultID uint `gorm:"not null;index"`

	// Physical identity
	//
	// SerialNumber is the canonical production identity.
	SerialNumber string `gorm:"size:255;not null;uniqueIndex"`

	// HardwareID is optional because not every product exposes
	// a hardware identifier during manufacturing.
	HardwareID string `gorm:"size:255;index"`

	// Customer ownership.
	//
	// This is the UUID of the user who activated the device.
	//
	// It is intentionally nullable because a manufactured device
	// does not belong to an end user yet.
	CustomerUUID *uuid.UUID `gorm:"type:uuid;index"`

	// IoT runtime state.
	//
	// This field represents connectivity/runtime state only.
	// It does not represent inventory or shipment state.
	Status DeviceStatus `gorm:"size:50;not null;index"`

	// ActivatedAt records when the device was first activated
	// by an end user.
	ActivatedAt *time.Time

	// LastOnlineAt records the last time the device was observed
	// online by the IoT infrastructure.
	LastOnlineAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Device) TableName() string {
	return "devices"
}
