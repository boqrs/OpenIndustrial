package iec104

import (
	"encoding/binary"
	"fmt"
	"math"
)

// ParseASDU decodes a raw byte slice into a structured ASDU.
func ParseASDU(data []byte) (*ASDU, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("invalid ASDU length: got %d, want at least 6", len(data))
	}

	asdu := &ASDU{
		TypeID:        TypeID(data[0]),
		VSQ:           data[1],
		Cause:         CauseOfTransmission(binary.LittleEndian.Uint16(data[2:4])),
		CommonAddress: binary.LittleEndian.Uint16(data[4:6]),
	}

	numObjects := int(asdu.VSQ & 0x7F)     // Number of objects is in the lower 7 bits
	isSequential := (asdu.VSQ & 0x80) != 0 // SQ bit: 0 = individual, 1 = sequential addresses

	var err error
	objectData := data[6:]

	// The parsing logic depends heavily on the TypeID.
	switch asdu.TypeID {
	case M_SP_NA_1: // Single-point information
		asdu.Objects, err = parseInformationObjects(objectData, numObjects, isSequential, parseSinglePoint)
	case M_DP_NA_1: // Double-point information
		asdu.Objects, err = parseInformationObjects(objectData, numObjects, isSequential, parseDoublePoint)
	case M_ME_NA_1: // Measured value, normalized value
		asdu.Objects, err = parseInformationObjects(objectData, numObjects, isSequential, parseNormalizedValue)
	case M_ME_NC_1: // Measured value, short floating point number
		asdu.Objects, err = parseInformationObjects(objectData, numObjects, isSequential, parseShortFloat)
	default:
		return nil, fmt.Errorf("unsupported TypeID: %d", asdu.TypeID)
	}

	if err != nil {
		return nil, err
	}

	return asdu, nil
}

// a generic parser function type
type objectParser func(data []byte) (interface{}, int, error)

// parseInformationObjects is a generic helper to parse a sequence of information objects.
func parseInformationObjects(data []byte, numObjects int, isSequential bool, parser objectParser) ([]InformationObject, error) {
	objects := make([]InformationObject, 0, numObjects)
	offset := 0

	var firstIOA uint32
	if isSequential {
		if len(data) < 3 {
			return nil, fmt.Errorf("not enough data for sequential IOA")
		}
		// IOA is 3 bytes in the standard
		firstIOA = uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16
		offset += 3
	}

	for i := 0; i < numObjects; i++ {
		var ioa uint32
		if !isSequential {
			if len(data[offset:]) < 3 {
				return nil, fmt.Errorf("not enough data for individual IOA")
			}
			ioa = uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16
			offset += 3
		} else {
			ioa = firstIOA + uint32(i)
		}

		value, consumed, err := parser(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse object %d: %w", i, err)
		}
		offset += consumed

		objects = append(objects, InformationObject{
			IOA:   ioa,
			Value: value,
			// Timestamp parsing would be more complex, involving CP56Time2a
		})
	}
	return objects, nil
}

// parseSinglePoint parses a single-point information value (1 byte).
func parseSinglePoint(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("not enough data for single point value")
	}
	// Value is in the first bit (0=OFF, 1=ON)
	return (data[0] & 0x01) == 1, 1, nil
}

// parseDoublePoint parses a double-point information value (1 byte).
func parseDoublePoint(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("not enough data for double point value")
	}
	// Value is in the two least significant bits (00: intermediate, 01: OFF, 10: ON, 11: intermediate)
	return data[0] & 0x03, 1, nil
}

// parseNormalizedValue parses a measured value, normalized without quality descriptor (2 bytes).
func parseNormalizedValue(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("not enough data for normalized value")
	}
	// Signed 16-bit value, scaled from -1.0 to +1.0
	value := int16(binary.LittleEndian.Uint16(data))
	return float64(value) / 32768.0, 2, nil
}

// parseShortFloat parses a short floating point value (4 bytes).
func parseShortFloat(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("not enough data for short float value")
	}
	bits := binary.LittleEndian.Uint32(data)
	return math.Float32frombits(bits), 4, nil
}