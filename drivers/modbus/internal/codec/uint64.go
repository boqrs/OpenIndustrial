package codec

import (
	"fmt"
)

func init() {
	RegisterDecoder(DataTypeUint64, decodeUint64)
}

func decodeUint64(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("not enough bytes for uint64")
	}

	var u64 uint64
	switch wordOrder {
	case HighWordFirst: // [HW1 LW1 HW2 LW2]
		w1, err := getUint16(data[0:2], byteOrder)
		if err != nil {
			return nil, err
		}
		w2, err := getUint16(data[2:4], byteOrder)
		if err != nil {
			return nil, err
		}
		w3, err := getUint16(data[4:6], byteOrder)
		if err != nil {
			return nil, err
		}
		w4, err := getUint16(data[6:8], byteOrder)
		if err != nil {
			return nil, err
		}
		u64 = uint64(w1)<<48 | uint64(w2)<<32 | uint64(w3)<<16 | uint64(w4)
	case LowWordFirst: // [LW1 HW1 LW2 HW2]
		w1, err := getUint16(data[0:2], byteOrder)
		if err != nil {
			return nil, err
		}
		w2, err := getUint16(data[2:4], byteOrder)
		if err != nil {
			return nil, err
		}
		w3, err := getUint16(data[4:6], byteOrder)
		if err != nil {
			return nil, err
		}
		w4, err := getUint16(data[6:8], byteOrder)
		if err != nil {
			return nil, err
		}
		u64 = uint64(w2)<<48 | uint64(w1)<<32 | uint64(w4)<<16 | uint64(w3)
	default:
		return nil, fmt.Errorf("unsupported word order: %s", wordOrder)
	}

	return u64, nil
}