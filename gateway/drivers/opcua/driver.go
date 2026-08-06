package opcua

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type DriverState int

const (
	DriverStateStopped DriverState = iota
	DriverStateStarting
	DriverStateRunning
	DriverStateStopping
)

type DriverConfig struct {
	Connection     ConnectionConfig
	Subscription   SubscriptionConfig
	Poll           PollConfig
	Mappings       []NodeMapping
	CollectionMode string // "poll" or "subscribe"
}

type Driver struct {
	name      string
	config    DriverConfig
	adapter   Adapter
	poller    *Poller
	cache     Cache
	publisher Publisher
	sampleCh  chan Sample
	ctx       context.Context
	cancel    context.CancelFunc
	state     DriverState
	wg        sync.WaitGroup
	mu        sync.Mutex
}

func NewDriver(name string, config DriverConfig, publisher Publisher) (*Driver, error) {
	ctx, cancel := context.WithCancel(context.Background())

	adapter := NewOpcuaClientAdapter()

	if publisher == nil {
		publisher = NewNoOpPublisher()
	}

	driver := &Driver{
		name:      name,
		config:    config,
		adapter:   adapter,
		cache:     NewMemoryCache(),
		publisher: publisher,
		sampleCh:  make(chan Sample, 100),
		ctx:       ctx,
		cancel:    cancel,
		state:     DriverStateStopped,
	}

	if config.CollectionMode == "poll" {
		poller, err := NewPollerBuilder().
			WithAdapter(adapter).
			WithConfig(config.Poll).
			WithMappings(config.Mappings).
			WithSampleChan(driver.sampleCh).
			Build()
		if err != nil {
			return nil, fmt.Errorf("failed to build poller: %w", err)
		}
		driver.poller = poller
	}

	return driver, nil
}

func (d *Driver) Start() error {
	d.mu.Lock()
	if d.state != DriverStateStopped {
		d.mu.Unlock()
		return fmt.Errorf("driver is not in a stopped state")
	}
	d.state = DriverStateStarting
	d.mu.Unlock()

	log.Printf("Starting driver '%s'...", d.name)

	if err := d.adapter.Connect(d.ctx, d.config.Connection); err != nil {
		d.mu.Lock()
		d.state = DriverStateStopped
		d.mu.Unlock()
		return fmt.Errorf("failed to connect adapter: %w", err)
	}

	d.wg.Add(1)
	go d.consume()

	if d.config.CollectionMode == "subscribe" {
		err := d.adapter.Subscribe(d.ctx, d.config.Subscription, d.config.Mappings, d.sampleCh)
		if err != nil {
			d.Stop()
			return fmt.Errorf("failed to start subscription mode: %w", err)
		}
		log.Println("Driver started in subscription mode.")
	} else if d.poller != nil {
		d.poller.Start(d.ctx)
		log.Println("Driver started in polling mode.")
	} else {
		d.Stop()
		return fmt.Errorf("no collection mode configured")
	}

	d.mu.Lock()
	d.state = DriverStateRunning
	d.mu.Unlock()

	log.Printf("Driver '%s' started successfully.", d.name)
	return nil
}

func (d *Driver) Stop() error {
	d.mu.Lock()
	if d.state != DriverStateRunning && d.state != DriverStateStarting {
		d.mu.Unlock()
		return fmt.Errorf("driver is not in a running state")
	}
	d.state = DriverStateStopping
	d.mu.Unlock()

	log.Printf("Stopping driver '%s'...", d.name)

	d.cancel()

	if d.poller != nil {
		d.poller.Stop()
	}

	d.wg.Wait()

	if err := d.adapter.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting adapter: %v", err)
	}

	d.mu.Lock()
	d.state = DriverStateStopped
	d.mu.Unlock()

	log.Printf("Driver '%s' stopped.", d.name)
	return nil
}

func (d *Driver) consume() {
	defer d.wg.Done()
	log.Println("Consumer started.")
	for {
		select {
		case sample, ok := <-d.sampleCh:
			if !ok {
				log.Println("Sample channel closed, consumer stopping.")
				return
			}
			d.cache.Set(sample)
			if err := d.publisher.Publish(sample); err != nil {
				log.Printf("Error publishing sample: %v", err)
			}
		case <-d.ctx.Done():
			log.Println("Context cancelled, consumer stopping.")
			return
		}
	}
}

func (d *Driver) GetValue(id string) (Sample, bool) {
	return d.cache.Get(id)
}

func (d *Driver) WriteValue(id string, value interface{}) error {
	if d.state != DriverStateRunning {
		return fmt.Errorf("driver is not running")
	}
	mapping, ok := d.findMappingByID(id)
	if !ok {
		return fmt.Errorf("no mapping found for id: %s", id)
	}
	sample := Sample{
		ID:     id,
		NodeID: mapping.NodeID,
		Value:  value,
	}
	return d.adapter.WriteNode(d.ctx, sample)
}

func (d *Driver) findMappingByID(id string) (NodeMapping, bool) {
	for _, m := range d.config.Mappings {
		if m.ID == id {
			return m, true
		}
	}
	return NodeMapping{}, false
}