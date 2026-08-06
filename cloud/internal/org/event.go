package org

import (
	"time"

	"github.com/google/uuid"
)

// OrgCreatedEvent is published when a new organization is successfully created.
type OrgCreatedEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	OrgID     uuid.UUID `json:"org_id"`
	OrgName   string    `json:"org_name"`
	Timestamp time.Time `json:"timestamp"`
}

// NewOrgCreatedEvent creates a new OrgCreatedEvent.
func NewOrgCreatedEvent(org *Org) *OrgCreatedEvent {
	return &OrgCreatedEvent{
		EventID:   uuid.New(),
		OrgID:     org.ID,
		OrgName:   org.Name,
		Timestamp: time.Now().UTC(),
	}
}