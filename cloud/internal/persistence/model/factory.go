package model

import "time"

// Factory represents a factory belonging to a tenant.
type Factory struct {
	ID uint `gorm:"primaryKey"`

	// Tenant is the security and data-isolation boundary.
	TenantID uint `gorm:"not null;uniqueIndex:idx_factory_tenant_code,priority:1;index"`

	// Each factory is associated with exactly one Resource.
	ResourceID uint `gorm:"not null;uniqueIndex"`

	// Factory code is unique within a tenant.
	Code string `gorm:"type:varchar(100);not null;uniqueIndex:idx_factory_tenant_code,priority:2"`
    Name       string `gorm:"size:128;not null"`

	Address  string `gorm:"type:text"`
	Timezone string `gorm:"type:varchar(100);not null;default:'UTC'"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

