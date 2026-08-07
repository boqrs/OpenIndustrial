package identity

import "time"

// User represents a human user in the system.
// It holds authentication credentials and personal information.
// A User's affiliation to an organization is managed via Membership.
type User struct {
	ID           string
	Username     string // Unique login name
	Email        string
	Phone        string
	PasswordHash string
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}