package bacnet

import (
	"context"
	"fmt"
	"log"

	"github.com/OpenGongChang/OpenIndustrial/gateway/model/point"
	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/driver"
)

const DriverType = "BACnet"

type Driver struct {
	*driver.BaseDriver

	config Config

	ctx driver.Context // Use the new driver context

	adapter Adapter

	poller *Poller

	cancel context.CancelFunc
}

func NewDriver(config Config) (*Driver, error) {

	if err := config.Validate(); err != nil {
		return nil, err
	}

	d := &Driver{
		BaseDriver: driver.NewBaseDriver(config.Name, DriverType),
		config:     config,
	}

	return d, nil
}

func (d *Driver) Init(ctx driver.Context) error {

	d.ctx = ctx // Store the new context

	switch d.config.Connection.Mode {

	case ConnectionModeIP:
		d.adapter = NewGoBACnetAdapter() // 修复点：将 NewIPAdapter() 替换为 NewGoBACnetAdapter()

	case ConnectionModeMSTP:
		return ErrUnsupportedMode

	default:
		return ErrUnsupportedMode
	}

	if err := d.adapter.Connect(context.Background(), d.config.Connection); err != nil {
		return err
	}

	poller, err := NewPoller(d.adapter, d.config)
	if err != nil {
		return err
	}

	d.poller = poller

	return d.BaseDriver.Init(ctx)
}

func (d *Driver) Start(ctx context.Context) error {

	if d.poller == nil {
		return ErrNotInitialized
	}

	driverCtx, cancel := context.WithCancel(ctx)

	d.cancel = cancel

	d.poller.Start(driverCtx)

	go d.publishLoop(driverCtx)

	log.Printf("[%s] started", d.Name())

	return nil
}

func (d *Driver) Stop(ctx context.Context) error {

	if d.cancel != nil {
		d.cancel()
	}

	if d.poller != nil {
		d.poller.Stop()
	}

	if d.adapter != nil {
		_ = d.adapter.Disconnect(ctx)
	}

	log.Printf("[%s] stopped", d.Name())

	return nil
}

func (d *Driver) publishLoop(ctx context.Context) {

	for {

		select {

		case <-ctx.Done():
			return

		case sample, ok := <-d.poller.Samples():

			if !ok {
				return
			}

			// Convert to standard driver.Event and submit
			event := driver.NewDataEvent(
				d.Name(),
				"", // BACnet doesn't have a simple DeviceID concept, using driver name or empty for now.
				sample.ID,
				sample.Value,
				toPointQuality(sample.Quality),
				sample.Timestamp,
			)
			d.ctx.Submit(event)
		}
	}
}

// toPointQuality converts bacnet.Quality (string) to point.Quality (uint8).
func toPointQuality(q Quality) point.Quality {
	switch q {
	case QualityGood:
		return point.QualityGood
	case QualityBad:
		return point.QualityBad
	case QualityUncertain, QualityDisconnected, QualityNotSupported:
		return point.QualityUncertain
	default:
		return point.QualityUncertain
	}
}


func (d *Driver) Read(pointID string) (any, error) {

	if d.poller == nil {
		return nil, ErrNotInitialized
	}

	return nil, fmt.Errorf("Read() not implemented, use Poller cache")
}

func (d *Driver) Write(pointID string, value any) error {

	for _, node := range d.config.NodeMappings {

		if node.ID != pointID {
			continue
		}

		return d.adapter.WriteProperty(
			context.Background(),
			node,
			value,
		)
	}

	return ErrObjectNotFound
}