package model

import "time"

type ProductionPlanAllocation struct {
	ID uint `gorm:"primaryKey"`

	ProductionPlanID uint `gorm:"not null;index"`
	SalesOrderItemID uint `gorm:"not null;index"`

	AllocatedQuantity int64 `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ProductionPlanAllocation) TableName() string {
	return "production_plan_allocations"
}
