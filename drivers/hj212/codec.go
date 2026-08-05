package hj212

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/sigurn/crc16"
)

var crcTable = crc16.MakeTable(crc16.CRC16_MODBUS)

// encodeFrame takes the core data payload (the "CP" part), wraps it with the HJ/T 212 frame structure,
// calculates length and CRC, and returns the full message bytes.
func EncodeFrame(cpData []byte) ([]byte, error) {
	var buf bytes.Buffer

	// Header
	buf.WriteString("##")

	// Length (4 digits, padded with 0)
	lengthStr := fmt.Sprintf("%04d", len(cpData))
	buf.WriteString(lengthStr)

	// Data
	buf.Write(cpData)

	// CRC (4 hex characters)
	crc := crc16.Checksum(cpData, crcTable)
	crcStr := fmt.Sprintf("%04X", crc)
	buf.WriteString(crcStr)

	// Footer
	buf.WriteString("&&")

	return buf.Bytes(), nil
}

// decodeFrame takes a raw byte slice from the wire, validates the frame structure,
// checks length and CRC, and returns the clean data payload (the "CP" part).
func DecodeFrame(data []byte) ([]byte, error) {
	if len(data) < 12 { // Minimum length: ##LLLL......CRC&
		return nil, ErrInvalidLength
	}

	// 1. Check Header and Footer
	if !bytes.HasPrefix(data, []byte("##")) {
		return nil, ErrInvalidHeader
	}
	// Note: We don't check for the "&&" footer immediately, as it might be part of a larger buffer.
	// Instead, we use the length field to determine the expected end of the message.

	// 2. Parse Length
	length, err := strconv.Atoi(string(data[2:6]))
	if err != nil {
		return nil, ErrInvalidLength
	}

	// 3. Check if the buffer contains the full message
	expectedTotalLength := 2 + 4 + length + 4 + 2 // ## + LLLL + Data + CRC + &&
	if len(data) < expectedTotalLength {
		return nil, ErrInvalidLength // Not enough data yet
	}

	// 4. Extract segments
	fullMessage := data[:expectedTotalLength]
	cpData := fullMessage[6 : 6+length]
	crcStr := string(fullMessage[6+length : 6+length+4])
	footer := fullMessage[6+length+4:]

	// 5. Validate Footer
	if !bytes.Equal(footer, []byte("&&")) {
		return nil, ErrInvalidFooter
	}

	// 6. Validate CRC
	expectedCRC, err := strconv.ParseUint(crcStr, 16, 16)
	if err != nil {
		return nil, ErrCRCMismatch
	}
	calculatedCRC := crc16.Checksum(cpData, crcTable)
	if uint16(expectedCRC) != calculatedCRC {
		return nil, ErrCRCMismatch
	}

	return cpData, nil
}