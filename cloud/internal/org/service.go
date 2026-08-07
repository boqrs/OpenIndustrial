package org

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrgType = errors.New("invalid organization type")
)

// Service provides organization-related business logic.
type Service struct {
	repo Repository
	// eventBus publisher
}

// NewService creates a new organization service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateOrganization creates a new organization.
func (s *Service) CreateOrganization(ctx context.Context, name, description string, orgType OrgType) (*Organization, error) {
	if !orgType.Valid() {
		return nil, ErrInvalidOrgType
	}
	org := NewOrganization(name, description, orgType)

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}
	// TODO: publish event
	return org, nil
}

// GetOrganization retrieves an organization by its ID.
func (s *Service) GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error) {
	return s.repo.Get(ctx, id)
}