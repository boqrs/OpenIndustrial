package hj212

import "time"

// DataSegment represents the core data part of an HJ/T 212 message.
// It contains the key-value pairs found in the "CP" field.
type DataSegment struct {
	QN string // Request Number
	ST string // System Type
	CN string // Command Number
	PW string // Password
	MN string // Monitor Node ID

	// Pollutants holds the key-value data for each pollutant code.
	// e.g., "w01018-Rtd": "56.3"
	Pollutants map[string]string
}

// PointMapping maps a user-defined ID to a specific HJ/T 212 pollutant code.
type PointMapping struct {
	ID   string `yaml:"id"`   // User-friendly ID, e.g., "sewage_cod"
	Code string `yaml:"code"` // HJ/T 212 pollutant code, e.g., "w01018"
	Unit string `yaml:"unit"` // e.g., "mg/L"
}

// Sample represents a single data point collected from a device.
// This is the unified data model for all drivers.
type Sample struct {
	PointID   string
	Value     interface{}
	Timestamp time.Time
	Quality   string // "good", "bad", etc.
	Source    string // e.g., "hj212"
}