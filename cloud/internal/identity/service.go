package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service provides business logic for identity management.
type Service struct {
	repo Repository
}

// NewService creates a new identity Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateUser creates a new user.
func (s *Service) CreateUser(ctx context.Context, user *User) error {
	user.ID = "user_" + uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Status = StatusActive

	return s.repo.CreateUser(ctx, user)
}

// JoinOrganization creates a membership linking a user to an organization.
func (s *Service) JoinOrganization(ctx context.Context, userID string, orgID string) error {
	member := &Membership{
		ID:        "member_" + uuid.NewString(),
		UserID:    userID,
		OrgID:     orgID,
		Status:    StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateMembership(ctx, member)
}