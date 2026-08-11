package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant represents an isolated workspace in the system.
// It is the root entity for all other resources within a tenant's scope.
type Tenant struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier for the tenant.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	Name   string `gorm:"type:varchar(100);not null;unique"`
	Status string `gorm:"type:varchar(20);default:'active'"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName explicitly sets the table name for the Tenant model.
func (Tenant) TableName() string {
	return "tenants"
}