package product

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
)


type AttributeResponse struct {
	ID          uint `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label,omitempty"`
	Description string    `json:"description,omitempty"`
	DataType    string    `json:"data_type"`
	Unit        string    `json:"unit,omitempty"`
	Required    bool      `json:"required"`
}


type AttributeDefinitionResponse struct {
	ID          uint `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label,omitempty"`
	Description string    `json:"description,omitempty"`
	DataType    string    `json:"data_type"`
	Unit        string    `json:"unit,omitempty"`
	Required    bool      `json:"required"`
}

type ProductModelResponse struct {
	ID          uint                     `json:"id"`
	ResourceID  uint                     `json:"resource_id"`
	Name        string                        `json:"name"`
	Code        string                        `json:"code"`
	Version     string                        `json:"version"`
	Category    string                        `json:"category"`
	Description string                        `json:"description,omitempty"`
	Status      string                        `json:"status"`
	Attributes  []AttributeResponse `json:"attributes,omitempty"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
}

type ProductDetailResponse struct {
	ID          uint                     `json:"id"`
	ResourceID  uint                     `json:"resource_id"`
	Name        string                        `json:"name"`
	Code        string                        `json:"code"`
	Version     string                        `json:"version"`
	Category    string                        `json:"category"`
	Description string                        `json:"description,omitempty"`
	Status      string                        `json:"status"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
	Attribute []*model.ResourceAttribute	`json:"attribute,omitempty"`
	AttributeDefinition []*model.AttributeDefinition `json:"attribute_definition,omitempty"`
}

type CreateProductModelResponse struct {
	ID          uint                     `json:"id"`
	ResourceID  uint                     `json:"resource_id"`
	Name        string                        `json:"name"`
	Code        string                        `json:"code"`
	Version     string                        `json:"version"`
	Category    string                        `json:"category"`
	Description string                        `json:"description,omitempty"`
	Status      string                        `json:"status"`
	Attributes  []AttributeDefinitionResponse `json:"attributes,omitempty"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
}

type UpdateProductModelResponse struct {
	ID          uint                     `json:"id"`
	ResourceID  uint                     `json:"resource_id"`
	Name        string                        `json:"name"`
	Code        string                        `json:"code"`
	Version     string                        `json:"version"`
	Category    string                        `json:"category"`
	Description string                        `json:"description,omitempty"`
	Status      string                        `json:"status"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
}



type ProductModelListResponse struct {
	Items    []*ProductModelResponse `json:"items"`
	pkg.PageBaseResp
}