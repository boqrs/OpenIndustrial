package resource

import "errors"

var (
	ErrNotFound    = errors.New("resource not found")
	ErrOrgRequired = errors.New("organization ID is required")
)