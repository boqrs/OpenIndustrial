package param

import (
	"time"
	"github.com/google/uuid"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)


type DeviceTypeResponse struct {
	ID          uuid.UUID                `json:"id"`
	ResourceID  uuid.UUID                `json:"resource_id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Category    model.DeviceTypeCategory `json:"category"`
	Description string                   `json:"description,omitempty"`
	Enabled     bool                     `json:"enabled"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type DeviceResponse struct {
	ID               uuid.UUID              `json:"id"`
	ResourceID       uuid.UUID              `json:"resource_id"`
	DeviceType       DeviceTypeResponse     `json:"device_type"`
	Name             string                 `json:"name"`
	Code             string                 `json:"code,omitempty"`
	Status           string                 `json:"status"`
	ParentResourceID *uuid.UUID             `json:"parent_resource_id,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}


type DeviceConnectionResponse struct {
	ID               uint                   `json:"id"`
	SourceResourceID uuid.UUID              `json:"source_resource_id"`
	TargetResourceID uuid.UUID              `json:"target_resource_id"`
	ConnectionType   string                 `json:"connection_type"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type DeviceTopologyResponse struct {
	Device      DeviceResponse             `json:"device"`
	Children    []DeviceResponse           `json:"children"`
	Connections []DeviceConnectionResponse `json:"connections"`
}