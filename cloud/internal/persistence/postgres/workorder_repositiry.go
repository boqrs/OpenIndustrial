package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	//"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/workorder"
	"github.com/boqrs/nexus/database"
	"github.com/google/uuid"
)

type WorkOrderRepository struct {
	db *database.DBProvider
}


func NewWorkOrderRepository(db *database.DBProvider) *WorkOrderRepository {
	return &WorkOrderRepository{
		db: db,
	}
}

func (r *WorkOrderRepository) Create(ctx context.Context, entity *model.WorkOrder) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *WorkOrderRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.WorkOrder, error) {
	var entity model.WorkOrder
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&entity).Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, workorder.ErrWorkOrderNotFound
		// }
		return nil, err
	}
	return &entity, nil
}

func (r *WorkOrderRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*model.WorkOrder, error) {
	var entity model.WorkOrder
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&entity).Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, workorder.ErrWorkOrderNotFound
		// }
		return nil, err
	}
	return &entity, nil
}

func (r *WorkOrderRepository) List(ctx context.Context, tenantID uuid.UUID, productID *uint, offset, limit int) ([]*model.WorkOrder, error) {
	var entities []*model.WorkOrder
	query := r.db.Get().WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("planned_start_at ASC")

	// if status != nil {
	// 	query = query.Where("status = ?", *status)
	// }

	if productID != nil {
		query = query.Where("production_plan_id = ?", *productID)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *WorkOrderRepository) Update(ctx context.Context, entity *model.WorkOrder) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *WorkOrderRepository) SumPlannedQuantityByPlanID(ctx context.Context, tenantID uuid.UUID, productionPlanID uint) (int64, error) {
	var total int64
	err := r.db.Get().WithContext(ctx).
		Model(&model.WorkOrder{}).
		Where("tenant_id = ? AND production_plan_id = ?", tenantID, productionPlanID).
		Where("status <> ?", model.WorkOrderStatusCancelled).
		Select("COALESCE(SUM(planned_quantity), 0)").
		Scan(&total).
		Error

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *WorkOrderRepository) Delete(ctx context.Context, tenantID uuid.UUID, id uint) error {
	return r.db.Get().WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&model.WorkOrder{}, "id = ?", id).
		Error
}

func (r *WorkOrderRepository) Count(ctx context.Context, tenantID uuid.UUID, productID uint) (int64, error) {
	var count int64
	if err := r.db.Get().WithContext(ctx).
		Model(&model.WorkOrder{}).
		Where("tenant_id = ?", tenantID).
		Where("product_id = ?", productID).
		Count(&count).
		Error; err != nil{
			return 0, err
		}

		return count, nil
}