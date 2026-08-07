package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service provides business logic for managing gateways.
type Service struct {
	repo Repository
}

// NewService creates a new gateway service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// RegisterGateway creates a new gateway from a registration request.
// It assigns a name and saves the new gateway.
func (s *Service) RegisterGateway(ctx context.Context, model string) (*Gateway, error) {
	// For example, generate a default name or use a different logic
	name := "Gateway " + uuid.New().String()[:8]
	gw := NewGateway(name, "", model, "") // resourceID can be assigned later
	return s.repo.Create(ctx, gw)
}

// ListGateways retrieves all gateways.
func (s *Service) ListGateways(ctx context.Context) ([]*Gateway, error) {
	return s.repo.List(ctx)
}

// GetGateway retrieves a gateway by its ID.
func (s *Service) GetGateway(ctx context.Context, id uuid.UUID) (*Gateway, error) {
	return s.repo.Get(ctx, id)
}

// Heartbeat updates the gateway's last seen timestamp.
func (s *Service) Heartbeat(ctx context.Context, id uuid.UUID) error {
	gw, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	gw.Status = "Online"
	gw.LastHeartbeat = time.Now().UTC()
	return s.repo.Update(ctx, gw)
}

// UpdateGateway updates an existing gateway.
func (s *Service) UpdateGateway(ctx context.Context, gw *Gateway) error {
	return s.repo.Update(ctx, gw)
}