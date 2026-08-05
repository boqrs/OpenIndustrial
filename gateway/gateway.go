package gateway

import (
	"context"
	"log"
	"sync"

	"github.com/OpenGongChang/OpenIndustrial/runtime/driver"
)

// Gateway is the core of the industrial edge gateway.
type Gateway struct {
	config    *GatewayConfig
	registry  *Registry
	collector *Collector
	cache     *Cache
	// add other components like publisher etc.

	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
}

// NewGateway creates a new gateway instance.
func NewGateway(config *GatewayConfig) *Gateway {
	ctx, cancel := context.WithCancel(context.Background())
	cache := NewCache()
	return &Gateway{
		config:    config,
		registry:  NewRegistry(),
		collector: NewCollector(cache),
		cache:     cache,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Registry returns the driver registry.
func (g *Gateway) Registry() *Registry {
	return g.registry
}

// Collector returns the gateway's data collector.
func (g *Gateway) Collector() *Collector {
	return g.collector
}

// Start initializes and starts the gateway and all its drivers.
func (g *Gateway) Start() error {
	g.mu.Lock()
	if g.running {
		g.mu.Unlock()
		return nil // already running
	}
	g.running = true
	g.mu.Unlock()

	log.Println("Starting gateway...")

	// Start collector
	go g.collector.run(g.ctx)

	// Create and start drivers based on config
	for _, driverConfig := range g.config.Drivers {
		d, err := g.registry.Create(driverConfig)
		if err != nil {
			log.Printf("Error creating driver %s: %v", driverConfig.ID, err)
			continue
		}

		// The context for the driver should be derived from the gateway's context
		// but we need a runtime context object. Let's create a simple one for now.
		driverCtx := g.createDriverContext(d)
		if err := d.Init(driverCtx); err != nil {
			log.Printf("Error initializing driver %s: %v", d.Name(), err)
			continue
		}

		if err := d.Start(g.ctx); err != nil {
			log.Printf("Error starting driver %s: %v", d.Name(), err)
			continue
		}
		log.Printf("Driver %s started.", d.Name())
	}

	log.Println("Gateway started successfully.")
	return nil
}

// Stop gracefully shuts down the gateway.
func (g *Gateway) Stop() error {
	g.mu.Lock()
	if !g.running {
		g.mu.Unlock()
		return nil // not running
	}
	g.running = false
	g.mu.Unlock()

	log.Println("Stopping gateway...")
	g.cancel() // Signal all components to stop

	// Stop all drivers
	for _, d := range g.registry.List() {
		if err := d.Stop(context.Background()); err != nil { // Use a new context for stopping
			log.Printf("Error stopping driver %s: %v", d.Name(), err)
		}
	}

	log.Println("Gateway stopped.")
	return nil
}

// This is a placeholder. We need to define a proper runtime context.
// For now, it connects the driver's output to the collector.
func (g *Gateway) createDriverContext(d driver.Driver) driver.Context {
	// This is the crucial link!
	// We need a way for the driver to submit events.
	// The context is the perfect place for that.
	return &driverContext{
		collector: g.collector,
	}
}

// driverContext implements driver.Context
type driverContext struct {
	collector *Collector
}

func (c *driverContext) Submit(event driver.Event) {
	c.collector.Submit(event)
}