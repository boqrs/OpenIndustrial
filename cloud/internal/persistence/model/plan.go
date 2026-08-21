package model

import (
	"time"

	"github.com/google/uuid"
)

// ProductionPlanStatus defines the lifecycle of a production plan.
type ProductionPlanStatus string

const (
	ProductionPlanStatusDraft      ProductionPlanStatus = "draft"
	ProductionPlanStatusReleased   ProductionPlanStatus = "released"
	ProductionPlanStatusInProgress ProductionPlanStatus = "in_progress"
	ProductionPlanStatusCompleted  ProductionPlanStatus = "completed"
	ProductionPlanStatusCancelled  ProductionPlanStatus = "cancelled"
)

// ProductionPlan represents a high-level plan to produce a certain quantity of a product within a specific timeframe.
type ProductionPlan struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID        uuid.UUID `gorm:"type:uuid;not null;index"`
	PlanNo          string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_tenant_plan_no"`
	ProductID       uuid.UUID `gorm:"type:uuid;not null;index"`
	FactoryID       uuid.UUID `gorm:"type:uuid;not null;index"`
	PlannedQuantity int       `gorm:"not null"`
	PlannedStartAt  time.Time `gorm:"not null;index"`
	PlannedEndAt    time.Time `gorm:"not null"`
	Status          ProductionPlanStatus `gorm:"type:varchar(32);not null;default:'draft';index"`
	Description     string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (ProductionPlan) TableName() string {
	return "production_plans"
}