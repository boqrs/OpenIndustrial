package modbus

import (
	"fmt"

	modbusClient "github.com/goburrow/modbus"
)

// newTCPClientHandler creates and configures a Modbus TCP client handler.
func newTCPClientHandler(cfg ConnectionConfig) (ModbusConnectionHandler, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("modbus TCP address cannot be empty")
	}

	handler := modbusClient.NewTCPClientHandler(cfg.Address)
	handler.Timeout = cfg.Timeout
	// handler.IdleTimeout = cfg.IdleTimeout // Removed as IdleTimeout is not in ConnectionConfig
	handler.SlaveId = cfg.SlaveID

	return handler, nil
}