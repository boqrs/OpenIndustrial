package productinstance

import (
	"time"
)

// ProductInstance represents a physical product.
type ProductInstance struct {
	ID string

	// Serial Number
	SN string

	// Product model ID
	ProductID string

	// Owner organization
	OrgID string

	// Current lifecycle state
	State string

	// Current location type
	LocationType string

	// Current location id
	LocationID string

	// Current owner
	OwnerOrgID string

	CreatedAt time.Time

	UpdatedAt time.Time
}