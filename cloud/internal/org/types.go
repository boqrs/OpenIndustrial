package org

// OrgStatus represents the status of an organization.
type OrgStatus string

const (
	OrgStatusActive    OrgStatus = "active"
	OrgStatusInactive  OrgStatus = "inactive"
	OrgStatusSuspended OrgStatus = "suspended"
)

// OrgType represents the type of an organization.
type OrgType string

const (
	OrgTypeTenant   OrgType = "tenant"
	OrgTypePartner  OrgType = "partner"
	OrgTypeInternal OrgType = "internal"
)

// Valid checks if the org type is valid.
func (ot OrgType) Valid() bool {
	switch ot {
	case OrgTypeTenant, OrgTypePartner, OrgTypeInternal:
		return true
	default:
		return false
	}
}