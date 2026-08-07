package resource

import "time"

// Resource represents a node in the industrial resource hierarchy.
type Resource struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        Type                   `json:"type"`
	ParentID    string                 `json:"parent_id,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Relation represents a directed link between two resources.
type Relation struct {
	ID         string                 `json:"id"`
	FromID     string                 `json:"from_id"` // The source resource of the relation
	ToID       string                 `json:"to_id"`   // The target resource of the relation
	Type       RelationType           `json:"type"`    // The type of the relation
	Direction  RelationDirection      `json:"direction"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}