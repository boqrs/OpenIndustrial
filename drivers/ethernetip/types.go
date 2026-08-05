package ethernetip

import "time"

// Sample represents a single data point collected from a device.
// This is the unified data model for all drivers.
type Sample struct {
	PointID   string
	Value     interface{}
	Timestamp time.Time
	Quality   Quality
	Source    string // e.g., "ethernetip"
}

// Quality defines the quality of a data point.
type Quality string

const (
	QualityGood         Quality = "good"
	QualityBad          Quality = "bad"
	QualityUncertain    Quality = "uncertain"
	QualityDisconnected Quality = "disconnected"
)

// PointMapping defines how a high-level PointID maps to a specific
// location within an EtherNet/IP device.
type PointMapping struct {
	ID   string
	Type PointType

	// For CIP Object access
	Class     uint16
	Instance  uint16
	Attribute uint16

	// For PLC Tag access
	Tag string

	DataType DataType
	Writable bool
}

// PointType specifies the addressing mode for a point.
type PointType string

const (
	PointTypeCIPObject PointType = "cip"
	PointTypePLCTag    PointType = "tag"
)

// DataType defines the data type of a tag or attribute.
type DataType string

const (
	// Basic Types
	TypeBOOL  DataType = "BOOL"  // Boolean
	TypeSINT  DataType = "SINT"  // 8-bit signed integer
	TypeINT   DataType = "INT"   // 16-bit signed integer
	TypeDINT  DataType = "DINT"  // 32-bit signed integer
	TypeLINT  DataType = "LINT"  // 64-bit signed integer
	TypeUSINT DataType = "USINT" // 8-bit unsigned integer
	TypeUINT  DataType = "UINT"  // 16-bit unsigned integer
	TypeUDINT DataType = "UDINT" // 32-bit unsigned integer
	TypeULINT DataType = "ULINT" // 64-bit unsigned integer
	TypeREAL  DataType = "REAL"  // 32-bit floating point
	TypeLREAL DataType = "LREAL" // 64-bit floating point
	TypeSTRING DataType = "STRING" // String
)