package role

// Type distinguishes between system-defined and custom-defined roles.
type Type string

const (
	// TypeSystem indicates a role that is predefined by the platform.
	TypeSystem Type = "system"
	// TypeCustom indicates a role that is created by a user within an organization.
	TypeCustom Type = "custom"
)