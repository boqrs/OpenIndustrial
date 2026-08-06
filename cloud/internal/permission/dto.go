package permission

import (
	"time"
)

// PermissionResponse defines the structure for a response containing permission details.
type PermissionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToPermissionResponse converts a Permission entity to a PermissionResponse DTO.
func ToPermissionResponse(p *Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

// ToPermissionListResponse converts a slice of Permission entities to a slice of PermissionResponse DTOs.
func ToPermissionListResponse(permissions []*Permission) []*PermissionResponse {
	res := make([]*PermissionResponse, len(permissions))
	for i, p := range permissions {
		res[i] = ToPermissionResponse(p)
	}
	return res
}