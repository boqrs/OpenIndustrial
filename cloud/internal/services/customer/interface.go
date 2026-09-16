package customer

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context, customer *model.Customer) error

	GetByID(ctx context.Context, id uint) (*model.Customer, error)

	GetByCode(ctx context.Context, code string) (*model.Customer, error)

	List(ctx context.Context, offset, limit int) ([]*model.Customer, int64, error)

	Update(ctx context.Context, customer *model.Customer) error
}

type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*Response, error)

	GetByID(ctx context.Context, id uint) (*Response, error)

	GetByCode(ctx context.Context, code string) (*Response, error)

	List(ctx context.Context, req *ListRequest) (*ListResponse, error)

	Update(ctx context.Context, id uint, req *UpdateRequest) (*Response, error)

	Deactivate(ctx context.Context, id uint) error
}
