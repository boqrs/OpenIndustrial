package model

import (
	"time"

	"github.com/google/uuid"
)

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
}