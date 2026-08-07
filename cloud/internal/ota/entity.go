package ota

import "time"

// Firmware represents a specific version of the gateway firmware.
type Firmware struct {
	// ID is the unique identifier for the firmware.
	ID string

	// Version is the version string (e.g., "v1.2.0").
	Version string

	// URL is the download link for the firmware package.
	URL string

	// Hash is the SHA256 hash of the firmware file for integrity checking.
	Hash string

	// Signature is the cryptographic signature of the hash for authenticity.
	Signature string

	// CreatedAt is the timestamp when the firmware was uploaded.
	CreatedAt time.Time
}