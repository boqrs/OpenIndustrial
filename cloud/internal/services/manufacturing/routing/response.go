
package routing

import (
	"encoding/json"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type OperationResponse struct {
	ID                      uint            `json:"id"`
	RoutingID               uint            `json:"routingId"`
	Code                    string          `json:"code"`
	Name                    string          `json:"name"`
	Description             string          `json:"description"`
	WorkstationID           *uint           `json:"workstationId"`
	Sequence                int             `json:"sequence"`
	StandardDurationSeconds int64           `json:"standardDurationSeconds"`
	Required                bool            `json:"required"`
	Parameters              json.RawMessage `json:"parameters"`
	CreatedAt               time.Time       `json:"createdAt"`
	UpdatedAt               time.Time       `json:"updatedAt"`
}

type RoutingResponse struct {
	ID          uint                `json:"id"`
	ResourceID  uint                `json:"resourceId"`
	ProductID   uint                `json:"productId"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Version     int                 `json:"version"`
	Description string              `json:"description"`
	Status      model.RoutingStatus `json:"status"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}