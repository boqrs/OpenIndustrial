package permission

import (
	"context"
)

// Service provides business logic for managing permissions and policies.
type Service struct {
	repo Repository
}

// NewService creates a new permission service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListAll retrieves all available permissions.
func (s *Service) ListAll(ctx context.Context) ([]*Permission, error) {
	return s.repo.ListAllPermissions(ctx)
}

// GetPoliciesForRole retrieves all policies for a given role.
func (s *Service) GetPoliciesForRole(ctx context.Context, roleID string) ([]*Policy, error) {
	return s.repo.GetPoliciesForSubject(ctx, roleID)
}