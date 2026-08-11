package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a named set of permissions.
// This is the GORM model for the `roles` table.
type Role struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_role_name"`
	Description string    `gorm:"type:text"`
	IsSystem    bool      `gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// --- GORM Relationships ---

	// Permissions is a many-to-many relationship.
	// GORM will automatically create a join table named `role_permissions`.
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

// TableName explicitly sets the table name for the Role model.
func (Role) TableName() string {
	return "roles"
}