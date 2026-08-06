package device

import (
	"context"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the device domain.
type Service struct {
	repo Repository
}

// NewService creates a new device service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateDevice handles the business logic of creating a new device.
func (s *Service) CreateDevice(ctx context.Context, orgID, gatewayID uuid.UUID, name, model string) (*Device, error) {
	dev, err := NewDevice(orgID, gatewayID, name, model)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, dev); err != nil {
		return nil, err
	}

	return dev, nil
}

// GetDeviceByID retrieves a device by its ID.
func (s *Service) GetDeviceByID(ctx context.Context, orgID, deviceID uuid.UUID) (*Device, error) {
	return s.repo.FindByID(ctx, orgID, deviceID)
}

// ListDevicesForGateway lists all devices connected to a specific gateway.
func (s *Service) ListDevicesForGateway(ctx context.Context, orgID, gatewayID uuid.UUID) ([]*Device, error) {
	return s.repo.ListByGateway(ctx, orgID, gatewayID)
}