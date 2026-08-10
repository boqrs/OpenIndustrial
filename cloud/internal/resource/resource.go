package resource

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Properties is a flexible JSONB type for resource-specific attributes.
type Properties map[string]interface{}

// Value implements the driver.Valuer interface, allowing this type to be written to the database.
func (p Properties) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface, allowing this type to be read from the database.
func (p *Properties) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &p)
}

// Resource represents a generic entity in the system, like a device, product, or supplier.
type Resource struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	TenantID   uuid.UUID  `db:"tenant_id" json:"tenant_id"`
	Type       string     `db:"type" json:"type"`
	Name       string     `db:"name" json:"name"`
	Properties Properties `db:"properties" json:"properties,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

// ResourceRelation defines a typed link between two resources.
type ResourceRelation struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	TenantID         uuid.UUID  `db:"tenant_id" json:"tenant_id"`
	SourceResourceID uuid.UUID  `db:"source_resource_id" json:"source_resource_id"`
	TargetResourceID uuid.UUID  `db:"target_resource_id" json:"target_resource_id"`
	RelationType     string     `db:"relation_type" json:"relation_type"`
	Properties       Properties `db:"properties" json:"properties,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// ResourceRepository defines the interface for resource-related database operations.
// This is the contract that our persistence layer (e.g., PostgreSQL) must fulfill.
type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *Resource) error
	GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*Resource, error)
	UpdateResource(ctx context.Context, resource *Resource) error
	DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error
	ListResourcesByType(ctx context.Context, tenantID uuid.UUID, resourceType string, offset, limit int) ([]Resource, error)

	CreateRelation(ctx context.Context, relation *ResourceRelation) error
	GetRelationsBySource(ctx context.Context, tenantID, sourceID uuid.UUID) ([]ResourceRelation, error)
	GetRelationsByTarget(ctx context.Context, tenantID, targetID uuid.UUID) ([]ResourceRelation, error)
	DeleteRelation(ctx context.Context, tenantID, relationID uuid.UUID) error
}