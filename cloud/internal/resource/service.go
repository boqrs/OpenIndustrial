package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service provides use cases for the resource domain.
type Service struct {
	resourceRepo ResourceRepository
	groupRepo    GroupRepository
}

// NewService creates a new resource service.
func NewService(resourceRepo ResourceRepository, groupRepo GroupRepository) *Service {
	return &Service{
		resourceRepo: resourceRepo,
		groupRepo:    groupRepo,
	}
}

// CreateProductParams defines the parameters for creating a new product (Resource).
// These fields now directly match the Resource struct.
type CreateProductParams struct {
	Name         string
	Description  string
	Type         string
	SerialNumber string
	OwnerGroupID uuid.UUID
}

// CreateProduct handles the business logic for creating a new product.
// It creates the resource and associates it with an owner group.
func (s *Service) CreateProduct(ctx context.Context, tenantID uuid.UUID, params CreateProductParams) (*Resource, error) {
	// In a real app, this should be a single transaction.

	// 1. Create the resource using the corrected fields.
	resource := &Resource{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         params.Name,
		Description:  params.Description,  // CORRECTED
		Type:         params.Type,
		SerialNumber: params.SerialNumber, // CORRECTED
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.resourceRepo.CreateResource(ctx, resource); err != nil {
		return nil, err
	}

	// 2. Associate the new resource with its owner group.
	if err := s.groupRepo.AddResourceToGroup(ctx, tenantID, resource.ID, params.OwnerGroupID); err != nil {
		// If this fails, we should ideally roll back the resource creation.
		return nil, err
	}

	return resource, nil
}

// ListUserGroups retrieves all groups that the given user is a member of.
func (s *Service) ListUserGroups(ctx context.Context, tenantID, userID uuid.UUID) ([]*Group, error) {
	groups, err := s.groupRepo.ListGroupsByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	// Ensure we always return a non-nil slice for JSON marshalling, which is good practice.
	if groups == nil {
		return []*Group{}, nil
	}
	return groups, nil
}