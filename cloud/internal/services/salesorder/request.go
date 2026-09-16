package salesorder

import "time"

type CreateRequest struct {
	OrderNo     string       `json:"orderNo"`
	CustomerID  uint         `json:"customerID"`
	OrderDate   time.Time    `json:"orderDate"`
	Description string       `json:"description"`
	Items       []CreateItem `json:"items"`
}

type CreateItem struct {
	ProductID       uint  `json:"productID"`
	OrderedQuantity int64 `json:"orderedQuantity"`
}

type UpdateRequest struct {
	OrderDate   *time.Time `json:"orderDate,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type ListRequest struct {
	Status *string `form:"status"`
	Offset int     `form:"offset"`
	Limit  int     `form:"limit"`
}
