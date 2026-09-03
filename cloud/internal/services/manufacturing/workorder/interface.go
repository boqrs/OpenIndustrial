package workorder

import (
	"context"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context, workOrder *model.WorkOrder) error
	GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.WorkOrder, error)
	List(ctx context.Context, tenantID uuid.UUID, productID *uint, offset, limit int) ([]*model.WorkOrder, error)
	Count(ctx context.Context, tenantID uuid.UUID, productID uint) (int64, error)
	Update(ctx context.Context, workOrder *model.WorkOrder) error
	SumQuantityByPlanID(ctx context.Context, tenantID uuid.UUID, productionPlanID uint) (int64, error)
}
type Service interface {
	Create(ctx context.Context, tenantID uuid.UUID, req *CreateRequest) (*Response, error)
	GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*Response, error)
	List(ctx context.Context, req *ListRequest) (*ListResp, error)
	Update(ctx context.Context, tenantID uuid.UUID, id uint, req *UpdateRequest) (*Response, error)
	Release(ctx context.Context, tenantID uuid.UUID, id uint) error
	Start(ctx context.Context, tenantID uuid.UUID, id uint) error
	Complete(ctx context.Context, tenantID uuid.UUID, id uint) error
	Cancel(ctx context.Context, tenantID uuid.UUID, id uint) error
}
