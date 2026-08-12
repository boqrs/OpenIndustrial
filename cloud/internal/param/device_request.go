package param

import (

	"github.com/google/uuid"
)

// AttributeDefinitionParam defines the structure for creating a new attribute definition
// as part of a product model.
type AttributeDefinitionParam struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	DataType   string `json:"value_type" binding:"required,oneof=string,int,float,bool"`
	Unit        string `json:"unit"`
	Label string `json:"label"`
}

// CreateProductModelRequest defines the parameters for creating a new product model,
// including its attribute definitions.
type CreateProductModelRequest struct {
	Name       string                              `json:"name" binding:"required,min=2,max=100"`
	Attributes map[string]AttributeDefinitionParam `json:"attributes"` // Key is the attribute key
	TenantID   uuid.UUID                           `json:"-"`
}

// RegisterDeviceRequest defines the parameters for registering a new device instance
// (for both IoT products and factory assets).
type RegisterDeviceRequest struct {
	ProductModelID uuid.UUID              `json:"product_model_id" binding:"required"`
	InstanceName   string                 `json:"instance_name" binding:"required,min=2,max=100"`
	Attributes     map[string]interface{} `json:"attributes"`
	TenantID       uuid.UUID              `json:"-"`
	SerialNumber string `json:"serial_number"`
}