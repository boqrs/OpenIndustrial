package permission

import "errors"

var (
	ErrPermissionNameRequired = errors.New("permission name is required")
	ErrPermissionNotFound     = errors.New("permission not found")
)