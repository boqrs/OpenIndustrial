package permission

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a single, granular action that can be permitted or denied.
// e.g., "users.create", "devices.read", "workorders.start"
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"` // The permission string itself
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewPermission creates a new Permission entity.
// In a real system, these might be seeded from a static list rather than created via API.
func NewPermission(name, description string) (*Permission, error) {
	if name == "" {
		return nil, ErrPermissionNameRequired
	}
	return &Permission{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}, nil
}