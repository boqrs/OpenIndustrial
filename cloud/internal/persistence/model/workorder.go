package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkOrderStatus string

const (
	WorkOrderStatusDraft      WorkOrderStatus = "draft"
	WorkOrderStatusReleased   WorkOrderStatus = "released"
	WorkOrderStatusInProgress WorkOrderStatus = "in_progress"
	WorkOrderStatusCompleted  WorkOrderStatus = "completed"
	WorkOrderStatusCancelled  WorkOrderStatus = "cancelled"
)

func (s WorkOrderStatus) String() string {
	return string(s)
}

func (s WorkOrderStatus) IsValid() bool {
	switch s {
	case WorkOrderStatusDraft, WorkOrderStatusReleased, WorkOrderStatusInProgress, WorkOrderStatusCompleted, WorkOrderStatusCancelled:
		return true
	default:
		return false
	}
}

type WorkOrder struct {
	ID uint `gorm:"primaryKey"`
	ResourceID uint `gorm:"not null;index"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductionPlanID uint `gorm:"not null;index"`
	FactoryID          uint  `gorm:"not null;index"`
	ProductionLineID   uint  `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`
	BOMID uint `gorm:"not null;index"`
	RoutingID uint `gorm:"not null;index"`
	Code string `gorm:"type:varchar(100);not null"`
	PlannedQuantity int64 `gorm:"not null;default:0"`
	Priority int `gorm:"not null;default:0"`
	DueDate *time.Time
	Status WorkOrderStatus `gorm:"type:varchar(50);not null;default:'draft';index"`
	StartedAt *time.Time
	CompletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}