package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"github.com/boqrs/nexus/database"
)

type BomRepository struct {
	db *database.DBProvider
}

func NewBOMRepository(db *database.DBProvider) *BomRepository {
	return &BomRepository{
		db: db,
	}
}

func (r *BomRepository) Create(ctx context.Context,bom *model.BOM) error {
	return r.db.Get().
		WithContext(ctx).
		Create(bom).
		Error
}

func (r *BomRepository) GetByID(ctx context.Context,tenantID uuid.UUID,id uint) (*model.BOM, error) {
	var bom model.BOM

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND id = ?",
			tenantID,
			id,
		).
		First(&bom).
		Error

	if err != nil {
		return nil, err
	}

	return &bom, nil
}

func (r *BomRepository) GetByNoVersion(ctx context.Context,tenantID uuid.UUID,bomNo string,version int) (*model.BOM, error) {
	var bom model.BOM

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND bom_no = ? AND version = ?",
			tenantID,
			bomNo,
			version,
		).
		First(&bom).
		Error

	if err != nil {
		return nil, err
	}

	return &bom, nil
}

func (r *BomRepository) List(ctx context.Context,tenantID uuid.UUID,productID uuid.UUID,offset int,limit int) ([]*model.BOM, int64, error) {
	var boms []*model.BOM
	var total int64

	query := r.db.Get().
		WithContext(ctx).
		Model(&model.BOM{}).
		Where(
			"tenant_id = ? AND product_id = ?",
			tenantID,
			productID,
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Order("version DESC").
		Find(&boms).
		Error

	if err != nil {
		return nil, 0, err
	}

	return boms, total, nil
}

func (r *BomRepository) Update(ctx context.Context,bom *model.BOM) error {
	return r.db.Get().
		WithContext(ctx).
		Save(bom).
		Error
}

func (r *BomRepository) CreateItems(ctx context.Context,items []*model.BOMItem) error {
	if len(items) == 0 {
		return nil
	}

	return r.db.Get().
		WithContext(ctx).
		Create(&items).
		Error
}

func (r *BomRepository) GetItems(ctx context.Context,tenantID uuid.UUID,bomID uint) ([]*model.BOMItem, error) {
	var items []*model.BOMItem

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND bom_id = ?",
			tenantID,
			bomID,
		).
		Order("sequence ASC, id ASC").
		Find(&items).
		Error

	return items, err
}

func (r *BomRepository) DeleteItems(ctx context.Context,tenantID uuid.UUID,bomID uint) error {
	return r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND bom_id = ?",
			tenantID,
			bomID,
		).
		Delete(&model.BOMItem{}).
		Error
}