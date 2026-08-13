package postgres

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/factory"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"gorm.io/gorm"
)

// FactoryRepository implements the factory.Repository interface using GORM.
type FactoryRepository struct {
	db *gorm.DB
}

// NewFactoryRepository creates a new instance of FactoryRepository.
func NewFactoryRepository(db *gorm.DB) factory.Repository {
	return &FactoryRepository{db: db}
}

// Create inserts a new factory record into the database.
func (r *FactoryRepository) Create(ctx context.Context, f *model.Factory) error {
	return r.db.WithContext(ctx).Create(f).Error
}