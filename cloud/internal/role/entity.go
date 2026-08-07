package role

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/permission"
	"github.com/google/uuid"
)

// Role represents a collection of permissions.
type Role struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	// A role is a collection of permissions.
	Permissions []*permission.Permission `json:"permissions"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// NewRole creates a new Role.
func NewRole(orgID uuid.UUID, name, description string) *Role {
	now := time.Now().UTC()
	return &Role{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        name,
		Description: description,
		Permissions: make([]*permission.Permission, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}