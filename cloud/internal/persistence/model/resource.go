package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// StatusActive means the resource is active and operational.
	StatusActive = "active"
	// StatusInactive means the resource is not currently in use.
	StatusInactive = "inactive"
	// StatusArchived means the resource is archived for historical purposes.
	StatusArchived = "archived"
	// StatusPending means the resource is awaiting some action (e.g., approval).
	StatusPending = "pending"

	// StatusProvisioned means a device has been registered in the factory but not yet activated by an end-user.
	StatusProvisioned = "PROVISIONED"
	// StatusOnboarded means a device has been successfully activated by an end-user.
	StatusOnboarded = "ONBOARDED"
	// StatusOffline means a device is currently not connected.
	StatusOffline = "OFFLINE"
	// StatusDecommissioned means a device has been permanently taken out of service.
	StatusDecommissioned = "DECOMMISSIONED"
)

// Resource is the GORM model for the 'resources' table, reflecting the final architectural decision.
// It uses an auto-incrementing integer ID as the primary key for internal use,
// and a separate UUID field as the public-facing business identifier.
// This model should only be used within the persistence layer.
type Resource struct {
	// ID is the internal, auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the external-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:uuid_generate_v4()"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Renamed fields to avoid SQL keyword conflicts, with explicit column mapping.
	ResourceType   string `gorm:"column:resource_type;type:varchar(100);not null;index"`
	ResourceName   string `gorm:"column:resource_name;type:varchar(255);not null"`
	ResourceStatus string `gorm:"column:resource_status;type:varchar(50);not null;default:'active'"`

	Code     *string `gorm:"type:varchar(100);uniqueIndex"`
	Metadata []byte  `gorm:"type:jsonb"`
	Version  int     `gorm:"column:record_version;not null;default:1"`

	// ParentID is a nullable foreign key to the internal 'ID' field.
	ParentID uuid.UUID `gorm:"index"`

	// OwnerGroupID references the UUID of a group.
	OwnerGroupID *uuid.UUID `gorm:"type:uuid;index"`

	// Standard timestamp fields
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for the Resource model.
func (Resource) TableName() string {
	return "resources"
}