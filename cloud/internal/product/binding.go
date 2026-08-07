package product

import "time"

// DeviceBinding links a physical device to a product instance (SN).
// This is crucial for customer-facing IoT features.
type DeviceBinding struct {
	ID                string
	DeviceID          string    // The ID of the physical device (e.g., a mainboard)
	ProductInstanceID string    // The product SN this device is part of
	StartAt           time.Time
	EndAt             *time.Time
}