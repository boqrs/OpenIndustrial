package asset

import "errors"

var (
	ErrAssetSNRequired = errors.New("asset serial number is required")
	ErrAssetNotFound   = errors.New("asset not found")
)