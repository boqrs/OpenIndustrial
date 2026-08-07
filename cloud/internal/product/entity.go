package product

import (
	"time"

	"github.com/google/uuid"
)

// Product defines the core entity for a product model.
type Product struct {
	ID          string            `json:"id"`
	OrgID       string            `json:"org_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Spec        map[string]string `json:"spec"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// LifecycleEvent represents a significant event in a product instance's life.
type LifecycleEvent struct {
	ProductInstanceID string
	Type              string // e.g., "MANUFACTURED", "SOLD", "ACTIVATED", "DECOMMISSIONED"
	Timestamp         time.Time
	Metadata          map[string]interface{}
}

// NewProduct creates a new Product entity.
func NewProduct(orgID, name, description string, spec map[string]string) (*Product, error) {
	now := time.Now().UTC()
	return &Product{
		ID:          uuid.NewString(),
		OrgID:       orgID,
		Name:        name,
		Description: description,
		Spec:        spec,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}