package product

import (
	"time"
)

// CreateProductRequest defines the request for creating a product.
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Spec        string `json:"spec"`
	Description string `json:"description"`
}

// ProductResponse is the DTO for a product.
type ProductResponse struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Spec        string    `json:"spec"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToProductResponse converts a Product entity to a DTO.
func ToProductResponse(product *Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID.String(),
		OrgID:       product.OrgID.String(),
		Name:        product.Name,
		Code:        product.Code,
		Spec:        product.Spec,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToProductListResponse converts a slice of Product entities to a slice of DTOs.
func ToProductListResponse(products []*Product) []*ProductResponse {
	responses := make([]*ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToProductResponse(p)
	}
	return responses
}