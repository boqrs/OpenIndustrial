package param

import (
	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)
/****************************************************************************/
type AttributeDefinitionRequest struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	DataType    string `json:"data_type"`
	Unit        string `json:"unit,omitempty"`
	Required    bool   `json:"required"`
}

type CreateDeviceTypeRequest struct {
	Name        string                              `json:"name"`
	Code        string                              `json:"code"`
	Category    model.DeviceTypeCategory            `json:"category"`
	Description string                              `json:"description,omitempty"`
	Attributes  map[string]AttributeDefinitionRequest `json:"attributes,omitempty"`
}

type UpdateDeviceTypeRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

// --- Device DTOs ---

type CreateDeviceRequest struct {
	DeviceTypeID     uuid.UUID              `json:"device_type_id"`
	Name             string                 `json:"name"`
	Code             string                 `json:"code,omitempty"`
	ParentResourceID *uuid.UUID             `json:"parent_resource_id,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
}

type UpdateDeviceRequest struct {
	Name             *string                `json:"name,omitempty"`
	Code             *string                `json:"code,omitempty"`
	ParentResourceID *uuid.UUID             `json:"parent_resource_id,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
}



// --- Topology DTOs ---

type AttachDeviceRequest struct {
	ParentResourceID uuid.UUID `json:"parent_resource_id"`
}

type ConnectDeviceRequest struct {
	TargetResourceID uuid.UUID              `json:"target_resource_id"`
	ConnectionType   string                 `json:"connection_type"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}