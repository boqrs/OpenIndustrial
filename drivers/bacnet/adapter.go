package bacnet

import (
	"context"
	"sync"
)

// Adapter 定义 BACnet 通信适配器。
// Driver、Poller 永远只依赖 Adapter，而不依赖具体 BACnet SDK。
type Adapter interface {

	// 建立连接
	Connect(ctx context.Context, cfg ConnectionConfig) error

	// 断开连接
	Disconnect(ctx context.Context) error

	// 当前是否已连接
	IsConnected() bool

	// 批量读取属性
	ReadProperties(ctx context.Context, nodes []NodeMapping) ([]Sample, error)

	// 写属性
	WriteProperty(ctx context.Context, node NodeMapping, value any) error

	// 当前连接信息
	ConnectionInfo() ConnectionInfo
}

type ConnectionInfo struct {
	Mode        ConnectionMode
	Address     string
	Port        uint16
	DeviceID    uint32
	Connected   bool
	Description string
}

// BaseAdapter 提供 Adapter 公共能力。
// 真正的 BACnet SDK Adapter 建议直接匿名嵌入 BaseAdapter。
type BaseAdapter struct {
	mu   sync.RWMutex
	info ConnectionInfo
}

func NewBaseAdapter() *BaseAdapter {
	return &BaseAdapter{}
}

func (b *BaseAdapter) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.info.Connected
}

func (b *BaseAdapter) ConnectionInfo() ConnectionInfo {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.info
}

func (b *BaseAdapter) setConnected(info ConnectionInfo) {
	b.mu.Lock()
	defer b.mu.Unlock()

	info.Connected = true
	b.info = info
}

func (b *BaseAdapter) setDisconnected() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.info.Connected = false
}