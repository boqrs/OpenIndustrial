package org

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents a tenant or a partner in the system.
type Organization struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Type        OrgType    `json:"type"`
	Status      OrgStatus  `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewOrganization creates a new Organization.
func NewOrganization(name, description string, orgType OrgType) *Organization {
	return &Organization{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Type:        orgType,
		Status:      OrgStatusActive, // Default status
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}