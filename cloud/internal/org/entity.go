package org

import (
	"time"

	"github.com/google/uuid"
)

// Org represents a single organization or tenant in the system.
// It is the top-level container for all resources.
type Org struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewOrg creates a new Organization entity.
func NewOrg(name string) (*Org, error) {
	// Basic validation
	if name == "" {
		return nil, ErrOrgNameRequired
	}

	now := time.Now().UTC()
	return &Org{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}