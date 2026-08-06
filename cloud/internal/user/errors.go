package user

import "errors"

var (
	ErrUserOrgIDRequired  = errors.New("user must be associated with an organization")
	ErrUsernameRequired   = errors.New("username is required")
	ErrUserEmailRequired  = errors.New("user email is required")
	ErrUserPasswordRequired = errors.New("user password is required")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
)