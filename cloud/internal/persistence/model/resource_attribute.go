package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// ResourceAttribute is the GORM model for the 'resource_attributes' table.
// It stores the actual value of a defined attribute for a specific resource instance.
// Following the new architecture decision, all foreign keys are now UUIDs.
type ResourceAttribute struct {
	// ID is the internal, auto-incrementing primary key. This model is not
	// exposed directly via API, so it does not need a public-facing UUID.
	ID uint `gorm:"primaryKey"`

	// Foreign key to the 'resources' table's public-facing UUID.
	ResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_resource_attr_def"`

	// Foreign key to the 'attribute_definitions' table's public-facing UUID.
	AttributeDefinitionID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_resource_attr_def"`

	// The actual value of the attribute, stored in a flexible JSONB format.
	Value []byte `gorm:"type:jsonb"`

	// Standard timestamp fields
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships for GORM to preload data if needed.
	// Note: GORM needs to know which fields to use for the join.
	Resource            Resource            `gorm:"references:UUID;foreignKey:ResourceID"`
	AttributeDefinition AttributeDefinition `gorm:"references:UUID;foreignKey:AttributeDefinitionID"`
}

// TableName specifies the table name for the ResourceAttribute model.
func (ResourceAttribute) TableName() string {
	return "resource_attributes"
}

// SetValue serializes the provided value into JSON and stores it in the Value field.
func (ra *ResourceAttribute) SetValue(value interface{}) {
	// We don't need to handle nil, as json.Marshal will correctly produce "null".
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		// In case of a marshaling error, we can store a JSON null or an error string.
		// Storing null is safer.
		ra.Value = []byte("null")
		return
	}
	ra.Value = jsonBytes
}

// GetValueColumns returns the single column name that holds the attribute value.
func (ra *ResourceAttribute) GetValueColumns() []string {
	return []string{"value"}
}