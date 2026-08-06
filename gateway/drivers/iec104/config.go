package iec104

import "time"

// ConnectionConfig defines the configuration for connecting to an IEC 104 slave (server).
type ConnectionConfig struct {
	// Host is the IP address or hostname of the IEC 104 slave.
	Host string
	// Port is the TCP port, usually 2404.
	Port int

	// Timeout for connection and various operations.
	Timeout time.Duration

	// CommonAddress is the address of the station (or ASDU).
	// Can be 1 or 2 bytes, but we'll use uint16.
	CommonAddress uint16

	// GeneralInterrogationInterval is the interval for performing a general interrogation (C_IC_NA_1).
	// If zero, no periodic interrogation is performed after the initial one.
	GeneralInterrogationInterval time.Duration

	// Time synchronization settings
	TimeSyncInterval time.Duration
}