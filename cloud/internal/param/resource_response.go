package param

import (
	"time"

	"github.com/google/uuid"

)

// Resource is the standard data transfer object for a single resource response.
// It defines the JSON structure that clients will receive.
type Resource struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Type          string     `json:"type"`
	Name          string     `json:"name"`
	Code          *string    `json:"code,omitempty"`
	Status        string     `json:"status"`
	Metadata      []byte     `json:"metadata,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	RecordVersion int        `json:"record_version"`
	ParentID      *uuid.UUID `json:"parent_id,omitempty"`
	OwnerGroupID  *uuid.UUID `json:"owner_group_id,omitempty"`
}