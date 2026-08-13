package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user account within a specific tenant.
type User struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserType string    `gorm:"type:varchar(20);not null"`

	// Profile stores user-specific information like name, avatar, etc.
	// Using json.RawMessage with a `jsonb` type is flexible.
	Profile json.RawMessage `gorm:"type:jsonb"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// --- GORM Relationships ---

	// A user belongs to a tenant.
	Tenant Tenant `gorm:"foreignKey:TenantID;references:UUID"`

	// A user can have multiple principals (login methods).
	Principals []*Principal `gorm:"foreignKey:UserID;references:UUID"`

	// A user can belong to many groups.
	Groups []*Group `gorm:"many2many:user_groups;"`

	// A user can have many roles.
	Roles []*Role `gorm:"many2many:user_roles;"`
}

// TableName explicitly sets the table name for the User model.
func (User) TableName() string {
	return "users"
}

// Group represents a collection of users for permissioning.
// This is the GORM model, with a separate auto-incrementing primary key (ID)
// and a public-facing business key (UUID).
type Group struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// --- GORM Relationships ---
	// A group belongs to a tenant.
	// Tenant Tenant `gorm:"foreignKey:TenantID;references:UUID"` // Uncomment if Tenant model exists and has UUID as a business key.
}

// TableName explicitly sets the table name for the Group model.
func (Group) TableName() string {
	return "groups"
}

// Tenant represents an isolated workspace in the system.
// It is the root entity for all other resources within a tenant's scope.
type Tenant struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier for the tenant.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	Name   string `gorm:"type:varchar(100);not null;unique"`
	Status string `gorm:"type:varchar(20);default:'active'"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName explicitly sets the table name for the Tenant model.
func (Tenant) TableName() string {
	return "tenants"
}