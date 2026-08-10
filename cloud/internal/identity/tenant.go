package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system.
type Tenant struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// TenantRepository defines the interface for tenant persistence.
type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *Tenant) error
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}