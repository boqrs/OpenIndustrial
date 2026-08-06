package productinstance

import "errors"

var (
	ErrNotFound = errors.New(
		"product instance not found",
	)

	ErrDuplicateSN = errors.New(
		"serial number already exists",
	)
)