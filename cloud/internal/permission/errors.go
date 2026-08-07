package permission

import "errors"

var (
	ErrPolicyNotFound = errors.New("policy not found")
	ErrForbidden      = errors.New("action forbidden")
)