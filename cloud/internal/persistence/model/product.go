package model

import (
	"time"

)

type ProductModel struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	// ResourceID references resources.uuid.
	ResourceID uint `gorm:"not null;index"`

	// Code identifies the product model family.
	// Version distinguishes immutable model definitions.
	Code    string `gorm:"type:varchar(100);not null"`
	Version string `gorm:"type:varchar(50);not null"`

	// Product category, for example:
	// robot / plc / gateway / machine / product_iot.
	Category string `gorm:"type:varchar(100);not null;index"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}