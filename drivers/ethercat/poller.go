package ethercat

import (
	"context"
	"log"
	"time"
)

// Poller is responsible for the high-frequency cyclic task of exchanging PDO data.
type Poller struct {
	adapter  Adapter
	pdoMap   []PDOMapping
	cycle    time.Duration
	cache    *Cache
	output   chan<- Sample
	decoder  *pdoDecoder // The decoder will be responsible for parsing the raw PDO bytes.
}

// NewPoller creates a new Poller instance.
func NewPoller(adapter Adapter, pdoMap []PDOMapping, cycle time.Duration, cache *Cache, output chan<- Sample) *Poller {
	return &Poller{
		adapter:  adapter,
		pdoMap:   pdoMap,
		cycle:    cycle,
		cache:    cache,
		output:   output,
		decoder:  newPdoDecoder(pdoMap),
	}
}

// Start begins the real-time polling loop.
// This loop must run with high priority and minimal jitter.
func (p *Poller) Start(ctx context.Context) {
	log.Printf("Starting EtherCAT poller with a cycle time of %s...", p.cycle)
	ticker := time.NewTicker(p.cycle)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.executeCycle(ctx)
		case <-ctx.Done():
			log.Println("Stopping EtherCAT poller.")
			return
		}
	}
}

// executeCycle performs a single, time-critical PDO exchange cycle.
func (p *Poller) executeCycle(ctx context.Context) {
	// 1. Write output PDOs to slaves
	// TODO: Get output data from a command queue or cache
	outputData := make(map[uint16][]byte) // Placeholder
	if err := p.adapter.WritePDOs(ctx, outputData); err != nil {
		log.Printf("Error writing PDOs: %v", err)
		// In a real-world scenario, this might trigger a state change to error.
		return
	}

	// 2. Read input PDOs from slaves
	rawData, err := p.adapter.ReadPDOs(ctx)
	if err != nil {
		log.Printf("Error reading PDOs: %v", err)
		return
	}

	// 3. Decode the raw PDO data into Samples
	if len(rawData) > 0 {
		samples := p.decoder.decode(rawData)
		timestamp := time.Now()

		// 4. Update cache and send samples to the output channel
		for _, sample := range samples {
			sample.Timestamp = timestamp
			sample.Quality = QualityGood // Assume good quality if we got this far
			p.cache.Set(sample)
			p.output <- sample
		}
	}
}