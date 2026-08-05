package codec

func init() {
	RegisterDecoder(DataTypeString, decodeString)
}

func decodeString(data []byte, byteOrder Endian, wordOrder WordOrder) (interface{}, error) {
	// Assuming UTF-8 encoding for simplicity.
	// In Modbus, strings are often ASCII or specific vendor encodings.
	// This might need to be configurable or more robust in a real-world scenario.
	if len(data) == 0 {
		return "", nil // Return empty string for empty data
	}
	return string(data), nil
}