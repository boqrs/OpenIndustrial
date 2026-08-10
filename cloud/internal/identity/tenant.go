package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Tenant represents a customer organization in the system.
// It is the highest level of data isolation.
type Tenant struct {
	ID        uuid.UUID
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TenantRepository defines the interface for database operations on tenants.
type TenantRepository interface {
	// CreateTenant creates a new tenant in the database.
	CreateTenant(ctx context.Context, tenant *Tenant) error

	// GetTenantByID retrieves a tenant by its unique ID.
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)

	// UpdateTenant updates an existing tenant's information.
	UpdateTenant(ctx context.Context, tenant *Tenant) error
}