package codec

import (
	"encoding/binary"
	"fmt"
)

// Decode uses the registered decoder for the given DataType to convert raw bytes into a Go type.
// It takes the raw byte slice, the target DataType, byte order, and word order as input.
func Decode(data []byte, dataType DataType, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	decoderFn, err := GetDecoder(dataType)
	if err != nil {
		return nil, err
	}
	return decoderFn(data, byteOrder, wordOrder)
}

// getUint16 is a helper to extract a uint16 from a byte slice based on byteOrder.
func getUint16(b []byte, byteOrder Endian) (uint16, error) {
	if len(b) < 2 {
		return 0, fmt.Errorf("not enough bytes for uint16")
	}
	switch byteOrder {
	case BigEndian:
		return binary.BigEndian.Uint16(b), nil
	case LittleEndian:
		return binary.LittleEndian.Uint16(b), nil
	default:
		return 0, fmt.Errorf("unsupported byte order: %s", byteOrder)
	}
}