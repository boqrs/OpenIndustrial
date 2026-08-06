package role

import "errors"

var (
	ErrRoleNameRequired = errors.New("role name is required")
	ErrRoleNotFound     = errors.New("role not found")
)