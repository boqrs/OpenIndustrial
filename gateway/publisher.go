package gateway

import (
	"context"
)

// Publisher is the interface for components that send data out of the gateway.
// Examples include MQTT publishers, HTTP publishers, etc.
type Publisher interface {
	// Init initializes the publisher with its configuration and a reference to the collector.
	Init(config PublisherConfig, collector *Collector) error

	// Start begins the publishing process.
	Start(ctx context.Context) error

	// Stop gracefully stops the publisher.
	Stop(ctx context.Context) error

	// ID returns the unique identifier of the publisher instance.
	ID() string
}

// PublisherConfig holds the configuration for a single publisher instance.
// This structure is analogous to driver.Config.
type PublisherConfig struct {
	ID       string                 `mapstructure:"id"`
	Type     string                 `mapstructure:"type"`
	Settings map[string]interface{} `mapstructure:"settings"`
}