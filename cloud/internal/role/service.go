package role

import (
	"context"

	"github.com/google/uuid"
)

// Service provides role-related business logic.
type Service struct {
	repo Repository
	// permissionService permission.Service // To validate permission IDs
}

// NewService creates a new role service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateRole creates a new role.
func (s *Service) CreateRole(ctx context.Context, orgID uuid.UUID, req *CreateRoleRequest) (*Role, error) {
	role := NewRole(orgID, req.Name, req.Description)

	// In a real system, you would fetch the permission objects from the DB
	// based on req.PermissionIDs and assign them to role.Permissions.
	// For now, we'll leave it empty.

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

// GetRoleByID retrieves a role by its ID.
func (s *Service) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*Role, error) {
	return s.repo.GetRoleByID(ctx, roleID)
}

// ListRolesForOrg lists all roles for a given organization.
func (s *Service) ListRolesForOrg(ctx context.Context, orgID uuid.UUID) ([]*Role, error) {
	return s.repo.ListRolesByOrg(ctx, orgID)
}