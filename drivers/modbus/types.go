package modbus

import (
	"fmt"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/modbus/internal/codec" // Import internal/codec for decoding
)

// RegisterType 表示 Modbus 数据区类型。
type RegisterType uint8

const (
	// Coil (0xxxx)
	RegisterCoil RegisterType = iota

	// Discrete Input (1xxxx)
	RegisterDiscreteInput

	// Input Register (3xxxx)
	RegisterInputRegister

	// Holding Register (4xxxx)
	RegisterHoldingRegister
)

func (t RegisterType) String() string {
	switch t {
	case RegisterCoil:
		return "coil"

	case RegisterDiscreteInput:
		return "discrete_input"

	case RegisterInputRegister:
		return "input_register"

	case RegisterHoldingRegister:
		return "holding_register"

	default:
		return "unknown"
	}
}

// NodeMapping
//
// Runtime 永远只认识 Node。
// Modbus Driver 负责 Node 与 Register 的映射。
type NodeMapping struct {

	// Runtime Node ID
	PointID string `json:"pointId" yaml:"pointId"`

	// Holding Register
	Register RegisterType `json:"register" yaml:"register"`

	// 起始地址
	Address uint16 `json:"address" yaml:"address"`

	// 数据类型 (使用 internal/codec 中的 DataType)
	DataType codec.DataType `json:"dataType" yaml:"dataType"`

	// 占用寄存器数量
	Length uint16 `json:"length" yaml:"length"`

	// 字节序 (使用 internal/codec 中的 Endian)
	ByteOrder codec.Endian `json:"byteOrder" yaml:"byteOrder"`

	// 字序 (使用 internal/codec 中的 WordOrder)
	WordOrder codec.WordOrder `json:"wordOrder" yaml:"wordOrder"`

	// 是否允许写
	Writable bool `json:"writable" yaml:"writable"`

	// 工程量缩放
	Scale float64 `json:"scale" yaml:"scale"`

	// 工程量偏移
	Offset float64 `json:"offset" yaml:"offset"`

	// 描述
	Description string `json:"description" yaml:"description"`
}

// Validate 验证 NodeMapping 的有效性。
func (nm *NodeMapping) Validate() error {
	if nm.PointID == "" {
		return fmt.Errorf("%w: node mapping PointID cannot be empty", ErrInvalidConfig)
	}
	if nm.Register == 0 {
		return fmt.Errorf("%w: node mapping %s register type cannot be empty", ErrInvalidConfig, nm.PointID)
	}
	if nm.DataType == "" {
		return fmt.Errorf("%w: node mapping %s data type cannot be empty", ErrInvalidConfig, nm.PointID)
	}
	// Length is automatically derived, so no need to validate here unless it's explicitly set and needs bounds checking.

	return nil
}

// RegisterCount 返回当前 NodeMapping 占用的寄存器数量。
//
// 注意：
// 这里不需要用户配置 Length。
// Runtime 根据 DataType 自动推导。
func (m NodeMapping) RegisterCount() uint16 {

	switch m.DataType {

	case codec.DataTypeBool:
		return 1

	case codec.DataTypeInt16:
		return 1

	case codec.DataTypeUint16:
		return 1

	case codec.DataTypeInt32:
		return 2

	case codec.DataTypeUint32:
		return 2

	case codec.DataTypeFloat32:
		return 2

	case codec.DataTypeInt64:
		return 4

	case codec.DataTypeUint64:
		return 4

	case codec.DataTypeFloat64:
		return 4

	default:
		return 1
	}
}

// Batch
//
// Optimizer 根据多个 NodeMapping
// 自动生成批量读取任务。
type Batch struct {

	Register RegisterType

	Start uint16

	Count uint16

	Points []NodeMapping
}

// Sample represents a collected data point.
// This is the unified output format for the Runtime.
type Sample struct {
	PointID   string    `json:"point_id"`
	Value     any       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Quality   Quality   `json:"quality"`
}

// Quality represents the quality of a collected sample.
type Quality string

const (
	Good         Quality = "good"
	Bad          Quality = "bad"
	Uncertain    Quality = "uncertain"
	Disconnected Quality = "disconnected"
)