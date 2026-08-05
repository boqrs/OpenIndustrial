package codec

import (
	"encoding/binary"
	"fmt"
)

func init() {
	RegisterDecoder(DataTypeInt16, decodeInt16)
}

func decodeInt16(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("not enough bytes for int16")
	}

	var val int16
	switch byteOrder {
	case BigEndian:
		val = int16(binary.BigEndian.Uint16(data))
	case LittleEndian:
		val = int16(binary.LittleEndian.Uint16(data))
	default:
		return nil, fmt.Errorf("unsupported byte order: %s", byteOrder)
	}
	return val, nil
}