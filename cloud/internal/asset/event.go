package asset

import (
	"time"

	"github.com/google/uuid"
)

// AssetCreatedEvent is published when a new asset is provisioned in the system.
type AssetCreatedEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	AssetID   uuid.UUID `json:"asset_id"`
	SN        string    `json:"sn"`
	OrgID     uuid.UUID `json:"org_id"`
	ProductID uuid.UUID `json:"product_id"`
	Timestamp time.Time `json:"timestamp"`
}

// NewAssetCreatedEvent creates a new AssetCreatedEvent.
func NewAssetCreatedEvent(asset *Asset) *AssetCreatedEvent {
	return &AssetCreatedEvent{
		EventID:   uuid.New(),
		AssetID:   asset.ID,
		SN:        asset.SN,
		OrgID:     asset.OrgID,
		ProductID: asset.ProductID,
		Timestamp: time.Now().UTC(),
	}
}

// AssetStateChangedEvent is a generic event for when an asset's status changes.
type AssetStateChangedEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	AssetID   uuid.UUID `json:"asset_id"`
	SN        string    `json:"sn"`
	OrgID     uuid.UUID `json:"org_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
	// Context can hold extra data, e.g., WorkOrderID, ShippingID, CustomerID
	Context map[string]interface{} `json:"context,omitempty"`
}