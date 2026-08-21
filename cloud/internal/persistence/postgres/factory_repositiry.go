package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"

)

// factoryRepository implements the factory.Repository interface using GORM.
type factoryRepository struct {
	db *database.DBProvider
}

// NewFactoryRepository creates a new GORM-based repository for factories.
func NewFactoryRepository(db *database.DBProvider) factory.Repository {
	return &factoryRepository{db: db}
}

// Create inserts a new factory record into the database.
func (r *factoryRepository) Create(ctx context.Context, entity *model.Factory) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

// GetByUUID retrieves a factory by its business UUID.
func (r *factoryRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*model.Factory, error) {
	var f model.Factory
	err := r.db.Get().WithContext(ctx).Where("uuid = ?", id).First(&f).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}
	return &f, nil
}

// GetByResourceID retrieves a factory by its associated resource ID.
func (r *factoryRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.Factory, error) {
	var f model.Factory
	err := r.db.Get().WithContext(ctx).Where("resource_id = ?", resourceID).First(&f).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}
	return &f, nil
}

// GetByCode retrieves a factory by its unique code.
// Note: The method signature in the interface had a typo 'Cortex', corrected to 'Context'.
func (r *factoryRepository) GetByCode(ctx context.Context, code string) (*model.Factory, error) {
	var f model.Factory
	err := r.db.Get().WithContext(ctx).Where("code = ?", code).First(&f).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, factory.ErrFactoryNotFound
		}
		return nil, err
	}
	return &f, nil
}

// Update saves changes to an existing factory record.
func (r *factoryRepository) Update(ctx context.Context, entity *model.Factory) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

// Delete removes a factory record from the database by its business UUID.
func (r *factoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Get().WithContext(ctx).Where("uuid = ?", id).Delete(&model.Factory{}).Error
}