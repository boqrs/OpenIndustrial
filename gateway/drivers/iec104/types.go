package iec104

import "time"

// Sample represents a single data point collected from a device.
// This is the unified data model for all drivers.
type Sample struct {
	PointID   string
	Value     interface{}
	Timestamp time.Time
	Quality   Quality
	Source    string // e.g., "iec104"
}

// Quality defines the quality of a data point.
type Quality string

const (
	QualityGood      Quality = "good"
	QualityBad       Quality = "bad"
	QualityUncertain Quality = "uncertain"
)

// PointMapping defines the mapping between a user-defined PointID and
// the specific IEC 104 Information Object Address (IOA).
type PointMapping struct {
	ID            string
	IOA           uint32
	TypeID        TypeID
	CommonAddress uint16
	// DataType can be used for hints, though TypeID often implies it.
	DataType string
}

// ASDU (Application Service Data Unit) is the main data container in IEC 104.
type ASDU struct {
	TypeID        TypeID
	VSQ           uint8 // Variable Structure Qualifier
	Cause         CauseOfTransmission
	CommonAddress uint16
	Objects       []InformationObject
}

// InformationObject represents a single data point within an ASDU.
type InformationObject struct {
	IOA       uint32
	Value     interface{}
	Timestamp time.Time
	Quality   byte // The raw quality byte from the protocol
}

// TypeID represents the type of information in the ASDU.
type TypeID uint8

// CauseOfTransmission indicates the reason for the data transmission.
type CauseOfTransmission uint16

// Common constants for TypeID (a subset)
const (
	// Monitored direction (slave to master)
	M_SP_NA_1 TypeID = 1  // Single-point information
	M_DP_NA_1 TypeID = 3  // Double-point information
	M_ME_NA_1 TypeID = 9  // Measured value, normalized value
	M_ME_NB_1 TypeID = 11 // Measured value, scaled value
	M_ME_NC_1 TypeID = 13 // Measured value, short floating point number
	M_IT_NA_1 TypeID = 34 // Integrated totals

	// Controlled direction (master to slave)
	C_SC_NA_1 TypeID = 45 // Single command
	C_DC_NA_1 TypeID = 46 // Double command
	C_SE_NC_1 TypeID = 50 // Set-point command, short floating point number
	C_IC_NA_1 TypeID = 100 // Interrogation command
)

