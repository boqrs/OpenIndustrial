package opcua

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Poller struct {
	adapter   Adapter
	optimizer *Optimizer
	config    PollConfig
	mappings  []NodeMapping
	sampleCh  chan<- Sample
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

type PollerBuilder struct {
	adapter   Adapter
	optimizer *Optimizer
	config    PollConfig
	mappings  []NodeMapping
	sampleCh  chan<- Sample
}

func NewPollerBuilder() *PollerBuilder {
	return &PollerBuilder{}
}

func (b *PollerBuilder) WithAdapter(adapter Adapter) *PollerBuilder {
	b.adapter = adapter
	return b
}

func (b *PollerBuilder) WithOptimizer(optimizer *Optimizer) *PollerBuilder {
	b.optimizer = optimizer
	return b
}

func (b *PollerBuilder) WithConfig(config PollConfig) *PollerBuilder {
	b.config = config
	return b
}

func (b *PollerBuilder) WithMappings(mappings []NodeMapping) *PollerBuilder {
	b.mappings = mappings
	return b
}

func (b *PollerBuilder) WithSampleChan(sampleCh chan<- Sample) *PollerBuilder {
	b.sampleCh = sampleCh
	return b
}

func (b *PollerBuilder) Build() (*Poller, error) {
	if b.adapter == nil {
		return nil, errors.New("poller builder: adapter is required")
	}
	if b.sampleCh == nil {
		return nil, errors.New("poller builder: sample channel is required")
	}
	if b.config.Interval <= 0 {
		return nil, errors.New("poller builder: poll interval must be positive")
	}
	if b.optimizer == nil {
		log.Println("Poller builder: optimizer not provided, creating a new one.")
		b.optimizer = NewOptimizer()
	}

	return &Poller{
		adapter:   b.adapter,
		optimizer: b.optimizer,
		config:    b.config,
		mappings:  b.mappings,
		sampleCh:  b.sampleCh,
		stopCh:    make(chan struct{}),
	}, nil
}

func (p *Poller) Start(ctx context.Context) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ticker := time.NewTicker(p.config.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.poll(ctx)
			case <-p.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (p *Poller) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

func (p *Poller) poll(ctx context.Context) {
	if !p.adapter.IsConnected() {
		log.Println("Poll skipped: adapter is not connected")
		return
	}

	samples, err := p.adapter.ReadNodes(ctx, p.mappings)
	if err != nil {
		log.Printf("Error during poll: %v", err)
		return
	}

	for _, sample := range samples {
		if optimizedSample, ok := p.optimizer.LastValueFilter(sample); ok {
			select {
			case p.sampleCh <- optimizedSample:
			case <-ctx.Done():
				return
			default:
				log.Printf("Sample channel is full, dropping sample for node %s", optimizedSample.NodeID)
			}
		}
	}
}