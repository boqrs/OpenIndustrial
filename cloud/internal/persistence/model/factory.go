package model

import "time"

// Factory represents a factory belonging to a tenant.
// TenantID defines the security boundary.
// ResourceID links the factory to its corresponding Resource.
type Factory struct {
	// Internal database primary key.
	ID uint `gorm:"primaryKey"`

	// Tenant that owns this factory.
	// Factory codes are unique within a tenant, not globally.
	TenantID uint `gorm:"not null;uniqueIndex:idx_factory_tenant_code,priority:1;index"`

	// Corresponding factory resource.
	// One Resource can be associated with at most one Factory.
	ResourceID uint `gorm:"not null;uniqueIndex"`

	// Factory code, unique within the same tenant.
	Code string `gorm:"type:varchar(100);not null;uniqueIndex:idx_factory_tenant_code,priority:2"`

	// Physical address of the factory.
	Address string `gorm:"type:text"`

	// Timezone of the factory's location.
	Timezone string `gorm:"type:varchar(100);not null;default:'UTC'"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
