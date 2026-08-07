package resource

// Resolver provides an interface to query the resource graph.
// It is used by the authorization engine to determine resource hierarchy and relationships.
type Resolver interface {
	GetResource(id string) (*Resource, error)
	// IsParent checks if a resource is a direct or indirect parent of another resource.
	IsParent(parentID string, childID string) (bool, error)
}