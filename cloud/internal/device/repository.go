package device

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrDeviceNotFound is returned when a device is not found.
	ErrDeviceNotFound = errors.New("device not found")
)

// Repository defines the persistence interface for device related entities.
type Repository interface {
	Create(ctx context.Context, device *Device) error
	FindByID(ctx context.Context, orgID, deviceID uuid.UUID) (*Device, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Device, error)
	SaveTelemetry(ctx context.Context, deviceID uuid.UUID, ts time.Time, points map[string]interface{}) error
}

// memoryRepository is an in-memory implementation of the Repository interface.
type memoryRepository struct {
	mu        sync.RWMutex
	devices   map[string]*Device
	telemetry map[string][]map[string]interface{}
}

// NewMemoryRepository creates a new in-memory device repository.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		devices:   make(map[string]*Device),
		telemetry: make(map[string][]map[string]interface{}),
	}
}

// Create saves a new device to the repository.
func (r *memoryRepository) Create(ctx context.Context, device *Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices[device.ID.String()] = device
	return nil
}

// FindByID retrieves a device by its ID.
func (r *memoryRepository) FindByID(ctx context.Context, orgID, deviceID uuid.UUID) (*Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if dev, ok := r.devices[deviceID.String()]; ok && dev.OrgID == orgID {
		return dev, nil
	}
	return nil, ErrDeviceNotFound
}

// ListByOrg lists all devices for a given organization.
func (r *memoryRepository) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*Device
	for _, dev := range r.devices {
		if dev.OrgID == orgID {
			results = append(results, dev)
		}
	}
	return results, nil
}

// SaveTelemetry saves telemetry data for a device.
func (r *memoryRepository) SaveTelemetry(ctx context.Context, deviceID uuid.UUID, ts time.Time, points map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	dataPoint := map[string]interface{}{
		"timestamp": ts,
	}
	for k, v := range points {
		dataPoint[k] = v
	}
	r.telemetry[deviceID.String()] = append(r.telemetry[deviceID.String()], dataPoint)
	return nil
}