package hj212

import "context"

// Adapter is the interface that abstracts the underlying communication layer (e.g., TCP, Serial).
// It provides a unified way for the driver to interact with different transport protocols.
type Adapter interface {
	// Connect establishes a connection to the device or gateway.
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// Disconnect closes the connection.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the connection is currently active.
	IsConnected() bool

	// ReadDataSegment blocks and waits for the next complete data segment from the device.
	// In HJ/T 212, devices often send data proactively, so this method is designed
	// to listen for and decode these incoming messages.
	ReadDataSegment(ctx context.Context) (*DataSegment, error)

	// SendCommand sends a command or request to the device.
	// For example, requesting real-time data (CN=2011).
	SendCommand(ctx context.Context, segment *DataSegment) error
}