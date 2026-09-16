package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"gorm.io/gorm/clause"
)

type salesOrderRepository struct {
	db *database.DBProvider
}

func NewSalesOrderRepository(db *database.DBProvider) *salesOrderRepository {
	return &salesOrderRepository{
		db: db,
	}
}

func (r *salesOrderRepository) Create(
	ctx context.Context,
	order *model.SalesOrder,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(order).
		Error
}

func (r *salesOrderRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.SalesOrder, error) {
	var order model.SalesOrder

	err := r.db.Get().
		WithContext(ctx).
		First(&order, id).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *salesOrderRepository) GetByOrderNo(
	ctx context.Context,
	orderNo string,
) (*model.SalesOrder, error) {
	var order model.SalesOrder

	err := r.db.Get().
		WithContext(ctx).
		Where("order_no = ?", orderNo).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *salesOrderRepository) List(
	ctx context.Context,
	status *model.SalesOrderStatus,
	offset,
	limit int,
) ([]*model.SalesOrder, int64, error) {
	db := r.db.Get().WithContext(ctx)

	query := db.Model(&model.SalesOrder{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []*model.SalesOrder

	if err := query.
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).
		Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil

}

func (r *salesOrderRepository) Update(
	ctx context.Context,
	order *model.SalesOrder,
) error {
	return r.db.Get().
		WithContext(ctx).
		Model(&model.SalesOrder{}).
		Where("id = ?", order.ID).
		Updates(map[string]any{
			"order_date":  order.OrderDate,
			"description": order.Description,
			"status":      order.Status,
			"updated_at":  order.UpdatedAt,
		}).
		Error
}

func (r *salesOrderRepository) CreateItem(
	ctx context.Context,
	item *model.SalesOrderItem,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(item).
		Error
}

func (r *salesOrderRepository) CreateItemsTx(
	ctx context.Context,
	items []*model.SalesOrderItem,
) error {
	if len(items) == 0 {
		return nil
	}

	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Create(&items).
		Error

}

func (r *salesOrderRepository) DeleteItemsTx(
	ctx context.Context,
	orderID uint,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Where("sales_order_id = ?", orderID).
		Delete(&model.SalesOrderItem{}).
		Error
}

func (r *salesOrderRepository) GetItemByID(
	ctx context.Context,
	id uint,
) (*model.SalesOrderItem, error) {
	var item model.SalesOrderItem

	err := r.db.Get().
		WithContext(ctx).
		First(&item, id).
		Error

	if err != nil {
		return nil, err
	}

	return &item, nil

}

func (r *salesOrderRepository) ListItems(
	ctx context.Context,
	orderID uint,
) ([]*model.SalesOrderItem, error) {
	var items []*model.SalesOrderItem

	err := r.db.Get().
		WithContext(ctx).
		Where("sales_order_id = ?", orderID).
		Order("id ASC").
		Find(&items).
		Error

	return items, err

}

func (r *salesOrderRepository) GetByIDForUpdateTx(
	ctx context.Context,
	id uint,
) (*model.SalesOrder, error) {
	var order model.SalesOrder

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		First(&order, id).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *salesOrderRepository) UpdateTx(
	ctx context.Context,
	order *model.SalesOrder,
) error {
	return dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Model(&model.SalesOrder{}).
		Where("id = ?", order.ID).
		Updates(map[string]any{
			"status": order.Status,
		}).
		Error
}
