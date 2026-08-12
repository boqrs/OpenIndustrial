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
// AttributeDefinition defines the blueprint for an attribute that can be associated
// with a ProductModel. It specifies the name, data type, and other constraints
// for an attribute, but does not hold a value itself.
type AttributeDefinition struct {
	// ID is the internal, auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique identifier for the attribute definition.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	// ProductModelID is a foreign key linking this attribute definition to a specific
	// product model (which is a type of Resource). This creates a "blueprint" of
	// attributes for all device instances of that model.
	ResourceID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Name is the machine-readable key for the attribute (e.g., "motor_speed").
	// It should be unique per product model.
	Name string `gorm:"type:varchar(255);not null;uniqueIndex:idx_attr_def_model_name"`

	// Description provides a human-readable explanation of the attribute's purpose.
	Description string `gorm:"type:text"`

	// DataType specifies the value type for this attribute (e.g., "string", "int", "float", "bool").
	DataType AttributeValueType `gorm:"type:varchar(50);not null"`

	// Unit specifies the physical unit for numerical data types (e.g., "RPM", "°C", "Pa").
	// It provides essential context for interpreting the attribute's value.
	Unit string `gorm:"type:varchar(50)"`

	// Label provides a user-friendly name for display purposes (e.g., "Motor Speed").
	Label string `gorm:"type:varchar(255);not null"`

	// Standard timestamp fields
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// GORM relationship for preloading.
	Resource Resource `gorm:"references:UUID;foreignKey:ProductModelID"`
}

// TableName specifies the table name for the AttributeDefinition model.
func (AttributeDefinition) TableName() string {
	return "attribute_definitions"
}