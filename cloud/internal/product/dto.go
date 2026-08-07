package product

import "time"

// CreateProductRequest defines the request body for creating a new product.
type CreateProductRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Spec        map[string]string `json:"spec"`
}

// ProductResponse defines the standard response for a product.
type ProductResponse struct {
	ID          string            `json:"id"`
	OrgID       string            `json:"org_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Spec        map[string]string `json:"spec"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ToProductResponse converts a Product entity to a ProductResponse.
func ToProductResponse(product *Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		OrgID:       product.OrgID,
		Name:        product.Name,
		Description: product.Description,
		Spec:        product.Spec,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToProductListResponse converts a slice of Product entities to a slice of ProductResponse.
func ToProductListResponse(products []*Product) []*ProductResponse {
	res := make([]*ProductResponse, len(products))
	for i, p := range products {
		res[i] = ToProductResponse(p)
	}
	return res
}