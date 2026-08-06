package org

import "time"

// CreateOrgRequest defines the structure for a request to create a new organization.
type CreateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateOrgRequest defines the structure for a request to update an organization.
type UpdateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

// OrgResponse defines the structure for a response containing organization details.
// It might be a subset of the Org entity or include additional computed fields.
type OrgResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToOrgResponse converts an Org entity to an OrgResponse DTO.
func ToOrgResponse(org *Org) *OrgResponse {
	return &OrgResponse{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
}