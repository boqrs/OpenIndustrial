package profinet

import (
	"context"
	"log"
	"time"
)

// Poller is responsible for periodically collecting data from the adapter
// and updating the cache.
type Poller struct {
	adapter  Adapter
	cache    *Cache
	interval time.Duration
	devices  []DeviceConfig // List of devices to poll
}

// NewPoller creates a new Poller instance.
func NewPoller(adapter Adapter, cache *Cache, devices []DeviceConfig, interval time.Duration) *Poller {
	return &Poller{
		adapter:  adapter,
		cache:    cache,
		devices:  devices,
		interval: interval,
	}
}

// Start begins the polling loop. It runs until the context is cancelled.
func (p *Poller) Start(ctx context.Context) {
	log.Println("Starting PROFINET poller...")
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.poll()
		case <-ctx.Done():
			log.Println("Stopping PROFINET poller.")
			return
		}
	}
}

// poll performs a single polling cycle.
func (p *Poller) poll() {
	samples, err := p.collect()
	if err != nil {
		log.Printf("Error collecting PROFINET samples: %v", err)
		return
	}

	for _, s := range samples {
		p.cache.Update(s)
	}
	log.Printf("Polled and updated %d samples.", len(samples))
}

// collect fetches data from all configured devices and converts it into Sample format.
// NOTE: This is a placeholder implementation. The actual logic will depend on GSDML parsing
// and the data layout of each device.
func (p *Poller) collect() ([]Sample, error) {
	var allSamples []Sample

	for _, device := range p.devices {
		// 1. Read cyclic input data for the device.
		inputData, err := p.adapter.ReadInputs(context.Background(), device.StationName)
		if err != nil {
			log.Printf("Failed to read inputs for device %s: %v", device.StationName, err)
			// Create a "bad quality" sample for the whole device to indicate a problem.
			allSamples = append(allSamples, Sample{
				PointID:   device.StationName,
				Value:     nil,
				Timestamp: time.Now(),
				Quality:   QualityBad,
				Source:    "profinet",
			})
			continue
		}

		// 2. Parse the inputData based on the device's GSDML and ModuleConfig.
		// This is where the raw bytes are transformed into meaningful values.
		// For now, we'll just create a single sample representing the raw byte array.
		// TODO: Implement proper GSDML-based parsing.
		allSamples = append(allSamples, Sample{
			PointID:   device.StationName + ".Inputs",
			Value:     inputData,
			Timestamp: time.Now(),
			Quality:   QualityGood,
			Source:    "profinet",
		})
	}

	return allSamples, nil
}