package routing

import (
	//"github.com/google/uuid"

)

type CreateRoutingRequest struct {
	ProductID   uint   `json:"productId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"isDefault"`
}

type UpdateRoutingRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsDefault   *bool   `json:"isDefault,omitempty"`
}

type CreateOperationRequest struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	WorkCenterID   uint   `json:"workCenterId"` // Assuming WorkCenter is also a resource with uint ID
	Sequence       int    `json:"sequence"`
	SetupTime      int    `json:"setupTime"`
	ProcessingTime int    `json:"processingTime"`
}

type UpdateOperationRequest struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	//WorkCenterID   *uint   `json:"workCenterId,omitempty"`
	Sequence       *int    `json:"sequence,omitempty"`
	//SetupTime      *int    `json:"setupTime,omitempty"`
	//ProcessingTime *int    `json:"processingTime,omitempty"`
}