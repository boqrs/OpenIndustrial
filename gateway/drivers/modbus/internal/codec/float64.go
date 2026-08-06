package codec

import (
	"fmt"
	"math"
)

func init() {
	RegisterDecoder(DataTypeFloat64, decodeFloat64)
}

func decodeFloat64(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("not enough bytes for float64")
	}

	var u64 uint64
	switch wordOrder {
	case HighWordFirst: // [HW1 LW1 HW2 LW2] -> [HW1 LW1] [HW2 LW2]
		// Assuming 4 words for float64
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

		// Reconstruct the 64-bit value based on word order
		// For HighWordFirst, it means the most significant 32-bit part comes first,
		// and within that, the high word comes first.
		// This is a common interpretation for Modbus float64.
		// (w1 w2) (w3 w4) -> (MSB 32-bit) (LSB 32-bit)
		// MSB 32-bit: w1 << 16 | w2
		// LSB 32-bit: w3 << 16 | w4
		u64 = uint64(w1)<<48 | uint64(w2)<<32 | uint64(w3)<<16 | uint64(w4)

	case LowWordFirst: // [LW1 HW1 LW2 HW2] -> [LW1 HW1] [LW2 HW2]
		// Similar logic for LowWordFirst
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

		// (w2 w1) (w4 w3) -> (MSB 32-bit) (LSB 32-bit)
		// MSB 32-bit: w2 << 16 | w1
		// LSB 32-bit: w4 << 16 | w3
		u64 = uint64(w2)<<48 | uint64(w1)<<32 | uint64(w4)<<16 | uint64(w3)

	default:
		return nil, fmt.Errorf("unsupported word order: %s", wordOrder)
	}

	return math.Float64frombits(u64), nil
}