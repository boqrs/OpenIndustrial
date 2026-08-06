package user

import (
	"context"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the user domain.
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
func (s *Service) CreateUser(ctx context.Context, orgID uuid.UUID, name, email, password string) (*User, error) {
	// In a real app, you might check if the email is already taken.
	user, err := NewUser(orgID, name, email, password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// AuthenticateUser verifies a user's credentials and returns the user if successful.
func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}