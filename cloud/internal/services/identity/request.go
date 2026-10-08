package identity

import "github.com/google/uuid"

// =====================================================
// Authentication
// =====================================================

type LoginRequest struct {
	TenantCode string `json:"tenant_code" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
}

type LogoutRequest struct {
	UserID       uuid.UUID `json:"user_id" validate:"required"`
	TenantID     uint      `json:"tenant_id" validate:"required"`
	RefreshToken string    `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// =====================================================
// Access Request
// =====================================================

// RequestAccessRequest is submitted by a user who wants to join a tenant.
//
// It does NOT create a User.
// It only notifies the tenant administrator.
type RequestAccessRequest struct {
	TenantCode string `json:"tenant_code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
}

// =====================================================
// Invitation
// =====================================================

type InviteUserRequest struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	RoleID uint   `json:"role_id" validate:"required"`
}

type AcceptInvitationRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name"`
}

// =====================================================
// User Management
// =====================================================

type CreateUserRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	UUID uuid.UUID `json:"uuid"`

	Name string `json:"name" validate:"required"`

	Email string `json:"email" validate:"required,email"`

	RoleID uint `json:"role_id" validate:"required"`

	// UserType and Status are intentionally ignored by the service.
	// Normal user creation always creates an employee in invited state.
	UserType string `json:"user_type"`

	Status string `json:"status"`
}

type GetUserRequest struct {
	TenantID uint      `json:"tenant_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}

type ListUsersRequest struct {
	TenantID uint `json:"tenant_id" validate:"required"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	Status   string `json:"status"`
	UserType string `json:"user_type"`
	Keyword  string `json:"keyword"`
}

type UpdateUserRequest struct {
	TenantID uint      `json:"tenant_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`

	Name string `json:"name"`
}

type DisableUserRequest struct {
	TenantID uint      `json:"tenant_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}

type EnableUserRequest struct {
	TenantID uint      `json:"tenant_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}

// =====================================================
// Password
// =====================================================

type UpdatePasswordRequest struct {
	TenantID uint      `json:"tenant_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`

	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	TenantID uint      `json:"tenant_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`

	NewPassword string `json:"new_password" validate:"required,min=8"`

	// OperatorID is populated from JWT by the handler.
	OperatorID uuid.UUID `json:"-"`
}
