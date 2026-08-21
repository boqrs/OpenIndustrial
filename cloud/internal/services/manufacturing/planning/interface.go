package planning

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// Repository defines the persistence interface for production plans.
type Repository interface {
	Create(ctx context.Context, entity *model.ProductionPlan) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.ProductionPlan, error)
	GetByPlanNo(ctx context.Context, tenantID uuid.UUID, planNo string) (*model.ProductionPlan, error)
	List(ctx context.Context, tenantID uuid.UUID, status *model.ProductionPlanStatus) ([]*model.ProductionPlan, error)
	Update(ctx context.Context, entity *model.ProductionPlan) error
}

// Service defines the business logic for managing production plans.
type Service interface {
	CreateProductionPlan(ctx context.Context, req *CreateProductionPlanRequest) (*ProductionPlanResponse, error)
	GetProductionPlanByID(ctx context.Context, id uuid.UUID) (*ProductionPlanResponse, error)
	ListProductionPlans(ctx context.Context, status *model.ProductionPlanStatus) ([]*ProductionPlanResponse, error)
	UpdateProductionPlan(ctx context.Context, id uuid.UUID, req *UpdateProductionPlanRequest) (*ProductionPlanResponse, error)
	ReleaseProductionPlan(ctx context.Context, id uuid.UUID) error
	CancelProductionPlan(ctx context.Context, id uuid.UUID) error
}