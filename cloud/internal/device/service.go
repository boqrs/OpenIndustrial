package device

import (
	"context"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
)

// Service defines the business logic for the device domain.
type Service interface {
	CreateDevice(ctx context.Context, orgID, gatewayID, name, model string) (*Device, error)
	GetDeviceByID(ctx context.Context, orgID, deviceID string) (*Device, error)
	ListDevicesByOrg(ctx context.Context, orgID string) ([]*Device, error)
	ListDevicesForGateway(ctx context.Context, orgID, gatewayID string) ([]*Device, error)
	RecordTelemetry(ctx context.Context, deviceID string, ts time.Time, points map[string]interface{}) error
}

type service struct {
	repo     Repository
	eventBus event.Bus
}

// NewService creates a new device service.
func NewService(repo Repository, bus event.Bus) Service {
	return &service{
		repo:     repo,
		eventBus: bus,
	}
}

// CreateDevice creates a new device.
func (s *service) CreateDevice(ctx context.Context, orgID, gatewayID, name, model string) (*Device, error) {
	d, err := NewDevice(orgID, gatewayID, name, model)
	if err != nil {
		return nil, err
	}
	return d, s.repo.Create(ctx, d)
}

// GetDeviceByID retrieves a device by its ID.
func (s *service) GetDeviceByID(ctx context.Context, orgID, deviceID string) (*Device, error) {
	return s.repo.FindByID(ctx, orgID, deviceID)
}

// ListDevicesByOrg lists all devices for an organization.
func (s *service) ListDevicesByOrg(ctx context.Context, orgID string) ([]*Device, error) {
	return s.repo.ListByOrg(ctx, orgID)
}

// ListDevicesForGateway lists all devices for a specific gateway.
func (s *service) ListDevicesForGateway(ctx context.Context, orgID, gatewayID string) ([]*Device, error) {
	return s.repo.ListByGateway(ctx, orgID, gatewayID)
}

// RecordTelemetry records new telemetry data for a device.
func (s *service) RecordTelemetry(ctx context.Context, deviceID string, ts time.Time, points map[string]interface{}) error {
	return s.repo.SaveTelemetry(ctx, deviceID, ts, points)
}