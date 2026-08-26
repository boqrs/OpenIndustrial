package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductionExecutionStatus string

const (
	ProductionExecutionStatusPending    ProductionExecutionStatus = "pending"
	ProductionExecutionStatusInProgress ProductionExecutionStatus = "in_progress"
	ProductionExecutionStatusCompleted  ProductionExecutionStatus = "completed"
	ProductionExecutionStatusFailed     ProductionExecutionStatus = "failed"
	ProductionExecutionStatusCancelled  ProductionExecutionStatus = "cancelled"
)

// ProductionExecution is a Resource-backed manufacturing entity.
//
// ResourceUUID is the public identity of the execution.
// ID is only used for internal database relations.
type ProductionExecution struct {
	ID uint `gorm:"primaryKey"`

	ResourceUUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	WorkOrderID uint `gorm:"not null;index"`

	ProductID uint `gorm:"not null;index"`

	RoutingID uint `gorm:"not null;index"`

	// RoutingVersion is the routing version actually used
	// by this execution.
	RoutingVersion int `gorm:"not null"`

	DeviceID *uint `gorm:"index"`

	Status ProductionExecutionStatus `gorm:"type:varchar(50);not null;default:'pending';index"`

	StartedAt *time.Time

	CompletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ProductionExecution) TableName() string {
	return "production_executions"
}

type ExecutionOperationStatus string

const (
	ExecutionOperationStatusPending    ExecutionOperationStatus = "pending"
	ExecutionOperationStatusInProgress ExecutionOperationStatus = "in_progress"
	ExecutionOperationStatusCompleted  ExecutionOperationStatus = "completed"
	ExecutionOperationStatusSkipped    ExecutionOperationStatus = "skipped"
	ExecutionOperationStatusFailed     ExecutionOperationStatus = "failed"
)

// ExecutionOperation is an execution-time child entity.
// It is not a Resource and therefore owns its own UUID.
type ExecutionOperation struct {
	ID uint `gorm:"primaryKey"`

	UUID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()"`

	ExecutionID uint `gorm:"not null;index"`

	RoutingOperationID *uint `gorm:"index"`

	Sequence int `gorm:"not null"`

	Code string `gorm:"type:varchar(100);not null"`

	Name string `gorm:"type:varchar(255);not null"`

	Description string `gorm:"type:text"`

	WorkstationID *uint `gorm:"index"`

	Status ExecutionOperationStatus `gorm:"type:varchar(50);not null;default:'pending';index"`

	StartedAt *time.Time

	CompletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ExecutionOperation) TableName() string {
	return "execution_operations"
}