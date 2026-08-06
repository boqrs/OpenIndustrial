package simulator

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/model/point"
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/driver"
	"github.com/go-viper/mapstructure/v2"
)

const DriverType = "simulator"

// Driver is a simple driver that simulates data generation.
type Driver struct {
	*driver.BaseDriver
	ctx    driver.Context
	config Config
	cancel context.CancelFunc
}

// PointConfig defines the configuration for a single simulated point.
type PointConfig struct {
	PointID   string        `mapstructure:"pointId"`
	ValueType string        `mapstructure:"valueType"` // "float", "bool", "int", "sine"
	Interval  time.Duration `mapstructure:"interval"`
	// For float/int
	Min float64 `mapstructure:"min"`
	Max float64 `mapstructure:"max"`
	// For sine wave
	Base      float64 `mapstructure:"base"`
	Amplitude float64 `mapstructure:"amplitude"`
	Period    time.Duration `mapstructure:"period"`
}

// DeviceConfig defines a set of points for a simulated device.
type DeviceConfig struct {
	DeviceID string        `mapstructure:"deviceId"`
	Points   []PointConfig `mapstructure:"points"`
}

// Config holds the configuration for the multi-device simulator driver.
type Config struct {
	Devices []DeviceConfig `mapstructure:"devices"`
}

// NewDriver creates a new simulator driver.
func NewDriver(cfg driver.Config) (driver.Driver, error) {
	var simConfig Config
	// Use a mapstructure decoder for robust config parsing
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &simConfig,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mapstructure decoder: %w", err)
	}

	if err := decoder.Decode(cfg.Settings); err != nil {
		return nil, fmt.Errorf("failed to decode simulator settings: %w", err)
	}

	d := &Driver{
		BaseDriver: driver.NewBaseDriver(cfg.ID, DriverType),
		config:     simConfig,
	}
	return d, nil
}

// Init initializes the driver.
func (d *Driver) Init(ctx driver.Context) error {
	d.ctx = ctx
	return d.BaseDriver.Init(ctx)
}

// Start begins the data simulation for all configured points.
func (d *Driver) Start(ctx context.Context) error {
	driverCtx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	log.Printf("[%s] simulator driver started.", d.Name())

	for _, device := range d.config.Devices {
		for _, point := range device.Points {
			go d.pointSimulationLoop(driverCtx, device, point)
		}
	}

	return nil
}

// Stop halts all data simulation loops.
func (d *Driver) Stop(ctx context.Context) error {
	if d.cancel != nil {
		d.cancel()
	}
	log.Printf("[%s] simulator driver stopped.", d.Name())
	return nil
}

// pointSimulationLoop is the main loop for generating data for a single point.
func (d *Driver) pointSimulationLoop(ctx context.Context, device DeviceConfig, p PointConfig) {
	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	// For sine wave calculation
	startTime := time.Now()

	log.Printf("[%s] starting simulation for point %s on device %s (every %s)", d.Name(), p.PointID, device.DeviceID, p.Interval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var value any
			switch p.ValueType {
			case "float":
				value = p.Min + rand.Float64()*(p.Max-p.Min)
			case "int":
				value = int64(p.Min + rand.Float64()*(p.Max-p.Min))
			case "bool":
				value = rand.Intn(2) == 1
			case "sine":
				elapsed := time.Since(startTime).Seconds()
				periodSeconds := p.Period.Seconds()
				if periodSeconds == 0 {
					periodSeconds = 60 // Avoid division by zero
				}
				radian := 2 * math.Pi * (elapsed / periodSeconds)
				value = p.Base + p.Amplitude*math.Sin(radian)
			default:
				log.Printf("[%s] unknown value type for point %s: %s", d.Name(), p.PointID, p.ValueType)
				continue
			}

			event := driver.NewDataEvent(
				d.Name(),
				device.DeviceID,
				p.PointID,
				value,
				point.QualityGood,
				time.Now(),
			)

			// log.Printf("[%s] submitting event: %+v", d.Name(), event)
			d.ctx.Submit(event)
		}
	}
}