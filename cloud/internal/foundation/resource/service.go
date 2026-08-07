package resource

import "context"

// Service defines the business logic interface for managing resources and the graph.
type Service interface {
	CreateResource(ctx context.Context, namespaceID string, parentID *string, name string, resType ResourceType, metadata map[string]string) (*Resource, error)
	GetResource(ctx context.Context, namespaceID, id string) (*Resource, error)
	UpdateResourceMetadata(ctx context.Context, namespaceID, id string, metadata map[string]string) (*Resource, error)
	DeleteResource(ctx context.Context, namespaceID, id string) error

	AddRelation(ctx context.Context, namespaceID, sourceID, targetID string, relType RelationType, properties map[string]string) (*Relation, error)
	RemoveRelation(ctx context.Context, namespaceID, relationID string) error

	GetChildren(ctx context.Context, namespaceID, parentID string) ([]*Resource, error)
	GetRelations(ctx context.Context, namespaceID, resourceID string) ([]*Relation, error)
}