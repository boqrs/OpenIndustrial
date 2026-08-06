package modbus

import (
	"context"
	"errors" // Import errors for ErrInvalidConfig
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/drivers/modbus/internal/codec" // Import internal/codec for decoding

	modbusClient "github.com/goburrow/modbus"
)

// ErrInvalidConfig is returned when the Modbus configuration is invalid.
var ErrInvalidConfig = errors.New("invalid modbus config")

// Adapter 定义 Modbus 协议访问能力。
//
// Adapter 屏蔽所有 Modbus 协议细节：
//
//	Coil
//	Discrete Input
//	Holding Register
//	Input Register
//
// Runtime 永远不会直接接触 Goburrow SDK。
//
// Adapter 是整个 Modbus Driver 的唯一协议入口。
type Adapter interface {

	// Connect 建立连接。
	//
	// Connect 应当支持重复调用。
	// 如果已经连接，应直接返回 nil。
	Connect(ctx context.Context) error

	// Close 关闭连接。
	//
	// Close 应当支持重复调用。
	Close() error

	// Connected 返回当前连接状态。
	Connected() bool

	// Read 读取一个 PointMapping。
	//
	// Adapter 内部负责：
	//
	//	RegisterType
	//	Address
	//	Codec
	//	Endian
	//
	// 最终统一返回 Sample。
	Read(ctx context.Context, mapping NodeMapping) (Sample, error)

	// ReadBatch 批量读取。
	//
	// ReadBatch 可以根据地址连续性自动优化读取次数。
	ReadBatch(ctx context.Context, mappings []NodeMapping) ([]Sample, error)

	// Write 写入一个 Point。
	Write(ctx context.Context, mapping NodeMapping, value any) error
}

// ModbusConnectionHandler extends modbus.ClientHandler with Connect and Close methods.
// Concrete implementations like *modbus.TCPClientHandler and *modbus.RTUClientHandler
// already provide these methods.
type ModbusConnectionHandler interface {
	modbusClient.ClientHandler
	Connect() error
	Close() error
}

// ModbusAdapter 是 Adapter 接口的 Modbus 实现。
type ModbusAdapter struct {
	config      *Config
	handler     ModbusConnectionHandler    // Modbus 客户端处理器 (TCP 或 RTU)
	client      modbusClient.Client        // Modbus 客户端
	mu          sync.Mutex                 // 保护连接状态
	isConnected bool
	optimizer   *Optimizer                 // Modbus 请求优化器
}

// NewModbusAdapter 创建一个新的 ModbusAdapter 实例。
func NewModbusAdapter(cfg *Config) (Adapter, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	adapter := &ModbusAdapter{
		config:    cfg,
		optimizer: NewOptimizer(), // 初始化优化器
	}

	var err error
	// 根据连接模式初始化 handler
	switch cfg.Connection.Mode {
	case ModeTCP:
		adapter.handler, err = newTCPClientHandler(cfg.Connection)
	case ModeRTU:
		adapter.handler, err = newRTUClientHandler(cfg.Connection)
	default:
		return nil, fmt.Errorf("unsupported connection mode: %s", cfg.Connection.Mode)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create modbus client handler: %w", err)
	}

	adapter.client = modbusClient.NewClient(adapter.handler)
	return adapter, nil
}

// Connect implements the Adapter interface.
func (a *ModbusAdapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isConnected {
		return nil // Already connected
	}

	err := a.handler.Connect()
	if err != nil {
		a.isConnected = false
		return fmt.Errorf("modbus connect failed: %w", err)
	}
	a.isConnected = true
	return nil
}

// Close implements the Adapter interface.
func (a *ModbusAdapter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isConnected {
		return nil // Already closed or not connected
	}

	err := a.handler.Close()
	if err != nil {
		return fmt.Errorf("modbus close failed: %w", err)
	}
	a.isConnected = false
	return nil
}

// Connected implements the Adapter interface.
func (a *ModbusAdapter) Connected() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isConnected
}

// Read implements the Adapter interface.
func (a *ModbusAdapter) Read(ctx context.Context, mapping NodeMapping) (Sample, error) {
	if !a.Connected() {
		return Sample{PointID: mapping.PointID, Quality: Disconnected, Timestamp: time.Now(), Value: nil}, fmt.Errorf("modbus not connected")
	}

	var results []byte
	var err error

	// Perform Modbus read based on RegisterType
	switch mapping.Register {
	case RegisterCoil:
		results, err = a.client.ReadCoils(mapping.Address, mapping.Length)
	case RegisterDiscreteInput:
		results, err = a.client.ReadDiscreteInputs(mapping.Address, mapping.Length)
	case RegisterHoldingRegister:
		results, err = a.client.ReadHoldingRegisters(mapping.Address, mapping.Length)
	case RegisterInputRegister:
		results, err = a.client.ReadInputRegisters(mapping.Address, mapping.Length)
	default:
		return Sample{PointID: mapping.PointID, Quality: Bad, Timestamp: time.Now(), Value: nil}, fmt.Errorf("unsupported register type: %s", mapping.Register)
	}

	if err != nil {
		return Sample{PointID: mapping.PointID, Quality: Bad, Timestamp: time.Now(), Value: nil}, fmt.Errorf("modbus read failed for point %s: %w", mapping.PointID, err)
	}

	// Decode the raw bytes using the internal codec
	decodedValue, err := codec.Decode(results, mapping.DataType, mapping.ByteOrder, mapping.WordOrder)
	if err != nil {
		return Sample{PointID: mapping.PointID, Quality: Bad, Timestamp: time.Now(), Value: nil}, fmt.Errorf("failed to decode value for point %s: %w", mapping.PointID, err)
	}

	// Apply scaling and offset if defined
	if mapping.Scale != 0 || mapping.Offset != 0 {
		switch v := decodedValue.(type) {
		case int16:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case uint16:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case int32:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case uint32:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case int64:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case uint64:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case float32:
			decodedValue = float64(v)*mapping.Scale + mapping.Offset
		case float64:
			decodedValue = v*mapping.Scale + mapping.Offset
		// Add other numeric types if necessary
		}
	}

	return Sample{
		PointID:   mapping.PointID,
		Value:     decodedValue,
		Timestamp: time.Now(),
		Quality:   Good,
	}, nil
}

// ReadBatch implements the Adapter interface.
// ReadBatch 从 Modbus 设备批量读取多个点的值。
func (a *ModbusAdapter) ReadBatch(ctx context.Context, mappings []NodeMapping) ([]Sample, error) {
	if !a.Connected() {
		return nil, fmt.Errorf("modbus not connected")
	}

	optimizedBatches := a.optimizer.Optimize(mappings)
	allSamples := make([]Sample, 0, len(mappings))
	samplesMap := make(map[string]Sample)

	for _, batch := range optimizedBatches {
		var results []byte
		var err error
		switch batch.Register {
		case RegisterCoil:
			results, err = a.client.ReadCoils(batch.Start, batch.Count)
		case RegisterDiscreteInput:
			results, err = a.client.ReadDiscreteInputs(batch.Start, batch.Count)
		case RegisterHoldingRegister:
			results, err = a.client.ReadHoldingRegisters(batch.Start, batch.Count)
		case RegisterInputRegister:
			results, err = a.client.ReadInputRegisters(batch.Start, batch.Count)
		default:
			err = fmt.Errorf("unsupported register type for batch: %s", batch.Register)
		}

		if err != nil {
			for _, m := range batch.Points {
				samplesMap[m.PointID] = Sample{PointID: m.PointID, Quality: Bad, Timestamp: time.Now(), Value: nil}
			}
			log.Printf("Modbus ReadBatch: Failed to read batch for register type %s, start %d, count %d: %v", batch.Register, batch.Start, batch.Count, err)
			continue
		}

		// Implement efficient decoding from batch results
		for _, m := range batch.Points {
			var pointData []byte
			var decodeErr error
			var decodedValue interface{}

			switch m.Register {
			case RegisterCoil, RegisterDiscreteInput:
				// For coils/discrete inputs, results are bit-packed.
				// Calculate the bit offset within the results byte slice.
				relativeAddress := m.Address - batch.Start
				byteIndex := relativeAddress / 8
				bitOffset := relativeAddress % 8

				if int(byteIndex) < len(results) {
					// Extract the bit value
					bitValue := (results[byteIndex] >> bitOffset) & 0x01
					if bitValue == 1 {
						pointData = []byte{0x01} // Represent true as 0x01 for codec
					} else {
						pointData = []byte{0x00} // Represent false as 0x00 for codec
					}
					decodedValue, decodeErr = codec.Decode(pointData, codec.DataTypeBool, m.ByteOrder, m.WordOrder)
				} else {
					decodeErr = fmt.Errorf("bit index out of bounds for coil/discrete input")
				}

			case RegisterHoldingRegister, RegisterInputRegister:
				// For registers, results are byte-packed (2 bytes per register).
				// Calculate the byte offset within the results byte slice.
				relativeAddress := m.Address - batch.Start
				byteOffset := relativeAddress * 2 // Each register is 2 bytes
				numBytes := m.RegisterCount() * 2 // Total bytes for this point

				if int(byteOffset+numBytes) <= len(results) {
					pointData = results[byteOffset : byteOffset+numBytes]
					decodedValue, decodeErr = codec.Decode(pointData, m.DataType, m.ByteOrder, m.WordOrder)
				} else {
					decodeErr = fmt.Errorf("byte index out of bounds for register")
				}
			}

			if decodeErr != nil {
				samplesMap[m.PointID] = Sample{PointID: m.PointID, Quality: Bad, Timestamp: time.Now(), Value: nil}
				log.Printf("Modbus ReadBatch: Failed to decode point %s (Address: %d, DataType: %s): %v", m.PointID, m.Address, m.DataType, decodeErr)
			} else {
				// Apply scaling and offset
				finalValue := applyScalingAndOffset(decodedValue, m.Scale, m.Offset)
				samplesMap[m.PointID] = Sample{PointID: m.PointID, Quality: Good, Timestamp: time.Now(), Value: finalValue}
			}
		}
	}

	for _, m := range mappings {
		if sample, ok := samplesMap[m.PointID]; ok {
			allSamples = append(allSamples, sample)
		} else {
			// Should not happen if all mappings are covered by batches, but as a fallback
			allSamples = append(allSamples, Sample{PointID: m.PointID, Quality: Bad, Timestamp: time.Now(), Value: nil})
		}
	}
	return allSamples, nil
}

// applyScalingAndOffset applies scale and offset to a decoded value.
// It handles various numeric types.
func applyScalingAndOffset(value interface{}, scale float64, offset float64) interface{} {
	if scale == 1.0 && offset == 0.0 {
		return value
	}

	switch v := value.(type) {
	case int16:
		return float64(v)*scale + offset
	case uint16:
		return float64(v)*scale + offset
	case int32:
		return float64(v)*scale + offset
	case uint32:
		return float64(v)*scale + offset
	case int64:
		return float64(v)*scale + offset
	case uint64:
		return float64(v)*scale + offset
	case float32:
		return float64(v)*scale + offset
	case float64:
		return v*scale + offset
	default:
		return value // Cannot apply to non-numeric types
	}
}


// Write implements the Adapter interface.
func (a *ModbusAdapter) Write(ctx context.Context, mapping NodeMapping, value any) error {
	if !a.Connected() {
		return fmt.Errorf("modbus not connected")
	}

	// TODO: Implement Modbus write logic.
	// This will involve encoding the value to bytes using internal/codec (if applicable)
	// and then calling appropriate a.client.Write... functions.
	return fmt.Errorf("Write not implemented yet")
}