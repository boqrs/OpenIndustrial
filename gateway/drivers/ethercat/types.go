package ethercat

import "time"

// Sample represents a single data point collected from a device.
// This is the unified data model for all drivers.
type Sample struct {
	PointID   string
	Value     interface{}
	Timestamp time.Time
	Quality   Quality
	Source    string // e.g., "ethercat"
}

// Quality defines the quality of a data point.
type Quality string

const (
	QualityGood         Quality = "good"
	QualityBad          Quality = "bad"
	QualityUncertain    Quality = "uncertain"
	QualityDisconnected Quality = "disconnected"
)

// SlaveInfo represents a discovered EtherCAT slave device.
type SlaveInfo struct {
	Index       uint16
	VendorID    uint32
	ProductCode uint32
	Revision    uint32
	Name        string
}

// SDORequest represents a request to read from or write to an SDO object.
type SDORequest struct {
	Slave    uint16
	Index    uint16
	SubIndex uint8
	Value    interface{} // Used for writing, ignored for reading
	Timeout  time.Duration
}

// SDOResponse represents the data returned from an SDO read request.
type SDOResponse struct {
	Data []byte
}

// PDOEntry defines a single entry within a Process Data Object (PDO).
type PDOEntry struct {
	Name      string
	Index     uint16
	SubIndex  uint8
	Offset    uint // Offset in bits within the PDO buffer
	BitLength uint
	DataType  DataType
}

// PDOMapping defines the PDO configuration for a specific slave.
type PDOMapping struct {
	Slave     uint16
	Direction PDODirection
	Entries   []PDOEntry
}

// PDODirection specifies if the PDO is for input (Tx) or output (Rx).
type PDODirection string

const (
	PDOInput  PDODirection = "input"  // Data from slave to master (TxPDO)
	PDOOutput PDODirection = "output" // Data from master to slave (RxPDO)
)

// DataType defines the data type of a PDO/SDO entry.
type DataType string

const (
	TypeBOOL  DataType = "BOOL"
	TypeSINT  DataType = "SINT"
	TypeINT   DataType = "INT"
	TypeDINT  DataType = "DINT"
	TypeLINT  DataType = "LINT"
	TypeUSINT DataType = "USINT"
	TypeUINT  DataType = "UINT"
	TypeUDINT DataType = "UDINT"
	TypeULINT DataType = "ULINT"
	TypeREAL  DataType = "REAL"
	TypeLREAL DataType = "LREAL"
)