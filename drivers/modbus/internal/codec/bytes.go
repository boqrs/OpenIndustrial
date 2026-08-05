package codec

import "fmt"

func init() {
	RegisterDecoder(DataTypeBytes, decodeBytes)
}

func decodeBytes(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	// For DataTypeBytes, we simply return the raw byte slice.
	// No byte order or word order conversion is typically needed unless specified otherwise.
	if len(data) == 0 {
		return nil, fmt.Errorf("empty byte slice for bytes type")
	}
	return data, nil
}