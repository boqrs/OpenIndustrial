package product

import (
	"time"

	"github.com/google/uuid"
)


type AttributeDefinitionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label,omitempty"`
	Description string    `json:"description,omitempty"`
	DataType    string    `json:"data_type"`
	Unit        string    `json:"unit,omitempty"`
	Required    bool      `json:"required"`
}

type ProductModelResponse struct {
	ID          uuid.UUID                     `json:"id"`
	ResourceID  uuid.UUID                     `json:"resource_id"`
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

type ProductModelListResponse struct {
	Items    []*ProductModelResponse `json:"items"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Total    int64                   `json:"total"`
}