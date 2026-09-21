package model

import "time"

type SalesOrderItem struct {
	ID uint `gorm:"primaryKey"`

	SalesOrderID uint `gorm:"not null;index"`

	ProductID uint `gorm:"not null;index"`

	OrderedQuantity int64 `gorm:"not null"`

	// ReservedQuantity represents the quantity already allocated
	// to created / in-transit shipments but not yet delivered.
	ReservedQuantity int64 `gorm:"not null;default:0"`

	// FulfilledQuantity represents the quantity already delivered
	// to the customer.
	FulfilledQuantity int64 `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SalesOrderItem) TableName() string {
	return "sales_order_items"
}
