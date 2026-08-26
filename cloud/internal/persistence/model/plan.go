package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductionPlanStatus string

const (
	ProductionPlanStatusDraft      ProductionPlanStatus = "draft"
	ProductionPlanStatusReleased   ProductionPlanStatus = "released"
	ProductionPlanStatusInProgress ProductionPlanStatus = "in_progress"
	ProductionPlanStatusCompleted  ProductionPlanStatus = "completed"
	ProductionPlanStatusCancelled  ProductionPlanStatus = "cancelled"
)

// ProductionPlan is a Resource-backed manufacturing entity.
//
// ResourceUUID is the public identity of the plan.
// ID is only used for internal database relations.
type ProductionPlan struct {
	ID uint `gorm:"primaryKey"`

	ResourceUUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	PlanNo string `gorm:"type:varchar(100);not null"`

	ProductID uint `gorm:"not null;index"`

	FactoryID uint `gorm:"not null;index"`

	PlannedQuantity int64 `gorm:"not null"`

	PlannedStartAt time.Time `gorm:"not null;index"`

	PlannedEndAt time.Time `gorm:"not null"`

	Status ProductionPlanStatus `gorm:"type:varchar(32);not null;default:'draft';index"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ProductionPlan) TableName() string {
	return "production_plans"
}