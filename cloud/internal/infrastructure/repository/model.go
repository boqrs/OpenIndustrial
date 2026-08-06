package repository

import "time"

type ProductInstanceModel struct {
	ID        string    `gorm:"primaryKey"`
	SN        string    `gorm:"uniqueIndex"`
	ProductID string
	OrgID     string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ProductInstanceModel) TableName() string {
	return "product_instances"
}

type LifecycleEventModel struct {
	ID                string `gorm:"primaryKey"`
	ProductInstanceID string
	EventType         string
	FromState         string
	ToState           string
	Source            string
	Operator          string
	Payload           []byte
	CreatedAt         time.Time
}

func (LifecycleEventModel) TableName() string {
	return "lifecycle_events"
}