package salesorder

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Response struct {
	ID uint `json:"id"`

	CustomerID uint `json:"customerID"`

	OrderNo string `json:"orderNo"`

	Status model.SalesOrderStatus `json:"status"`

	OrderDate time.Time `json:"orderDate"`

	Description string `json:"description"`

	Items []*ItemResponse `json:"items"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ItemResponse struct {
	ID uint `json:"id"`

	SalesOrderID uint `json:"salesOrderID"`

	ProductID uint `json:"productID"`

	OrderedQuantity int64 `json:"orderedQuantity"`

	ReservedQuantity int64 `json:"reservedQuantity"`

	FulfilledQuantity int64 `json:"fulfilledQuantity"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListResponse struct {
	Items []*Response `json:"items"`

	Total int64 `json:"total"`
}