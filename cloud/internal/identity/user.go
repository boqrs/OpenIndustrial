package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserType defines the type of a user.
type UserType string

const (
	UserTypeAdmin  UserType = "admin"
	UserTypeMember UserType = "member"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID       `db:"id"`
	TenantID  uuid.UUID       `db:"tenant_id"`
	UserType  string          `db:"user_type"`
	Profile   json.RawMessage `db:"profile"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

// Principal represents a way to authenticate a user.
type Principal struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	TenantID   uuid.UUID `db:"tenant_id"`
	Provider   string    `db:"provider"`
	Identifier string    `db:"identifier"`
	Credential string    `db:"credential"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type ListUsersRepoParams struct {
	Limit  int
	Offset int
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*User, error)
	CreatePrincipal(ctx context.Context, principal *Principal) error
	GetPrincipal(ctx context.Context, tenantID uuid.UUID, provider, identifier string) (*Principal, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID, params ListUsersRepoParams) ([]*User, error)
}