package model

import (
	"time"

	"github.com/google/uuid"
)

type ExecutionResultStatus string

const (
	ExecutionResultStatusDraft     ExecutionResultStatus = "draft"
	ExecutionResultStatusConfirmed ExecutionResultStatus = "confirmed"
	ExecutionResultStatusCancelled ExecutionResultStatus = "cancelled"
)

// ExecutionResult represents the final production result of one execution.
//
// An Execution has exactly one ExecutionResult in the current design.
// The result owns production quantities and final confirmation state.
// Detailed per-unit production data remains in ExecutionOperation.Result.
type ExecutionResult struct {
	ID uint `gorm:"primaryKey"`
	TenantID uuid.UUID `gorm:"ype:uuid;not null;index"`
	ExecutionID uint `gorm:"not null;uniqueIndex"`
	WorkOrderID uint `gorm:"not null;index"`
	ProducedQuantity  int64 `gorm:"not null;default:0"`
	QualifiedQuantity int64 `gorm:"not null;default:0"`
	RejectedQuantity  int64 `gorm:"not null;default:0"`
	Status ExecutionResultStatus `gorm:"type:varchar(50);not null;default:'draft';index"`
	ConfirmedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ExecutionResult) TableName() string {
	return "execution_results"
}