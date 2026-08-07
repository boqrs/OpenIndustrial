package resource

import "time"

// Resource represents a unique, addressable entity within the system, designed for a multi-tenant environment.
type Resource struct {
	ID          string            `json:"id"`
	NamespaceID string            `json:"namespace_id"`
	Type        ResourceType      `json:"type"`
	Name        string            `json:"name"`
	ParentID    *string           `json:"parent_id,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}