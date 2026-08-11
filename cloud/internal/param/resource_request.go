package param

import (
	"github.com/google/uuid"
)

// CreateResource is the authoritative parameter structure for creating a resource.
// It is used by the API layer for binding and passed directly to the Service layer.
type CreateResource struct {
	// Fields from the request body
	Type         string     `json:"type" binding:"required,resourcetype"`
	Name         string     `json:"name" binding:"required,min=2,max=100"`
	Code         *string    `json:"code,omitempty"`
	Status       string     `json:"status" binding:"required,oneof=active inactive archived pending"`
	Metadata     []byte     `json:"metadata,omitempty"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	OwnerGroupID *uuid.UUID `json:"owner_group_id,omitempty"`

	// TenantID is populated by the handler from the request context, not the body.
	TenantID uuid.UUID `json:"-"`
}

// UpdateResource is the authoritative parameter structure for updating a resource.
type UpdateResource struct {
	// Fields from the request body
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Code     *string `json:"code,omitempty"`
	Status   string `json:"status" binding:"required,oneof=active inactive archived pending"`
	Metadata []byte `json:"metadata,omitempty"`
	Version  int    `json:"version" binding:"required,gt=0"`

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
