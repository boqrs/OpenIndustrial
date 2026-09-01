package model

import (
	"time"


	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents an action that can be performed, e.g., "users:create".
// This is the GORM model for the `permissions` table.
// Permissions are typically seeded and not managed via a regular API,
// so we use a simple auto-incrementing primary key.
type Permission struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// Name is the unique identifier for the permission, e.g., "resources:create".
	Name string `gorm:"type:varchar(255);unique;not null"`

	Description string `gorm:"type:text"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName explicitly sets the table name for the Permission model.
func (Permission) TableName() string {
	return "permissions"
}

// Principal represents a way to authenticate a user (e.g., username/password, OAuth).
type Principal struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`
	// UUID is the public-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`
	// Provider indicates the authentication method, e.g., "password", "google", "wechat".
	Provider string `gorm:"type:varchar(50);not null;index:idx_provider_identifier,unique"`
	// Identifier is the unique ID within the provider, e.g., username, email, or openid.
	Identifier string `gorm:"type:varchar(255);not null;index:idx_provider_identifier,unique"`
	// Credential stores the secret, e.g., a hashed password.
	Credential string `gorm:"type:text"`
	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// --- GORM Relationships ---
	// A principal belongs to a user.
	User User `gorm:"foreignKey:UserID;references:UUID"`
}

// TableName explicitly sets the table name for the Principal model.
func (Principal) TableName() string {
	return "principals"
}

// Role represents a named set of permissions.
// This is the GORM model for the `roles` table.
type Role struct {
	// ID is the internal auto-incrementing primary key.
	ID uint `gorm:"primaryKey"`

	// UUID is the public-facing, unique business identifier.
	UUID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`

	TenantID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_role_name"`
	Description string    `gorm:"type:text"`
	IsSystem    bool      `gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// --- GORM Relationships ---

	// Permissions is a many-to-many relationship.
	// GORM will automatically create a join table named `role_permissions`.
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

// TableName explicitly sets the table name for the Role model.
func (Role) TableName() string {
	return "roles"
}