package identity

import (
	"context"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
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
	ListUsers(ctx context.Context, tenantID uuid.UUID, params ListUsersRepoReq) ([]*model.User, error)
}

// Service defines the interface for the identity service.
// This is the contract that the rest of the application will use.
type Service interface {
	RegisterNewTenant(ctx context.Context, req *RegisterTenantRequest) (*RegisterTenantResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetCurrentUser(ctx context.Context, tenantID, userID uuid.UUID) (*GetCurrentUserResponse, error)
	CreateUser(ctx context.Context, tenantID uuid.UUID, req *CreateUserRequest) (*CreateUserResponse, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID, req *ListUsersRequest) ([]*UserResponse, error)
	UpdateUser(ctx context.Context, tenantID, userID uuid.UUID, req *UpdateUserRequest) error
	DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]*RoleResponse, error)
	AssignRoleToUser(ctx context.Context, tenantID, userID uuid.UUID, req *AssignRoleToUserRequest) error
	ListUserGroups(ctx context.Context, tenantID, userID uuid.UUID) ([]*GroupResponse, error)
}
