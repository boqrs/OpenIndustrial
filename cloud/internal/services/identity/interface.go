package identity

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// =====================================================
// Repository
// =====================================================

// UserRepository defines persistence operations for the
// identity domain.
//
// Identity persistence includes:
//
//   - Tenant
//   - User
//   - Role
//   - Permission
//   - Invitation
//   - Principal
//
// Tenant is the isolation boundary of the identity domain.
//
// A user belongs to exactly one tenant and exactly one role.
// Permissions are granted to roles through role_permissions.
type UserRepository interface {

	// =====================================================
	// Tenant
	// =====================================================

	// GetTenantByCode retrieves a tenant by its unique code.
	//
	// Tenant code is used by the login flow to identify
	// the factory before authenticating the user.
	GetTenantByCode(
		ctx context.Context,
		code string,
	) (*model.Tenant, error)

	// GetTenantByID retrieves a tenant by its database ID.
	GetTenantByID(
		ctx context.Context,
		id uint,
	) (*model.Tenant, error)

	// =====================================================
	// User
	// =====================================================

	// CreateUser creates a user inside a tenant.
	//
	// RoleID must reference a role belonging to the same tenant.
	CreateUser(
		ctx context.Context,
		user *model.User,
	) error

	// GetUserByID retrieves a user by tenant ID and user UUID.
	//
	// TenantID is required to enforce tenant isolation.
	GetUserByID(
		ctx context.Context,
		tenantID uint,
		userID uuid.UUID,
	) (*model.User, error)

	// GetUserByEmail retrieves a user by tenant ID and email.
	//
	// Email uniqueness is tenant-scoped.
	GetUserByEmail(
		ctx context.Context,
		tenantID uint,
		email string,
	) (*model.User, error)

	// ListUsers retrieves users belonging to a tenant.
	ListUsers(
		ctx context.Context,
		tenantID uint,
		limit int,
		offset int,
		status string,
		userType string,
		keyword string,
	) ([]*model.User, error)

	// UpdateUser updates an existing user.
	UpdateUser(
		ctx context.Context,
		user *model.User,
	) error

	// =====================================================
	// Role
	// =====================================================

	// GetRole retrieves a role by tenant and name.
	//
	// Role names are unique within a tenant.
	GetRole(
		ctx context.Context,
		tenantID uint,
		name string,
	) (*model.Role, error)

	// GetRoleByID retrieves a role by tenant and role ID.
	//
	// TenantID is intentionally required to prevent a role
	// belonging to another tenant from being used.
	GetRoleByID(
		ctx context.Context,
		tenantID uint,
		roleID uint,
	) (*model.Role, error)

	// CheckRolePermission checks whether the specified role
	// belonging to the specified tenant has a permission.
	//
	// The permission relationship is:
	//
	// Role
	//   -> RolePermission
	//      -> Permission
	CheckRolePermission(
		ctx context.Context,
		tenantID uint,
		roleID uint,
		permissionName string,
	) (bool, error)

	// =====================================================
	// Permission
	// =====================================================

	// CreatePermission creates a system-defined permission.
	//
	// Permissions are global and are not tenant-specific.
	// This is primarily intended for bootstrap/system data.
	CreatePermission(
		ctx context.Context,
		permission *model.Permission,
	) error

	// GetPermission retrieves a permission by resource and action.
	//
	// The resulting permission name is:
	//
	//     resourceKey:action
	GetPermission(
		ctx context.Context,
		resourceKey string,
		action string,
	) (*model.Permission, error)

	// =====================================================
	// Invitation
	// =====================================================

	// CreateInvitation creates a user invitation.
	CreateInvitation(
		ctx context.Context,
		invitation *model.UserInvitation,
	) error

	// GetInvitationByToken retrieves an invitation by its hashed token.
	GetInvitationByToken(
		ctx context.Context,
		tokenHash string,
	) (*model.UserInvitation, error)

	// GetInvitationByID retrieves an invitation by tenant and UUID.
	GetInvitationByID(
		ctx context.Context,
		tenantID uint,
		id uuid.UUID,
	) (*model.UserInvitation, error)

	// UpdateInvitation updates an existing invitation.
	UpdateInvitation(
		ctx context.Context,
		invitation *model.UserInvitation,
	) error

	// =====================================================
	// Principal
	// =====================================================

	// CreatePrincipal creates an authentication principal.
	CreatePrincipal(
		ctx context.Context,
		principal *model.Principal,
	) error

	// GetPrincipal retrieves a principal by tenant, provider,
	// and identifier.
	GetPrincipal(
		ctx context.Context,
		tenantID uint,
		provider string,
		identifier string,
	) (*model.Principal, error)

	// UpdatePrincipal updates an existing authentication principal.
	UpdatePrincipal(
		ctx context.Context,
		principal *model.Principal,
	) error

	GetAdminByTenantID(
		ctx context.Context,
		tenantID uint,
	) (*model.User, error)

	GetUserStats(
		ctx context.Context,
		tenantID uint,
	) (
		*UserStats,
		error,
	)

	ListRoles(
		ctx context.Context,
		tenantID uint,
	) ([]*model.Role, error)
}

// =====================================================
// Service
// =====================================================

// Service
//
// Identity business service.
//
// Responsibility:
//
// - Authentication
// - User lifecycle
// - Invitation
// - Password management
type Service interface {

	// =====================================================
	// Authentication
	// =====================================================

	// 用户登录
	//
	// tenant_code + email + password
	//
	Login(
		ctx context.Context,
		req LoginRequest,
	) (
		*LoginResponse,
		error,
	)

	// 用户退出登录
	//
	// 第一阶段:
	// refresh token失效
	//
	// 后续:
	// token blacklist
	//
	Logout(
		ctx context.Context,
		req LogoutRequest,
	) error

	// 刷新Token
	//
	RefreshToken(
		ctx context.Context,
		req RefreshTokenRequest,
	) (
		*LoginResponse,
		error,
	)

	// =====================================================
	// Invitation
	// =====================================================

	// 管理员邀请用户
	//
	// Tenant Admin
	//
	// User(invited)
	//
	// Invitation
	//
	InviteUser(
		ctx context.Context,
		tenantID uint,
		operatorID uuid.UUID,
		req InviteUserRequest,
	) error
	// 接受邀请完成注册
	//
	// Token
	// Password
	// Principal
	// Active User
	//
	AcceptInvitation(
		ctx context.Context,
		req AcceptInvitationRequest,
	) error

	GetInvitation(
		ctx context.Context,
		tenantID uint,
		id uuid.UUID,
	) (
		*model.UserInvitation,
		error,
	)

	// =====================================================
	// User Management
	// =====================================================

	CreateUser(
		ctx context.Context,
		req CreateUserRequest,
	) (
		*model.User,
		error,
	)

	GetUser(
		ctx context.Context,
		tenantID uint,
		userID uuid.UUID,
	) (
		*model.User,
		error,
	)

	GetUserByEmail(
		ctx context.Context,
		tenantID uint,
		email string,
	) (
		*model.User,
		error,
	)

	ListUsers(
		ctx context.Context,
		req ListUsersRequest,
	) (
		[]*model.User,
		error,
	)

	UpdateUser(
		ctx context.Context,
		req UpdateUserRequest,
	) error

	DisableUser(
		ctx context.Context,
		req DisableUserRequest,
	) error

	EnableUser(
		ctx context.Context,
		req EnableUserRequest,
	) error

	// =====================================================
	// Password
	// =====================================================

	UpdatePassword(
		ctx context.Context,
		req UpdatePasswordRequest,
	) error

	ResetPassword(
		ctx context.Context,
		req ResetPasswordRequest,
	) error
	// =====================================================
	// Tenant
	// =====================================================

	GetTenant(
		ctx context.Context,
		tenantID uint,
	) (
		*model.Tenant,
		error,
	)

	GetTenantByCode(
		ctx context.Context,
		code string,
	) (
		*model.Tenant,
		error,
	)

	RequestAccess(
		ctx context.Context,
		req RequestAccessRequest,
	) error

	GetUserStats(
		ctx context.Context,
		tenantID uint,
	) (
		*UserStats,
		error,
	)

	ListRoles(
		ctx context.Context,
		tenantID uint,
	) ([]*model.Role, error)
}
