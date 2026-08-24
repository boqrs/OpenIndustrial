package workorder

import (
	"context"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context,entity *model.WorkOrder) error
	GetByID(ctx context.Context,tenantID uuid.UUID,id uuid.UUID) (*model.WorkOrder, error)
	GetByOrderNo(ctx context.Context,tenantID uuid.UUID,orderNo string) (*model.WorkOrder, error)
	List(ctx context.Context,tenantID uuid.UUID,status *model.WorkOrderStatus,productionPlanID *uuid.UUID) ([]*model.WorkOrder, error)
	Update(ctx context.Context,entity *model.WorkOrder) error

	/*
		SumPlannedQuantity returns the total planned quantity
		of work orders belonging to a production plan.

		Only active work orders are counted.
		Cancelled work orders are excluded.
	*/
	SumPlannedQuantity(ctx context.Context,tenantID uuid.UUID,productionPlanID uuid.UUID) (int, error)
}

type Service interface {
	CreateWorkOrder(ctx context.Context,req *CreateWorkOrderRequest) (*WorkOrderResponse, error)
	GetWorkOrder(ctx context.Context,id uuid.UUID) (*WorkOrderResponse, error)
	ListWorkOrders(ctx context.Context,status *model.WorkOrderStatus,productionPlanID *uuid.UUID) ([]*WorkOrderResponse, error)
	UpdateWorkOrder(ctx context.Context,id uuid.UUID,req *UpdateWorkOrderRequest) (*WorkOrderResponse, error)
	ReleaseWorkOrder(ctx context.Context,id uuid.UUID) error
	StartWorkOrder(ctx context.Context,id uuid.UUID) error
	CancelWorkOrder(ctx context.Context,id uuid.UUID) error
}