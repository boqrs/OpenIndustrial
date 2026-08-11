package resource

import (
	"time"

	"github.com/google/uuid"
)

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

// AttributeDefinition defines a type of attribute that can be associated with a resource.
type AttributeDefinition struct {
	ID          uuid.UUID          `db:"id"`
	TenantID    uuid.UUID          `db:"tenant_id"`
	Key         string             `db:"key"`
	Name        string             `db:"name"`
	Description *string            `db:"description"`
	ValueType   AttributeValueType `db:"value_type"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
}

// ResourceAttribute represents a specific attribute value for a given resource.
type ResourceAttribute struct {
	ResourceID    uuid.UUID   `db:"resource_id"`
	AttributeID   uuid.UUID   `db:"attribute_id"`
	ValueString   *string     `db:"value_string"`
	ValueText     *string     `db:"value_text"`
	ValueInteger  *int64      `db:"value_integer"`
	ValueFloat    *float64    `db:"value_float"`
	ValueBoolean  *bool       `db:"value_boolean"`
	ValueDateTime *time.Time  `db:"value_datetime"`
	ValueJSON     []byte      `db:"value_json"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
}