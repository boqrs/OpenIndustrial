package resource

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Resource represents a fundamental entity in the industrial ecosystem.
type Resource struct {
	ID           uuid.UUID  `db:"id"`
	TenantID     uuid.UUID  `db:"tenant_id"`
	Type         string     `db:"type"`
	Name         string     `db:"name"`
	Code         *string    `db:"code"`
	Status       string     `db:"status"`
	Metadata     []byte     `db:"metadata"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
	ParentID     *uuid.UUID `db:"parent_id"`
	OwnerGroupID *uuid.UUID `db:"owner_group_id"`
}

// ResourceRelation defines the relationship between two resources.
type ResourceRelation struct {
	FromID       uuid.UUID `db:"from_id"`
	ToID         uuid.UUID `db:"to_id"`
	RelationType string    `db:"relation_type"`
	CreatedAt    time.Time `db:"created_at"`
}

// ResourceRepository defines the persistence interface for Resource entities.
type ResourceRepository interface {
	CreateResource(ctx context.Context, resource *Resource) error
	GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*Resource, error)
	UpdateResource(ctx context.Context, resource *Resource) error
	DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error
	ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*Resource, error)
	CheckUserInSameGroupAsResource(ctx context.Context, userID, resourceID uuid.UUID) (bool, error)
	CreateResourceRelation(ctx context.Context, relation *ResourceRelation) error
	ListResourceRelations(ctx context.Context, resourceID uuid.UUID) ([]*ResourceRelation, error)
}

// AttributeDefinitionRepository defines the persistence interface for AttributeDefinition entities.
type AttributeDefinitionRepository interface {
	CreateAttributeDefinition(ctx context.Context, def *AttributeDefinition) error
	GetAttributeDefinitionByID(ctx context.Context, tenantID, defID uuid.UUID) (*AttributeDefinition, error)
	GetAttributeDefinitionByKey(ctx context.Context, tenantID uuid.UUID, key string) (*AttributeDefinition, error)
	ListAttributeDefinitions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*AttributeDefinition, error)
	UpdateAttributeDefinition(ctx context.Context, def *AttributeDefinition) error
	DeleteAttributeDefinition(ctx context.Context, tenantID, defID uuid.UUID) error
}

// ResourceAttributeRepository defines the persistence interface for ResourceAttribute values.
type ResourceAttributeRepository interface {
	SetAttribute(ctx context.Context, attr *ResourceAttribute) error
	SetAttributes(ctx context.Context, attrs []*ResourceAttribute) error
	GetAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) (*ResourceAttribute, error)
	GetAttributesByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*ResourceAttribute, error)
	DeleteAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) error
}