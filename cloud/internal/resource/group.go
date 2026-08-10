package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Group represents a collection of users or resources for permission management (ABAC).
type Group struct {
	ID          uuid.UUID `db:"id" json:"id"`
	TenantID    uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// GroupRepository defines the interface for group-related database operations.
type GroupRepository interface {
	CreateGroup(ctx context.Context, group *Group) error
	GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*Group, error)
	AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	AddResourceToGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error
	RemoveResourceFromGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error
}

// AuthorizationRepository defines the interface for complex authorization checks.
type AuthorizationRepository interface {
	// CheckUserInSameGroupAsResource is the core of our ABAC instance-level check.
	// It verifies if a user and a resource share at least one common group.
	CheckUserInSameGroupAsResource(ctx context.Context, tenantID, userID, resourceID uuid.UUID) (bool, error)
}