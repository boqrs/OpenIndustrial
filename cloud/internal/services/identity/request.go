package identity

import (
	"github.com/google/uuid"
)

// =====================================================
// Authentication
// =====================================================

// LoginRequest
//
// 用户登录
//
// tenant_code + email + password
type LoginRequest struct {

	// 工厂编码
	TenantCode string `json:"tenant_code" validate:"required"`

	// 用户邮箱
	Email string `json:"email" validate:"required,email"`

	// 密码
	Password string `json:"password" validate:"required"`
}

// LogoutRequest
//
// 用户退出登录
type LogoutRequest struct {

	// 当前用户ID
	UserID uuid.UUID `json:"user_id" validate:"required"`

	// 当前Tenant
	TenantID uint `json:"tenant_id" validate:"required"`

	// Refresh Token
	//
	// 用于注销refresh token
	//
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRequest
//
// 使用refresh token获取新的token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// =====================================================
// User Invitation
// =====================================================

// InviteUserRequest
//
// 管理员邀请用户
type InviteUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	// 创建用户姓名
	Name string `json:"name" validate:"required"`

	// 用户邮箱
	Email string `json:"email" validate:"required,email"`

	// 创建人
	//
	// 当前管理员用户
	//
	CreatedBy uint `json:"created_by" validate:"required"`
}

// AcceptInvitationRequest
//
// 用户通过邮件完成注册
type AcceptInvitationRequest struct {
	Token string `json:"token" validate:"required"`

	Password string `json:"password" validate:"required,min=8"`

	Name string `json:"name"`
}

// =====================================================
// User Management
// =====================================================

// CreateUserRequest
//
// 内部创建用户
//
// 一般由Invitation流程调用
type CreateUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UUID uuid.UUID `json:"uuid"`

	Name string `json:"name" validate:"required"`

	Email string `json:"email" validate:"required,email"`

	UserType string `json:"user_type"`

	Status string `json:"status"`
}

// GetUserRequest
//
// 查询单个用户
type GetUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// ListUsersRequest
//
// 用户列表
type ListUsersRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	// 分页
	//
	Limit int `json:"limit"`

	Offset int `json:"offset"`

	// 可选过滤

	Status string `json:"status"`

	UserType string `json:"user_type"`

	Keyword string `json:"keyword"`
}

// UpdateUserRequest
//
// 修改用户基础信息
type UpdateUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UserID uuid.UUID `json:"user_id" validate:"required"`

	Name string `json:"name"`

	Email string `json:"email"`
}

// DisableUserRequest
//
// 禁用用户
type DisableUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// EnableUserRequest
//
// 激活用户
type EnableUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// =====================================================
// Password
// =====================================================

// UpdatePasswordRequest
//
// 用户自己修改密码
type UpdatePasswordRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UserID uuid.UUID `json:"user_id" validate:"required"`

	OldPassword string `json:"old_password" validate:"required"`

	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResetPasswordRequest
//
// 管理员重置用户密码
type ResetPasswordRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UserID uuid.UUID `json:"user_id" validate:"required"`

	NewPassword string `json:"new_password" validate:"required,min=8"`

	OperatorID uint `json:"operator_id" validate:"required"`
}
