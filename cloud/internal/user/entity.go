package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
// Each user belongs to a single organization.
type User struct {
	ID             uuid.UUID `json:"id"`
	OrgID          uuid.UUID `json:"org_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"` // Never expose this field in JSON responses
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NewUser creates a new User entity.
func NewUser(orgID uuid.UUID, name, email, password string) (*User, error) {
	// Basic validation
	if name == "" {
		return nil, ErrUserNameRequired
	}
	if email == "" {
		return nil, ErrUserEmailRequired
	}
	if password == "" {
		return nil, ErrUserPasswordRequired
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:             uuid.New(),
		OrgID:          orgID,
		Name:           name,
		Email:          email,
		PasswordHash:   string(hashedPassword),
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// CheckPassword verifies if the provided password matches the user's hashed password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}