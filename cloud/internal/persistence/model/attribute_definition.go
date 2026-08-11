package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AttributeValueType defines the possible data types for an attribute.
type AttributeValueType string

const (
	AttributeValueTypeString   AttributeValueType = "string"
	AttributeValueTypeText     AttributeValueType = "text"
	AttributeValueTypeInteger  AttributeValueType = "integer"
	AttributeValueTypeFloat    AttributeValueType = "float"
	AttributeValueTypeBoolean  AttributeValueType = "boolean"
	AttributeValueTypeDateTime AttributeValueType = "datetime"
	AttributeValueTypeJSON     AttributeValueType = "json"
)

// AttributeDefinition is the GORM model for the 'attribute_definitions' table.
// It defines the schema for a custom attribute that can be attached to resources.
// It uses the dual ID (internal auto-increment + external UUID) pattern because
// it's an independent business entity that may be exposed via its own API endpoint.
type AttributeDefinition struct {
	// ID is the internal, auto-incrementing primary key for database performance.
	ID uint `gorm:"primaryKey"`

	// UUID is the external-facing, unique business identifier for API security and consistency.
	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:uuid_generate_v4()"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	// The unique key or name for the attribute, e.g., "color", "weight".
	// It's unique per tenant.
	Name string `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_attribute_name"`

	// A human-readable label for the attribute, e.g., "Product Color".
	Label string `gorm:"type:varchar(255);not null"`

	// The data type of the attribute's value.
	DataType AttributeValueType `gorm:"type:varchar(50);not null"`

	// Indicates if this attribute is required for resources of a certain type.
	Required bool `gorm:"not null;default:false"`

	// A default value for the attribute, stored as JSON text.
	DefaultValue *string `gorm:"type:text"`

	// Standard timestamp fields
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


// TableName specifies the table name for the AttributeDefinition model.
func (AttributeDefinition) TableName() string {
	return "attribute_definitions"
}