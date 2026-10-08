package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant represents an isolated organization boundary.
//
// Tenant is created by platform bootstrap,
// not by user registration.
type Tenant struct {
	ID uint `gorm:"primaryKey"`

	UUID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`

	Name string `gorm:"type:varchar(128);not null"`

	Code string `gorm:"type:varchar(64);uniqueIndex;not null"`

	Status string `gorm:"type:varchar(32);not null;default:'active'"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Tenant) TableName() string {
	return "tenants"
}

const (
	UserStatusInvited = "invited"

	UserStatusActive = "active"

	UserStatusDisabled = "disabled"
)

const (
	UserTypeAdmin = "admin"

	UserTypeEmployee = "employee"
)

// User represents a user identity inside a tenant.
//
// A user is always owned by exactly one tenant.
type User struct {
	ID uint `gorm:"primaryKey"`

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_uuid"`

	TenantID uint `gorm:"not null;index"`

	Tenant Tenant `gorm:"foreignKey:TenantID"`

	RoleID uint `gorm:"not null;index"`

	Role Role `gorm:"foreignKey:RoleID"`

	Email string `gorm:"type:varchar(255);not null"`

	Name string `gorm:"type:varchar(128);not null"`

	UserType string `gorm:"type:varchar(32);not null;default:'employee'"`

	Status string `gorm:"type:varchar(32);not null;default:'invited'"`

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {

	return "users"

}

type UserInvitation struct {
	ID uint `gorm:"primaryKey"`

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	TenantID uint `gorm:"not null;index"`

	UserID uint `gorm:"not null;index"`

	User User `gorm:"foreignKey:UserID"`

	Email string `gorm:"type:varchar(255);not null"`

	TokenHash string `gorm:"type:varchar(255);not null;index"`

	ExpiresAt time.Time

	UsedAt *time.Time

	CreatedBy uint `gorm:"not null"`

	CreatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserInvitation) TableName() string {

	return "user_invitations"

}

const (
	PrincipalProviderPassword = "password"

	PrincipalProviderOAuth = "oauth"
)

type Principal struct {
	ID uint `gorm:"primaryKey"`

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	TenantID uint `gorm:"not null;index;uniqueIndex:idx_principal_identity"`

	UserID uint `gorm:"not null;index"`

	User User `gorm:"foreignKey:UserID"`

	Provider string `gorm:"type:varchar(32);not null;uniqueIndex:idx_principal_identity"`

	Identifier string `gorm:"type:varchar(255);not null;uniqueIndex:idx_principal_identity"`

	SecretHash string `gorm:"type:text"`

	Status string `gorm:"type:varchar(32);not null;default:'active'"`

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Principal) TableName() string {

	return "principals"

}
