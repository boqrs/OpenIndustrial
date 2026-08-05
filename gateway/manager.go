package gateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/OpenGongChang/OpenIndustrial/runtime/driver"
)

type DriverManager struct {
	mu      sync.RWMutex
	drivers map[string]driver.Driver
}

func NewDriverManager() *DriverManager {
	return &DriverManager{
		drivers: make(
			map[string]driver.Driver,
		),
	}
}

func (m *DriverManager) Register(
	d driver.Driver,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := d.Name()

	if _, ok := m.drivers[id]; ok {
		return fmt.Errorf(
			"driver %s already exists",
			id,
		)
	}

	m.drivers[id] = d
	return nil
}

func (m *DriverManager) Get(
	name string,
) (driver.Driver, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	d, ok := m.drivers[name]
	return d, ok
}

func (m *DriverManager) StartAll(
	ctx context.Context,
) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, d := range m.drivers {
		if err := d.Start(ctx); err != nil {
			return fmt.Errorf(
				"start driver %s failed:%w",
				d.Name(),
				err,
			)
		}
	}
	return nil
}

func (m *DriverManager) StopAll(
	ctx context.Context,
) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, d := range m.drivers {
		d.Stop(ctx)
	}
	return nil
}