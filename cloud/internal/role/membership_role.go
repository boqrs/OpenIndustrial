package role

// Binding represents the link between a membership and a role.
// This entity establishes that a user, through their membership in an organization,
// is assigned a specific role.
type Binding struct {
	ID           string
	MembershipID string
	RoleID       string
}