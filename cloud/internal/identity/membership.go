package identity

import "time"

// Membership represents a user's affiliation with an organization.
// It links a User to an Org, establishing that the user is a member of that organization.
type Membership struct {
	ID string

	UserID string
	OrgID  string

	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}