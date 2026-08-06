package product

import "errors"

var (
	ErrProductNameRequired  = errors.New("product name is required")
	ErrProductModelRequired = errors.New("product model is required")
	ErrProductNotFound      = errors.New("product not found")
)