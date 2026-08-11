package identity

import (
	"context"

	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)

// GroupRepository defines the interface for accessing group data.
type GroupRepository interface {
	CreateGroup(ctx context.Context, group *model.Group) error
	GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*model.Group, error)
	AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error
	ListGroupsByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Group, error)
	// Note: The methods for adding/removing resources from groups are intentionally
	// left out here. That logic belongs to a higher-level service or the resource kernel
	// itself, which would hold a reference to a group ID.
}


// PermissionRepository defines the interface for permission-related database operations.
type PermissionRepository interface {
	// CheckPermissionForUser checks if a user has a specific permission through their roles.
	CheckPermissionForUser(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error)

	// CreatePermission adds a new permission to the database.
	CreatePermission(ctx context.Context, p *model.Permission) error

	// GetPermission retrieves a permission by its key and action.
	GetPermission(ctx context.Context, resourceKey, action string) (*model.Permission, error)

	// ListPermissionsByRole retrieves all permissions associated with a specific role.
	ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*model.Permission, error)
}

// RoleRepository defines the interface for role persistence.
type RoleRepository interface {
	CreateRole(ctx context.Context, role *model.Role) error
	GetRoleByID(ctx context.Context, tenantID, id uuid.UUID) (*model.Role, error)
	GetRoleByName(ctx context.Context, tenantID uuid.UUID, name string) (*model.Role, error)
	AddUserToRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error
}

// TenantRepository defines the interface for tenant persistence.
type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *model.Tenant) error
	GetTenantByID(ctx context.Context, id uuid.UUID) (*model.Tenant, error)
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*model.User, error)
	CreatePrincipal(ctx context.Context, principal *model.Principal) error
	GetPrincipal(ctx context.Context, tenantID uuid.UUID, provider, identifier string) (*model.Principal, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID, params param.ListUsersRepoReq) ([]*model.User, error)
}