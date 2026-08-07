package org

import (
	"time"

	"github.com/google/uuid"
)

// CreateOrgRequest defines the request for creating an organization.
type CreateOrgRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Type        OrgType `json:"type" binding:"required"`
}

// OrgResponse is the DTO for an organization.
type OrgResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        OrgType   `json:"type"`
	Status      OrgStatus `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToOrgResponse converts an Organization entity to a DTO.
func ToOrgResponse(org *Organization) *OrgResponse {
	return &OrgResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		Type:        org.Type,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}