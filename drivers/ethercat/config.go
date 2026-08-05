package ethercat

import "time"

// ConnectionConfig defines the configuration for the EtherCAT master.
type ConnectionConfig struct {
	// Interface is the name of the network interface to use (e.g., "eth0", "enp3s0").
	// This is a mandatory field.
	Interface string

	// CycleTime is the desired cyclic task interval for PDO exchange.
	// Common values are 1ms, 2ms, 4ms.
	CycleTime time.Duration

	// Timeout for various operations like SDO access.
	Timeout time.Duration
}