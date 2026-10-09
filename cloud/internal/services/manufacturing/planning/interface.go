package planning

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// Repository defines the persistence interface for production plans.
type Repository interface {
	Create(ctx context.Context, entity *model.ProductionPlan) error
	GetByID(ctx context.Context, tenantID uint, id uint) (*model.ProductionPlan, error)
	GetByPlanNo(ctx context.Context, tenantID uint, planNo string) (*model.ProductionPlan, error)
	List(ctx context.Context, tenantID uint, status *model.ProductionPlanStatus) ([]*model.ProductionPlan, error)
	Update(ctx context.Context, entity *model.ProductionPlan) error
	GetByIDForUpdateTx(ctx context.Context, tenantID uint, id uint) (*model.ProductionPlan, error)
	UpdateTx(ctx context.Context, entity *model.ProductionPlan) error
}

// Service defines the business logic for managing production plans.
type Service interface {
	CreateProductionPlan(ctx context.Context, req *CreateProductionPlanRequest) (*ProductionPlanResponse, error)
	GetProductionPlanByID(ctx context.Context, id uint) (*ProductionPlanResponse, error)
	ListProductionPlans(ctx context.Context, status *model.ProductionPlanStatus) ([]*ProductionPlanResponse, error)
	UpdateProductionPlan(ctx context.Context, id uint, req *UpdateProductionPlanRequest) (*ProductionPlanResponse, error)
	ReleaseProductionPlan(ctx context.Context, id uint) error
	CancelProductionPlan(ctx context.Context, id uint) error
}
