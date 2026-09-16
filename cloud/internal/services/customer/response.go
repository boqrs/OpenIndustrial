package customer

import (
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Response struct {
	ID uint `json:"id"`

	Code string             `json:"code"`
	Name string             `json:"name"`
	Type model.CustomerType `json:"type"`

	ContactName  string `json:"contactName"`
	ContactEmail string `json:"contactEmail"`
	ContactPhone string `json:"contactPhone"`
	Address      string `json:"address"`

	Status model.CustomerStatus `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListResponse struct {
	Items []*Response `json:"items"`
	Total int64       `json:"total"`
}
