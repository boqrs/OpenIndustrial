package resource

import "context"

// Repository defines the persistence interface for resources and their relations.
type Repository interface {
	CreateResource(ctx context.Context, resource *Resource) error
	GetResourceByID(ctx context.Context, namespaceID, id string) (*Resource, error)
	UpdateResource(ctx context.Context, resource *Resource) error
	DeleteResource(ctx context.Context, namespaceID, id string) error
	ListResourcesByParentID(ctx context.Context, namespaceID, parentID string) ([]*Resource, error)

	CreateRelation(ctx context.Context, relation *Relation) error
	DeleteRelation(ctx context.Context, namespaceID, id string) error
	FindRelations(ctx context.Context, namespaceID, resourceID string) ([]*Relation, error)
}