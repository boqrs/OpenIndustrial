package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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