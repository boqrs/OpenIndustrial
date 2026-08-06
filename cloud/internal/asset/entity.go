package asset

import (
	"time"

	"github.com/google/uuid"
)

// Asset represents a unique, physical instance of a product, identified by a Serial Number (SN).
// This is the central entity for lifecycle tracking.
type Asset struct {
	ID        uuid.UUID `json:"id"`
	SN        string    `json:"sn"`        // Serial Number, should be unique within an organization
	OrgID     uuid.UUID `json:"org_id"`
	ProductID uuid.UUID `json:"product_id"` // Links to the product blueprint
	Status    string    `json:"status"`    // e.g., "In_Production", "In_Stock", "In_Transit", "In_Use", "Decommissioned"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewAsset creates a new Asset entity.
func NewAsset(orgID, productID uuid.UUID, sn string) (*Asset, error) {
	if sn == "" {
		return nil, ErrAssetSNRequired
	}

	now := time.Now().UTC()
	return &Asset{
		ID:        uuid.New(),
		SN:        sn,
		OrgID:     orgID,
		ProductID: productID,
		Status:    "provisioned", // Initial status
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}