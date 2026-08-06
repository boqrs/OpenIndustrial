package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/drivers/bacnet"
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/eventbus"
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/object"
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/registry"
)

// MockRegistry 模拟 registry.Registry 接口
type MockRegistry struct{}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{}
}

func (m *MockRegistry) Add(obj object.Object) error {
	log.Printf("MockRegistry: Add object %s (mocked)", obj.GetID())
	return nil
}

func (m *MockRegistry) Get(id string) (object.Object, bool) {
	log.Printf("MockRegistry: Get object %s (mocked)", id)
	return nil, false
}

func (m *MockRegistry) Remove(id string) error {
	log.Printf("MockRegistry: Remove object %s (mocked)", id)
	return nil
}

func (m *MockRegistry) List() []object.Object {
	log.Println("MockRegistry: List objects (mocked)")
	return []object.Object{}
}

// MockEventBus 模拟 EventBus
type MockEventBus struct{}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{}
}

func (m *MockEventBus) Publish(event eventbus.Event) {
	log.Printf("EventBus Published: Topic=%s, Data=%+v", event.Topic, event.Data)
}

func (m *MockEventBus) Subscribe(topic string, handler eventbus.EventHandler) func() {
	log.Printf("EventBus Subscribed to Topic: %s", topic)
	return func() {
		log.Printf("EventBus Unsubscribed from Topic: %s (mocked)", topic)
	}
}

// MockRuntimeContext 模拟 OpenIndustrial 的运行时上下文
type MockRuntimeContext struct {
	eventBus *MockEventBus
	registry *MockRegistry
}

func NewMockRuntimeContext() *MockRuntimeContext {
	return &MockRuntimeContext{
		eventBus: NewMockEventBus(),
		registry: NewMockRegistry(),
	}
}

func (m *MockRuntimeContext) EventBus() eventbus.EventBus {
	return m.eventBus
}

func (m *MockRuntimeContext) Registry() registry.Registry {
	return m.registry
}

func (m *MockRuntimeContext) Logger() *log.Logger {
	return log.Default()
}

func main() {
	log.Println("Starting BACnet example...")

	// 1. 配置 BACnet 驱动
	// 请根据您的实际 BACnet 设备信息修改以下配置
	cfg := bacnet.Config{
		Name: "MyBACnetDriver",
		Connection: bacnet.ConnectionConfig{
			Mode:          bacnet.ConnectionModeIP,
			DeviceAddress: "192.168.1.100", // 替换为您的 BACnet 设备 IP 地址
			Port:          47808,
			DeviceID:      123, // 替换为您的 BACnet 设备实例号 (0 表示自动发现第一个)
		},
		NodeMappings: []bacnet.NodeMapping{
			{
				ID:         "analog_input_0",
				ObjectType: bacnet.ObjectTypeAnalogInput,
				Instance:   0,
				PropertyID: bacnet.PropertyIDPresentValue,
				DataType:   bacnet.DataTypeFloat, // Corrected from DataTypeReal
			},
			{
				ID:         "binary_output_1",
				ObjectType: bacnet.ObjectTypeBinaryOutput,
				Instance:   1,
				PropertyID: bacnet.PropertyIDPresentValue,
				DataType:   bacnet.DataTypeBoolean,
			},
			// 添加更多您需要读取的节点
		},
		Poll: bacnet.PollConfig{ // Corrected from PollConfig
			Interval: 5 * time.Second, // 每 5 秒轮询一次
		},
	}

	// 2. 创建驱动实例
	driver, err := bacnet.NewDriver(cfg)
	if err != nil {
		log.Fatalf("Failed to create BACnet driver: %v", err)
	}

	// // 3. 模拟运行时上下文
	// mockRuntime := NewMockRuntimeContext()

	// // 4. 初始化驱动
	// if err := driver.Init(mockRuntime); err != nil {
	// 	log.Fatalf("Failed to initialize BACnet driver: %v", err)
	// }

	// 5. 启动驱动
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := driver.Start(ctx); err != nil {
		log.Fatalf("Failed to start BACnet driver: %v", err)
	}

	log.Println("BACnet driver started. Press Ctrl+C to stop.")

	// 6. 示例：写入一个值 (如果需要)
	// 注意：写入功能目前在 adapter_gobacnet.go 中标记为 ErrNotImplemented
	// 如果 gobacnet 库支持写入，并且您实现了 WriteProperty，可以取消注释以下代码
	/*
		time.Sleep(10 * time.Second) // 等待一段时间，确保驱动已连接并开始轮询
		log.Println("Attempting to write to binary_output_1...")
		writeErr := driver.Write("binary_output_1", true) // 尝试将 BinaryOutput 1 设置为 true
		if writeErr != nil {
			log.Printf("Failed to write to binary_output_1: %v", writeErr)
		} else {
			log.Println("Successfully sent write command to binary_output_1 (check device for actual change).")
		}
	*/

	// 7. 监听操作系统信号，实现优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // 阻塞直到接收到信号

	log.Println("Stopping BACnet driver...")

	// 8. 停止驱动
	if err := driver.Stop(context.Background()); err != nil {
		log.Fatalf("Failed to stop BACnet driver: %v", err)
	}

	log.Println("BACnet example finished.")
}