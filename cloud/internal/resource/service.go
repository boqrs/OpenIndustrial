package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service provides business logic for resources.
type Service struct {
	resourceRepo ResourceRepository
	groupRepo    GroupRepository
	authzRepo    AuthorizationRepository
}

// NewService creates a new resource service.
func NewService(
	resourceRepo ResourceRepository,
	groupRepo GroupRepository,
	authzRepo AuthorizationRepository,
) *Service {
	return &Service{
		resourceRepo: resourceRepo,
		groupRepo:    groupRepo,
		authzRepo:    authzRepo,
	}
}

// CreateProductParams defines the parameters for creating a product.
type CreateProductParams struct {
	Name         string
	Description  string
	Type         string
	SerialNumber string
	OwnerGroupID uuid.UUID
}

// CreateProduct creates a new product resource and associates it with an owner group.
func (s *Service) CreateProduct(ctx context.Context, tenantID uuid.UUID, params CreateProductParams) (*Resource, error) {
	// In a real transaction, you'd wrap these calls.
	resource := &Resource{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         params.Name,
		Description:  params.Description,
		Type:         params.Type,
		SerialNumber: params.SerialNumber,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.resourceRepo.CreateResource(ctx, resource); err != nil {
		return nil, err
	}

	// Associate the resource with its owner group.
	if err := s.groupRepo.AddResourceToGroup(ctx, tenantID, resource.ID, params.OwnerGroupID); err != nil {
		// In a real scenario, you might want to roll back the resource creation here.
		return nil, err
	}

	return resource, nil
}

// ListUserGroups lists all groups a user is a member of.
func (s *Service) ListUserGroups(ctx context.Context, tenantID, userID uuid.UUID) ([]*Group, error) {
	groups, err := s.groupRepo.ListGroupsByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	// Ensure we always return a non-nil slice
	if groups == nil {
		return []*Group{}, nil
	}
	return groups, nil
}

// GetResource retrieves a single resource by its ID.
// It will perform authorization checks in the future.
func (s *Service) GetResource(ctx context.Context, tenantID, resourceID uuid.UUID) (*Resource, error) {
	// TODO: Add authorization check here.
	// For example: s.authzRepo.CheckUserPermissionForResource(ctx, userID, resourceID, "read")
	return s.resourceRepo.GetResourceByID(ctx, tenantID, resourceID)
}

// ListResources retrieves a list of resources for a tenant.
// It will perform authorization filtering in the future.
func (s *Service) ListResources(ctx context.Context, tenantID uuid.UUID) ([]*Resource, error) {
	// TODO: Add authorization filtering here.
	// The current implementation returns all resources for the tenant.
	// A future implementation would filter based on the user's group memberships.
	return s.resourceRepo.ListResources(ctx, tenantID)
}