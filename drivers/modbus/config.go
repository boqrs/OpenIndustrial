package modbus

import (
	"fmt"
	"time"
)

// Mode 表示 Modbus 通讯模式。
type Mode string

const (
	ModeTCP Mode = "tcp"
	ModeRTU Mode = "rtu"
)

// ConnectionConfig 描述 Modbus 连接配置。
type ConnectionConfig struct {

	// TCP
	Address string `json:"address" yaml:"address"`

	// RTU
	SerialPort string `json:"serialPort" yaml:"serialPort"`

	BaudRate int `json:"baudRate" yaml:"baudRate"`

	DataBits int `json:"dataBits" yaml:"dataBits"`

	Parity string `json:"parity" yaml:"parity"`

	StopBits int `json:"stopBits" yaml:"stopBits"`

	SlaveID byte `json:"slaveId" yaml:"slaveId"`

	Mode Mode `json:"mode" yaml:"mode"`

	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

// PollConfig 定义采样策略。
type PollConfig struct {

	// 默认采样周期
	Interval time.Duration `json:"interval" yaml:"interval"`

	// 每次批量读取最大寄存器数量
	MaxBatchSize uint16 `json:"maxBatchSize" yaml:"maxBatchSize"`
}

// Config 为 Modbus Driver 配置。
type Config struct {
	Name string `json:"name" yaml:"name"`

	Connection ConnectionConfig `json:"connection" yaml:"connection"`

	Poll PollConfig `json:"poll" yaml:"poll"`

	NodeMappings []NodeMapping `json:"nodeMappings" yaml:"nodeMappings"` // 重命名为 NodeMappings

}

// Validate 验证 Modbus 配置的有效性。
func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("%w: driver name cannot be empty", ErrInvalidConfig)
	}
	if c.Connection.Mode == "" {
		return fmt.Errorf("%w: connection mode cannot be empty", ErrInvalidConfig)
	}
	if c.Poll.Interval <= 0 {
		return fmt.Errorf("%w: poll interval must be greater than 0", ErrInvalidConfig)
	}
	if len(c.NodeMappings) == 0 {
		return fmt.Errorf("%w: node mappings cannot be empty", ErrInvalidConfig)
	}

	for _, nm := range c.NodeMappings {
		if err := nm.Validate(); err != nil {
			return fmt.Errorf("%w: invalid node mapping %s: %w", ErrInvalidConfig, nm.PointID, err)
		}
	}

	return nil
}