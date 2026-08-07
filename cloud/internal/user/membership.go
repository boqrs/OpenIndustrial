package user

// Membership (or UserOrg) represents the relationship between a User and an Organization.
// It signifies that a user is a member of a specific organization.
type Membership struct {
	// ID is the unique identifier for the membership record.
	ID string

	// UserID is the ID of the user.
	UserID string

	// OrgID is the ID of the organization.
	OrgID string

	// Status indicates the status of the membership (e.g., active, pending, disabled).
	Status string
}