package bom

import (
	"time"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Response struct {
	ID uint `json:"id"`

	TenantID uuid.UUID `json:"tenant_id"`

	ProductID uuid.UUID `json:"product_id"`

	BOMNo string `json:"bom_no"`

	Version int `json:"version"`

	Status model.BOMStatus `json:"status"`

	Description string `json:"description"`

	Items []ItemResponse `json:"items"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

type ItemResponse struct {
	ID uint `json:"id"`

	MaterialID uint `json:"material_id"`

	Quantity float64 `json:"quantity"`

	Unit string `json:"unit"`

	Sequence int `json:"sequence"`

	OperationCode string `json:"operation_code"`

	IsOptional bool `json:"is_optional"`

	Description string `json:"description"`
}