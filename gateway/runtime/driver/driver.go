package driver

import (
	"context"
	"fmt"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/service"
)

// Context provides the driver with access to the gateway's runtime.
// It's the bridge between a driver and the core system.
type Context interface {
	// Submit sends an event to the gateway's collector.
	Submit(event Event)
}

// Driver defines the interface for a generic OpenIndustrial driver.
// Drivers are essentially services that interact with external systems or implement specific functionalities.
type Driver interface {
	service.Service
	// Init initializes the driver with the given runtime context.
	// This method should be called once after the driver is created and before it's started.
	Init(ctx Context) error
	// Type returns the type of the driver (e.g., "modbus", "opcua", "simulator").
	Type() string
}

// BaseDriver provides a basic implementation for the Driver interface,
// allowing embedding in custom driver structs to reduce boilerplate.
type BaseDriver struct {
	*service.BaseService
	driverType string
	ctx        Context // The runtime context provided during Init
}

// NewBaseDriver creates a new BaseDriver instance.
func NewBaseDriver(name, driverType string) *BaseDriver {
	return &BaseDriver{
		BaseService: service.NewBaseService(name),
		driverType:  driverType,
	}
}

// Init initializes the BaseDriver with the runtime context.
func (bd *BaseDriver) Init(ctx Context) error {
	if bd.ctx != nil {
		return fmt.Errorf("driver '%s' already initialized", bd.Name())
	}
	bd.ctx = ctx
	return nil
}

// Type returns the type of the driver.
func (bd *BaseDriver) Type() string {
	return bd.driverType
}

// Start and Stop methods are inherited from BaseService,
// custom drivers should override them as needed.
func (bd *BaseDriver) Start(ctx context.Context) error {
	// Default: do nothing, just return nil
	return nil
}

func (bd *BaseDriver) Stop(ctx context.Context) error {
	// Default: do nothing, just return nil
	return nil
}