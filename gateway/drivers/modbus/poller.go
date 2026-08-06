package modbus

import (
	"context"
	"log"
	"time"
)

// Poller is responsible for periodically polling Modbus points.
type Poller struct {
	adapter  Adapter
	config   PollConfig
	mappings []NodeMapping
	stopChan chan struct{}
	doneChan chan struct{}
}

// NewPoller creates a new Poller instance.
func NewPoller(adapter Adapter, config PollConfig, mappings []NodeMapping) *Poller {
	return &Poller{
		adapter:  adapter,
		config:   config,
		mappings: mappings,
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// Start begins the polling loop.
func (p *Poller) Start(ctx context.Context) {
	log.Printf("Modbus Poller started with interval: %s", p.config.Interval)
	ticker := time.NewTicker(p.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Perform polling
			p.poll(ctx)
		case <-p.stopChan:
			log.Println("Modbus Poller stopped.")
			close(p.doneChan)
			return
		case <-ctx.Done():
			log.Println("Modbus Poller stopped due to context cancellation.")
			close(p.doneChan)
			return
		}
	}
}

// Stop gracefully stops the polling loop.
func (p *Poller) Stop() {
	close(p.stopChan)
	<-p.doneChan // Wait for the polling loop to finish
}

// poll performs a single polling cycle.
func (p *Poller) poll(ctx context.Context) {
	if !p.adapter.Connected() {
		log.Println("Poller: Modbus adapter not connected, skipping poll.")
		return
	}

	// Use ReadBatch for efficiency, which will eventually be optimized by the Optimizer
	samples, err := p.adapter.ReadBatch(ctx, p.mappings)
	if err != nil {
		log.Printf("Poller: Failed to read batch: %v", err)
		// Handle error, e.g., update quality of all points to Bad/Uncertain
		for i := range samples {
			samples[i].Quality = Bad
			samples[i].Timestamp = time.Now()
		}
	}

	// Process samples (e.g., send to a data pipeline, log, etc.)
	for _, sample := range samples {
		log.Printf("Poller: Collected Sample - PointID: %s, Value: %v, Quality: %s", sample.PointID, sample.Value, sample.Quality)
		// In a real application, these samples would be sent to the Runtime's data processing pipeline.
	}
}