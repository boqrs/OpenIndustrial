package model

import "time"

type SalesOrderStatus string

const (
	SalesOrderStatusDraft     SalesOrderStatus = "draft"
	SalesOrderStatusConfirmed SalesOrderStatus = "confirmed"
	SalesOrderStatusCompleted SalesOrderStatus = "completed"
	SalesOrderStatusCancelled SalesOrderStatus = "cancelled"
)

type SalesOrder struct {
	ID uint `gorm:"primaryKey"`

	CustomerID uint `gorm:"not null;index"`

	OrderNo string `gorm:"type:varchar(100);not null;uniqueIndex"`

	Status SalesOrderStatus `gorm:"type:varchar(32);not null;default:'draft';index"`

	OrderDate time.Time `gorm:"not null;index"`

	Description string `gorm:"type:text"`

	Items []*SalesOrderItem `gorm:"foreignKey:SalesOrderID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SalesOrder) TableName() string {
	return "sales_orders"
}
