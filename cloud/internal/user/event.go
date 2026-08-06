package user

import (
	"time"

	"github.com/google/uuid"
)

// UserCreatedEvent is published when a new user is successfully created.
type UserCreatedEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	OrgID     uuid.UUID `json:"org_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserCreatedEvent creates a new UserCreatedEvent.
func NewUserCreatedEvent(user *User) *UserCreatedEvent {
	return &UserCreatedEvent{
		EventID:   uuid.New(),
		UserID:    user.ID,
		OrgID:     user.OrgID,
		Email:     user.Email,
		Timestamp: time.Now().UTC(),
	}
}