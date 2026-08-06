package ethernetip

import "context"

// Adapter is the core interface for interacting with an EtherNet/IP device.
// It abstracts the underlying messaging protocol (Explicit vs. Implicit).
type Adapter interface {
	// Connect establishes a connection to the device and registers a session.
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// Disconnect unregisters the session and closes the connection.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the adapter is connected and has a valid session.
	IsConnected() bool

	// Read reads a set of points (tags or CIP objects) from the device.
	// This is typically used for polled, non-real-time data.
	Read(ctx context.Context, points []PointMapping) ([]Sample, error)

	// Write writes a value to a single point on the device.
	Write(ctx context.Context, point PointMapping, value interface{}) error

	// Subscribe sets up a subscription for real-time data.
	// For Explicit mode, this might start a polling loop internally.
	// For Implicit mode, this would set up a UDP listener for Assembly data.
	Subscribe(ctx context.Context, mappings []PointMapping, ch chan<- Sample) error
}