package identity

// Status defines the lifecycle status of an identity-related entity.
type Status string

const (
	StatusActive   Status = "active"
	StatusDisabled Status = "disabled"
	StatusLocked   Status = "locked"
)