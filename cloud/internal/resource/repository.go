package resource

import "context"

// Repository defines the interface for resource and relation storage.
type Repository interface {
	CreateResource(ctx context.Context, resource *Resource) (*Resource, error)
	GetResource(ctx context.Context, id string) (*Resource, error)
	UpdateResource(ctx context.Context, resource *Resource) (*Resource, error)
	DeleteResource(ctx context.Context, id string) error
	CreateRelation(ctx context.Context, relation *Relation) (*Relation, error)
	ListRelations(ctx context.Context, fromID, toID, relType string) ([]*Relation, error)
}