package dlt645

import (
	"encoding/binary"
	"encoding/hex"
)

// buildReadFrame creates a DL/T 645 frame for a read command.
// Address is a 12-char BCD string (e.g., "000000123456").
// DI is the 4-byte Data Identifier.
func buildReadFrame(address string, di uint32) ([]byte, error) {
	addrBytes, err := bcdToBytes(address)
	if err != nil {
		return nil, err
	}

	// Data field for read command: DI + 33h
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, di)
	for i := range data {
		data[i] += 0x33
	}

	frame := &Frame{
		Address: addrBytes,
		Control: 0x01, // Read data command
		Data:    data,
	}

	return encodeFrame(frame)
}

// decodeResponseFrame parses a raw byte slice into a Frame structure and validates it.
func decodeResponseFrame(data []byte) (*Frame, error) {
	if len(data) < 12 {
		return nil, ErrInvalidFrameLength
	}
	if data[0] != 0x68 || data[2] != 0x68 || data[len(data)-1] != 0x16 {
		return nil, ErrInvalidFrameStartOrEnd
	}

	// Verify checksum
	length := int(data[9])
	if len(data) != 12+length {
		return nil, ErrInvalidFrameLength
	}
	calculatedCS := checksum(data[0 : 10+length])
	expectedCS := data[10+length]
	if calculatedCS != expectedCS {
		return nil, ErrChecksumMismatch
	}

	frame := &Frame{
		Address: data[1:7],
		Control: data[8],
		Data:    data[10 : 10+length],
	}

	return frame, nil
}

// parseData extracts the value from the data payload of a response frame.
func parseData(data []byte, point PointMapping) (interface{}, error) {
	// Subtract 0x33 from each byte to get the original value
	payload := make([]byte, len(data)-4) // First 4 bytes are the DI
	for i := 4; i < len(data); i++ {
		payload[i-4] = data[i] - 0x33
	}

	var val float64

	switch point.DataType {
	case DataTypeBCD:
		// BCD decoding needs a specific implementation based on the number of bytes.
		// This is a simplified example for a 4-byte BCD.
		if len(payload) < 4 {
			return nil, ErrDataLengthMismatch
		}
		val = bcdToFloat(payload[0:4])
	default:
		return nil, ErrUnsupportedDataType
	}

	// Apply scaling factor
	if point.Scale != 0 {
		val *= point.Scale
	}

	return val, nil
}

// --- Utility Functions ---

// checksum calculates the sum of bytes modulo 256.
func checksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return sum
}

// bcdToBytes converts a 12-character BCD string address to a 6-byte slice.
func bcdToBytes(s string) ([]byte, error) {
	if len(s) != 12 {
		return nil, ErrInvalidAddressFormat
	}
	// DL/T 645 addresses are sent in reverse order (little-endian).
	reversedStr, err := hex.DecodeString(s)
	if err != nil {
		return nil, ErrInvalidAddressFormat
	}
	// Reverse the byte array
	for i, j := 0, len(reversedStr)-1; i < j; i, j = i+1, j-1 {
		reversedStr[i], reversedStr[j] = reversedStr[j], reversedStr[i]
	}
	return reversedStr, nil
}

// bcdToFloat is a simplified BCD to float converter.
// A robust implementation would handle various BCD formats.
func bcdToFloat(bcd []byte) float64 {
	var result float64
	for i := len(bcd) - 1; i >= 0; i-- {
		high := float64(bcd[i] >> 4)
		low := float64(bcd[i] & 0x0F)
		result = result*100 + high*10 + low
	}
	return result
}