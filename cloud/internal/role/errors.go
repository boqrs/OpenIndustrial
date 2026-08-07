package role

import "errors"

var (
	ErrRoleNotFound = errors.New("role not found")
	ErrOrgRequired  = errors.New("organization ID is required")
)