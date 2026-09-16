package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/customer"
	"github.com/boqrs/nexus/database"
)

type customerRepository struct {
	db *database.DBProvider
}

func NewCustomerRepository(
	db *database.DBProvider,
) customer.Repository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) Create(
	ctx context.Context,
	entity *model.Customer,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(entity).
		Error
}

func (r *customerRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Customer, error) {
	var entity model.Customer

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *customerRepository) GetByCode(
	ctx context.Context,
	code string,
) (*model.Customer, error) {
	var entity model.Customer

	err := r.db.Get().
		WithContext(ctx).
		Where("code = ?", code).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *customerRepository) List(
	ctx context.Context,
	offset,
	limit int,
) ([]*model.Customer, int64, error) {
	var (
		entities []*model.Customer
		total    int64
	)

	db := r.db.Get().WithContext(ctx)

	if err := db.
		Model(&model.Customer{}).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&entities).
		Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *customerRepository) Update(
	ctx context.Context,
	entity *model.Customer,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(entity).
		Error
}

var _ customer.Repository = (*customerRepository)(nil)
