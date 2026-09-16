package model

import "time"

type SalesOrderItem struct {
	ID uint `gorm:"primaryKey"`

	SalesOrderID uint `gorm:"not null;index"`

	ProductID uint `gorm:"not null;index"`

	OrderedQuantity int64 `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SalesOrderItem) TableName() string {
	return "sales_order_items"
}
