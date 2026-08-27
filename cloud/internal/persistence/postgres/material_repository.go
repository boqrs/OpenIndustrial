package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"github.com/boqrs/nexus/database"

)


type materialRepository struct {
	db *database.DBProvider
}

// NewMaterialRepository creates a new GORM-based implementation of MaterialRepository.
func NewMaterialRepository(db *database.DBProvider) *materialRepository {
	return &materialRepository{db: db}
}

func (r *materialRepository) Create(ctx context.Context, material *model.Material) error {
	return r.db.Get().WithContext(ctx).Create(material).Error
}

func (r *materialRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.Material, error) {
	var material model.Material
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&material).
		Error
	return &material, err
}

func (r *materialRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*model.Material, error) {
	var material model.Material
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&material).
		Error
	return &material, err
}

func (r *materialRepository) List(ctx context.Context, tenantID uuid.UUID, offset int, limit int) ([]*model.Material, int64, error) {
	var materials []*model.Material
	var total int64

	query := r.db.Get().WithContext(ctx).Model(&model.Material{}).Where("tenant_id = ?", tenantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&materials).
		Error
	
	return materials, total, err
}

func (r *materialRepository) Update(ctx context.Context, material *model.Material) error {
	return r.db.Get().WithContext(ctx).Save(material).Error
}

func (r *materialRepository) Delete(ctx context.Context, tenantID uuid.UUID, id uint) error {
	return r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&model.Material{}).
		Error
}