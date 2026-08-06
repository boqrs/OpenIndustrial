package org

import (
	"context"
	"time"
)

// Service encapsulates the business logic for the organization domain.
type Service struct {
	repo Repository
}

// NewService creates a new organization service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateOrg handles the business logic of creating a new organization.
func (s *Service) CreateOrg(ctx context.Context, name string) (*Org, error) {
	org, err := NewOrg(name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrgByID retrieves an organization by its ID.
func (s *Service) GetOrgByID(ctx context.Context, id string) (*Org, error) {
	// In a real application, you might add more logic here,
	// like checking permissions, caching, etc.
	return s.repo.FindByID(ctx, id)
}

// UpdateOrg updates an organization's information.
func (s *Service) UpdateOrg(ctx context.Context, id string, name string) (*Org, error) {
	org, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, ErrOrgNameRequired
	}

	org.Name = name
	org.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// DeleteOrg removes an organization.
func (s *Service) DeleteOrg(ctx context.Context, id string) error {
	// You might want to check for existing users or resources in the org
	// before allowing deletion. This is where that business logic would go.
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}