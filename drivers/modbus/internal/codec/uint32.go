package codec

import (
	"fmt"
)

func init() {
	RegisterDecoder(DataTypeUint32, decodeUint32)
}

func decodeUint32(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("not enough bytes for uint32")
	}

	var u32 uint32
	switch wordOrder {
	case HighWordFirst: // [HW LW]
		hw, err := getUint16(data[0:2], byteOrder)
		if err != nil {
			return nil, err
		}
		lw, err := getUint16(data[2:4], byteOrder)
		if err != nil {
			return nil, err
		}
		u32 = uint32(hw)<<16 | uint32(lw)
	case LowWordFirst: // [LW HW]
		lw, err := getUint16(data[0:2], byteOrder)
		if err != nil {
			return nil, err
		}
		hw, err := getUint16(data[2:4], byteOrder)
		if err != nil {
			return nil, err
		}
		u32 = uint32(hw)<<16 | uint32(lw)
	default:
		return nil, fmt.Errorf("unsupported word order: %s", wordOrder)
	}

	return u32, nil
}