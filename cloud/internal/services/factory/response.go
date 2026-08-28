package factory

import (
	"time"
)

type FactoryResponse struct {
	ID uint `json:"id"`

	ResourceID uint `json:"resource_id"`

	Name string `json:"name"`

	Code string `json:"code"`

	Address string `json:"address,omitempty"`

	Timezone string `json:"timezone"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

type TopologyNodeResponse struct {
	ResourceID uint `json:"resource_id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Status string `json:"status"`
	ParentResourceID *uint `json:"parent_resource_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type FactoryTopologyResponse struct {
	Factory FactoryResponse `json:"factory"`

	Nodes []TopologyNodeResponse `json:"nodes"`
}