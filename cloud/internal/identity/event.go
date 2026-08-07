package identity

import "time"

// Event represents a domain event related to the identity context.
type Event struct {
	Type      string
	UserID    string
	OrgID     string
	Timestamp time.Time
}

const (
	EventUserCreated       = "identity.user.created"
	EventMembershipCreated = "identity.membership.created"
)