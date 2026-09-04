package postgres

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"github.com/boqrs/nexus/database"
)

type Repository struct {
	db *database.DBProvider
}

func NewRepository(	db *database.DBProvider) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	entity *model.ExecutionResult,
) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *Repository) GetByID(
	ctx context.Context,
	tenantID uuid.UUID,
	id uint,
) (*model.ExecutionResult, error) {
	var entity model.ExecutionResult

	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *Repository) GetByExecutionID(
	ctx context.Context,
	tenantID uuid.UUID,
	executionID uint,
) (*model.ExecutionResult, error) {
	var entity model.ExecutionResult

	err := r.db.Get().WithContext(ctx).
		Where(
			"tenant_id = ? AND execution_id = ?",
			tenantID,
			executionID,
		).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *Repository) Update(
	ctx context.Context,
	entity *model.ExecutionResult,
) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}