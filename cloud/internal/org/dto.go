package org

import "github.com/google/uuid"

// CreateOrganizationRequest defines the structure for a request to create a new organization.
type CreateOrganizationRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     OrgType `json:"type" binding:"required"`
	ParentID string  `json:"parent_id,omitempty"`
}

// OrganizationResponse defines the structure for a standard organization API response.
type OrganizationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      OrgType   `json:"type"`
	ParentID  string    `json:"parent_id,omitempty"`
	Status    OrgStatus `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// ToOrganizationResponse converts an Organization entity to an OrganizationResponse DTO.
func ToOrganizationResponse(org *Organization) *OrganizationResponse {
	return &OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Type:      org.Type,
		ParentID:  org.ParentID,
		Status:    org.Status,
		CreatedAt: org.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: org.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}