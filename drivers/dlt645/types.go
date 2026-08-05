package dlt645

import "time"

// ProtocolVersion defines the DL/T 645 protocol version.
type ProtocolVersion string

const (
	DLT645_1997 ProtocolVersion = "DL/T 645-1997"
	DLT645_2007 ProtocolVersion = "DL/T 645-2007"
)

// DataType defines the type of data for a point.
type DataType string

const (
	DataTypeUint32  DataType = "uint32"
	DataTypeInt32   DataType = "int32"
	DataTypeFloat32 DataType = "float32"
	DataTypeBCD     DataType = "bcd"
)

// Meter represents a physical electricity meter.
type Meter struct {
	Address      string          `yaml:"address"`
	Protocol     ProtocolVersion `yaml:"protocol"`
	Manufacturer string          `yaml:"manufacturer"`
}

// PointMapping maps a user-defined ID to a DL/T 645 Data Identifier (DI).
type PointMapping struct {
	ID       string   `yaml:"id"`
	DI       uint32   `yaml:"di"`
	DataType DataType `yaml:"type"`
	Scale    float64  `yaml:"scale"`
}

// Frame represents a DL/T 645 communication frame.
// As per the spec: 68 | Address | 68 | Control | Length | Data | CS | 16
type Frame struct {
	Address []byte // 6-byte BCD encoded address
	Control byte
	Data    []byte
}

// Sample represents a single data point collected from a device.
// This is the unified data model for all drivers.
type Sample struct {
	PointID   string
	Value     interface{}
	Timestamp time.Time
	Quality   Quality
	Source    string // e.g., "dlt645"
}

// Quality defines the quality of a data point.
type Quality string

const (
	QualityGood      Quality = "good"
	QualityBad       Quality = "bad"
	QualityUncertain Quality = "uncertain"
)

// WriteRequest represents a write operation to a device.
type WriteRequest struct {
	Address string
	Point   PointMapping
	Value   interface{}
}