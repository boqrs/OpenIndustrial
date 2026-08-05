package dlt645

import "time"

// ConnectionType defines the physical layer connection type.
type ConnectionType string

const (
	Serial ConnectionType = "serial"
	TCP    ConnectionType = "tcp"
)

// ConnectionConfig holds the configuration for connecting to a device or gateway.
type ConnectionConfig struct {
	// Type specifies the connection type: "serial" or "tcp".
	Type ConnectionType `yaml:"type"`

	// Address is the device address. For serial, it's the port (e.g., "/dev/ttyUSB0" on Linux or "COM3" on Windows).
	// For TCP, it's the host and port (e.g., "192.168.1.10:502").
	Address string `yaml:"address"`

	// Timeout for read/write operations.
	Timeout time.Duration `yaml:"timeout"`

	// Serial port specific settings.
	BaudRate int    `yaml:"baudRate"`
	DataBits int    `yaml:"dataBits"`
	StopBits int    `yaml:"stopBits"`
	Parity   string `yaml:"parity"` // "N" (None), "E" (Even), "O" (Odd)
}