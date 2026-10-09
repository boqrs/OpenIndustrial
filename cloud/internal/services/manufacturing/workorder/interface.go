package workorder

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context, workOrder *model.WorkOrder) error
	CreateTx(ctx context.Context, workOrder *model.WorkOrder) error
	GetByID(ctx context.Context, tenantID uint, id uint) (*model.WorkOrder, error)
	List(ctx context.Context, tenantID uint, productID *uint, offset, limit int) ([]*model.WorkOrder, error)
	GetByIDForUpdateTx(ctx context.Context, tenantID uint, id uint) (*model.WorkOrder, error)
	Count(ctx context.Context, tenantID uint, productID uint) (int64, error)
	Update(ctx context.Context, workOrder *model.WorkOrder) error
	UpdateTx(ctx context.Context, workOrder *model.WorkOrder) error
	SumQuantityByPlanID(ctx context.Context, tenantID uint, productionPlanID uint) (int64, error)
}
type Service interface {
	Create(ctx context.Context, tenantID uint, req *CreateRequest) (*Response, error)
	GetByID(ctx context.Context, tenantID uint, id uint) (*Response, error)
	List(ctx context.Context, req *ListRequest) (*ListResp, error)
	Update(ctx context.Context, tenantID uint, id uint, req *UpdateRequest) (*Response, error)
	Release(ctx context.Context, tenantID uint, id uint) error
	Start(ctx context.Context, tenantID uint, id uint) error
	Complete(ctx context.Context, tenantID uint, id uint) error
	Cancel(ctx context.Context, tenantID uint, id uint) error
}
