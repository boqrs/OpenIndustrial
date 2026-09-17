package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
)

type productionPlanAllocationRepository struct {
	db *database.DBProvider
}

func NewProductionPlanAllocationRepository(
	db *database.DBProvider,
) *productionPlanAllocationRepository {
	return &productionPlanAllocationRepository{
		db: db,
	}
}

func (r *productionPlanAllocationRepository) CreateTx(
	ctx context.Context,
	allocation *model.ProductionPlanAllocation,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(allocation).Error
}

func (r *productionPlanAllocationRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.ProductionPlanAllocation, error) {
	var allocation model.ProductionPlanAllocation

	err := r.db.Get().
		WithContext(ctx).
		First(&allocation, id).Error
	if err != nil {
		return nil, err
	}

	return &allocation, nil
}

func (r *productionPlanAllocationRepository) ListBySalesOrderItemID(
	ctx context.Context,
	salesOrderItemID uint,
) ([]*model.ProductionPlanAllocation, error) {
	var allocations []*model.ProductionPlanAllocation

	err := r.db.Get().
		WithContext(ctx).
		Where("sales_order_item_id = ?", salesOrderItemID).
		Order("id ASC").
		Find(&allocations).Error
	if err != nil {
		return nil, err
	}

	return allocations, nil
}

func (r *productionPlanAllocationRepository) ListByProductionPlanID(
	ctx context.Context,
	productionPlanID uint,
) ([]*model.ProductionPlanAllocation, error) {
	var allocations []*model.ProductionPlanAllocation

	err := r.db.Get().
		WithContext(ctx).
		Where("production_plan_id = ?", productionPlanID).
		Order("id ASC").
		Find(&allocations).Error
	if err != nil {
		return nil, err
	}

	return allocations, nil
}

func (r *productionPlanAllocationRepository) SumBySalesOrderItemID(
	ctx context.Context,
	salesOrderItemID uint,
) (int64, error) {
	var total int64

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ProductionPlanAllocation{}).
		Where("sales_order_item_id = ?", salesOrderItemID).
		Select("COALESCE(SUM(allocated_quantity), 0)").
		Scan(&total).
		Error

	return total, err
}

func (r *productionPlanAllocationRepository) SumByProductionPlanID(
	ctx context.Context,
	productionPlanID uint,
) (int64, error) {
	var total int64

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.ProductionPlanAllocation{}).
		Where("production_plan_id = ?", productionPlanID).
		Select("COALESCE(SUM(allocated_quantity), 0)").
		Scan(&total).
		Error

	return total, err
}
