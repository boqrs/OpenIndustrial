package hj212

import (
	"context"
	"fmt"
	"sync"
)

// Driver is the main entry point and controller for the HJ/T 212 driver.
type Driver struct {
	config     *ConnectionConfig
	adapter    Adapter
	poller     *Poller
	cache      *Cache
	sampleChan chan Sample
	stopOnce   sync.Once
}

// NewDriver creates a new HJ/T 212 driver instance.
func NewDriver(cfg *ConnectionConfig, points []PointMapping) (*Driver, error) {
	// Default to the TCP adapter for standard use.
	adapter := NewTCPAdapter()
	return NewDriverWithAdapter(cfg, points, adapter)
}

// NewDriverWithAdapter creates a new driver with a specific adapter, useful for testing.
func NewDriverWithAdapter(cfg *ConnectionConfig, points []PointMapping, adapter Adapter) (*Driver, error) {
	cache := NewCache()
	sampleChan := make(chan Sample, 100)
	poller := NewPoller(adapter, points, sampleChan)

	return &Driver{
		config:     cfg,
		adapter:    adapter,
		poller:     poller,
		cache:      cache,
		sampleChan: sampleChan,
	}, nil
}

// Start connects to the device and begins the listening process.
func (d *Driver) Start() error {
	if err := d.adapter.Connect(context.Background(), *d.config); err != nil {
		return err
	}

	// Start the background goroutine that moves samples from the poller to the cache.
	go d.processSamples()

	// Start the poller's listening loop.
	d.poller.Start()

	return nil
}

// Stop disconnects from the device and stops the listening process.
func (d *Driver) Stop() {
	d.stopOnce.Do(func() {
		d.poller.Stop()
		d.adapter.Disconnect(context.Background())
		close(d.sampleChan)
	})
}

// processSamples is a background task that reads samples from the poller's
// channel and updates the internal cache.
func (d *Driver) processSamples() {
	for sample := range d.sampleChan {
		d.cache.Set(sample.PointID, sample)
	}
}

// Read retrieves the latest value for a given point ID from the cache.
func (d *Driver) Read(pointID string) (Sample, error) {
	sample, found := d.cache.Get(pointID)
	if !found {
		return Sample{}, fmt.Errorf("point ID '%s' not found in cache", pointID)
	}
	return sample, nil
}

// Write sends a command to the device. (Future implementation)
func (d *Driver) Write(pointID string, value interface{}) error {
	// This would involve mapping the pointID to a command (CN) and value,
	// creating a DataSegment, and using adapter.SendCommand.
	return fmt.Errorf("write operation not yet implemented")
}