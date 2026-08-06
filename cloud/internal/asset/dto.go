package asset

import (
	"time"
)

// CreateAssetRequest defines the structure for a request to create a new asset.
type CreateAssetRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	SN        string `json:"sn" binding:"required"`
}

// AssetResponse defines the structure for a response containing asset details.
type AssetResponse struct {
	ID        string    `json:"id"`
	SN        string    `json:"sn"`
	OrgID     string    `json:"org_id"`
	ProductID string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToAssetResponse converts an Asset entity to an AssetResponse DTO.
func ToAssetResponse(asset *Asset) *AssetResponse {
	return &AssetResponse{
		ID:        asset.ID.String(),
		SN:        asset.SN,
		OrgID:     asset.OrgID.String(),
		ProductID: asset.ProductID.String(),
		Status:    asset.Status,
		CreatedAt: asset.CreatedAt,
	}
}

// ToAssetListResponse converts a slice of Asset entities to a slice of AssetResponse DTOs.
func ToAssetListResponse(assets []*Asset) []*AssetResponse {
	res := make([]*AssetResponse, len(assets))
	for i, asset := range assets {
		res[i] = ToAssetResponse(asset)
	}
	return res
}