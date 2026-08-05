package profinet

import (
	"context"
)

// Adapter is the core interface for interacting with a PROFINET network.
// It abstracts the underlying implementation (e.g., CGO bridge to p-net).
type Adapter interface {
	// Connect initializes the PROFINET controller on a specific network interface.
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// Disconnect closes the connection and releases all resources.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the adapter is connected and running.
	IsConnected() bool

	// Discover scans the network for PROFINET devices.
	Discover(ctx context.Context) ([]DeviceInfo, error)

	// ConnectDevice establishes a connection with a specific IO Device.
	ConnectDevice(ctx context.Context, device DeviceInfo) error

	// ReadInputs reads the last received cyclic input data for a device.
	ReadInputs(ctx context.Context, deviceID string) ([]byte, error)

	// WriteOutputs sends cyclic output data to a device.
	WriteOutputs(ctx context.Context, deviceID string, data []byte) error

	// ReadRecord performs an acyclic read of a specific record from a device.
	ReadRecord(ctx context.Context, req RecordRequest) (RecordResponse, error)

	// Subscribe provides a channel for receiving a stream of unified Sample data.
	Subscribe(ctx context.Context, ch chan<- Sample) error
}