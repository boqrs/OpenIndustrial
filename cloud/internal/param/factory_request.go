package param

import (
	"github.com/google/uuid"
)

type CreateFactoryRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Address  string `json:"address,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

type UpdateFactoryRequest struct {
	Name     *string `json:"name,omitempty"`
	Code     *string `json:"code,omitempty"`
	Address  *string `json:"address,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

type CreateTopologyNodeRequest struct {
	FactoryID uuid.UUID `json:"factory_id"`

	ParentResourceID *uuid.UUID `json:"parent_resource_id,omitempty"`

	Type string `json:"type"`

	Name string `json:"name"`

	Code string `json:"code,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateTopologyNodeRequest struct {
	Name     *string                `json:"name,omitempty"`
	Code     *string                `json:"code,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MoveTopologyNodeRequest struct {
	ResourceID       uuid.UUID `json:"resource_id"`
	ParentResourceID *uuid.UUID `json:"parent_resource_id,omitempty"`
}
