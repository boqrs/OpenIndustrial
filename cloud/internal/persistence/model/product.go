package model

import (
	"time"

	"github.com/google/uuid"
)

// ProductModel describes the static definition of a product/device model.
//
// Resource stores the common identity:
//   - UUID
//   - TenantID
//   - Name
//   - Type
//   - Status
//   - Metadata
//
// ProductModel stores product-domain-specific information.
type ProductModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	// ResourceID references resources.uuid.
	ResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	// Code identifies the product model family.
	// Version distinguishes immutable model definitions.
	Code    string `gorm:"type:varchar(100);not null"`
	Version string `gorm:"type:varchar(50);not null"`

	// Product category, for example:
	// robot / plc / gateway / machine / product_iot.
	Category string `gorm:"type:varchar(100);not null;index"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}