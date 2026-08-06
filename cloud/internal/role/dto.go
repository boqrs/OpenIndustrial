package role

import (
	"time"
)

// CreateRoleRequest defines the structure for a request to create a new role.
type CreateRoleRequest struct {
	Name        string       `json:"name" binding:"required"`
	Permissions []Permission `json:"permissions"`
}

// RoleResponse defines the structure for a response containing role details.
type RoleResponse struct {
	ID          string       `json:"id"`
	OrgID       string       `json:"org_id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ToRoleResponse converts a Role entity to a RoleResponse DTO.
func ToRoleResponse(role *Role) *RoleResponse {
	return &RoleResponse{
		ID:          role.ID.String(),
		OrgID:       role.OrgID.String(),
		Name:        role.Name,
		Permissions: role.Permissions,
		CreatedAt:   role.CreatedAt,
	}
}

// ToRoleListResponse converts a slice of Role entities to a slice of RoleResponse DTOs.
func ToRoleListResponse(roles []*Role) []*RoleResponse {
	res := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		res[i] = ToRoleResponse(role)
	}
	return res
}