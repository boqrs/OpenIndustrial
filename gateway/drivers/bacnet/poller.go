// drivers/bacnet/poller.go
package bacnet

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Poller 负责定时从 BACnet 设备读取属性值。
// 它将获取到的数据通过优化器处理后，发布到样本通道。
type Poller struct {
	adapter   Adapter
	config    Config // 使用整个 Config，因为它包含 NodeMappings 和 PollConfig
	optimizer *Optimizer
	stopCh    chan struct{}
	wg        sync.WaitGroup
	sampleCh  chan Sample // 用于发布处理后的样本
}

// NewPoller 创建一个新的 Poller 实例。
// 它接收一个已经初始化并连接的 Adapter 实例。
func NewPoller(adapter Adapter, config Config) (*Poller, error) {
	if adapter == nil {
		return nil, fmt.Errorf("BACnet adapter cannot be nil")
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid BACnet poller configuration: %w", err)
	}

	return &Poller{
		adapter:   adapter,
		config:    config,
		optimizer: NewOptimizer(), // 初始化优化器
		stopCh:    make(chan struct{}),
		sampleCh:  make(chan Sample, 100), // 缓冲通道
	}, nil
}

// Start 启动 Poller 的数据采集循环。
func (p *Poller) Start(ctx context.Context) {
	p.wg.Add(1)
	go p.run(ctx)
}

// Stop 停止 Poller 的数据采集循环并等待其完成。
func (p *Poller) Stop() {
	close(p.stopCh)
	p.wg.Wait()
	close(p.sampleCh) // 关闭样本通道
}

// Samples 返回一个只读通道，用于接收处理后的样本。
func (p *Poller) Samples() <-chan Sample {
	return p.sampleCh
}

// run 是 Poller 的主循环，定时进行轮询。
func (p *Poller) run(ctx context.Context) {
	defer p.wg.Done()

	log.Printf("BACnet Poller: Starting polling mode for device %s with interval %s",
		p.config.Connection.DeviceAddress, p.config.Poll.Interval)
	ticker := time.NewTicker(p.config.Poll.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.poll(ctx)
		case <-p.stopCh:
			log.Println("BACnet Poller: Stopping polling due to stop signal.")
			return
		case <-ctx.Done():
			log.Println("BACnet Poller: Stopping polling due to context cancellation.")
			return
		}
	}
}

// poll 执行一次轮询操作，读取 BACnet 属性值。
func (p *Poller) poll(ctx context.Context) {
	if !p.adapter.IsConnected() {
		log.Printf("BACnet Poller: Adapter not connected, skipping poll cycle.")
		return
	}

	// 批量读取所有配置的节点属性
	samples, err := p.adapter.ReadProperties(ctx, p.config.NodeMappings)
	if err != nil {
		log.Printf("BACnet Poller: Error reading properties: %v", err)
		// 如果读取失败，为所有节点发布 Bad 质量的样本
		for _, nm := range p.config.NodeMappings {
			p.publishSample(Sample{
				ID:         nm.ID,
				ObjectType: nm.ObjectType,
				Instance:   nm.Instance,
				PropertyID: nm.PropertyID,
				Value:      nil,
				Timestamp:  time.Now(),
				Quality:    QualityBad,
			})
		}
		return
	}

	// 通过优化器处理原始样本
	optimizedSamples := p.optimizer.ProcessSamples(samples)
	for _, sample := range optimizedSamples {
		p.publishSample(sample)
	}
}

// publishSample 将样本发布到样本通道。
func (p *Poller) publishSample(sample Sample) {
	select {
	case p.sampleCh <- sample:
		// 样本已发送
	default:
		log.Printf("BACnet Poller: Sample channel is full, dropping sample for ID: %s", sample.ID)
	}
}