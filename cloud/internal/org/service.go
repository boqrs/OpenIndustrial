package org

import "context"

// Service provides business logic for organization management.
// It orchestrates the application's business rules and relies on the repository for data access.
type Service struct {
	repo Repository
}

// NewService creates a new organization service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateOrganization handles the business logic of creating a new organization.
func (s *Service) CreateOrganization(ctx context.Context, name string, orgType OrgType, parentID string) (*Organization, error) {
	// 1. Use the factory to create a new Organization entity
	org, err := NewOrganization(name, orgType, parentID)
	if err != nil {
		return nil, err
	}

	// 2. (Optional) Perform additional business rule validations here.
	// For example, check if an organization with the same name already exists.

	// 3. Persist the new organization using the repository
	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}

	// 4. (Optional) Dispatch a domain event, e.g., "OrganizationCreated".

	return org, nil
}

// GetOrganization retrieves an organization by its ID.
func (s *Service) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	return s.repo.FindByID(ctx, id)
}