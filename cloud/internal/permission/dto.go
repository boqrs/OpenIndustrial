package permission

import (
	"time"
)

// PermissionResponse is the DTO for a permission.
type PermissionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToPermissionResponse converts a Permission entity to a DTO.
func ToPermissionResponse(p *Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

// ToPermissionListResponse converts a slice of Permission entities to a slice of DTOs.
func ToPermissionListResponse(permissions []*Permission) []*PermissionResponse {
	list := make([]*PermissionResponse, len(permissions))
	for i, p := range permissions {
		list[i] = ToPermissionResponse(p)
	}
	return list
}