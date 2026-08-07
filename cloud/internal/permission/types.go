package permission

// ScopeType defines the range to which a policy applies.
type ScopeType string

const (
	// ScopeOrganization applies the policy to all resources within an entire organization.
	ScopeOrganization ScopeType = "organization"
	// ScopeResource applies the policy to a single, specific resource.
	ScopeResource ScopeType = "resource"
	// ScopeSubTree applies the policy to a resource and all its descendants in the resource graph.
	ScopeSubTree ScopeType = "subtree"
)

// Effect determines whether a policy grants (allow) or revokes (deny) a permission.
type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)