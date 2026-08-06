package product

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product definition or blueprint (SKU).
// It is a template for creating physical assets (Product Instances or SNs).
type Product struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`      // Belongs to which organization
	Name        string    `json:"name"`       // e.g., "Ender-3 S1 Pro"
	Description string    `json:"description"`
	Model       string    `json:"model"`       // e.g., "CR-10 SE"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewProduct creates a new Product entity.
func NewProduct(orgID uuid.UUID, name, description, model string) (*Product, error) {
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if model == "" {
		return nil, ErrProductModelRequired
	}

	now := time.Now().UTC()
	return &Product{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        name,
		Description: description,
		Model:       model,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}