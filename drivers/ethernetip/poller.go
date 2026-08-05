package ethernetip

import (
	"context"
	"log"
	"time"
)

// Poller is responsible for periodically collecting data from the adapter.
type Poller struct {
	adapter  Adapter
	points   []PointMapping
	interval time.Duration
	output   chan<- Sample
}

// NewPoller creates a new Poller instance.
func NewPoller(adapter Adapter, points []PointMapping, interval time.Duration, output chan<- Sample) *Poller {
	return &Poller{
		adapter:  adapter,
		points:   points,
		interval: interval,
		output:   output,
	}
}

// Start begins the polling loop. It runs until the context is cancelled.
func (p *Poller) Start(ctx context.Context) {
	log.Println("Starting EtherNet/IP poller...")
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Perform an initial poll immediately.
	p.poll(ctx)

	for {
		select {
		case <-ticker.C:
			p.poll(ctx)
		case <-ctx.Done():
			log.Println("Stopping EtherNet/IP poller.")
			return
		}
	}
}

// poll performs a single polling cycle: reads data from the adapter and sends it to the output channel.
func (p *Poller) poll(ctx context.Context) {
	if !p.adapter.IsConnected() {
		log.Println("Adapter not connected, skipping poll cycle.")
		return
	}

	samples, err := p.adapter.Read(ctx, p.points)
	if err != nil {
		log.Printf("Error reading from EtherNet/IP adapter: %v", err)
		// Even on error, the adapter returns bad quality samples, so we send them.
	}

	if len(samples) > 0 {
		log.Printf("Polled %d samples from EtherNet/IP device.", len(samples))
		for _, s := range samples {
			p.output <- s
		}
	}
}