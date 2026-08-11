package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Group represents a collection of users for permissioning.
type Group struct {
	ID          uuid.UUID `db:"id" json:"id"`
	TenantID    uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// GroupRepository defines the interface for accessing group data.
type GroupRepository interface {
	CreateGroup(ctx context.Context, group *Group) error
	GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*Group, error)
	AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	ListGroupsByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*Group, error)
	// Note: The methods for adding/removing resources from groups are intentionally
	// left out here. That logic belongs to a higher-level service or the resource kernel
	// itself, which would hold a reference to a group ID.
}