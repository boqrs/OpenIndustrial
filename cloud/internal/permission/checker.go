package permission

// Checker defines the core interface for the permission engine.
// It is responsible for answering the fundamental question: "Can this user do this action on this resource?"
type Checker interface {
	// Check verifies if a user has the permission to perform an action on a specific resource.
	// It encapsulates the entire logic of checking memberships, roles, role-permissions, and resource ACLs.
	Check(userID string, action string, resourceID string) (bool, error)
}