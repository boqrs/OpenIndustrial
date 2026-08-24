package model

import (
	"time"

	"github.com/google/uuid"
)

type WorkOrderStatus string

const (
	WorkOrderStatusDraft      WorkOrderStatus = "draft"
	WorkOrderStatusReleased   WorkOrderStatus = "released"
	WorkOrderStatusInProgress WorkOrderStatus = "in_progress"
	WorkOrderStatusCompleted  WorkOrderStatus = "completed"
	WorkOrderStatusCancelled  WorkOrderStatus = "cancelled"
)

type WorkOrder struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	OrderNo string `gorm:"type:varchar(100);not null"`

	ProductionPlanID uuid.UUID `gorm:"type:uuid;not null;index"`

	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`

	FactoryID uuid.UUID `gorm:"type:uuid;not null;index"`

	PlannedQuantity int `gorm:"not null"`

	CompletedQuantity int `gorm:"not null;default:0"`

	PlannedStartAt time.Time `gorm:"not null;index"`

	PlannedEndAt time.Time `gorm:"not null"`

	Status WorkOrderStatus `gorm:"type:varchar(32);not null;default:'draft';index"`

	Priority int `gorm:"not null;default:0"`

	Description string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}