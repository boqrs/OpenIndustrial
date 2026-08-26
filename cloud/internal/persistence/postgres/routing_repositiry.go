package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
)

type routingRepository struct {
	db *database.DBProvider
}

func NewRoutingRepository(db *database.DBProvider) *routingRepository {
	return &routingRepository{db: db}
}

func (r *routingRepository) CreateRouting(ctx context.Context, entity *model.Routing) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *routingRepository) GetRouting(ctx context.Context, tenantID, id uuid.UUID) (*model.Routing, error) {
	var entity model.Routing
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) GetRoutingByNameAndVersion(ctx context.Context, tenantID, productID uuid.UUID, name, version string) (*model.Routing, error) {
	var entity model.Routing
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND product_id = ? AND name = ? AND version = ?", tenantID, productID, name, version).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not finding a record is not an error for this existence check
		}
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) ListRoutings(ctx context.Context, tenantID uuid.UUID, productID *uuid.UUID, status *model.RoutingStatus) ([]*model.Routing, error) {
	var entities []*model.Routing
	query := r.db.Get().WithContext(ctx).Where("tenant_id = ?", tenantID)
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	err := query.Order("created_at desc").Find(&entities).Error
	return entities, err
}

func (r *routingRepository) UpdateRouting(ctx context.Context, entity *model.Routing) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *routingRepository) DeactivateOtherRoutings(ctx context.Context, tenantID, productID, exceptRoutingID uuid.UUID) error {
	return r.db.Get().WithContext(ctx).
		Model(&model.Routing{}).
		Where("tenant_id = ? AND product_id = ? AND id != ? AND status = ?", tenantID, productID, exceptRoutingID, model.RoutingStatusActive).
		Update("status", model.RoutingStatusInactive).
		Error
}

func (r *routingRepository) CountOperations(ctx context.Context, tenantID, routingID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Get().WithContext(ctx).Model(&model.RoutingOperation{}).Where("tenant_id = ? AND routing_id = ?", tenantID, routingID).Count(&count).Error
	return count, err
}

func (r *routingRepository) CreateOperation(ctx context.Context, entity *model.RoutingOperation) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *routingRepository) GetOperation(ctx context.Context, tenantID, routingID, operationID uuid.UUID) (*model.RoutingOperation, error) {
	var entity model.RoutingOperation
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND routing_id = ? AND id = ?", tenantID, routingID, operationID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) GetOperationByCode(ctx context.Context, tenantID, routingID uuid.UUID, code string) (*model.RoutingOperation, error) {
	var entity model.RoutingOperation
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND routing_id = ? AND code = ?", tenantID, routingID, code).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not finding a record is not an error for this existence check
		}
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) GetOperationBySequence(ctx context.Context, tenantID, routingID uuid.UUID, sequence int) (*model.RoutingOperation, error) {
	var entity model.RoutingOperation
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND routing_id = ? AND sequence = ?", tenantID, routingID, sequence).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not finding a record is not an error for this existence check
		}
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) ListOperations(ctx context.Context, tenantID, routingID uuid.UUID) ([]*model.RoutingOperation, error) {
	var entities []*model.RoutingOperation
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND routing_id = ?", tenantID, routingID).Order("sequence asc").Find(&entities).Error
	return entities, err
}

func (r *routingRepository) UpdateOperation(ctx context.Context, entity *model.RoutingOperation) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *routingRepository) DeleteOperation(ctx context.Context, tenantID, routingID, operationID uuid.UUID) error {
	result := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND routing_id = ? AND id = ?", tenantID, routingID, operationID).Delete(&model.RoutingOperation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return error if no rows were deleted
	}
	return nil
}

// GetByID implements routing.Repository.
func (r *routingRepository) GetByID(ctx context.Context, tx *gorm.DB, tenantID, id uuid.UUID) (*model.Routing, error) {
	var entity model.Routing

	// Use the provided transaction 'tx' to ensure atomicity.
	db := tx.WithContext(ctx)

	// Preload operations and order them by sequence.
	err := db.
		Preload("Operations", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Where("tenant_id = ?", tenantID).
		First(&entity, "id = ?", id).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, routing.ErrRoutingNotFound
	// }
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
