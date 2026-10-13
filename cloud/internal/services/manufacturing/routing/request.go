package routing

import "encoding/json"

type CreateRoutingRequest struct {
	ProductID   uint   `json:"productId" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoutingRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateOperationRequest struct {
	Code                    string          `json:"code" binding:"required"`
	Name                    string          `json:"name" binding:"required"`
	Description             string          `json:"description"`
	WorkstationID           *uint           `json:"workstationId"`
	Sequence                int             `json:"sequence" binding:"required,gt=0"`
	StandardDurationSeconds int64           `json:"standardDurationSeconds"`
	Required                *bool           `json:"required"`
	Parameters              json.RawMessage `json:"parameters"`
}

type UpdateOperationRequest struct {
	Code                    *string         `json:"code,omitempty"`
	Name                    *string         `json:"name,omitempty"`
	Description             *string         `json:"description,omitempty"`
	WorkstationID           *uint           `json:"workstationId,omitempty"`
	Sequence                *int            `json:"sequence,omitempty"`
	StandardDurationSeconds *int64          `json:"standardDurationSeconds,omitempty"`
	Required                *bool           `json:"required,omitempty"`
	Parameters              json.RawMessage `json:"parameters,omitempty"`
}
