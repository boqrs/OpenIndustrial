package resource

import (
	"context"
	"fmt"
)

// GraphService provides methods for interacting with the resource graph.
type GraphService struct {
	repo Repository
}

// NewGraphService creates a new GraphService.
func NewGraphService(repo Repository) *GraphService {
	return &GraphService{repo: repo}
}

// AddRelation adds a new relation between two resources.
func (s *GraphService) AddRelation(ctx context.Context, fromID, toID string, relType RelationType) (*Relation, error) {
	// In a real implementation, you might want to check if the resources exist.
	relation := &Relation{
		FromID: fromID,
		ToID:   toID,
		Type:   relType,
	}
	return s.repo.CreateRelation(ctx, relation)
}

// GetChildren returns all direct children of a given resource.
func (s *GraphService) GetChildren(ctx context.Context, resourceID string) ([]*Resource, error) {
	relations, err := s.repo.ListRelations(ctx, resourceID, "", "")
	if err != nil {
		return nil, err
	}

	var children []*Resource
	for _, relation := range relations {
		// Assuming "contains" is the relation type for parent-child.
		if relation.Type == Contains && relation.FromID == resourceID {
			child, err := s.repo.GetResource(ctx, relation.ToID)
			if err != nil {
				// Handle error: log it, skip it, etc.
				continue
			}
			children = append(children, child)
		}
	}
	return children, nil
}

// FindPaths finds all paths between two resources.
// This is a simplified implementation and might not be efficient for large graphs.
func (s *GraphService) FindPaths(ctx context.Context, startID, endID string) ([][]*Resource, error) {
	// This requires a more complex graph traversal algorithm (like DFS or BFS)
	// which is beyond the scope of this basic example.
	return nil, fmt.Errorf("path finding not implemented")
}