package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserOrgIDRequired is returned when creating a user without an organization ID.
	ErrUserOrgIDRequired = errors.New("user must be associated with an organization")
)

// User represents a user account in the system.
type User struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User instance.
func NewUser(orgID uuid.UUID, username, email, firstName, lastName string) (*User, error) {
	if orgID == uuid.Nil {
		return nil, ErrUserOrgIDRequired
	}
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New(),
		OrgID:     orgID,
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true, // Users are active by default
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetPassword hashes the user's password and sets it.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the user's hashed password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}