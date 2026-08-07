package org

import "errors"

var (
	// ErrNotFound indicates that a requested organization was not found.
	ErrNotFound = errors.New("organization not found")

	// ErrInvalidType indicates that the provided organization type is not valid.
	ErrInvalidType = errors.New("invalid organization type")
)