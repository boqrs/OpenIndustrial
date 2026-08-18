package param

import (
	"github.com/google/uuid"
)

// CreateResource is the authoritative parameter structure for creating a resource.
// It is used by the API layer for binding and passed directly to the Service layer.

// CreateResource defines the parameters for creating a new resource.
type CreateResource struct {
	Type         string                 `json:"type" binding:"required,resourcetype"`
	Name         string                 `json:"name" binding:"required,min=2,max=100"`
	Code         *string                `json:"code,omitempty"`
	// Status must be one of the predefined resource statuses.
	Status       string                 `json:"status" binding:"required,oneof=active inactive archived pending PROVISIONED ONBOARDED OFFLINE DECOMMISSIONED"`
	Metadata     []byte                 `json:"metadata,omitempty"`
	ParentID     *uuid.UUID             `json:"parent_id,omitempty"`
	OwnerGroupID *uuid.UUID             `json:"owner_group_id,omitempty"`
	//Attributes   map[string]interface{} `json:"attributes,omitempty"`
	TenantID     uuid.UUID              `json:"-"`
}


// UpdateResource is the authoritative parameter structure for updating a resource.
type UpdateResource struct {
	// Fields from the request body
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Code     *string `json:"code,omitempty"`
	Status   string `json:"status" binding:"required,oneof=active inactive archived pending"`
	Metadata []byte `json:"metadata,omitempty"`
	Version  int    `json:"version" binding:"required,gt=0"`
	ParentID uuid.UUID `json:"parent_id,omitempty"`

	// Populated by the handler from the URL and context.
	TenantID   uuid.UUID `json:"-"`
	ResourceID uuid.UUID `json:"-"`
}

// CreateProduct is the authoritative parameter structure for creating a product.
// This definition is based on the authoritative version you provided.
type CreateProduct struct {
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	Type         string    `json:"type" binding:"required"`
	SerialNumber string    `json:"serial_number"`
	OwnerGroupID uuid.UUID `json:"owner_group_id" binding:"required"`

	// Populated by the handler from the context.
	TenantID uuid.UUID `json:"-"`
}

type CreateRelation struct {
	// FromResourceID is the UUID of the source resource in the relation.
	FromResourceID uuid.UUID `json:"from_resource_id" binding:"required"`

	// ToResourceID is the UUID of the target resource in the relation.
	ToResourceID uuid.UUID `json:"to_resource_id" binding:"required"`

	// RelationType describes the nature of the relationship (e.g., "IS_INSTANCE_OF").
	RelationType string `json:"relation_type" binding:"required"`
}
