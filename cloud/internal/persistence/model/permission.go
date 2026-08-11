package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents an action that can be performed, e.g., "users:create".
// This is the GORM model for the `permissions` table.
// Permissions are typically seeded and not managed via a regular API,
// so we use a simple auto-incrementing primary key.
type Permission struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// Name is the unique identifier for the permission, e.g., "resources:create".
	Name string `gorm:"type:varchar(255);unique;not null"`

	Description string `gorm:"type:text"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName explicitly sets the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}