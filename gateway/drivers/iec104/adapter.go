package iec104

import "context"

// Adapter is the interface that abstracts the underlying IEC 104 client implementation.
type Adapter interface {
	// Connect establishes a TCP connection and starts the IEC 104 communication sequence (STARTDT).
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// Disconnect sends a STOPDT message and closes the TCP connection.
	Disconnect(ctx context.Context) error

	// IsConnected returns true if the TCP connection is alive and STARTDT has been confirmed.
	IsConnected() bool

	// Read performs a targeted read, which is not a standard IEC 104 operation.
	// This might be implemented by sending a specific interrogation command.
	// For most use cases, Subscribe is preferred.
	Read(ctx context.Context, points []PointMapping) ([]Sample, error)

	// Write sends a command or set-point to the slave device.
	// It constructs and sends a C_..._NA_1 type ASDU.
	Write(ctx context.Context, point PointMapping, value interface{}) error

	// Subscribe initiates the data transfer and listens for incoming ASDUs.
	// It performs an initial General Interrogation (C_IC_NA_1) to get all data points,
	// and then continuously processes spontaneously generated messages from the slave.
	// The received data is parsed into InformationObjects and sent to the provided channel.
	Subscribe(ctx context.Context, ch chan<- InformationObject) error
}