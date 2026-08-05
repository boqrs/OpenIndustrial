package ethernetip

import (
	"context"
	"log"
)

const (
	// internalChannelSize is the size of the channels connecting poller, optimizer, and driver.
	internalChannelSize = 100
)

// Driver is the main entry point for the EtherNet/IP module.
// It encapsulates and manages the lifecycle of all components.
type Driver struct {
	config     *ConnectionConfig
	points     []PointMapping
	adapter    Adapter
	poller     *Poller
	optimizer  *Optimizer
	sampleChan chan Sample // The final output channel for the runtime
	cancel     context.CancelFunc
}

// NewDriver creates a new EtherNet/IP driver.
func NewDriver(cfg *ConnectionConfig, points []PointMapping) *Driver {
	// Based on the config, we select the appropriate adapter.
	// For now, we only have the CIP (Explicit) adapter.
	var adapter Adapter
	if cfg.Mode == ModeExplicit {
		adapter = NewCIPAdapter()
	} else {
		// In the future, we would instantiate NewImplicitAdapter() here.
		// For now, we'll just default to the explicit one.
		log.Printf("Warning: Mode '%s' not fully implemented. Defaulting to explicit.", cfg.Mode)
		adapter = NewCIPAdapter()
	}

	// Create the channel pipeline: Poller -> Optimizer -> Driver Output
	pollerToOptimizerChan := make(chan Sample, internalChannelSize)
	finalOutputChan := make(chan Sample, internalChannelSize)

	poller := NewPoller(adapter, points, cfg.Timeout, pollerToOptimizerChan) // Using Timeout as poll interval for now
	optimizer := NewOptimizer(pollerToOptimizerChan, finalOutputChan)

	return &Driver{
		config:     cfg,
		points:     points,
		adapter:    adapter,
		poller:     poller,
		optimizer:  optimizer,
		sampleChan: finalOutputChan,
	}
}

// Start connects the adapter and starts the polling and optimization goroutines.
func (d *Driver) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel

	log.Println("Starting EtherNet/IP driver...")
	if err := d.adapter.Connect(ctx, *d.config); err != nil {
		log.Printf("Failed to connect EtherNet/IP adapter: %v", err)
		cancel()
		return err
	}

	// Start the pipeline in reverse order of data flow.
	go d.optimizer.Start(ctx)
	go d.poller.Start(ctx)

	log.Println("EtherNet/IP driver started successfully.")
	return nil
}

// Stop stops all goroutines and disconnects the adapter.
func (d *Driver) Stop() error {
	log.Println("Stopping EtherNet/IP driver...")
	if d.cancel != nil {
		d.cancel()
	}

	if d.adapter != nil && d.adapter.IsConnected() {
		if err := d.adapter.Disconnect(context.Background()); err != nil {
			return err
		}
	}
	log.Println("EtherNet/IP driver stopped.")
	return nil
}

// Samples returns the final output channel for the runtime to consume.
func (d *Driver) Samples() <-chan Sample {
	return d.sampleChan
}