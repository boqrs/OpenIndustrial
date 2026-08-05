package dlt645

import (
	"context"
	"log"
	"sync"
	"time"
)

// Poller is responsible for periodically reading data from a list of meters.
type Poller struct {
	adapter    Adapter
	meters     []Meter
	points     []PointMapping
	interval   time.Duration
	sampleChan chan<- Sample
	wg         sync.WaitGroup
	cancel     context.CancelFunc
}

// NewPoller creates a new poller instance.
func NewPoller(adapter Adapter, meters []Meter, points []PointMapping, interval time.Duration, sampleChan chan<- Sample) *Poller {
	return &Poller{
		adapter:    adapter,
		meters:     meters,
		points:     points,
		interval:   interval,
		sampleChan: sampleChan,
	}
}

// Start begins the polling loop.
func (p *Poller) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg.Add(1)
	go p.pollLoop(ctx)
}

// Stop gracefully stops the polling loop.
func (p *Poller) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	p.wg.Wait()
}

// pollLoop is the main loop that triggers data reading at each interval.
func (p *Poller) pollLoop(ctx context.Context) {
	defer p.wg.Done()
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Run once immediately at the start
	p.pollAllMeters(ctx)

	for {
		select {
		case <-ticker.C:
			p.pollAllMeters(ctx)
		case <-ctx.Done():
			log.Println("DL/T 645 poller stopped.")
			return
		}
	}
}

// pollAllMeters iterates through all configured meters and reads their data points.
func (p *Poller) pollAllMeters(ctx context.Context) {
	if !p.adapter.IsConnected() {
		log.Println("DL/T 645 adapter is not connected, skipping poll cycle.")
		return
	}

	for _, meter := range p.meters {
		// A context for each meter's read operation
		readCtx, cancel := context.WithTimeout(ctx, p.interval)
		defer cancel()

		log.Printf("Polling meter: %s", meter.Address)
		samples, err := p.adapter.Read(readCtx, meter.Address, p.points)
		if err != nil {
			log.Printf("Failed to read from meter %s: %v", meter.Address, err)
			continue
		}

		for _, sample := range samples {
			p.sampleChan <- sample
		}
	}
}