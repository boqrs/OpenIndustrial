package model

import (
	"time"

	"github.com/google/uuid"
)

// Factory represents the factory entity in the database.
// It uses an auto-incrementing integer as the primary key for internal database use,
// and a separate UUID field as the public-facing business identifier.
type Factory struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing unique identifier for business operations.
	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:uuid_generate_v4()"`

	// ResourceID points to the corresponding entry in the resources table.
	ResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	// Code is a user-defined unique code for the factory.
	Code string `gorm:"type:varchar(100);not null;uniqueIndex"`

	// Address stores the physical address of the factory.
	Address string `gorm:"type:text"`

	// Timezone of the factory's location.
	Timezone string `gorm:"type:varchar(100);not null;default:'UTC'"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}