package identity

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// =====================================================
// Repository
// =====================================================

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

// UserRepository
//
// Identity domain persistence interface.
//
// 包含:
// Tenant
// User
// Invitation
// Principal
//
// Role / Permission 后续独立实现
type UserRepository interface {

	// =====================================================
	// Tenant
	// =====================================================

	// 根据Tenant Code查询工厂
	//
	// 登录入口使用
	//
	GetTenantByCode(
		ctx context.Context,
		code string,
	) (
		*model.Tenant,
		error,
	)

	// 根据Tenant ID查询工厂
	//
	GetTenantByID(
		ctx context.Context,
		id uint,
	) (
		*model.Tenant,
		error,
	)

	// =====================================================
	// User
	// =====================================================

	CreateUser(
		ctx context.Context,
		user *model.User,
	) error

	// Tenant隔离查询用户
	//
	// TenantID + User UUID
	//
	GetUserByID(
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

	UpdateUser(
		ctx context.Context,
		user *model.User,
	) error

	// 用户列表
	//
	// Tenant隔离
	//
	ListUsers(
		ctx context.Context,
		tenantID uint,
		limit int,
		offset int,
		status string,
		userType string,
		keyword string,
	) (
		[]*model.User,
		error,
	)

	// =====================================================
	// Invitation
	// =====================================================

	CreateInvitation(
		ctx context.Context,
		invitation *model.UserInvitation,
	) error

	GetInvitationByToken(
		ctx context.Context,
		tokenHash string,
	) (
		*model.UserInvitation,
		error,
	)

	GetInvitationByID(
		ctx context.Context,
		tenantID uint,
		id uuid.UUID,
	) (
		*model.UserInvitation,
		error,
	)

	UpdateInvitation(
		ctx context.Context,
		invitation *model.UserInvitation,
	) error

	// =====================================================
	// Principal
	// =====================================================

	CreatePrincipal(
		ctx context.Context,
		principal *model.Principal,
	) error

	GetPrincipal(
		ctx context.Context,
		tenantID uint,
		provider string,
		identifier string,
	) (
		*model.Principal,
		error,
	)

	UpdatePrincipal(
		ctx context.Context,
		principal *model.Principal,
	) error
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
}
