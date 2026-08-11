package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Group represents a collection of users for permissioning.
// This is the GORM model, with a separate auto-incrementing primary key (ID)
// and a public-facing business key (UUID).
type Group struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// --- GORM Relationships ---
	// A group belongs to a tenant.
	// Tenant Tenant `gorm:"foreignKey:TenantID;references:UUID"` // Uncomment if Tenant model exists and has UUID as a business key.
}

// TableName explicitly sets the table name for the Group model.
func (Group) TableName() string {
	return "groups"
}