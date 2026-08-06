package can

// ByteOrder defines the byte order for multi-byte signals.
type ByteOrder string

const (
	// BigEndian (Motorola)
	BigEndian ByteOrder = "BigEndian"
	// LittleEndian (Intel)
	LittleEndian ByteOrder = "LittleEndian"
)

// DataType defines the data type of a signal.
type DataType string

const (
	TypeSigned   DataType = "Signed"
	TypeUnsigned DataType = "Unsigned"
	TypeFloat32  DataType = "Float32"
	TypeFloat64  DataType = "Float64"
)