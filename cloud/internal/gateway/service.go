package gateway

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the gateway domain.
type Service struct {
	repo Repository
}

// NewService creates a new gateway service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// RegisterGateway handles the logic for a gateway registering itself with the cloud.
// It can be an idempotent operation.
func (s *Service) RegisterGateway(ctx context.Context, orgID uuid.UUID, gatewayID uuid.UUID, name string) (*Gateway, error) {
	gw, err := s.repo.FindByID(ctx, orgID, gatewayID)
	if err != nil {
		if err == ErrGatewayNotFound {
			// Gateway not found, create a new one
			newGw, createErr := NewGateway(orgID, name)
			if createErr != nil {
				return nil, createErr
			}
			newGw.ID = gatewayID // Use the ID provided by the gateway
			newGw.Status = "Online"
			
			if createErr := s.repo.Create(ctx, newGw); createErr != nil {
				return nil, createErr
			}
			return newGw, nil
		}
		return nil, err
	}

	// Gateway found, update its status and name if changed
	gw.Status = "Online"
	gw.LastSeenAt = time.Now().UTC()
	if name != "" && gw.Name != name {
		gw.Name = name
	}
	
	if err := s.repo.Update(ctx, gw); err != nil {
		return nil, err
	}
	return gw, nil
}

// HandleHeartbeat updates the gateway's last seen time to mark it as online.
func (s *Service) HandleHeartbeat(ctx context.Context, orgID, gatewayID uuid.UUID) error {
	gw, err := s.repo.FindByID(ctx, orgID, gatewayID)
	if err != nil {
		return err
	}

	gw.Status = "Online"
	gw.LastSeenAt = time.Now().UTC()
	return s.repo.Update(ctx, gw)
}