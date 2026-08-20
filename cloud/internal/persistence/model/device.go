package model

import (
	"time"

	"github.com/google/uuid"
)
/*
// DeviceTypeCategory defines the general category of a device type.
type DeviceTypeCategory string

const (
	DeviceTypeCategoryGateway    DeviceTypeCategory = "gateway"
	DeviceTypeCategoryController DeviceTypeCategory = "controller"
	DeviceTypeCategoryRobot      DeviceTypeCategory = "robot"
	DeviceTypeCategoryMachine    DeviceTypeCategory = "machine"
	DeviceTypeCategorySensor     DeviceTypeCategory = "sensor"
	DeviceTypeCategoryInspection DeviceTypeCategory = "inspection"
	DeviceTypeCategoryLogistics  DeviceTypeCategory = "logistics"
	DeviceTypeCategoryProductIoT DeviceTypeCategory = "product_iot"
)



// DeviceType is the business extension of a Resource whose
// Resource.Type is DEVICE_TYPE.
type DeviceType struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	ResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	Code string `gorm:"type:varchar(100);not null;uniqueIndex"`

	Category DeviceTypeCategory `gorm:"type:varchar(50);not null;index"`

	Description string `gorm:"type:text"`

	Enabled bool `gorm:"not null;default:true;index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Device is the business extension of a Resource whose
// Resource.Type is DEVICE.
type Device struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	ResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	DeviceTypeID uuid.UUID `gorm:"type:uuid;not null;index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}*/


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
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	ResourceID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	ProductModelID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Identity attributes
	SerialNumber string `gorm:"size:255;index"`
	HardwareID   string `gorm:"size:255;index"`

	// Runtime state
	Status DeviceStatus `gorm:"size:50;not null"`

	// Lifecycle timestamps
	ActivatedAt  *time.Time
	LastOnlineAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}