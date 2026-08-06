// drivers/bacnet/types.go
package bacnet

import (
	"fmt"
	"time"
)

// ObjectType represents a BACnet object type.
type ObjectType uint16

const (
	ObjectTypeAnalogInput        ObjectType = 0
	ObjectTypeAnalogOutput       ObjectType = 1
	ObjectTypeAnalogValue        ObjectType = 2
	ObjectTypeBinaryInput        ObjectType = 3
	ObjectTypeBinaryOutput       ObjectType = 4
	ObjectTypeBinaryValue        ObjectType = 5
	ObjectTypeMultiStateInput    ObjectType = 13
	ObjectTypeMultiStateOutput   ObjectType = 14
	ObjectTypeMultiStateValue    ObjectType = 19
	// Add more common object types as needed
)

// PropertyID represents a BACnet property identifier.
type PropertyID uint32

const (
	PropertyIDPresentValue PropertyID = 85
	PropertyIDStatusFlags  PropertyID = 111
	PropertyIDDescription  PropertyID = 28
	PropertyIDUnits        PropertyID = 117
	// Add more common property IDs as needed
)

// DataType represents the expected Go data type after reading a BACnet property.
// This helps in type assertion and validation.
type DataType string

const (
	DataTypeBoolean DataType = "Boolean"
	DataTypeInt     DataType = "Int" // For int32/int64
	DataTypeUint    DataType = "Uint" // For uint32/uint64
	DataTypeFloat   DataType = "Float" // For float32/float64
	DataTypeString  DataType = "String"
	DataTypeEnum    DataType = "Enum" // For multi-state values
	// Add more as needed, e.g., DateTime, Time, Date
)

// NodeMapping defines how an internal point ID maps to a BACnet Object and Property.
type NodeMapping struct {
	// ID is the internal unique identifier for this data point.
	ID string `json:"id" yaml:"id"`

	// ObjectType is the BACnet object type (e.g., AnalogInput).
	ObjectType ObjectType `json:"objectType" yaml:"objectType"`

	// Instance is the BACnet object instance number (e.g., 1 for AI:1).
	Instance uint32 `json:"instance" yaml:"instance"`

	// PropertyID is the BACnet property to read/write (e.g., Present_Value).
	PropertyID PropertyID `json:"propertyId" yaml:"propertyId"`

	// DataType is the expected Go data type of the property value.
	DataType DataType `json:"dataType" yaml:"dataType"`

	// Writable indicates if this property can be written to.
	Writable bool `json:"writable" yaml:"writable"`

	// Description provides a human-readable description of the node.
	Description string `json:"description" yaml:"description"`
}

// Validate checks if the NodeMapping is valid.
func (nm *NodeMapping) Validate() error {
	if nm.ID == "" {
		return fmt.Errorf("node mapping ID cannot be empty")
	}
	// Basic validation for ObjectType and PropertyID (can be expanded)
	if nm.ObjectType > 2000 { // Arbitrary upper limit for common types
		return fmt.Errorf("invalid object type %d for ID '%s'", nm.ObjectType, nm.ID)
	}
	if nm.PropertyID == 0 { // PropertyID 0 is usually reserved or invalid for Present_Value etc.
		return fmt.Errorf("invalid property ID %d for ID '%s'", nm.PropertyID, nm.ID)
	}
	if nm.DataType == "" {
		return fmt.Errorf("node mapping DataType cannot be empty for ID '%s'", nm.ID)
	}
	switch nm.DataType {
	case DataTypeBoolean, DataTypeInt, DataTypeUint, DataTypeFloat, DataTypeString, DataTypeEnum:
		// Valid
	default:
		return fmt.Errorf("unsupported DataType '%s' for ID '%s'", nm.DataType, nm.ID)
	}
	return nil
}

// Sample represents a collected data point from a BACnet property.
// This is the unified output format for the Runtime.
type Sample struct {
	ID         string    `json:"id"`          // Internal ID from NodeMapping
	ObjectType ObjectType `json:"object_type"` // Original BACnet Object Type
	Instance   uint32    `json:"instance"`    // Original BACnet Object Instance
	PropertyID PropertyID `json:"property_id"` // Original BACnet Property ID
	Value      any       `json:"value"`       // Decoded value
	Timestamp  time.Time `json:"timestamp"`   // Timestamp from the BACnet device or local
	Quality    Quality   `json:"quality"`     // Quality of the sample
}

// Quality represents the quality of a collected sample.
type Quality string

const (
	QualityGood        Quality = "good"
	QualityBad         Quality = "bad"
	QualityUncertain   Quality = "uncertain"
	QualityDisconnected Quality = "disconnected"
	QualityNotSupported Quality = "not_supported" // For features not supported by device/client
)
