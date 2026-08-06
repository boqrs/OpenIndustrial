package hj212

import (
	"context"
	"log"
	"sync"
	"time"
)

// Poller is responsible for continuously listening for data from the adapter.
// In the context of HJ/T 212, it acts as a listener for proactively sent data.
type Poller struct {
	adapter    Adapter
	mappings   []PointMapping
	sampleChan chan<- Sample
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewPoller creates a new Poller.
func NewPoller(adapter Adapter, mappings []PointMapping, sampleChan chan<- Sample) *Poller {
	return &Poller{
		adapter:    adapter,
		mappings:   mappings,
		sampleChan: sampleChan,
		stopChan:   make(chan struct{}),
	}
}

// Start begins the polling (listening) process in a background goroutine.
func (p *Poller) Start() {
	p.wg.Add(1)
	go p.listenLoop()
}

// Stop signals the listening loop to terminate and waits for it to finish.
func (p *Poller) Stop() {
	close(p.stopChan)
	p.wg.Wait()
}

// listenLoop is the main loop that continuously reads data from the adapter.
func (p *Poller) listenLoop() {
	defer p.wg.Done()
	log.Println("Poller: Starting listening loop...")

	for {
		select {
		case <-p.stopChan:
			log.Println("Poller: Stopping listening loop.")
			return
		default:
			// Blocking call to read the next data segment.
			segment, err := p.adapter.ReadDataSegment(context.Background())
			if err != nil {
				// In a real application, handle different errors differently.
				// e.g., reconnect on connection errors.
				log.Printf("Poller: Error reading data segment: %v", err)
				continue // Or implement a backoff strategy
			}

			log.Printf("Poller: Received data segment from MN %s", segment.MN)
			p.processDataSegment(segment)
		}
	}
}

// processDataSegment converts a DataSegment into multiple Samples and sends them to the channel.
func (p *Poller) processDataSegment(segment *DataSegment) {
	// Create a map for quick lookup of pollutant codes to their values.
	// e.g., "w01018-Rtd" -> "56.3"
	// We need to handle different suffixes like -Rtd, -Avg, -Flag etc.
	// This is a simplified example that assumes a "-Rtd" suffix or no suffix.
	for _, point := range p.mappings {
		// Try to find the value for the point's code with common suffixes.
		var rawValue string
		var found bool

		// Check for real-time data, then other common types
		suffixes := []string{"-Rtd", "-Avg", "-Min", "-Max", ""}
		for _, suffix := range suffixes {
			val, ok := segment.Pollutants[point.Code+suffix]
			if ok {
				rawValue = val
				found = true
				break
			}
		}

		if !found {
			continue
		}

		// TODO: Convert rawValue (string) to the appropriate data type based on mapping.
		// For now, we'll just pass it as a string.
		sample := Sample{
			PointID:   point.ID,
			Value:     rawValue,
			Timestamp: time.Now(), // TODO: Use QN field from segment for a more accurate timestamp
			Quality:   "good",
			Source:    "hj212",
		}
		p.sampleChan <- sample
	}
}