package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Factory represents a factory in the industrial system.
// It extends the base Resource model with factory-specific attributes.
type Factory struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ResourceID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex"` // Foreign key to the Resource
	Resource   Resource       `gorm:"foreignKey:ResourceID"`          // Belongs to Resource
	Code       string         `gorm:"type:varchar(100);uniqueIndex"`  // Factory-specific code, e.g., "F001"
	Address    string         `gorm:"type:varchar(255)"`              // Physical address of the factory
	Status     string         `gorm:"type:varchar(50)"`               // e.g., Active, Inactive, UnderConstruction
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}