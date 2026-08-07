package permission

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for storing permissions and policies.
type Repository interface {
	// Permission methods
	GetPermission(ctx context.Context, id uuid.UUID) (*Permission, error)
	ListAllPermissions(ctx context.Context) ([]*Permission, error)

	// Policy methods
	AddPolicy(ctx context.Context, policy *Policy) error
	RemovePolicy(ctx context.Context, policy *Policy) error
	GetPoliciesForSubject(ctx context.Context, subject string) ([]*Policy, error)
}