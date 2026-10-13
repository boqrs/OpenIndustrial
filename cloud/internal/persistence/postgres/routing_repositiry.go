package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"gorm.io/gorm/clause"
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

// ProductBelongsToTenant checks the product model through its backing Resource.
// ProductBelongsToTenant checks product ownership through its backing Resource.
func (r *routingRepository) ProductBelongsToTenant(
	ctx context.Context,
	tenantID, productID uint,
) (bool, error) {
	var count int64

	err := r.db.Get().
		WithContext(ctx).
		Table("product_models AS p").
		Joins("JOIN resources AS res ON res.id = p.resource_id").
		Where(
			"p.id = ? AND res.tenant_id = ? AND res.deleted_at IS NULL",
			productID,
			tenantID,
		).
		Count(&count).Error

	return count > 0, err
}

func (r *routingRepository) GetRoutingByID(ctx context.Context, tenantID, id uint) (*model.Routing, error) {
	var entity model.Routing
	err := r.db.Get().WithContext(ctx).
		Model(&model.Routing{}).
		Joins("JOIN resources AS res ON res.id = routings.resource_id").
		Where("res.tenant_id = ? AND routings.id = ?", tenantID, id).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) GetRoutingByNameAndVersion(
	ctx context.Context,
	tenantID, productID uint,
	name string,
	version int,
) (*model.Routing, error) {
	var entity model.Routing
	err := r.db.Get().WithContext(ctx).
		Model(&model.Routing{}).
		Joins("JOIN resources AS res ON res.id = routings.resource_id").
		Where(
			"res.tenant_id = ? AND routings.product_id = ? AND routings.name = ? AND routings.version = ?",
			tenantID, productID, name, version,
		).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) ListRoutings(
	ctx context.Context,
	tenantID uint,
	productID *uint,
	status *model.RoutingStatus,
) ([]*model.Routing, error) {
	var entities []*model.Routing
	query := r.db.Get().WithContext(ctx).
		Model(&model.Routing{}).
		Joins("JOIN resources AS res ON res.id = routings.resource_id").
		Where("res.tenant_id = ?", tenantID)

	if productID != nil {
		query = query.Where("routings.product_id = ?", *productID)
	}
	if status != nil {
		query = query.Where("routings.status = ?", *status)
	}

	err := query.Order("routings.created_at DESC").Find(&entities).Error
	return entities, err
}

func (r *routingRepository) UpdateRouting(ctx context.Context, entity *model.Routing) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *routingRepository) DeactivateOtherRoutings(
	ctx context.Context,
	tenantID, productID, exceptRoutingID uint,
) error {
	db := r.db.Get().WithContext(ctx)
	resourceIDs := db.Table("resources").
		Select("id").
		Where("tenant_id = ?", tenantID)

	return db.Model(&model.Routing{}).
		Where("resource_id IN (?)", resourceIDs).
		Where("product_id = ? AND id <> ? AND status = ?",
			productID, exceptRoutingID, model.RoutingStatusActive).
		Update("status", model.RoutingStatusInactive).Error
}

func (r *routingRepository) CountOperations(ctx context.Context, tenantID, routingID uint) (int64, error) {
	var count int64
	err := r.db.Get().WithContext(ctx).
		Model(&model.RoutingOperation{}).
		Joins("JOIN routings AS rt ON rt.id = routing_operations.routing_id").
		Joins("JOIN resources AS res ON res.id = rt.resource_id").
		Where("res.tenant_id = ? AND rt.id = ?", tenantID, routingID).
		Count(&count).Error
	return count, err
}

func (r *routingRepository) CreateOperation(ctx context.Context, entity *model.RoutingOperation) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *routingRepository) GetOperation(
	ctx context.Context,
	tenantID, routingID, operationID uint,
) (*model.RoutingOperation, error) {
	var entity model.RoutingOperation
	err := r.db.Get().WithContext(ctx).
		Model(&model.RoutingOperation{}).
		Joins("JOIN routings AS rt ON rt.id = routing_operations.routing_id").
		Joins("JOIN resources AS res ON res.id = rt.resource_id").
		Where("res.tenant_id = ? AND rt.id = ? AND routing_operations.id = ?",
			tenantID, routingID, operationID).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) GetOperationByCode(
	ctx context.Context,
	tenantID, routingID uint,
	code string,
) (*model.RoutingOperation, error) {
	var entity model.RoutingOperation
	err := r.db.Get().WithContext(ctx).
		Model(&model.RoutingOperation{}).
		Joins("JOIN routings AS rt ON rt.id = routing_operations.routing_id").
		Joins("JOIN resources AS res ON res.id = rt.resource_id").
		Where("res.tenant_id = ? AND rt.id = ? AND routing_operations.code = ?",
			tenantID, routingID, code).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) GetOperationBySequence(
	ctx context.Context,
	tenantID, routingID uint,
	sequence int,
) (*model.RoutingOperation, error) {
	var entity model.RoutingOperation
	err := r.db.Get().WithContext(ctx).
		Model(&model.RoutingOperation{}).
		Joins("JOIN routings AS rt ON rt.id = routing_operations.routing_id").
		Joins("JOIN resources AS res ON res.id = rt.resource_id").
		Where("res.tenant_id = ? AND rt.id = ? AND routing_operations.sequence = ?",
			tenantID, routingID, sequence).
		First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) ListOperations(
	ctx context.Context,
	tenantID, routingID uint,
) ([]*model.RoutingOperation, error) {
	var entities []*model.RoutingOperation
	err := r.db.Get().WithContext(ctx).
		Model(&model.RoutingOperation{}).
		Joins("JOIN routings AS rt ON rt.id = routing_operations.routing_id").
		Joins("JOIN resources AS res ON res.id = rt.resource_id").
		Where("res.tenant_id = ? AND rt.id = ?", tenantID, routingID).
		Order("routing_operations.sequence ASC").
		Find(&entities).Error
	return entities, err
}

func (r *routingRepository) UpdateOperation(ctx context.Context, entity *model.RoutingOperation) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *routingRepository) DeleteOperation(
	ctx context.Context,
	tenantID, routingID, operationID uint,
) error {
	result := r.db.Get().WithContext(ctx).
		Where(
			"id = ? AND routing_id = ? AND routing_id IN (?)",
			operationID,
			routingID,
			r.db.Get().WithContext(ctx).
				Table("routings AS rt").
				Select("rt.id").
				Joins("JOIN resources AS res ON res.id = rt.resource_id").
				Where("res.tenant_id = ?", tenantID),
		).
		Delete(&model.RoutingOperation{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Legacy UUID-based lookup retained for callers that still use this method.
// Tenant ownership is still resolved through the backing Resource.
func (r *routingRepository) GetByID(
	ctx context.Context,
	tx *gorm.DB,
	tenantID uint,
	id uuid.UUID,
) (*model.Routing, error) {
	var entity model.Routing
	err := tx.WithContext(ctx).
		Model(&model.Routing{}).
		Joins("JOIN resources AS res ON res.id = routings.resource_id").
		Where("res.tenant_id = ? AND res.uuid = ?", tenantID, id).
		Preload("Operations", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *routingRepository) HasUnfinishedWorkOrders(
	ctx context.Context,
	tenantID, routingID uint,
) (bool, error) {
	var count int64
	err := r.db.Get().
		WithContext(ctx).
		Model(&model.WorkOrder{}).
		Where(
			"tenant_id = ? AND routing_id = ?",
			tenantID,
			routingID,
		).
		Where(
			"status IN ?",
			[]model.WorkOrderStatus{
				model.WorkOrderStatusDraft,
				model.WorkOrderStatusReleased,
				model.WorkOrderStatusInProgress,
			},
		).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *routingRepository) ActivateRoutingExclusive(
	ctx context.Context,
	tenantID uint,
	productID uint,
	routingID uint,
) (bool, error) {
	blocked := false

	err := r.db.Get().WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			// 1. 锁定产品记录，串行化同一产品的激活操作。
			var product struct {
				ID uint
			}

			err := tx.
				Table("product_models AS p").
				Select("p.id").
				Joins("JOIN resources AS res ON res.id = p.resource_id").
				Where(
					"p.id = ? AND res.tenant_id = ? AND res.deleted_at IS NULL",
					productID,
					tenantID,
				).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Take(&product).Error
			if err != nil {
				return err
			}

			// 2. 验证目标路线属于当前租户及指定产品。
			var target model.Routing

			err = tx.
				Model(&model.Routing{}).
				Joins("JOIN resources AS res ON res.id = routings.resource_id").
				Where(
					"res.tenant_id = ? AND routings.product_id = ? AND routings.id = ?",
					tenantID,
					productID,
					routingID,
				).
				Take(&target).Error
			if err != nil {
				return err
			}

			if target.Status == model.RoutingStatusArchived {
				return errors.New("cannot activate an archived routing")
			}

			if target.Status == model.RoutingStatusActive {
				return nil
			}

			// 3. 查询同一产品的其他活动路线。
			var activeRoutings []model.Routing

			err = tx.
				Model(&model.Routing{}).
				Joins("JOIN resources AS res ON res.id = routings.resource_id").
				Where(
					"res.tenant_id = ? AND routings.product_id = ? AND routings.status = ? AND routings.id <> ?",
					tenantID,
					productID,
					model.RoutingStatusActive,
					routingID,
				).
				Find(&activeRoutings).Error
			if err != nil {
				return err
			}

			// 4. 任何旧路线仍被未完成工单使用，都拒绝切换。
			for _, oldRouting := range activeRoutings {
				var count int64

				err = tx.
					Model(&model.WorkOrder{}).
					Where(
						"tenant_id = ? AND routing_id = ?",
						tenantID,
						oldRouting.ID,
					).
					Where(
						"status IN ?",
						[]model.WorkOrderStatus{
							model.WorkOrderStatusDraft,
							model.WorkOrderStatusReleased,
							model.WorkOrderStatusInProgress,
						},
					).
					Count(&count).Error
				if err != nil {
					return err
				}

				if count > 0 {
					blocked = true
					return nil
				}
			}

			// 5. 停用旧路线。
			for _, oldRouting := range activeRoutings {
				result := tx.
					Model(&model.Routing{}).
					Where(
						"id = ? AND status = ?",
						oldRouting.ID,
						model.RoutingStatusActive,
					).
					Update("status", model.RoutingStatusInactive)

				if result.Error != nil {
					return result.Error
				}
			}

			// 6. 激活目标路线。
			result := tx.
				Model(&model.Routing{}).
				Where(
					"id = ? AND status <> ?",
					routingID,
					model.RoutingStatusArchived,
				).
				Update("status", model.RoutingStatusActive)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected != 1 {
				return gorm.ErrRecordNotFound
			}

			return nil
		},
	)

	if err != nil {
		return false, err
	}

	return blocked, nil
}
