package gateway

import (
	"fmt"
	"sync"

	runtimeDriver "github.com/OpenGongChang/OpenIndustrial/gateway/runtime/driver"
)

// Factory creates driver instance.
type Factory func(
	config runtimeDriver.Config,
) (
	runtimeDriver.Driver,
	error,
)



// Registry manages all drivers.
//
// Responsibilities:
// - register driver factories
// - create drivers
// - manage lifecycle
type Registry struct {


	factories map[string]Factory


	drivers map[string]runtimeDriver.Driver


	mu sync.RWMutex

}



// NewRegistry creates driver registry.
func NewRegistry() *Registry {


	return &Registry{

		factories:make(
			map[string]Factory,
		),


		drivers:make(
			map[string]runtimeDriver.Driver,
		),

	}

}



// Register registers driver factory.
//
// Example:
//
// registry.Register(
//     "modbus",
//     modbus.NewDriver,
// )
//
func (r *Registry) Register(
	driverType string,
	factory Factory,
) error {


	r.mu.Lock()
	defer r.mu.Unlock()


	if _, exists := r.factories[driverType]; exists {

		return fmt.Errorf(
			"driver type already registered: %s",
			driverType,
		)

	}


	r.factories[driverType]=factory


	return nil
}



// Create creates driver instance.
func (r *Registry) Create(
	cfg runtimeDriver.Config,
) (
	runtimeDriver.Driver,
	error,
){


	r.mu.RLock()

	factory, ok :=
		r.factories[cfg.Type]


	r.mu.RUnlock()



	if !ok {

		return nil,
			fmt.Errorf(
				"driver factory not found: %s",
				cfg.Type,
			)

	}



	driver, err :=
		factory(cfg)


	if err != nil {

		return nil, err

	}



	r.mu.Lock()

	r.drivers[cfg.ID]=driver

	r.mu.Unlock()



	return driver,nil

}



// Get returns driver instance.
func (r *Registry) Get(
	id string,
)(
	runtimeDriver.Driver,
	bool,
){


	r.mu.RLock()
	defer r.mu.RUnlock()


	d,ok :=
		r.drivers[id]


	return d,ok

}



// List returns all drivers.
func (r *Registry) List() []runtimeDriver.Driver {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]runtimeDriver.Driver, 0, len(r.drivers))

	for _, d := range r.drivers {
		result = append(result, d)
	}

	return result
}