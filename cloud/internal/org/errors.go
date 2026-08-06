package org

import "errors"

var (
	ErrOrgNameRequired = errors.New("organization name is required")
	ErrInvalidOrgType  = errors.New("invalid organization type")
)