package model

import (
	"time"
	"encoding/json"


	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// StatusActive means the resource is active and operational.
	StatusActive = "active"
	// StatusInactive means the resource is not currently in use.
	StatusInactive = "inactive"
	// StatusArchived means the resource is archived for historical purposes.
	StatusArchived = "archived"
	// StatusPending means the resource is awaiting some action (e.g., approval).
	StatusPending = "pending"

	// StatusProvisioned means a device has been registered in the factory but not yet activated by an end-user.
	StatusProvisioned = "PROVISIONED"
	// StatusOnboarded means a device has been successfully activated by an end-user.
	StatusOnboarded = "ONBOARDED"
	// StatusOffline means a device is currently not connected.
	StatusOffline = "OFFLINE"
	// StatusDecommissioned means a device has been permanently taken out of service.
	StatusDecommissioned = "DECOMMISSIONED"
)

// Resource is the GORM model for the 'resources' table, reflecting the final architectural decision.
// It uses an auto-incrementing integer ID as the primary key for internal use,
// and a separate UUID field as the public-facing business identifier.
// This model should only be used within the persistence layer.
type Resource struct {
	// ID is the internal, auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the external-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:uuid_generate_v4()"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Renamed fields to avoid SQL keyword conflicts, with explicit column mapping.
	ResourceType   string `gorm:"column:resource_type;type:varchar(100);not null;index"`
	ResourceName   string `gorm:"column:resource_name;type:varchar(255);not null"`
	ResourceStatus string `gorm:"column:resource_status;type:varchar(50);not null;default:'active'"`

	Code     *string `gorm:"type:varchar(100);uniqueIndex"`
	Metadata []byte  `gorm:"type:jsonb"`
	Version  int     `gorm:"column:record_version;not null;default:1"`

	// ParentID is a nullable foreign key to the internal 'ID' field.
	ParentID uuid.UUID `gorm:"index"`

	// OwnerGroupID references the UUID of a group.
	OwnerGroupID *uuid.UUID `gorm:"type:uuid;index"`

	// Standard timestamp fields
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for the Resource model.
func (Resource) TableName() string {
	return "resources"
}

// ConnectionType is a string type that defines the semantics of a technical binding between resources.
type ConnectionType string

const (
	// ConnectionTypeConnectedThrough indicates that a resource (e.g., a device)
	// is connected or communicates through another resource (e.g., a gateway).
	ConnectionTypeConnectedThrough ConnectionType = "connected_through"

	// ConnectionTypeMonitoredBy indicates that a resource (e.g., an asset)
	// is being monitored by another resource (e.g., a sensor).
	ConnectionTypeMonitoredBy ConnectionType = "monitored_by"

	// ConnectionTypeControls indicates that a resource (e.g., a PLC)
	// is controlling another resource (e.g., a robotic arm).
	ConnectionTypeControls ConnectionType = "controls"

	// ConnectionTypePoweredBy indicates that a resource's power is supplied
	// by another resource (e.g., a specific power circuit).
	ConnectionTypePoweredBy ConnectionType = "powered_by"

	// ConnectionTypePairedWith indicates that two resources are functionally paired.
	ConnectionTypePairedWith ConnectionType = "paired_with"
)

// ResourceConnection defines a specific, runtime technical binding between two resources.
type ResourceConnection struct {
	ID       uint      `gorm:"primaryKey"`
	// SourceResourceID is the resource where the connection originates (e.g., a Device).
	SourceResourceID uuid.UUID `gorm:"type:uuid;not null;index:idx_connection_unique,priority:1"`
	SourceResource   Resource  `gorm:"foreignKey:SourceResourceID;references:UUID"`
	// TargetResourceID is the resource where the connection terminates (e.g., a Gateway).
	TargetResourceID uuid.UUID `gorm:"type:uuid;not null;index:idx_connection_unique,priority:2"`
	TargetResource   Resource  `gorm:"foreignKey:TargetResourceID;references:UUID"`
	// ConnectionType defines the semantics of the technical binding (e.g., "connected_through").
	ConnectionType ConnectionType `gorm:"type:varchar(100);not null;index:idx_connection_unique,priority:3"`
	// Metadata can store additional context about the connection itself.
	Metadata  []byte         `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


// TableName specifies the table name for the ResourceAttribute model.
func (ResourceConnection) TableName() string {
	return "resource_connections"
}


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