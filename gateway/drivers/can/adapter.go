package can

import "context"

// ConnectionConfig holds the configuration for connecting to a CAN interface.
type ConnectionConfig struct {
	Interface string `json:"interface"`
}

// Adapter is the interface for a CAN bus adapter.
// It provides a hardware-agnostic way to send and receive CAN frames.
type Adapter interface {
	// Connect establishes a connection to the CAN bus.
	Connect(ctx context.Context, cfg ConnectionConfig) error
	// Disconnect closes the connection to the CAN bus.
	Disconnect(ctx context.Context) error
	// IsConnected returns true if the adapter is connected to the bus.
	IsConnected() bool
	// Receive returns a read-only channel that emits received CAN frames.
	Receive() <-chan Frame
	// Send transmits a CAN frame to the bus.
	Send(ctx context.Context, frame Frame) error
}