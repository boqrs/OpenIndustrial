package param

import (
	"time"

	"github.com/google/uuid"
)

// FactoryResponse combines Resource and Factory information for API responses.
type FactoryResponse struct {
	ID         uuid.UUID `json:"id"` // This is the Resource ID
	TenantID   uuid.UUID `json:"tenant_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Code       string    `json:"code"`
	Address    string    `json:"address"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}