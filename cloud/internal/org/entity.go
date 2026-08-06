package org

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/shared"
	"github.com/google/uuid"
)

// OrgType defines the type of an organization.
type OrgType string

const (
	OrgTypeFactory  OrgType = "factory"
	OrgTypeCustomer OrgType = "customer"
	OrgTypePartner  OrgType = "partner"
	OrgTypePlatform OrgType = "platform"
)

// OrgStatus defines the status of an organization.
type OrgStatus string

const (
	OrgStatusActive   OrgStatus = "active"
	OrgStatusInactive OrgStatus = "inactive"
)

// Organization represents a single organization or tenant in the system.
// It is the top-level container for all resources and can be structured hierarchically.
type Organization struct {
	shared.BaseEntity
	Name     string    `json:"name"`
	Type     OrgType   `json:"type"`
	ParentID string    `json:"parent_id,omitempty"`
	Status   OrgStatus `json:"status"`
}

// NewOrganization creates a new Organization entity.
func NewOrganization(name string, orgType OrgType, parentID string) (*Organization, error) {
	if name == "" {
		return nil, ErrOrgNameRequired
	}

	now := time.Now().UTC()
	return &Organization{
		BaseEntity: shared.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:     name,
		Type:     orgType,
		ParentID: parentID,
		Status:   OrgStatusActive,
	}, nil
}