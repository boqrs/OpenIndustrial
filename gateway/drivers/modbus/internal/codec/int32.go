package codec

import (
	"fmt"
)

func init() {
	RegisterDecoder(DataTypeInt32, decodeInt32)
}

func decodeInt32(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("not enough bytes for int32")
	}

	var u32 uint32
	switch wordOrder {
	case HighWordFirst: // [HW LW]
		// data[0:2] is HW, data[2:4] is LW
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
		// data[0:2] is LW, data[2:4] is HW
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

	return int32(u32), nil
}