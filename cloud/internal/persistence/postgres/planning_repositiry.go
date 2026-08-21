package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"

)

// productionPlanRepository implements the planning.Repository interface.
type productionPlanRepository struct {
	db *database.DBProvider
}

// NewProductionPlanRepository creates a new GORM repository for production plans.
func NewProductionPlanRepository(db *database.DBProvider) planning.Repository {
	return &productionPlanRepository{db: db}
}

func (r *productionPlanRepository) Create(ctx context.Context, entity *model.ProductionPlan) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *productionPlanRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*model.ProductionPlan, error) {
	var entity model.ProductionPlan
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, planning.ErrProductionPlanNotFound
		}
		return nil, err
	}
	return &entity, nil
}

func (r *productionPlanRepository) GetByPlanNo(ctx context.Context, tenantID uuid.UUID, planNo string) (*model.ProductionPlan, error) {
	var entity model.ProductionPlan
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND plan_no = ?", tenantID, planNo).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, planning.ErrProductionPlanNotFound
		}
		return nil, err
	}
	return &entity, nil
}

func (r *productionPlanRepository) List(ctx context.Context, tenantID uuid.UUID, status *model.ProductionPlanStatus) ([]*model.ProductionPlan, error) {
	var entities []*model.ProductionPlan
	query := r.db.Get().WithContext(ctx).Where("tenant_id = ?", tenantID).Order("planned_start_at ASC")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *productionPlanRepository) Update(ctx context.Context, entity *model.ProductionPlan) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}