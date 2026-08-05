package dlt645

import (
	"bytes"
)

// encodeFrame serializes a Frame struct into a byte slice ready for transmission.
func encodeFrame(frame *Frame) ([]byte, error) {
	var buf bytes.Buffer

	// Start byte
	buf.WriteByte(0x68)

	// Address (6 bytes)
	buf.Write(frame.Address)

	// Second start byte
	buf.WriteByte(0x68)

	// Control code
	buf.WriteByte(frame.Control)

	// Data length
	buf.WriteByte(byte(len(frame.Data)))

	// Data
	buf.Write(frame.Data)

	// Checksum
	cs := checksum(buf.Bytes())
	buf.WriteByte(cs)

	// End byte
	buf.WriteByte(0x16)

	return buf.Bytes(), nil
}