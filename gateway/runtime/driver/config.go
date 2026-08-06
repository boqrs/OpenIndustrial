package driver

import "time"

// Config is the configuration for a driver instance.
type Config struct {
	ID   string `mapstructure:"id"`
	Type string `mapstructure:"type"`
	// Interval for data collection
	Interval time.Duration `mapstructure:"interval"`
	// Driver-specific configuration
	Settings map[string]interface{} `mapstructure:"settings"`
}