package model

import "time"

type CustomerType string

const (
	CustomerTypeBusiness    CustomerType = "business"
	CustomerTypeDistributor CustomerType = "distributor"
	CustomerTypeDealer      CustomerType = "dealer"
	CustomerTypeMarketplace CustomerType = "marketplace"
	CustomerTypeOEM         CustomerType = "oem"
)

type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
)

type Customer struct {
	ID uint `gorm:"primaryKey"`

	Code string       `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name string       `gorm:"type:varchar(255);not null"`
	Type CustomerType `gorm:"type:varchar(32);not null"`

	ContactName  string `gorm:"type:varchar(255)"`
	ContactEmail string `gorm:"type:varchar(255)"`
	ContactPhone string `gorm:"type:varchar(64)"`
	Address      string `gorm:"type:text"`

	Status CustomerStatus `gorm:"type:varchar(32);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Customer) TableName() string {
	return "customers"
}
