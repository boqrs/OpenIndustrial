package identity

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrMembershipNotFound = errors.New("membership not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)