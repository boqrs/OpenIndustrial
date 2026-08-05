package modbus

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	runtimeContext "github.com/OpenGongChang/OpenIndustrial/runtime/context"
	"github.com/OpenGongChang/OpenIndustrial/runtime/eventbus"
	"github.com/OpenGongChang/OpenIndustrial/runtime/object"
	"github.com/OpenGongChang/OpenIndustrial/runtime/registry"
)

// Driver 定义了 Modbus 驱动的接口。
type Driver interface {
	object.Object // 驱动也是一个可注册的对象

	Init(ctx runtimeContext.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	Read(pointID string) (Sample, error)
	Write(pointID string, value any) error
}

// modbusDriver 是 Driver 接口的实现。
type modbusDriver struct {
	config *Config
	adapter Adapter // Modbus 协议适配器

	eventBus eventbus.EventBus
	registry registry.Registry

	// 用于管理轮询 Goroutine 的生命周期
	pollingCtx    context.Context
	pollingCancel context.CancelFunc
	pollingWg     sync.WaitGroup

	mu sync.RWMutex // 保护驱动状态
}

// NewDriver 创建一个新的 Modbus 驱动实例。
func NewDriver(cfg Config) (Driver, error) {
	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid modbus driver config: %w", err)
	}

	// 创建 Modbus 适配器
	adapter, err := NewModbusAdapter(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create modbus adapter: %w", err)
	}

	return &modbusDriver{
		config:  &cfg,
		adapter: adapter,
	}, nil
}

// GetID implements object.Object.
func (d *modbusDriver) GetID() string {
	return d.config.Name
}

// GetKind implements object.Object.
func (d *modbusDriver) GetKind() object.Kind {
	return "ModbusDriver"
}

// Init implements Driver.
func (d *modbusDriver) Init(ctx runtimeContext.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.eventBus = ctx.EventBus()
	d.registry = ctx.Registry()

	// 将驱动自身注册到 Registry
	if err := d.registry.Add(d); err != nil {
		return fmt.Errorf("failed to register Modbus driver %s: %w", d.GetID(), err)
	}

	log.Printf("Modbus driver %s initialized.", d.GetID())
	return nil
}

// Start implements Driver.
func (d *modbusDriver) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.pollingCtx != nil {
		return fmt.Errorf("modbus driver %s already started", d.GetID())
	}

	// 连接 Modbus 设备
	if err := d.adapter.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to modbus device for driver %s: %w", d.GetID(), err)
	}
	log.Printf("Modbus driver %s connected to device.", d.GetID())

	d.pollingCtx, d.pollingCancel = context.WithCancel(ctx)
	d.pollingWg.Add(1)
	go d.startPolling() // 启动数据轮询 Goroutine

	log.Printf("Modbus driver %s started polling.", d.GetID())
	return nil
}

// Stop implements Driver.
func (d *modbusDriver) Stop(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.pollingCancel != nil {
		d.pollingCancel() // 发送停止信号给轮询 Goroutine
		d.pollingWg.Wait() // 等待轮询 Goroutine 退出
		d.pollingCancel = nil
		d.pollingCtx = nil
	}

	// 关闭 Modbus 连接
	if err := d.adapter.Close(); err != nil {
		log.Printf("Failed to close modbus connection for driver %s: %v", d.GetID(), err)
	}
	log.Printf("Modbus driver %s connection closed.", d.GetID())

	// 从 Registry 中注销驱动
	if d.registry != nil {
		if err := d.registry.Remove(d.GetID()); err != nil {
			log.Printf("Failed to unregister Modbus driver %s from registry: %v", d.GetID(), err)
		}
	}

	log.Printf("Modbus driver %s stopped.", d.GetID())
	return nil
}

// startPolling 启动一个 Goroutine 进行周期性数据轮询。
func (d *modbusDriver) startPolling() {
	defer d.pollingWg.Done()
	ticker := time.NewTicker(d.config.Poll.Interval)
	defer ticker.Stop()

	log.Printf("Modbus driver %s polling loop started with interval %s.", d.GetID(), d.config.Poll.Interval)

	for {
		select {
		case <-d.pollingCtx.Done():
			log.Printf("Modbus driver %s polling loop stopped.", d.GetID())
			return
		case <-ticker.C:
			d.performPolling()
		}
	}
}

// performPolling 执行一次数据轮询，读取所有配置的 Modbus 点位并发布到 EventBus。
func (d *modbusDriver) performPolling() {
	if !d.adapter.Connected() {
		log.Printf("Modbus driver %s not connected, skipping polling.", d.GetID())
		return
	}

	// 收集所有 PointMapping
	var allMappings []NodeMapping
	for _, node := range d.config.NodeMappings {
		allMappings = append(allMappings, node)
	}

	// 批量读取所有点位
	samples, err := d.adapter.ReadBatch(d.pollingCtx, allMappings)
	if err != nil {
		log.Printf("Modbus driver %s failed to read batch: %v", d.GetID(), err)
		// 可以在这里发布错误事件
		return
	}

	// 发布每个 Sample 到 EventBus
	for _, sample := range samples {
		event := eventbus.Event{
			Topic: fmt.Sprintf("modbus/%s/%s", d.GetID(), sample.PointID),
			Data:  sample,
		}
		d.eventBus.Publish(event)
		log.Printf("Modbus driver %s published sample for %s: %+v", d.GetID(), sample.PointID, sample.Value)
	}
}

// Read implements Driver.
func (d *modbusDriver) Read(pointID string) (Sample, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 查找对应的 PointMapping
	var targetMapping NodeMapping
	found := false
	for _, node := range d.config.NodeMappings {
		if node.PointID == pointID {
			targetMapping = node
			found = true
			break
		}
	}
	if !found {
		return Sample{}, fmt.Errorf("point ID %s not found in configuration", pointID)
	}

	// 调用适配器进行读取
	sample, err := d.adapter.Read(d.pollingCtx, targetMapping)
	if err != nil {
		return Sample{}, fmt.Errorf("failed to read point %s: %w", pointID, err)
	}
	return sample, nil
}

// Write implements Driver.
func (d *modbusDriver) Write(pointID string, value any) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 查找对应的 PointMapping
	var targetMapping NodeMapping
	found := false
	for _, node := range d.config.NodeMappings {
		if node.PointID == pointID {
			targetMapping = node
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("point ID %s not found in configuration", pointID)
	}

	// 调用适配器进行写入
	if err := d.adapter.Write(d.pollingCtx, targetMapping, value); err != nil {
		return fmt.Errorf("failed to write point %s: %w", pointID, err)
	}
	return nil
}