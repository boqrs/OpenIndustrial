package iec104

import "errors"

var (
	// ErrNotConnected is returned when an operation is attempted on a disconnected adapter.
	ErrNotConnected = errors.New("iec104: adapter is not connected")

	// ErrAlreadyConnected is returned when Connect is called on an already connected adapter.
	ErrAlreadyConnected = errors.New("iec104: adapter is already connected")

	// ErrInvalidAPCI is returned when a received frame has an invalid APCI format.
	ErrInvalidAPCI = errors.New("iec104: invalid APCI format")

	// ErrInvalidASDU is returned when a received ASDU cannot be parsed.
	ErrInvalidASDU = errors.New("iec104: invalid ASDU format")

	// ErrTimeout is returned when an operation does not complete within the specified timeout.
	ErrTimeout = errors.New("iec104: operation timed out")

	// ErrConnectionFailed is returned when the TCP connection cannot be established.
	ErrConnectionFailed = errors.New("iec104: connection failed")
)