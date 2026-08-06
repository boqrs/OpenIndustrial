package ethernetip

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Decode converts a byte slice from the PLC into a Go data type.
func Decode(data []byte, dataType DataType) (interface{}, error) {
	buf := bytes.NewReader(data)

	switch dataType {
	case TypeBOOL:
		var val int8
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			return nil, err
		}
		return val != 0, nil
	case TypeSINT:
		var val int8
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeINT:
		var val int16
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeDINT:
		var val int32
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeLINT:
		var val int64
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeUSINT:
		var val uint8
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeUINT:
		var val uint16
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeUDINT:
		var val uint32
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeULINT:
		var val uint64
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeREAL:
		var val float32
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeLREAL:
		var val float64
		err := binary.Read(buf, binary.LittleEndian, &val)
		return val, err
	case TypeSTRING:
		// Rockwell PLCs often use a specific STRING structure (LEN, DATA)
		// For now, we'll assume the raw bytes are the string, which might need refinement.
		return string(data), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDataType, dataType)
	}
}

// Encode converts a Go data type into a byte slice for writing to the PLC.
func Encode(value interface{}, dataType DataType) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch dataType {
	// We'll implement the most common types first.
	case TypeBOOL:
		val, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type for BOOL: expected bool, got %T", value)
		}
		if val {
			binary.Write(buf, binary.LittleEndian, int8(1))
		} else {
			binary.Write(buf, binary.LittleEndian, int8(0))
		}
	case TypeDINT:
		val, ok := value.(int32)
		if !ok {
			// Allow conversion from standard int
			if intVal, ok := value.(int); ok {
				val = int32(intVal)
			} else {
				return nil, fmt.Errorf("invalid type for DINT: expected int32 or int, got %T", value)
			}
		}
		binary.Write(buf, binary.LittleEndian, val)
	case TypeREAL:
		val, ok := value.(float32)
		if !ok {
			// Allow conversion from float64
			if f64Val, ok := value.(float64); ok {
				val = float32(f64Val)
			} else {
				return nil, fmt.Errorf("invalid type for REAL: expected float32 or float64, got %T", value)
			}
		}
		binary.Write(buf, binary.LittleEndian, val)
	default:
		return nil, fmt.Errorf("%w: encoding for %s not implemented yet", ErrUnsupportedDataType, dataType)
	}

	return buf.Bytes(), nil
}