package codec

import (
	"encoding/binary"
	"fmt"
)

func init() {
	RegisterDecoder(DataTypeUint16, decodeUint16)
}

func decodeUint16(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("not enough bytes for uint16")
	}

	var val uint16
	switch byteOrder {
	case BigEndian:
		val = binary.BigEndian.Uint16(data)
	case LittleEndian:
		val = binary.LittleEndian.Uint16(data)
	default:
		return nil, fmt.Errorf("unsupported byte order: %s", byteOrder)
	}
	return val, nil
}