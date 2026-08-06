package org

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/shared"

	"github.com/google/uuid"
)

// Org represents a single organization or tenant in the system.
// It is the top-level container for all resources.
type Org struct {
	shared.BaseEntity
	Name string `json:"name"`
}

// NewOrg creates a new Organization entity.
func NewOrg(name string) (*Org, error) {
	// Basic validation
	if name == "" {
		return nil, ErrOrgNameRequired
	}

	now := time.Now().UTC()
	return &Org{
		BaseEntity: shared.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name: name,
	}, nil
}