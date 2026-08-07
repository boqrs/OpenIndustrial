package permission

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a single, named permission.
// e.g., "products:read", "users:write"
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"` // The unique name of the permission
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Policy represents a rule for the authorization engine (e.g., Casbin).
// It's a tuple of (Subject, Object, Action).
// Example: (role:admin, /api/v1/products, POST)
type Policy struct {
	Subject string // e.g., a role ID or user ID
	Object  string // e.g., a resource path or object ID
	Action  string // e.g., "read", "write", "POST"
}