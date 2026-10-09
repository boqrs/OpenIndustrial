package allocation

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	CreateTx(
		ctx context.Context,
		allocation *model.ProductionPlanAllocation,
	) error

	GetByID(
		ctx context.Context,
		id uint,
	) (*model.ProductionPlanAllocation, error)

	ListBySalesOrderItemID(
		ctx context.Context,
		salesOrderItemID uint,
	) ([]*model.ProductionPlanAllocation, error)

	ListByProductionPlanID(
		ctx context.Context,
		productionPlanID uint,
	) ([]*model.ProductionPlanAllocation, error)

	SumBySalesOrderItemID(
		ctx context.Context,
		salesOrderItemID uint,
	) (int64, error)

	SumByProductionPlanID(
		ctx context.Context,
		productionPlanID uint,
	) (int64, error)
}

type Service interface {
	Create(
		ctx context.Context,
		tenantID uint,
		req *CreateRequest,
	) (*Response, error)

	GetByID(
		ctx context.Context,
		tenantID uint,
		id uint,
	) (*Response, error)

	ListBySalesOrderItemID(
		ctx context.Context,
		tenantID uint,
		salesOrderItemID uint,
	) ([]*Response, error)

	ListByProductionPlanID(
		ctx context.Context,
		tenantID uint,
		productionPlanID uint,
	) ([]*Response, error)

	GetAllocatedQuantityByProductionPlanID(
		ctx context.Context,
		tenantID uint,
		productionPlanID uint,
	) (int64, error)
}
