package codec

import (
	"fmt"
	"strconv"
)

func init() {
	RegisterDecoder(DataTypeBCD, decodeBCD)
}

func decodeBCD(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty byte slice for BCD type")
	}

	// BCD decoding typically involves converting each nibble (4 bits) into a decimal digit.
	// The byteOrder and wordOrder might influence how multi-byte BCD values are interpreted,
	// but for a simple BCD conversion, we'll assume a direct byte-by-byte processing.
	// This implementation will convert BCD bytes into a string representation of the decimal number.
	// If a numeric type (e.g., int, float) is required, further conversion would be needed.

	var bcdStr string
	for _, b := range data {
		// High nibble
		highNibble := (b >> 4) & 0x0F
		// Low nibble
		lowNibble := b & 0x0F

		if highNibble > 9 || lowNibble > 9 {
			return nil, fmt.Errorf("invalid BCD byte encountered: %x", b)
		}
		bcdStr += strconv.Itoa(int(highNibble))
		bcdStr += strconv.Itoa(int(lowNibble))
	}

	// Remove leading zeros if the number is not "0" itself
	if len(bcdStr) > 1 && bcdStr[0] == '0' {
		bcdStr = trimLeadingZeros(bcdStr)
	}

	return bcdStr, nil
}

// trimLeadingZeros removes leading zeros from a string, unless the string is "0".
func trimLeadingZeros(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			return s[i:]
		}
	}
	return "0" // All zeros, return "0"
}