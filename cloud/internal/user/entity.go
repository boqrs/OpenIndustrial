package user

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/shared"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserStatus defines the status of a user account.
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending" // e.g., waiting for email verification
)

// User represents a user account in the system.
// Each user must belong to an organization.
type User struct {
	shared.BaseEntity
	OrgID        uuid.UUID  `json:"org_id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"` // Never expose password hash
	Email        string     `json:"email"`
	Phone        string     `json:"phone,omitempty"`
	Status       UserStatus `json:"status"`
}

// NewUser creates a new User entity.
func NewUser(orgID uuid.UUID, username, password, email, phone string) (*User, error) {
	// Basic validation
	if orgID == uuid.Nil {
		return nil, ErrUserOrgIDRequired
	}
	if username == "" {
		return nil, ErrUsernameRequired
	}
	if email == "" {
		return nil, ErrUserEmailRequired
	}
	if password == "" {
		return nil, ErrUserPasswordRequired
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		BaseEntity: shared.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrgID:        orgID,
		Username:     username,
		PasswordHash: string(passwordHash),
		Email:        email,
		Phone:        phone,
		Status:       UserStatusActive, // Or UserStatusPending if email verification is needed
	}, nil
}

// CheckPassword compares a plaintext password with the user's hashed password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}