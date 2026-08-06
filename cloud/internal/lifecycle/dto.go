package lifecycle

import "github.com/google/uuid"

// --- Request DTOs ---

// CreateLifecycleInstanceRequest is the DTO for creating a new lifecycle instance for an entity.
type CreateLifecycleInstanceRequest struct {
	DefinitionName string    `json:"definition_name" binding:"required"`
	EntityType     string    `json:"entity_type" binding:"required"`
	EntityID       uuid.UUID `json:"entity_id" binding:"required"`
	InitialState   string    `json:"initial_state" binding:"required"`
}

// TriggerEventRequest is the DTO for triggering a new event in an entity's lifecycle.
type TriggerEventRequest struct {
	EntityType string    `json:"entity_type" binding:"required"`
	EntityID   uuid.UUID `json:"entity_id" binding:"required"`
	EventType  string    `json:"event_type" binding:"required"`
	FromState  string    `json:"from_state"`
	ToState    string    `json:"to_state" binding:"required"`
	Operator   string    `json:"operator"`
	Source     string    `json:"source"`
	Payload    string    `json:"payload"` // JSON string
}

// --- Response DTOs ---

// LifecycleInstanceResponse is the DTO for a lifecycle instance.
type LifecycleInstanceResponse struct {
	ID           uuid.UUID `json:"id"`
	EntityType   string    `json:"entity_type"`
	EntityID     uuid.UUID `json:"entity_id"`
	CurrentState string    `json:"current_state"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

// LifecycleEventResponse is the DTO for a single lifecycle event.
type LifecycleEventResponse struct {
	ID        uuid.UUID `json:"id"`
	FromState string    `json:"from_state"`
	ToState   string    `json:"to_state"`
	EventType string    `json:"event_type"`
	Operator  string    `json:"operator"`
	Source    string    `json:"source"`
	Payload   string    `json:"payload"`
	CreatedAt string    `json:"created_at"`
}

// TraceHistoryResponse is the DTO for the full trace history of an entity.
type TraceHistoryResponse struct {
	EntityType string                   `json:"entity_type"`
	EntityID   uuid.UUID                `json:"entity_id"`
	History    []LifecycleEventResponse `json:"history"`
}