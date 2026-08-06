package ethernetip

import "errors"

var (
	// ErrNotConnected is returned when an operation is attempted on a disconnected adapter.
	ErrNotConnected = errors.New("ethernet/ip adapter is not connected")

	// ErrAlreadyConnected is returned when Connect is called on an already connected adapter.
	ErrAlreadyConnected = errors.New("ethernet/ip adapter is already connected")

	// ErrSessionFailed is returned when session registration with the device fails.
	ErrSessionFailed = errors.New("failed to register session with the device")

	// ErrInvalidResponse is returned when the device sends a malformed or unexpected response.
	ErrInvalidResponse = errors.New("received invalid response from device")

	// ErrTagNotFound is returned when a requested tag does not exist on the PLC.
	ErrTagNotFound = errors.New("tag not found")

	// ErrUnsupportedDataType is returned when a data type is not supported by the codec.
	ErrUnsupportedDataType = errors.New("unsupported data type")
)