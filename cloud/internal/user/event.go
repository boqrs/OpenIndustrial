package user

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	// UserCreated is the event type for when a user is created.
	UserCreated = "user.created"
)

// UserCreatedPayload is the data for a UserCreated event.
type UserCreatedPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	OrgID     uuid.UUID `json:"org_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserCreatedEvent creates a new event payload for a created user.
func NewUserCreatedEvent(user *User) ([]byte, error) {
	payload := UserCreatedPayload{
		UserID:    user.ID,
		OrgID:     user.OrgID,
		Username:  user.Username,
		Email:     user.Email,
		Timestamp: user.CreatedAt,
	}
	return json.Marshal(payload)
}