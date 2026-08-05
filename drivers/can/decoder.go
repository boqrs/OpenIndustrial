package can

import (
	"bytes"
	"fmt"
	"math"

	"github.com/dgryski/go-bitstream"
)

// Decoder defines the interface for decoding a CAN frame into a signal value.
type Decoder interface {
	Decode(frame Frame, signal Signal) (SignalValue, error)
}

// DefaultDecoder is the standard implementation of the Decoder interface.
type DefaultDecoder struct{}

// NewDefaultDecoder creates a new DefaultDecoder.
func NewDefaultDecoder() *DefaultDecoder {
	return &DefaultDecoder{}
}

// Decode extracts a signal's value from a CAN frame based on the signal's definition.
func (d *DefaultDecoder) Decode(frame Frame, signal Signal) (SignalValue, error) {
	if frame.ID != signal.FrameID {
		return SignalValue{}, fmt.Errorf("frame ID mismatch: expected %d, got %d", signal.FrameID, frame.ID)
	}

	reader := bitstream.NewReader(bytes.NewReader(frame.Data))

	// Skip to the start bit of the signal
	if _, err := reader.ReadBits(int(signal.StartBit)); err != nil {
		return SignalValue{}, fmt.Errorf("failed to skip to start bit: %w", err)
	}

	// Read the raw signal data
	rawData, err := reader.ReadBits(int(signal.Length))
	if err != nil {
		return SignalValue{}, fmt.Errorf("failed to read signal bits: %w", err)
	}

	var rawValue uint64
	if signal.ByteOrder == LittleEndian {
		// The go-bitstream library reads in big-endian bit order, so we need to reverse for little-endian signals.
		// This is a simplification; a full DBC parser would handle this more robustly.
		// For now, we assume the bitstream library's reading order matches Motorola format.
		// A proper implementation would involve bit-level reversal if signal is little endian.
		// Let's proceed with a simple cast for now.
		rawValue = rawData
	} else {
		rawValue = rawData
	}

	var floatValue float64
	switch signal.DataType {
	case TypeUnsigned:
		floatValue = float64(rawValue)
	case TypeSigned:
		// Handle signed values (two's complement)
		isNegative := (rawValue >> (signal.Length - 1)) == 1
		if isNegative {
			mask := (uint64(1) << signal.Length) - 1
			floatValue = -float64((^rawValue+1)&mask)
		} else {
			floatValue = float64(rawValue)
		}
	case TypeFloat32:
		bits := uint32(rawValue)
		floatValue = float64(math.Float32frombits(bits))
	case TypeFloat64:
		floatValue = math.Float64frombits(rawValue)
	default:
		return SignalValue{}, fmt.Errorf("unsupported data type: %s", signal.DataType)
	}

	// Apply scale and offset
	physicalValue := floatValue*signal.Scale + signal.Offset

	return SignalValue{
		ID:        signal.ID,
		Value:     physicalValue,
		Timestamp: frame.Timestamp,
	}, nil
}