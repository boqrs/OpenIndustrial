package device

import (
	"time"

	"github.com/google/uuid"
)

// Device represents a physical or virtual device managed by the system.
// Its ID fields are now correctly typed as uuid.UUID to maintain consistency.
type Device struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}