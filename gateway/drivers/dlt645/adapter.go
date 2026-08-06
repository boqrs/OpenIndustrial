package dlt645

import "context"

// Adapter is the interface that abstracts the underlying communication layer (e.g., Serial, TCP).
// It provides a unified way for the driver to interact with different transport protocols.
type Adapter interface {
	// Connect establishes a connection to the device or gateway.
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// Disconnect closes the connection.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the connection is currently active.
	IsConnected() bool

	// Read sends a request to read data points from a specific meter address and returns the collected samples.
	// The adapter is responsible for framing the request, sending it, receiving the response,
	// and parsing it into the unified Sample model.
	Read(ctx context.Context, meterAddress string, points []PointMapping) ([]Sample, error)

	// Write sends a command to a specific meter.
	Write(ctx context.Context, req WriteRequest) error
}