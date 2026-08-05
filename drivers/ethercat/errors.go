package ethercat

import "errors"

var (
	// ErrNotConnected is returned when an operation is attempted on a disconnected adapter.
	ErrNotConnected = errors.New("ethercat: adapter is not connected")

	// ErrAlreadyConnected is returned when Connect is called on an already connected adapter.
	ErrAlreadyConnected = errors.New("ethercat: adapter is already connected")

	// ErrInterfaceNotSpecified is returned when the network interface is not provided in the config.
	ErrInterfaceNotSpecified = errors.New("ethercat: network interface not specified in config")

	// ErrSlaveNotFound is returned when an operation is requested for a non-existent slave.
	ErrSlaveNotFound = errors.New("ethercat: slave not found on the bus")

	// ErrSDOAbort is returned when a slave sends an SDO abort code.
	ErrSDOAbort = errors.New("ethercat: SDO operation aborted by slave")

	// ErrTimeout is returned when an operation does not complete within the specified timeout.
	ErrTimeout = errors.New("ethercat: operation timed out")
)