package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents a system-defined business capability.
//
// Permissions are global and are not tenant-specific.
// They should normally be created through system/bootstrap data
// rather than through a regular tenant API.
type Permission struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"type:varchar(128);uniqueIndex;not null"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Role represents a stable business role inside a tenant.
//
// A role is a predefined set of permissions. A user belongs to
// exactly one role.
//
// System roles are created by platform/bootstrap data and cannot
// normally be modified or deleted by tenant users.
type Role struct {
	ID uint `gorm:"primaryKey"`

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	TenantID uint `gorm:"not null;index"`

	Tenant Tenant `gorm:"foreignKey:TenantID"`

	Name string `gorm:"type:varchar(100);not null"`

	Description string `gorm:"type:text"`

	IsSystem bool `gorm:"not null;default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}

// RolePermission defines the permissions granted to a role.
//
// Permission itself is global, while the role belongs to a tenant.
type RolePermission struct {
	RoleID uint `gorm:"primaryKey"`

	PermissionID uint `gorm:"primaryKey"`

	CreatedAt time.Time

	Role       Role       `gorm:"foreignKey:RoleID"`
	Permission Permission `gorm:"foreignKey:PermissionID"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}