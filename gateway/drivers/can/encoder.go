package can

import (
	"bytes"
	"fmt"
	"math"

	"github.com/dgryski/go-bitstream"
)

// Encoder defines the interface for encoding a signal value into a CAN frame.
type Encoder interface {
	Encode(signal Signal, value any) (Frame, error)
}

// DefaultEncoder is the standard implementation of the Encoder interface.
type DefaultEncoder struct{}

// NewDefaultEncoder creates a new DefaultEncoder.
func NewDefaultEncoder() *DefaultEncoder {
	return &DefaultEncoder{}
}

// Encode converts a physical value into a CAN frame's data payload.
func (e *DefaultEncoder) Encode(signal Signal, value any) (Frame, error) {
	floatValue, ok := value.(float64)
	if !ok {
		// Attempt to convert from other numeric types
		switch v := value.(type) {
		case int:
			floatValue = float64(v)
		case int32:
			floatValue = float64(v)
		case int64:
			floatValue = float64(v)
		case float32:
			floatValue = float64(v)
		default:
			return Frame{}, fmt.Errorf("unsupported value type for encoding: %T", value)
		}
	}

	// Reverse the scaling and offset
	rawValueFloat := (floatValue - signal.Offset) / signal.Scale
	var rawValue uint64

	switch signal.DataType {
	case TypeUnsigned, TypeSigned:
		rawValue = uint64(rawValueFloat)
	case TypeFloat32:
		rawValue = uint64(math.Float32bits(float32(rawValueFloat)))
	case TypeFloat64:
		rawValue = math.Float64bits(rawValueFloat)
	default:
		return Frame{}, fmt.Errorf("unsupported data type for encoding: %s", signal.DataType)
	}

	// Create a buffer and a bitstream writer.
	buf := new(bytes.Buffer)
	writer := bitstream.NewWriter(buf)

	// A more robust implementation to handle bit-level placement.
	// We construct the full 64-bit data payload.
	var dataPayload uint64

	// Create a mask for the signal's length
	mask := (uint64(1) << signal.Length) - 1
	// Apply mask to ensure value fits
	rawValue &= mask

	// Shift the raw value to its start position.
	// Note: Bit numbering can be ambiguous (from LSB or MSB).
	// Assuming start bit is from LSB (e.g., bit 0).
	// A common convention is that bit 0 is the LSB of the first byte.
	// The go-bitstream library writes bits from LSB to MSB.
	// To place a signal, we need to consider the byte order.
	// This is a complex topic. For now, we'll use a simplified model
	// that may need adjustment based on specific DBC file conventions.
	// Let's assume a simple left shift for big-endian representation within the 64-bit word.
	dataPayload |= rawValue << signal.StartBit

	// Write the 64-bit payload to the buffer.
	// The library writes LSB first, so we might need to handle byte swapping for Big Endian.
	// Let's write as little-endian for now and see.
	err := writer.WriteBits(dataPayload, 64)
	if err != nil {
		return Frame{}, fmt.Errorf("bitstream.WriteBits failed: %w", err)
	}
	writer.Flush(bitstream.Zero)

	data := buf.Bytes()
	// Ensure data is 8 bytes long, padding if necessary.
	if len(data) < 8 {
		paddedData := make([]byte, 8)
		copy(paddedData, data)
		data = paddedData
	}

	frame := Frame{
		ID:       signal.FrameID,
		DLC:      8, // Assume 8 for now
		Data:     data,
		Extended: signal.FrameID > 0x7FF,
	}

	return frame, nil
}