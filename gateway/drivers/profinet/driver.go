package profinet

import (
	"context"
	"log"
)

// Driver is the main entry point for the PROFINET module.
// It encapsulates the adapter, poller, and cache, and manages their lifecycle.
type Driver struct {
	config  *ConnectionConfig
	adapter Adapter
	poller  *Poller
	cache   *Cache
	stop    context.CancelFunc
}

// NewDriver creates a new PROFINET driver.
func NewDriver(cfg *ConnectionConfig, devices []DeviceConfig) *Driver {
	// NewPnetAdapter will return either the real CGO adapter on Linux
	// or the stub adapter on other platforms, thanks to build tags.
	adapter := NewPnetAdapter()

	cache := NewCache()
	// The polling interval might come from config in the future.
	poller := NewPoller(adapter, cache, devices, cfg.CycleTime)

	return &Driver{
		config:  cfg,
		adapter: adapter,
		poller:  poller,
		cache:   cache,
	}
}

// SetAdapter allows injecting the adapter implementation after creation.
// This is useful because the CGO adapter might need complex initialization.
func (d *Driver) SetAdapter(adapter Adapter) {
	d.adapter = adapter
	d.poller.adapter = adapter // Ensure the poller also gets the real adapter.
}

// Start initializes the connection and begins polling.
func (d *Driver) Start() error {
	if d.adapter == nil {
		return ErrCGOInitializationFailed // Or a more specific "adapter not set" error.
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.stop = cancel

	log.Println("Starting PROFINET driver...")
	err := d.adapter.Connect(ctx, *d.config)
	if err != nil {
		log.Printf("Failed to connect PROFINET adapter: %v", err)
		return err
	}

	// Start the poller in a separate goroutine.
	go d.poller.Start(ctx)

	log.Println("PROFINET driver started successfully.")
	return nil
}

// Stop disconnects the adapter and stops the poller.
func (d *Driver) Stop() error {
	log.Println("Stopping PROFINET driver...")
	if d.stop != nil {
		d.stop()
	}

	if d.adapter != nil && d.adapter.IsConnected() {
		err := d.adapter.Disconnect(context.Background())
		if err != nil {
			return err
		}
	}
	log.Println("PROFINET driver stopped.")
	return nil
}

// GetCache returns the driver's data cache.
func (d *Driver) GetCache() *Cache {
	return d.cache
}