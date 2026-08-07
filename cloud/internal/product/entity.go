package product

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a type of product that the company manufactures.
type Product struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"` // A unique code for the product type
	Spec        string    `json:"spec"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewProduct creates a new Product entity.
func NewProduct(orgID uuid.UUID, name, code, spec, description string) *Product {
	now := time.Now().UTC()
	return &Product{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        name,
		Code:        code,
		Spec:        spec,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ProductInstance represents a single, unique manufactured item of a Product type.
type ProductInstance struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	OrgID           uuid.UUID `json:"org_id"`
	SerialNumber    string    `json:"serial_number"` // The unique serial number for this instance
	ManufacturedAt  time.Time `json:"manufactured_at"`
	LifecycleEvents []*LifecycleEvent `json:"lifecycle_events"`
}

// LifecycleEvent represents a significant event in a product instance's life.
type LifecycleEvent struct {
	ProductInstanceID uuid.UUID              `json:"product_instance_id"`
	Type              string                 `json:"type"` // e.g., "manufactured", "shipped", "installed"
	Timestamp         time.Time              `json:"timestamp"`
	Location          string                 `json:"location,omitempty"`
	Properties        map[string]interface{} `json:"properties,omitempty"`
}

// ProductionFinishedEvent is published when a work order for a product is completed.
type ProductionFinishedEvent struct {
	ProductInstanceID uuid.UUID `json:"product_instance_id"`
	SerialNumber      string    `json:"serial_number"`
	ProductID         uuid.UUID `json:"product_id"`
	OrgID             uuid.UUID `json:"org_id"`
	FinishedAt        time.Time `json:"finished_at"`
}