package ethernetip

import "time"

// ConnectionMode defines the primary communication method.
type ConnectionMode string

const (
	// ModeExplicit uses TCP for non-real-time read/write of tags (Explicit Messaging).
	ModeExplicit ConnectionMode = "explicit"
	// ModeImplicit uses UDP for real-time, cyclic I/O data exchange (Implicit Messaging).
	ModeImplicit ConnectionMode = "implicit"
)

// ConnectionConfig defines the parameters for connecting to an EtherNet/IP device.
type ConnectionConfig struct {
	// Host is the IP address of the PLC or device.
	Host string
	// Port for Explicit Messaging, typically 44818.
	Port int
	// Slot of the CPU in the chassis (e.g., for ControlLogix).
	Slot int
	// Timeout for network operations.
	Timeout time.Duration
	// Mode specifies whether to use Explicit or Implicit messaging.
	Mode ConnectionMode
}