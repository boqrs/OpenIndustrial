package org

import "errors"

var (
	// ErrOrgNameRequired is returned when an organization name is not provided.
	ErrOrgNameRequired = errors.New("organization name is required")
	// ErrOrgNotFound is returned when an organization is not found.
	ErrOrgNotFound = errors.New("organization not found")
)