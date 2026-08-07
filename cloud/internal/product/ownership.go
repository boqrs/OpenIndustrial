package product

import "time"

// Ownership records the transfer of a product instance to a customer.
type Ownership struct {
	ID                string
	ProductInstanceID string
	CustomerOrgID     string    // The organization ID of the customer
	UserID            string    // The specific user ID, if applicable
	StartAt           time.Time // When the ownership began
	EndAt             *time.Time // When the ownership ended (e.g., resale, decommission)
}