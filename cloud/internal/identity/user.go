package identity

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

// UserType defines the type for user roles.
type UserType string

// Defines the possible user types in the system.
const (
	UserTypeAdmin    UserType = "admin"    // 系统管理员 (Boss)
	UserTypeManager  UserType = "manager"  // 生产基地管理员 (厂长)
	UserTypeEmployee UserType = "employee" // 技术人员 (工程师)
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID       `db:"id" json:"id"`
	TenantID  uuid.UUID       `db:"tenant_id" json:"tenant_id"`
	UserType  string          `db:"user_type" json:"user_type"` // DB field remains string for flexibility
	Profile   json.RawMessage `db:"profile" json:"profile"`
	CreatedAt string          `db:"created_at" json:"created_at"`
	UpdatedAt string          `db:"updated_at" json:"updated_at"`
}

// Principal represents a login credential for a user.
type Principal struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	TenantID   uuid.UUID `db:"tenant_id"`
	Provider   string    `db:"provider"`
	Identifier string    `db:"identifier"`
	Credential string    `db:"credential"`
	CreatedAt  string    `db:"created_at"`
	UpdatedAt  string    `db:"updated_at"`
}

// ListUsersRepoParams defines parameters for listing users at the repository level.
type ListUsersRepoParams struct {
	Limit  int
	Offset int
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	CreatePrincipal(ctx context.Context, p *Principal) error
	GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*User, error)
	GetPrincipal(ctx context.Context, tenantID uuid.UUID, provider, identifier string) (*Principal, error)
	// 新增的方法
	ListUsers(ctx context.Context, tenantID uuid.UUID, params ListUsersRepoParams) ([]User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error
}