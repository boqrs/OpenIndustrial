package bom

type CreateRequest struct {
	ProductID uint `json:"product_id"`

	BOMNo string `json:"bom_no"`

	Version int `json:"version"`

	Description string `json:"description"`

	Items []ItemRequest `json:"items"`
}

type UpdateRequest struct {
	Description string `json:"description"`

	Items []ItemRequest `json:"items"`
}

type ItemRequest struct {
	MaterialID uint `json:"material_id"`

	Quantity float64 `json:"quantity"`

	Unit string `json:"unit"`

	Sequence int `json:"sequence"`

	OperationCode string `json:"operation_code"`

	IsOptional bool `json:"is_optional"`

	Description string `json:"description"`
}