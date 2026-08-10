package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Resource represents a digital twin of a physical or logical asset.
type Resource struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TenantID     uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	Type         string    `db:"type" json:"type"`
	SerialNumber string    `db:"serial_number" json:"serial_number"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// ResourceRelation defines a typed relationship between two resources.
type ResourceRelation struct {
	ID               uuid.UUID `db:"id" json:"id"`
	TenantID         uuid.UUID `db:"tenant_id" json:"tenant_id"`
	SourceResourceID uuid.UUID `db:"source_resource_id" json:"source_resource_id"`
	TargetResourceID uuid.UUID `db:"target_resource_id" json:"target_resource_id"`
	RelationType     string    `db:"relation_type" json:"relation_type"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// Group represents a collection of users or resources for permissioning.
type Group struct {
	ID          uuid.UUID `db:"id" json:"id"`
	TenantID    uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// --- Repositories: Interfaces for data access ---

// ResourceRepository defines the interface for accessing resource data.
type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *Resource) error
	GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*Resource, error)
	ListResources(ctx context.Context, tenantID uuid.UUID) ([]*Resource, error)
	// CreateResourceRelation is also needed if you manage relations here
	CreateResourceRelation(ctx context.Context, rel *ResourceRelation) error
}

// GroupRepository defines the interface for accessing group data.
type GroupRepository interface {
	CreateGroup(ctx context.Context, group *Group) error
	GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*Group, error)
	AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	AddResourceToGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error
	RemoveResourceFromGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error
	ListGroupsByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*Group, error)
}

// AuthorizationRepository defines the interface for authorization checks.
type AuthorizationRepository interface {
	CheckUserInSameGroupAsResource(ctx context.Context, tenantID, userID, resourceID uuid.UUID) (bool, error)
}