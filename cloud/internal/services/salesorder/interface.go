package salesorder

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type UnitOfWork interface {
	Execute(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}

type Repository interface {
	Create(
		ctx context.Context,
		order *model.SalesOrder,
	) error

	GetByID(
		ctx context.Context,
		id uint,
	) (*model.SalesOrder, error)

	GetByOrderNo(
		ctx context.Context,
		orderNo string,
	) (*model.SalesOrder, error)

	List(
		ctx context.Context,
		status *model.SalesOrderStatus,
		offset,
		limit int,
	) ([]*model.SalesOrder, int64, error)

	Update(
		ctx context.Context,
		order *model.SalesOrder,
	) error

	CreateItem(
		ctx context.Context,
		item *model.SalesOrderItem,
	) error

	CreateItemsTx(
		ctx context.Context,
		items []*model.SalesOrderItem,
	) error

	DeleteItemsTx(
		ctx context.Context,
		orderID uint,
	) error

	GetItemByID(
		ctx context.Context,
		id uint,
	) (*model.SalesOrderItem, error)

	ListItems(
		ctx context.Context,
		orderID uint,
	) ([]*model.SalesOrderItem, error)

	GetByIDForUpdateTx(
		ctx context.Context,
		id uint,
	) (*model.SalesOrder, error)

	UpdateTx(
		ctx context.Context,
		order *model.SalesOrder,
	) error

	GetItemByIDForUpdateTx(
		ctx context.Context,
		id uint,
	) (*model.SalesOrderItem, error)
}

type Service interface {
	Create(
		ctx context.Context,
		req *CreateRequest,
	) (*Response, error)

	GetByID(
		ctx context.Context,
		id uint,
	) (*Response, error)

	GetByOrderNo(
		ctx context.Context,
		orderNo string,
	) (*Response, error)

	List(
		ctx context.Context,
		req *ListRequest,
	) (*ListResponse, error)

	Update(
		ctx context.Context,
		id uint,
		req *UpdateRequest,
	) (*Response, error)

	Confirm(
		ctx context.Context,
		id uint,
	) error

	Cancel(
		ctx context.Context,
		id uint,
	) error
}