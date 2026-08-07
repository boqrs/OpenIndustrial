package org

import "github.com/google/uuid"

// OrgCreatedEvent is published when a new organization is created.
type OrgCreatedEvent struct {
	ID   uuid.UUID
	Name string
	Type OrgType
}

// OrgStatusChangedEvent is published when an organization's status changes.
type OrgStatusChangedEvent struct {
	ID        uuid.UUID
	OldStatus OrgStatus
	NewStatus OrgStatus
}