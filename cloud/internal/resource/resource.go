package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Resource represents a digital twin of a physical or logical asset in the system.
type Resource struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TenantID     uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`   // ADDED: This field was missing
	Type         string    `db:"type" json:"type"`
	SerialNumber string    `db:"serial_number" json:"serial_number"` // ADDED: This field was missing
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// ResourceRelation defines a typed relationship between two resources.
type ResourceRelation struct {
	ID               uuid.UUID `db:"id" json:"id"`
	TenantID         uuid.UUID `db:"tenant_id" json:"tenant_id"`
	SourceResourceID uuid.UUID `db:"source_resource_id" json:"source_resource_id"`
	TargetResourceID uuid.UUID `db:"target_resource_id" json:"target_resource_id"`
	RelationType     string    `db:"relation_type" json:"relation_type"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// ResourceRepository defines the persistence interface for Resources and their relations.
type ResourceRepository interface {
	CreateResource(ctx context.Context, res *Resource) error
	CreateResourceRelation(ctx context.Context, rel *ResourceRelation) error
}