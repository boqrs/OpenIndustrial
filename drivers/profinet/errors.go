package profinet

import "errors"

var (
	// ErrNotConnected is returned when an operation is attempted on a disconnected adapter.
	ErrNotConnected = errors.New("profinet adapter is not connected")

	// ErrAlreadyConnected is returned when Connect is called on an already connected adapter.
	ErrAlreadyConnected = errors.New("profinet adapter is already connected")

	// ErrDeviceNotFound is returned when an operation targets a device that is not found or connected.
	ErrDeviceNotFound = errors.New("device not found or session not established")

	// ErrCGOInitializationFailed is returned when the underlying C library (p-net) fails to initialize.
	ErrCGOInitializationFailed = errors.New("failed to initialize CGO layer for p-net")

	// ErrUnsupportedOnPlatform is returned when trying to use the adapter on a non-Linux platform.
	ErrUnsupportedOnPlatform = errors.New("profinet driver is only supported on Linux")
)