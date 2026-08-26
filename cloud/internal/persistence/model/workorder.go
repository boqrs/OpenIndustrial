package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkOrderStatus string

const (
	WorkOrderStatusDraft      WorkOrderStatus = "draft"
	WorkOrderStatusPlanned    WorkOrderStatus = "planned"
	WorkOrderStatusReleased   WorkOrderStatus = "released"
	WorkOrderStatusInProgress WorkOrderStatus = "in_progress"
	WorkOrderStatusCompleted  WorkOrderStatus = "completed"
	WorkOrderStatusCancelled  WorkOrderStatus = "cancelled"
)

// WorkOrder is a Resource-backed manufacturing entity.
//
// ResourceUUID is the public identity of the work order.
// ID is only used for internal database relations.
type WorkOrder struct {
	ID uint `gorm:"primaryKey"`

	ResourceUUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	ProductionPlanID uint `gorm:"not null;index"`

	ProductID uint `gorm:"not null;index"`

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

func (WorkOrder) TableName() string {
	return "work_orders"
}