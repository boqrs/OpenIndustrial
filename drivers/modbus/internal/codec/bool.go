package codec

import "fmt"

func init() {
	RegisterDecoder(DataTypeBool, decodeBool)
}

func decodeBool(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("not enough bytes for bool")
	}
	// Assuming a single byte represents a boolean, where 0 is false and non-zero is true.
	// This might need refinement based on specific Modbus coil/discrete input representation.
	return data[0] != 0, nil
}