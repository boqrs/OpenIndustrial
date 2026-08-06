package role

import (
	"context"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the role domain.
type Service struct {
	repo Repository
}

// NewService creates a new role service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateRole handles the business logic of creating a new role.
func (s *Service) CreateRole(ctx context.Context, orgID uuid.UUID, name string, permissions []Permission) (*Role, error) {
	role, err := NewRole(orgID, name, permissions)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID retrieves a role by its ID within a specific organization.
func (s *Service) GetRoleByID(ctx context.Context, orgID, roleID uuid.UUID) (*Role, error) {
	return s.repo.FindByID(ctx, orgID, roleID)
}

// ListRolesForOrg lists all roles for a given organization.
func (s *Service) ListRolesForOrg(ctx context.Context, orgID uuid.UUID) ([]*Role, error) {
	return s.repo.ListByOrg(ctx, orgID)
}