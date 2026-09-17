package allocation

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
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
		tenantID uuid.UUID,
		req *CreateRequest,
	) (*Response, error)

	GetByID(
		ctx context.Context,
		tenantID uuid.UUID,
		id uint,
	) (*Response, error)

	ListBySalesOrderItemID(
		ctx context.Context,
		tenantID uuid.UUID,
		salesOrderItemID uint,
	) ([]*Response, error)

	ListByProductionPlanID(
		ctx context.Context,
		tenantID uuid.UUID,
		productionPlanID uint,
	) ([]*Response, error)
}
