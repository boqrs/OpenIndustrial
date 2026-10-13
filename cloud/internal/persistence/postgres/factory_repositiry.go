package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	"github.com/boqrs/nexus/database"
)

type factoryRepository struct {
	db *database.DBProvider
}

func NewFactoryRepository(db *database.DBProvider) factory.Repository {
	return &factoryRepository{db: db}
}

func (r *factoryRepository) Create(
	ctx context.Context,
	entity *model.Factory,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(entity).Error
}

func (r *factoryRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Factory, error) {
	var entity model.Factory

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}

	return &entity, nil
}

func (r *factoryRepository) GetByResourceID(
	ctx context.Context,
	resourceID uint,
) (*model.Factory, error) {
	var entity model.Factory

	err := r.db.Get().
		WithContext(ctx).
		Where("resource_id = ?", resourceID).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}

	return &entity, nil
}

func (r *factoryRepository) GetByCode(
	ctx context.Context,
	tenantID uint,
	code string,
) (*model.Factory, error) {
	var entity model.Factory

	err := r.db.Get().
		WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}

	return &entity, nil
}

func (r *factoryRepository) Update(
	ctx context.Context,
	entity *model.Factory,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(entity).Error
}

func (r *factoryRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.Get().
		WithContext(ctx).
		Delete(&model.Factory{ID: id})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return factory.ErrFactoryNotFound
	}

	return nil
}

// ListByTenant returns factories belonging to the specified tenant.
func (r *factoryRepository) ListByTenant(
    ctx context.Context,
    tenantID uint,
) ([]*model.Factory, error) {
    var entities []*model.Factory

    err := r.db.Get().
        WithContext(ctx).
        Where("tenant_id = ?", tenantID).
        Order("id ASC").
        Find(&entities).Error
    if err != nil {
        return nil, err
    }

    return entities, nil
}