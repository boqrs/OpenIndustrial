package can

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// DriverState represents the operational state of the driver.
type DriverState int

const (
	DriverStateStopped DriverState = iota
	DriverStateStarting
	DriverStateRunning
	DriverStateStopping
)

// Driver is the main entry point and orchestrator for the CAN driver.
type Driver struct {
	name      string
	config    Config
	adapter   Adapter
	decoder   Decoder
	encoder   Encoder
	cache     Cache
	publisher Publisher
	optimizer *Optimizer
	poller    *Poller
	state     DriverState
	signals   map[uint32][]Signal // Map from FrameID to list of signals
	mu        sync.Mutex
	wg        sync.WaitGroup
	cancel    context.CancelFunc
}

// NewDriver creates a new CAN driver instance.
func NewDriver(config Config, publisher Publisher) (*Driver, error) {
	signalMap := make(map[uint32][]Signal)
	for _, s := range config.Signals {
		signalMap[s.FrameID] = append(signalMap[s.FrameID], s)
	}

	adapter := NewSocketCANAdapter()
	encoder := NewDefaultEncoder()

	d := &Driver{
		name:      config.Name,
		config:    config,
		adapter:   adapter,
		decoder:   NewDefaultDecoder(),
		encoder:   encoder,
		cache:     NewMemoryCache(),
		publisher: publisher,
		optimizer: NewOptimizer(),
		poller:    NewPoller(adapter, encoder, config.Signals, config.TxMessages),
		state:     DriverStateStopped,
		signals:   signalMap,
	}
	return d, nil
}

// Start connects the adapter and begins processing CAN frames.
func (d *Driver) Start(ctx context.Context) error {
	d.mu.Lock()
	if d.state != DriverStateStopped {
		d.mu.Unlock()
		return fmt.Errorf("driver is not in a stopped state")
	}
	d.state = DriverStateStarting
	d.mu.Unlock()

	ctx, d.cancel = context.WithCancel(ctx)

	if err := d.adapter.Connect(ctx, d.config.Connection); err != nil {
		d.mu.Lock()
		d.state = DriverStateStopped
		d.mu.Unlock()
		return fmt.Errorf("failed to connect adapter: %w", err)
	}

	d.wg.Add(1)
	go d.consume(ctx)

	if d.poller != nil {
		d.poller.Start(ctx)
	}

	d.mu.Lock()
	d.state = DriverStateRunning
	d.mu.Unlock()
	log.Printf("CAN driver '%s' started on interface '%s'", d.name, d.config.Connection.Interface)
	return nil
}

// Stop disconnects the adapter and stops all processing.
func (d *Driver) Stop(ctx context.Context) error {
	d.mu.Lock()
	if d.state != DriverStateRunning {
		d.mu.Unlock()
		return fmt.Errorf("driver is not in a running state")
	}
	d.state = DriverStateStopping
	d.mu.Unlock()

	if d.poller != nil {
		d.poller.Stop()
	}

	if d.cancel != nil {
		d.cancel()
	}

	d.wg.Wait()

	if err := d.adapter.Disconnect(ctx); err != nil {
		log.Printf("failed to disconnect adapter gracefully: %v", err)
	}

	d.mu.Lock()
	d.state = DriverStateStopped
	d.mu.Unlock()
	log.Printf("CAN driver '%s' stopped", d.name)
	return nil
}

// Read retrieves the latest value of a signal from the cache.
func (d *Driver) Read(id string) (SignalValue, bool) {
	return d.cache.Get(id)
}

// Write encodes a signal value and sends it as a CAN frame.
func (d *Driver) Write(ctx context.Context, id string, value any) error {
	var targetSignal Signal
	var found bool
	for _, s := range d.config.Signals {
		if s.ID == id {
			targetSignal = s
			found = true
			break
		}
	}

	if !found {
		return ErrSignalNotFound
	}

	frame, err := d.encoder.Encode(targetSignal, value)
	if err != nil {
		return fmt.Errorf("failed to encode signal %s: %w", id, err)
	}

	return d.adapter.Send(ctx, frame)
}

// consume is the main loop for processing incoming CAN frames.
func (d *Driver) consume(ctx context.Context) {
	defer d.wg.Done()
	frameChan := d.adapter.Receive()

	for {
		select {
		case frame, ok := <-frameChan:
			if !ok {
				log.Println("frame channel closed, stopping consumer")
				return
			}
			d.processFrame(frame)
		case <-ctx.Done():
			log.Println("context cancelled, stopping consumer")
			return
		}
	}
}

// processFrame decodes all relevant signals from a single CAN frame.
func (d *Driver) processFrame(frame Frame) {
	signals, ok := d.signals[frame.ID]
	if !ok {
		return // No signals defined for this frame ID
	}

	for _, signal := range signals {
		value, err := d.decoder.Decode(frame, signal)
		if err != nil {
			log.Printf("failed to decode signal %s: %v", signal.ID, err)
			continue
		}

		if d.optimizer.Optimize(value) {
			d.cache.Set(value)
			if d.publisher != nil {
				if err := d.publisher.Publish(value); err != nil {
					log.Printf("failed to publish signal %s: %v", signal.ID, err)
				}
			}
		}
	}
}