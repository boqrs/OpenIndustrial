package bacnet

import (
	"context"
	"fmt"
	//"net"
	"sync"
	"time"

	"github.com/alexbeltran/gobacnet"
	gb "github.com/alexbeltran/gobacnet/types"
)

// ErrNotImplemented 表示功能尚未实现。
var ErrNotImplemented = fmt.Errorf("not implemented")

// GoBACnetAdapter 实现了 BACnet 协议适配器接口。
type GoBACnetAdapter struct {
	*BaseAdapter // 嵌入 BaseAdapter 以管理连接状态

	mu sync.Mutex // 互斥锁，用于保护并发访问

	client *gobacnet.Client // gobacnet 客户端实例

	target gb.Address // 目标 BACnet 设备的地址
	device gb.Device  // 存储发现的 BACnet 设备信息 (已确保存在)
}

// NewGoBACnetAdapter 创建并返回一个新的 GoBACnetAdapter 实例。
func NewGoBACnetAdapter() *GoBACnetAdapter {
	return &GoBACnetAdapter{
		BaseAdapter: NewBaseAdapter(),
	}
}

// Connect 根据配置建立与 BACnet 设备的连接。
func (a *GoBACnetAdapter) Connect(
	ctx context.Context,
	cfg ConnectionConfig,
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.IsConnected() {
		return nil // 如果已经连接，则直接返回
	}

	if err := cfg.Validate(); err != nil {
		return err // 验证连接配置
	}

	switch cfg.Mode {
	case ConnectionModeIP:
		// 将 BACnet/IP 连接和设备发现逻辑委托给 connectIPWithDiscovery 函数
		return connectIPWithDiscovery(ctx, a, cfg)

	default:
		return ErrUnsupportedMode // 不支持的连接模式
	}
}

// Disconnect 关闭与 BACnet 设备的连接。
func (a *GoBACnetAdapter) Disconnect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.client != nil {
		a.client.Close() // 关闭 gobacnet 客户端
		a.client = nil
	}

	a.setDisconnected() // 更新连接状态为断开

	return nil
}

// ReadProperties 从 BACnet 设备读取多个属性。
func (a *GoBACnetAdapter) ReadProperties(
	ctx context.Context,
	nodes []NodeMapping,
) ([]Sample, error) {
	if !a.IsConnected() {
		return nil, ErrDisconnected // 如果未连接，则返回错误
	}

	samples := make([]Sample, 0, len(nodes))

	for _, node := range nodes {
		sample, err := a.readProperty(ctx, node) // 调用辅助方法读取单个属性
		if err != nil {
			// 如果读取失败，记录错误并添加一个质量为 Bad 的样本
			samples = append(samples, Sample{
				ID:         node.ID,
				ObjectType: node.ObjectType,
				Instance:   node.Instance,
				PropertyID: node.PropertyID,
				Quality:    QualityBad,
			})
			continue
		}
		samples = append(samples, sample)
	}

	return samples, nil
}

// WriteProperty 向 BACnet 设备写入单个属性。
func (a *GoBACnetAdapter) WriteProperty(
	ctx context.Context,
	node NodeMapping,
	value any,
) error {
	if !a.IsConnected() {
		return ErrDisconnected // 如果未连接，则返回错误
	}

	if !node.Writable {
		return ErrPropertyNotWritable // 如果属性不可写，则返回错误
	}

	return a.writeProperty(ctx, node, value) // 调用辅助方法写入单个属性
}

// readProperty 是一个辅助方法，用于从 BACnet 设备读取单个属性。
func (a *GoBACnetAdapter) readProperty(ctx context.Context, node NodeMapping) (Sample, error) {
	// 构建 ReadPropertyData 请求
	req := gb.ReadPropertyData{
		Object: gb.Object{
			ID: gb.ObjectID{
				Type:     gb.ObjectType(node.ObjectType),
				Instance: gb.ObjectInstance(node.Instance),
			},
			Properties: []gb.Property{
				{
					Type:       uint32(node.PropertyID),
					ArrayIndex: gb.ArrayAll,
				},
			},
		},
	}

	// 构建目标设备
	destDevice := gb.Device{Addr: a.target} // 假设 gb.Device 有一个 Address 字段

	value, err := a.client.ReadProperty(destDevice, req) // 移除 context.Context 参数
	if err != nil {
		return Sample{}, fmt.Errorf("read property %s:%d:%d failed: %w", node.ObjectType, node.Instance, node.PropertyID, err)
	}

	// 解码读取到的值
	v, err := decodeValue(value, node.DataType)
	if err != nil {
		return Sample{}, fmt.Errorf("decode value for %s:%d:%d failed: %w", node.ObjectType, node.Instance, node.PropertyID, err)
	}

	return Sample{
		ID:         node.ID,
		ObjectType: node.ObjectType,
		Instance:   node.Instance,
		PropertyID: node.PropertyID,
		Value:      v,
		Timestamp:  time.Now(),
		Quality:    QualityGood,
	}, nil
}

// writeProperty 是一个辅助方法，用于向 BACnet 设备写入单个属性。
func (a *GoBACnetAdapter) writeProperty(ctx context.Context, node NodeMapping, value any) error {
	// 当前使用的 gobacnet 库版本 (alexbeltran/gobacnet) 不支持 WriteProperty 方法。
	// 如果需要写入功能，可能需要：
	// 1. 寻找并切换到支持写入的 gobacnet 分叉或版本。
	// 2. 手动实现 BACnet 写入 APDU 的构建和发送。
	// 3. 暂时将此功能标记为未实现。
	//
	// 为了让代码通过编译，我们暂时返回 ErrNotImplemented。
	// 请与用户讨论如何处理 BACnet 写入功能。
	return ErrNotImplemented
}