package role

import (
	"time"
)

// CreateRoleRequest defines the request for creating a role.
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	// List of permission IDs to assign to the role.
	PermissionIDs []string `json:"permission_ids"`
}

// RoleResponse is the DTO for a role.
type RoleResponse struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToRoleResponse converts a Role entity to a DTO.
func ToRoleResponse(role *Role) *RoleResponse {
	return &RoleResponse{
		ID:          role.ID.String(),
		OrgID:       role.OrgID.String(),
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// ToRoleListResponse converts a slice of Role entities to a slice of DTOs.
func ToRoleListResponse(roles []*Role) []*RoleResponse {
	responses := make([]*RoleResponse, len(roles))
	for i, r := range roles {
		responses[i] = ToRoleResponse(r)
	}
	return responses
}