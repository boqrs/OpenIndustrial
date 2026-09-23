package model

import (
	"time"
	///"github.com/google/uuid"
)

// DeviceConnection represents runtime MQTT connection state.
//
// It belongs to Resource instead of Device because
// IoT identity is based on ResourceID.
type DeviceConnection struct {
	ID uint `gorm:"primaryKey"`

	ResourceID uint `gorm:"not null;uniqueIndex"`

	// Current connection state

	Status string `gorm:"size:50;not null"`

	ClientID string `gorm:"size:255;not null;index"`

	LastConnectedAt *time.Time

	LastDisconnectedAt *time.Time

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (DeviceConnection) TableName() string {
	return "device_connections"
}

// DeviceTopicPermission represents MQTT ACL.
//
// It does not belong to MQTT provider.
// It describes cloud authorization.
type DeviceTopicPermission struct {
	ID uint `gorm:"primaryKey"`

	ResourceID uint `gorm:"not null;index"`

	Topic string `gorm:"size:500;not null"`

	Permission string `gorm:"size:50;not null"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (DeviceTopicPermission) TableName() string {

	return "device_topic_permissions"

}
