// drivers/bacnet/config.go
package bacnet

import (
	"fmt"
	"net"
	"time"
)

// ConnectionMode defines the BACnet connection mode.
type ConnectionMode string

const (
	ConnectionModeIP   ConnectionMode = "ip"   // BACnet/IP (UDP)
	ConnectionModeMSTP ConnectionMode = "mstp" // BACnet MS/TP (Serial) - Note: MS/TP implementation is more complex
)

// ConnectionConfig defines the configuration for connecting to a BACnet device.
type ConnectionConfig struct {
	Mode ConnectionMode `json:"mode" yaml:"mode"` // Connection mode (e.g., "ip", "mstp")

	// BACnet/IP specific
	DeviceAddress string `json:"deviceAddress" yaml:"deviceAddress"` // IP address or hostname of the BACnet device
	Port          uint16 `json:"port" yaml:"port"`                   // UDP port, default 47808 (BACnet/IP standard port)
	NetworkNumber uint16 `json:"networkNumber" yaml:"networkNumber"` // BACnet network number (for routing)
	DeviceID      uint32 `json:"deviceId" yaml:"deviceId"`           // BACnet device instance ID to connect to

	// BACnet MS/TP specific (more complex, often requires a gateway)
	SerialPort string `json:"serialPort" yaml:"serialPort"` // Serial port path (e.g., "/dev/ttyUSB0")
	BaudRate   uint32 `json:"baudRate" yaml:"baudRate"`     // Baud rate for serial communication
	Master     bool   `json:"master" yaml:"master"`         // Is this device the MS/TP master?
	MaxMaster  uint8  `json:"maxMaster" yaml:"maxMaster"`   // Max master address
	MaxInfoFrames uint8 `json:"maxInfoFrames" yaml:"maxInfoFrames"` // Max info frames

	Timeout time.Duration `json:"timeout" yaml:"timeout"` // Timeout for BACnet requests
}

// Validate checks if the ConnectionConfig is valid.
func (cc *ConnectionConfig) Validate() error {
	if cc.Mode == "" {
		return fmt.Errorf("connection mode cannot be empty")
	}
	if cc.Timeout <= 0 {
		return fmt.Errorf("connection timeout must be greater than 0")
	}

	switch cc.Mode {
	case ConnectionModeIP:
		if cc.DeviceAddress == "" {
			return fmt.Errorf("device address cannot be empty for IP mode")
		}
		if cc.Port == 0 {
			cc.Port = 47808 // Default BACnet/IP port
		}
		// Basic IP address validation
		if net.ParseIP(cc.DeviceAddress) == nil {
			// Not a direct IP, could be hostname, will be resolved later
		}
		if cc.DeviceID == 0 {
			return fmt.Errorf("device ID cannot be 0 for IP mode")
		}
	case ConnectionModeMSTP:
		return fmt.Errorf("BACnet MS/TP mode is not yet fully supported in this driver example")
		// if cc.SerialPort == "" {
		// 	return fmt.Errorf("serial port cannot be empty for MS/TP mode")
		// }
		// if cc.BaudRate == 0 {
		// 	return fmt.Errorf("baud rate cannot be 0 for MS/TP mode")
		// }
		// if cc.DeviceID == 0 {
		// 	return fmt.Errorf("device ID cannot be 0 for MS/TP mode")
		// }
	default:
		return fmt.Errorf("unsupported connection mode: %s", cc.Mode)
	}
	return nil
}

// PollConfig defines configuration for polling BACnet properties.
type PollConfig struct {
	Interval time.Duration `json:"interval" yaml:"interval"` // How often to poll properties
}

// Validate checks if the PollConfig is valid.
func (pc *PollConfig) Validate() error {
	if pc.Interval <= 0 {
		return fmt.Errorf("poll interval must be greater than 0")
	}
	return nil
}

// Config defines the overall configuration for the BACnet driver.
type Config struct {
	Name string `json:"name" yaml:"name"` // Unique name for this BACnet driver instance

	Connection ConnectionConfig `json:"connection" yaml:"connection"`

	Poll PollConfig `json:"poll" yaml:"poll"`

	// NodeMappings define which BACnet objects/properties to monitor/control.
	NodeMappings []NodeMapping `json:"nodeMappings" yaml:"nodeMappings"`
}

// Validate checks if the overall Config is valid.
func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("config name cannot be empty")
	}
	if err := c.Connection.Validate(); err != nil {
		return fmt.Errorf("invalid connection config: %w", err)
	}
	if err := c.Poll.Validate(); err != nil {
		return fmt.Errorf("invalid poll config: %w", err)
	}
	if len(c.NodeMappings) == 0 {
		return fmt.Errorf("at least one node mapping must be provided")
	}
	for i, nm := range c.NodeMappings {
		if err := nm.Validate(); err != nil {
			return fmt.Errorf("invalid node mapping at index %d: %w", i, err)
		}
	}
	return nil
}
