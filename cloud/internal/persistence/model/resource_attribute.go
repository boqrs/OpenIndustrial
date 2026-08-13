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
    ID uint `gorm:"primaryKey"`

    ResourceID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_resource_attr_def"`

    AttributeDefinitionID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_resource_attr_def"`

    Value []byte `gorm:"type:jsonb;not null"`

    CreatedAt time.Time
    UpdatedAt time.Time

    Resource Resource `gorm:"foreignKey:ResourceID;references:UUID"`

    AttributeDefinition AttributeDefinition `gorm:"foreignKey:AttributeDefinitionID;references:ID"`
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