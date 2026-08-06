package profinet

import "time"

// ConnectionConfig defines the configuration for the PROFINET controller adapter.
type ConnectionConfig struct {
	// Interface is the name of the network interface to use (e.g., "eth0", "ens33").
	Interface string

	// LocalIP is the IP address of the controller on the PROFINET network.
	LocalIP string

	// CycleTime is the desired cyclic data exchange interval.
	// Common values are 1ms, 2ms, 4ms, 8ms.
	CycleTime time.Duration

	// Timeout for various network operations.
	Timeout time.Duration
}

// DeviceConfig defines the configuration for a specific IO Device to connect to.
type DeviceConfig struct {
	// Name is a user-defined name for the device.
	Name string

	// IP address of the device.
	IP string

	// StationName is the PROFINET station name of the device.
	StationName string

	// GSDML is the path or content of the GSDML file for this device.
	// This is crucial for understanding the device's modules and data layout.
	GSDML string

	// Modules defines the expected module configuration for the device.
	Modules []ModuleConfig
}

// ModuleConfig defines a single module within a device.
type ModuleConfig struct {
	// Slot number where the module is located.
	Slot uint16

	// SubSlot number within the slot.
	SubSlot uint16

	// InputSize is the size of the input data area in bytes.
	InputSize uint16

	// OutputSize is the size of the output data area in bytes.
	OutputSize uint16
}