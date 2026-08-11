package model

import (
	"time"
	
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// ResourceRelation is the GORM model for the 'resource_relations' table.
// It represents a directional, typed relationship between two resources.
// Following the new architecture decision, all foreign keys are now UUIDs.
type ResourceRelation struct {
	// ID is the internal, auto-incrementing primary key. This model is not
	// exposed directly via API, so it does not need a public-facing UUID.
	ID uint `gorm:"primaryKey"`

	// Foreign key for the source resource, referencing the public-facing UUID.
	FromResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_relation_uniqueness"`

	// Foreign key for the target resource, referencing the public-facing UUID.
	ToResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_relation_uniqueness"`

	// Describes the nature of the relationship, e.g., "depends_on", "contains".
	RelationType string `gorm:"type:varchar(100);not null;uniqueIndex:idx_relation_uniqueness"`

	// Standard timestamp fields
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// GORM relationships for preloading.
	FromResource Resource `gorm:"references:UUID;foreignKey:FromResourceID"`
	ToResource   Resource `gorm:"references:UUID;foreignKey:ToResourceID"`
}

// TableName specifies the table name for the ResourceRelation model.
func (ResourceRelation) TableName() string {
	return "resource_relations"
}