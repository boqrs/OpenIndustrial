package ethercat

import (
	"context"
	"log"
)

// Driver is the main entry point for the EtherCAT module.
// It encapsulates and manages the lifecycle of all components.
type Driver struct {
	config     *ConnectionConfig
	pdoMap     []PDOMapping
	adapter    Adapter
	poller     *Poller
	cache      *Cache
	sampleChan chan Sample
	cancel     context.CancelFunc
}

// NewDriver creates a new EtherCAT driver.
func NewDriver(cfg *ConnectionConfig, pdoMap []PDOMapping) *Driver {
	// The adapter_soem.go file will provide the concrete implementation for this interface.
	adapter := NewSOEMAdapter()

	// The cache will store the latest values of all points.
	cache := NewCache()

	// The final output channel for the runtime.
	sampleChan := make(chan Sample, 100)

	// The poller is responsible for the high-frequency cyclic task.
	poller := NewPoller(adapter, pdoMap, cfg.CycleTime, cache, sampleChan)

	return &Driver{
		config:     cfg,
		pdoMap:     pdoMap,
		adapter:    adapter,
		poller:     poller,
		cache:      cache,
		sampleChan: sampleChan,
	}
}

// Start connects the adapter, configures slaves, and starts the polling cycle.
func (d *Driver) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	log.Println("Starting EtherCAT driver...")

	// 1. Connect to the EtherCAT bus via the specified interface.
	if err := d.adapter.Connect(ctx, *d.config); err != nil {
		log.Printf("Failed to connect EtherCAT adapter: %v", err)
		cancel()
		return err
	}

	// 2. Scan for all slaves on the bus.
	slaves, err := d.adapter.ScanSlaves(ctx)
	if err != nil {
		log.Printf("Failed to scan for EtherCAT slaves: %v", err)
		d.adapter.Disconnect(ctx)
		cancel()
		return err
	}
	log.Printf("Found %d slaves on the bus.", len(slaves))
	for _, s := range slaves {
		log.Printf("  - Slave %d: %s (Vendor: 0x%X, Product: 0x%X)", s.Index, s.Name, s.VendorID, s.ProductCode)
	}

	// 3. Configure PDO mappings for the slaves.
	if err := d.adapter.ConfigurePDOs(ctx, d.pdoMap); err != nil {
		log.Printf("Failed to configure EtherCAT PDOs: %v", err)
		d.adapter.Disconnect(ctx)
		cancel()
		return err
	}
	log.Println("Successfully configured PDOs.")

	// 4. Start the high-frequency polling cycle.
	go d.poller.Start(ctx)

	log.Println("EtherCAT driver started successfully.")
	return nil
}

// Stop stops the poller and disconnects the adapter.
func (d *Driver) Stop() error {
	log.Println("Stopping EtherCAT driver...")
	if d.cancel != nil {
		d.cancel()
	}

	if d.adapter != nil && d.adapter.IsConnected() {
		if err := d.adapter.Disconnect(context.Background()); err != nil {
			return err
		}
	}
	log.Println("EtherCAT driver stopped.")
	return nil
}

// Samples returns the final output channel for the runtime to consume.
func (d *Driver) Samples() <-chan Sample {
	return d.sampleChan
}