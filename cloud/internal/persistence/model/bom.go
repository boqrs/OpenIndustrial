package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BOMStatus defines the lifecycle status of a BOM.
type BOMStatus string

const (
	BOMStatusDraft     BOMStatus = "draft"
	BOMStatusReleased  BOMStatus = "released"
	BOMStatusObsolete  BOMStatus = "obsolete"
	BOMStatusCancelled BOMStatus = "cancelled"
)

// BOM represents the header of a Bill of Materials.
// It defines the structure of a product for manufacturing.
type BOM struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	// ProductID is the domain ID of the product this BOM describes.
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Code is a human-readable identifier for the BOM, e.g., "MAIN-BOARD-V2".
	Code string `gorm:"type:varchar(100);not null"`

	// Version is the engineering version of this BOM.
	Version int `gorm:"not null"`

	// Status indicates the current lifecycle state of the BOM.
	Status BOMStatus `gorm:"type:varchar(32);not null;default:'draft';index"`

	// EffectiveFrom is the date when this BOM version becomes active.
	EffectiveFrom *time.Time `gorm:"index"`

	// EffectiveTo is the date when this BOM version is no longer active.
	EffectiveTo *time.Time `gorm:"index"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the database table name for the BOM model.
func (BOM) TableName() string {
	return "boms"
}

// BOMItem represents a single line item within a BOM.
// It specifies a material and its required quantity.
type BOMItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	// BOMID is the ID of the BOM header this item belongs to.
	BOMID uuid.UUID `gorm:"type:uuid;not null;index"`

	// MaterialID is the domain ID of the material required.
	MaterialID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Sequence defines the order or grouping of items in the BOM.
	Sequence int `gorm:"not null;default:0"`

	// Quantity is the amount of the material required.
	Quantity float64 `gorm:"type:numeric(20,6);not null"`

	// Unit is the unit of measure for the quantity (e.g., "pcs", "kg").
	Unit string `gorm:"type:varchar(32);not null"`

	// ScrapRate is the expected percentage of material that will be wasted.
	ScrapRate float64 `gorm:"type:numeric(10,6);not null;default:0"`

	// IsOptional indicates if this item is optional for the assembly.
	IsOptional bool `gorm:"not null;default:false"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the database table name for the BOMItem model.
func (BOMItem) TableName() string {
	return "bom_items"
}