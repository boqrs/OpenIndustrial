package model

import (
	"time"

	"gorm.io/gorm"
)

type RoutingStatus string

const (
	RoutingStatusDraft    RoutingStatus = "draft"
	RoutingStatusActive   RoutingStatus = "active"
	RoutingStatusInactive RoutingStatus = "inactive"
	RoutingStatusArchived RoutingStatus = "archived"
)

// Routing is a Resource-backed manufacturing entity.
//
// ResourceUUID is the public identity of the routing.
// ID is only used for internal database relations.
type Routing struct {
	ID uint `gorm:"primaryKey"`
	ResourceID uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`
	Code string `gorm:"type:varchar(100);not null"`
	Name string `gorm:"type:varchar(255);not null"`
	Version int `gorm:"not null;default:1"`
	Status RoutingStatus `gorm:"type:varchar(50);not null;default:'draft';index"`
	Description string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Routing) TableName() string {
	return "routings"
}

type RoutingOperation struct {
	ID uint `gorm:"primaryKey"`
	RoutingID uint `gorm:"not null;index"`
	Sequence int `gorm:"not null"`
	Code string `gorm:"type:varchar(100);not null"`
	Name string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	WorkstationID *uint `gorm:"index"`
	StandardDurationSeconds int64 `gorm:"not null;default:0"`
	Required bool `gorm:"not null;default:true"`
	Parameters []byte `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RoutingOperation) TableName() string {
	return "routing_operations"
}