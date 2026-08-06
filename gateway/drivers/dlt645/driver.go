package dlt645

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Driver is the main entry point for the DL/T 645 module.
type Driver struct {
	config     *ConnectionConfig
	adapter    Adapter
	poller     *Poller
	cache      *Cache
	sampleChan chan Sample
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewDriver creates a new DL/T 645 driver with the default serial adapter.
func NewDriver(cfg *ConnectionConfig, meters []Meter, points []PointMapping) (*Driver, error) {
	adapter := NewSerialAdapter()
	return NewDriverWithAdapter(cfg, meters, points, adapter)
}

// NewDriverWithAdapter creates a new DL/T 645 driver with a custom adapter.
// This is useful for testing or for implementing other transport layers.
func NewDriverWithAdapter(cfg *ConnectionConfig, meters []Meter, points []PointMapping, adapter Adapter) (*Driver, error) {
	cache := NewCache()
	sampleChan := make(chan Sample, 100)

	// The poller will read data and send it to the sampleChan.
	poller := NewPoller(adapter, meters, points, cfg.Timeout, sampleChan)

	return &Driver{
		config:     cfg,
		adapter:    adapter,
		poller:     poller,
		cache:      cache,
		sampleChan: sampleChan,
	}, nil
}

// Start connects the adapter and starts the polling and processing loops.
func (d *Driver) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	log.Println("Starting DL/T 645 driver...")

	if err := d.adapter.Connect(ctx, *d.config); err != nil {
		log.Printf("Failed to connect DL/T 645 adapter: %v", err)
		cancel()
		return err
	}
	log.Println("DL/T 645 adapter connected successfully.")

	// Start the background loop to process samples from the poller
	d.wg.Add(1)
	go d.processSamples(ctx)

	// Start the poller
	d.poller.Start()

	log.Println("DL/T 645 driver started successfully.")
	return nil
}

// Stop stops the poller and disconnects the adapter gracefully.
func (d *Driver) Stop() error {
	log.Println("Stopping DL/T 645 driver...")
	if d.cancel != nil {
		d.cancel()
	}

	// Stop the poller first
	if d.poller != nil {
		d.poller.Stop()
	}

	// Wait for the sample processing loop to finish
	d.wg.Wait()

	// Disconnect the adapter
	if d.adapter != nil && d.adapter.IsConnected() {
		if err := d.adapter.Disconnect(context.Background()); err != nil {
			return err
		}
	}
	log.Println("DL/T 645 driver stopped.")
	return nil
}

// processSamples runs a loop to read Samples from the channel and update the cache.
func (d *Driver) processSamples(ctx context.Context) {
	defer d.wg.Done()
	for {
		select {
		case sample := <-d.sampleChan:
			d.cache.Set(sample.PointID, sample)
			log.Printf("Updated cache for %s: %v", sample.PointID, sample.Value)
		case <-ctx.Done():
			log.Println("Sample processing loop stopped.")
			return
		}
	}
}

// Read returns the latest value for a specific point ID from the cache.
func (d *Driver) Read(pointID string) (Sample, error) {
	sample, found := d.cache.Get(pointID)
	if !found {
		return Sample{}, fmt.Errorf("point ID '%s' not found in cache", pointID)
	}
	return sample, nil
}

// Write sends a command to a specific point.
func (d *Driver) Write(pointID string, value interface{}) error {
	// This requires mapping the high-level pointID back to a meter address and PointMapping.
	// A more complex mapping structure would be needed for a robust implementation.
	return fmt.Errorf("write not yet implemented in driver layer")
}