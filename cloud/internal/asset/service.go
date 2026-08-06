package asset

import (
	"context"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the asset domain.
type Service struct {
	repo Repository
}

// NewService creates a new asset service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateAsset handles the business logic of creating a new asset.
func (s *Service) CreateAsset(ctx context.Context, orgID, productID uuid.UUID, sn string) (*Asset, error) {
	// Here you might add logic to check for SN uniqueness within the organization.
	asset, err := NewAsset(orgID, productID, sn)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// GetAssetBySN retrieves an asset by its serial number. This will be a very common operation.
func (s *Service) GetAssetBySN(ctx context.Context, orgID uuid.UUID, sn string) (*Asset, error) {
	return s.repo.FindBySN(ctx, orgID, sn)
}

// ListAssetsForOrg lists all assets for a given organization.
func (s *Service) ListAssetsForOrg(ctx context.Context, orgID uuid.UUID) ([]*Asset, error) {
	return s.repo.ListByOrg(ctx, orgID)
}