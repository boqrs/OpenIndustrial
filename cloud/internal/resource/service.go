package resource

import (
	"context"
)

// Service provides business logic for managing resources.
type Service struct {
	repo Repository
}

// NewService creates a new resource service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateResource creates a new resource.
func (s *Service) CreateResource(ctx context.Context, name, description string, resourceType Type, parentID *string) (*Resource, error) {
	resource := &Resource{
		Name:        name,
		Description: description,
		Type:        resourceType,
	}
	if parentID != nil {
		resource.ParentID = *parentID
	}

	return s.repo.CreateResource(ctx, resource)
}

// GetResourceTree retrieves a resource and its descendants as a tree structure.
// Note: This is a placeholder. A real implementation would be more complex.
func (s *Service) GetResourceTree(ctx context.Context, id string) (*Resource, error) {
	// The GetTree method was removed as it was not defined in the repository interface.
	// This method needs a proper implementation based on graph traversal (e.g., using ListRelations).
	// For now, it will just return the single resource.
	return s.repo.GetResource(ctx, id)
}

// UpdateResource updates an existing resource.
func (s *Service) UpdateResource(ctx context.Context, id string, name, description string) (*Resource, error) {
	resource, err := s.repo.GetResource(ctx, id)
	if err != nil {
		return nil, err
	}

	resource.Name = name
	resource.Description = description

	return s.repo.UpdateResource(ctx, resource)
}

// DeleteResource deletes a resource.
func (s *Service) DeleteResource(ctx context.Context, id string) error {
	return s.repo.DeleteResource(ctx, id)
}