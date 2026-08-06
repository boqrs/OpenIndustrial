package product

import (
	"time"
)

// CreateProductRequest defines the structure for a request to create a new product.
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Model       string `json:"model" binding:"required"`
}

// ProductResponse defines the structure for a response containing product details.
type ProductResponse struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Model       string    `json:"model"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToProductResponse converts a Product entity to a ProductResponse DTO.
func ToProductResponse(product *Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID.String(),
		OrgID:       product.OrgID.String(),
		Name:        product.Name,
		Description: product.Description,
		Model:       product.Model,
		CreatedAt:   product.CreatedAt,
	}
}

// ToProductListResponse converts a slice of Product entities to a slice of ProductResponse DTOs.
func ToProductListResponse(products []*Product) []*ProductResponse {
	res := make([]*ProductResponse, len(products))
	for i, product := range products {
		res[i] = ToProductResponse(product)
	}
	return res
}