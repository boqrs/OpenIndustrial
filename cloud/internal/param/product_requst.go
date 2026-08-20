package param

import (
	//"time"

	//"github.com/google/uuid"
)

type AttributeDefinitionRequest struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	DataType    string `json:"data_type"`
	Unit        string `json:"unit,omitempty"`
	Required    bool   `json:"required"`
}

type CreateProductModelRequest struct {
	Name        string                              `json:"name"`
	Code        string                              `json:"code"`
	Version     string                              `json:"version"`
	Category    string                              `json:"category"`
	Description string                              `json:"description,omitempty"`
	Attributes  map[string]AttributeDefinitionRequest `json:"attributes,omitempty"`
}

type UpdateProductModelRequest struct {
	Name        *string `json:"name,omitempty"`
	Category    *string `json:"category,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateAttributeDefinitionsRequest struct {
	Attributes map[string]AttributeDefinitionRequest `json:"attributes"`
}

type ListProductModelsRequest struct {
	Category string `json:"category,omitempty"`
	Status   string `json:"status,omitempty"`
	Code     string `json:"code,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}
