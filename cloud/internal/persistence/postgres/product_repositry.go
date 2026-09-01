package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	"github.com/boqrs/nexus/database"
)

// productRepository implements the product.Repository interface using GORM.
type ProductRepository struct {
	db *database.DBProvider
}

// NewProductRepository creates a new GORM-based repository for product models.
func NewProductRepository(	db *database.DBProvider) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, entity *model.ProductModel) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *ProductRepository) GetByID(ctx context.Context, id uint) (*model.ProductModel, error) {
	var pm model.ProductModel
	err := r.db.Get().WithContext(ctx).Where("id = ?", id).First(&pm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &pm, nil
}

func (r *ProductRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.ProductModel, error) {
	var pm model.ProductModel
	err := r.db.Get().WithContext(ctx).Where("resource_id = ?", resourceID).First(&pm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &pm, nil
}

func (r *ProductRepository) GetByCodeAndVersion(ctx context.Context, code string, version string) (*model.ProductModel, error) {
	var pm model.ProductModel
	err := r.db.Get().WithContext(ctx).Where("code = ? AND version = ?", code, version).First(&pm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &pm, nil
}

func (r *ProductRepository) List(ctx context.Context,req product.ListProductModelsRequest) ([]*model.ProductModel, int64, error){
	var items []*model.ProductModel
	var total int64

	query := r.db.Get().WithContext(ctx).Model(&model.ProductModel{})

	// Apply filters
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Code != "" {
		query = query.Where("code = ?", req.Code)
	}
	// Note: Filtering by status requires a JOIN with the resources table.
	if req.Status != "" {
		query = query.Joins("JOIN resources ON resources.uuid = product_models.resource_id").
			Where("resources.status = ?", req.Status)
	}

	// Get total count before applying pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (req.CurrentPage - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, entity *model.ProductModel) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Get().WithContext(ctx).Where("id = ?", id).Delete(&model.ProductModel{}).Error
}