package resource

import "time"

// Relation defines a directional link with properties between two resources in a multi-tenant environment.
type Relation struct {
	ID          string            `json:"id"`
	NamespaceID string            `json:"namespace_id"`
	SourceID    string            `json:"source_id"`
	TargetID    string            `json:"target_id"`
	Type        RelationType      `json:"type"`
	Properties  map[string]string `json:"properties,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}