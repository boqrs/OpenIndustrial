package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	"github.com/boqrs/nexus/database"
)

// factoryRepository implements the Factory Repository interface.
type factoryRepository struct {
	db *database.DBProvider
}

// NewFactoryRepository creates a new Factory Repository.
func NewFactoryRepository(db *database.DBProvider) factory.Repository {
	return &factoryRepository{db: db}
}

// Create inserts a new factory record.
func (r *factoryRepository) Create(
	ctx context.Context,
	entity *model.Factory,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(entity).Error
}

// GetByID retrieves a factory by its auto-increment primary key.
func (r *factoryRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Factory, error) {
	var f model.Factory

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&f).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}

	return &f, nil
}

// GetByResourceID retrieves a factory by its associated Resource ID.
func (r *factoryRepository) GetByResourceID(
	ctx context.Context,
	resourceID uint,
) (*model.Factory, error) {
	var f model.Factory

	err := r.db.Get().
		WithContext(ctx).
		Where("resource_id = ?", resourceID).
		First(&f).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}

	return &f, nil
}

// GetByCode retrieves a factory by its code.
// Tenant-scoped uniqueness will be implemented together with the
// Factory TenantID model change.
func (r *factoryRepository) GetByCode(
	ctx context.Context,
	code string,
) (*model.Factory, error) {
	var f model.Factory

	err := r.db.Get().
		WithContext(ctx).
		Where("code = ?", code).
		First(&f).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}

	return &f, nil
}

// Update saves changes to an existing factory.
func (r *factoryRepository) Update(
	ctx context.Context,
	entity *model.Factory,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(entity).Error
}

// Delete removes a factory by its primary key.
func (r *factoryRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	return r.db.Get().
		WithContext(ctx).
		Delete(&model.Factory{ID: id}).Error
}
