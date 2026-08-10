package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the Resource Kernel.
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

// CreateProductParams defines the parameters for creating a new product.
type CreateProductParams struct {
	TenantID     uuid.UUID
	Name         string
	Properties   Properties
	OwnerGroupID uuid.UUID // The ID of the group that will own this new product.
}

// CreateProduct is our first "template" business method.
// It handles the creation of a new product and associates it with an owner group.
func (s *Service) CreateProduct(ctx context.Context, params CreateProductParams) (*Resource, error) {
	// 1. Prepare the new resource object
	newProduct := &Resource{
		ID:         uuid.New(),
		TenantID:   params.TenantID,
		Type:       "product", // Hardcoded type for this specific business logic
		Name:       params.Name,
		Properties: params.Properties,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// TODO: Wrap the following two operations in a single database transaction.
	// This ensures that if adding the resource to a group fails, the resource creation is rolled back.

	// 2. Save the new resource to the database
	if err := s.resourceRepo.CreateResource(ctx, newProduct); err != nil {
		return nil, err
	}

	// 3. Associate the new resource with its owner group (core of ABAC)
	if err := s.groupRepo.AddResourceToGroup(ctx, params.TenantID, newProduct.ID, params.OwnerGroupID); err != nil {
		// In a real transactional implementation, the resource creation would be rolled back here.
		return nil, err
	}

	// 4. Return the newly created product
	return newProduct, nil
}