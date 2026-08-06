package permission

import (
	"context"
)

// Service encapsulates the business logic for the permission domain.
type Service struct {
	repo Repository
}

// NewService creates a new permission service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ListAllPermissions lists all available permissions in the system.
func (s *Service) ListAllPermissions(ctx context.Context) ([]*Permission, error) {
	return s.repo.ListAll(ctx)
}