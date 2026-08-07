package user

import (
	"context"

	"github.com/google/uuid"
)

// Service provides business logic for user management.
type Service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateUser handles the business logic of creating a new user.
func (s *Service) CreateUser(ctx context.Context, orgID uuid.UUID, username, email, password, firstName, lastName string) (*User, error) {
	// Here you might want to check if a user with the same email or username already exists.
	// existingUser, err := s.repo.FindByEmail(ctx, email)
	// if err == nil && existingUser != nil {
	// 	 return nil, ErrUserAlreadyExists
	// }

	user, err := NewUser(orgID, username, email, firstName, lastName)
	if err != nil {
		return nil, err
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate verifies a user's credentials and returns the user if successful.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		// If the error is 'not found', return a generic invalid credentials error.
		// This prevents attackers from enumerating existing user emails.
		return nil, ErrInvalidCredentials
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}