package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaterialType string

const (
	MaterialTypeRaw        MaterialType = "raw"
	MaterialTypeComponent  MaterialType = "component"
	MaterialTypeConsumable MaterialType = "consumable"
	MaterialTypePackaging  MaterialType = "packaging"
)

type Material struct {
	ID uint `gorm:"primaryKey"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Code string `gorm:"type:varchar(100);not null"`
	Name string `gorm:"type:varchar(255);not null"`

	MaterialType MaterialType `gorm:"type:varchar(32);not null"`
	Unit         string       `gorm:"type:varchar(32);not null"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Material) TableName() string {
	return "materials"
}