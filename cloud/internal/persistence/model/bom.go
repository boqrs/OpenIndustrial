package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BOMStatus string

const (
	BOMStatusDraft    BOMStatus = "draft"
	BOMStatusReleased BOMStatus = "released"
	BOMStatusObsolete BOMStatus = "obsolete"
)

type BOM struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID uint `gorm:"not null;index"`
	BOMNo string `gorm:"type:varchar(100);not null"`
	Version int `gorm:"not null"`
	Status BOMStatus `gorm:"type:varchar(32);not null;default:'draft';index"`
	Description string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (BOM) TableName() string {
	return "boms"
}

type BOMItem struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`
	BOMID uint `gorm:"not null;index"`
	MaterialID uint `gorm:"not null;index"`
	Quantity float64 `gorm:"type:numeric(20,6);not null"`
	Unit string `gorm:"type:varchar(32);not null"`
	Sequence int `gorm:"not null;default:0"`
	OperationCode string `gorm:"type:varchar(100)"`
	IsOptional bool `gorm:"not null;default:false"`
	Description string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (BOMItem) TableName() string {
	return "bom_items"
}
