package model

import (
	"time"
)

// Factory represents the factory entity in the database.
// It uses an auto-incrementing integer as the primary key for internal database use,
// and a separate UUID field as the public-facing business identifier.
type Factory struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`
	// ResourceID points to the corresponding entry in the resources table.
	ResourceID uint `gorm:"not null;index"`
	// Code is a user-defined unique code for the factory.
	Code string `gorm:"type:varchar(100);not null;uniqueIndex"`
	// Address stores the physical address of the factory.
	Address string `gorm:"type:text"`
	// Timezone of the factory's location.
	Timezone string `gorm:"type:varchar(100);not null;default:'UTC'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}