package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConnectionType is a string type that defines the semantics of a technical binding between resources.
type ConnectionType string

const (
	// ConnectionTypeConnectedThrough indicates that a resource (e.g., a device)
	// is connected or communicates through another resource (e.g., a gateway).
	ConnectionTypeConnectedThrough ConnectionType = "connected_through"

	// ConnectionTypeMonitoredBy indicates that a resource (e.g., an asset)
	// is being monitored by another resource (e.g., a sensor).
	ConnectionTypeMonitoredBy ConnectionType = "monitored_by"

	// ConnectionTypeControls indicates that a resource (e.g., a PLC)
	// is controlling another resource (e.g., a robotic arm).
	ConnectionTypeControls ConnectionType = "controls"

	// ConnectionTypePoweredBy indicates that a resource's power is supplied
	// by another resource (e.g., a specific power circuit).
	ConnectionTypePoweredBy ConnectionType = "powered_by"

	// ConnectionTypePairedWith indicates that two resources are functionally paired.
	ConnectionTypePairedWith ConnectionType = "paired_with"
)

// ResourceConnection defines a specific, runtime technical binding between two resources.
type ResourceConnection struct {
	ID       uint      `gorm:"primaryKey"`
	// SourceResourceID is the resource where the connection originates (e.g., a Device).
	SourceResourceID uuid.UUID `gorm:"type:uuid;not null;index:idx_connection_unique,priority:1"`
	SourceResource   Resource  `gorm:"foreignKey:SourceResourceID;references:UUID"`
	// TargetResourceID is the resource where the connection terminates (e.g., a Gateway).
	TargetResourceID uuid.UUID `gorm:"type:uuid;not null;index:idx_connection_unique,priority:2"`
	TargetResource   Resource  `gorm:"foreignKey:TargetResourceID;references:UUID"`
	// ConnectionType defines the semantics of the technical binding (e.g., "connected_through").
	ConnectionType ConnectionType `gorm:"type:varchar(100);not null;index:idx_connection_unique,priority:3"`
	// Metadata can store additional context about the connection itself.
	Metadata  []byte         `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


// TableName specifies the table name for the ResourceAttribute model.
func (ResourceConnection) TableName() string {
	return "resource_connections"
}