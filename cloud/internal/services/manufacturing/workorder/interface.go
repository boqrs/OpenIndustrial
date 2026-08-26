package workorder

import (
	"context"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context, entity *model.WorkOrder) error
	GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.WorkOrder, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*model.WorkOrder, error)
	List(ctx context.Context, tenantID uuid.UUID, status *model.WorkOrderStatus, productionPlanID *uint) ([]*model.WorkOrder, error)
	Update(ctx context.Context, entity *model.WorkOrder) error
	SumPlannedQuantityByPlanID(ctx context.Context, tenantID uuid.UUID, productionPlanID uint) (int64, error)
}

type Service interface {
	CreateWorkOrder(ctx context.Context, req *CreateWorkOrderRequest) (*WorkOrderResponse, error)
	GetWorkOrder(ctx context.Context, id uint) (*WorkOrderResponse, error)
	ListWorkOrders(ctx context.Context, status *model.WorkOrderStatus, productionPlanID *uint) ([]*WorkOrderResponse, error)
	UpdateWorkOrder(ctx context.Context, id uint, req *UpdateWorkOrderRequest) (*WorkOrderResponse, error)
	ReleaseWorkOrder(ctx context.Context, id uint) error
	StartWorkOrder(ctx context.Context, id uint) error
	CancelWorkOrder(ctx context.Context, id uint) error
}