package modbus

import (
	"fmt"

	modbusClient "github.com/goburrow/modbus"
)

// newRTUClientHandler creates and configures a Modbus RTU client handler.
func newRTUClientHandler(cfg ConnectionConfig) (ModbusConnectionHandler, error) {
	if cfg.SerialPort == "" {
		return nil, fmt.Errorf("modbus RTU serial port cannot be empty")
	}
	handler := modbusClient.NewRTUClientHandler(cfg.SerialPort)
	handler.Timeout = cfg.Timeout
	handler.SlaveId = cfg.SlaveID
	handler.BaudRate = cfg.BaudRate
	handler.DataBits = cfg.DataBits
	handler.Parity = cfg.Parity
	handler.StopBits = cfg.StopBits
	return handler, nil
}