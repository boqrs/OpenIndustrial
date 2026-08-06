package iec104

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// MarshalASDU encodes a structured ASDU into a raw byte slice.
func MarshalASDU(asdu *ASDU) ([]byte, error) {
	buf := new(bytes.Buffer)

	// TypeID (1 byte)
	buf.WriteByte(byte(asdu.TypeID))

	// VSQ (1 byte)
	buf.WriteByte(asdu.VSQ)

	// Cause of Transmission (2 bytes)
	binary.Write(buf, binary.LittleEndian, asdu.Cause)

	// Common Address (2 bytes)
	binary.Write(buf, binary.LittleEndian, asdu.CommonAddress)

	// Information Objects
	isSequential := (asdu.VSQ & 0x80) != 0
	if len(asdu.Objects) > 0 {
		if isSequential {
			// Write the first IOA
			ioa := asdu.Objects[0].IOA
			buf.WriteByte(byte(ioa & 0xFF))
			buf.WriteByte(byte((ioa >> 8) & 0xFF))
			buf.WriteByte(byte((ioa >> 16) & 0xFF))
		}

		for _, obj := range asdu.Objects {
			if !isSequential {
				// Write individual IOA for each object
				ioa := obj.IOA
				buf.WriteByte(byte(ioa & 0xFF))
				buf.WriteByte(byte((ioa >> 8) & 0xFF))
				buf.WriteByte(byte((ioa >> 16) & 0xFF))
			}

			// Marshal the object's value based on the ASDU TypeID
			valueBytes, err := marshalInformationObjectValue(asdu.TypeID, obj.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal value for IOA %d: %w", obj.IOA, err)
			}
			buf.Write(valueBytes)
		}
	}

	return buf.Bytes(), nil
}

// marshalInformationObjectValue encodes the value part of an Information Object.
func marshalInformationObjectValue(typeID TypeID, value interface{}) ([]byte, error) {
	switch typeID {
	case C_SC_NA_1: // Single command
		if v, ok := value.(bool); ok {
			if v {
				return []byte{0x01}, nil // ON
			} else {
				return []byte{0x00}, nil // OFF
			}
		}
		return nil, fmt.Errorf("invalid value type for C_SC_NA_1: expected bool, got %T", value)

	case C_IC_NA_1: // Interrogation command
		if v, ok := value.(int); ok && v == 20 {
			return []byte{20}, nil // QOI = 20
		}
		return nil, fmt.Errorf("invalid value for C_IC_NA_1: expected QOI=20, got %v", value)

	// TODO: Add cases for other command TypeIDs (e.g., C_SE_NC_1 for setpoints)

	default:
		return nil, fmt.Errorf("marshaling for TypeID %d is not supported", typeID)
	}
}