package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"gorm.io/gorm/clause"
)

type executionRepository struct {
	db *database.DBProvider
}

func NewExecutionRepository(db *database.DBProvider) *executionRepository {
	return &executionRepository{
		db: db,
	}
}

func (r *executionRepository) CreateExecution(ctx context.Context, entity *model.ProductionExecution, operations []*model.ExecutionOperation) error {
	return r.db.Get().WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(entity).Error; err != nil {
				return err
			}

			if len(operations) == 0 {
				return nil
			}

			for _, op := range operations {
				op.ExecutionID = entity.ID
			}

			if err := tx.Create(
				&operations,
			).Error; err != nil {
				return err
			}

			return nil
		})
}

func (r *executionRepository) createExecution(
	db *gorm.DB,
	entity *model.ProductionExecution,
	operations []*model.ExecutionOperation,
) error {
	if entity == nil {
		return gorm.ErrInvalidData
	}

	if err := db.Create(entity).Error; err != nil {
		return err
	}

	if len(operations) == 0 {
		return nil
	}

	for _, operation := range operations {
		if operation == nil {
			continue
		}

		operation.ExecutionID = entity.ID
	}

	if err := db.Create(&operations).Error; err != nil {
		return err
	}

	return nil
}

// CreateExecutionTx creates an execution and its operations using the
// transaction already stored in the context by UnitOfWork.
func (r *executionRepository) CreateExecutionTx(
	ctx context.Context,
	entity *model.ProductionExecution,
	operations []*model.ExecutionOperation,
) error {
	return r.createExecution(
		dbFromContext(ctx, r.db.Get()),
		entity,
		operations,
	)
}

func (r *executionRepository) GetExecutionByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.ProductionExecution, error) {
	var entity model.ProductionExecution

	err := r.db.Get().WithContext(ctx).
		Where(
			"tenant_id = ? AND id = ?",
			tenantID,
			id,
		).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *Repository) GetExecutionByIDForUpdateTx(
	ctx context.Context,
	tenantID uuid.UUID,
	id uint,
) (*model.ProductionExecution, error) {
	var entity model.ProductionExecution

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where(
			"tenant_id = ? AND id = ?",
			tenantID,
			id,
		).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *executionRepository) ListExecutions(ctx context.Context, tenantID uuid.UUID, workOrderID *uint, status *model.ProductionExecutionStatus) ([]*model.ProductionExecution, error) {
	var entities []*model.ProductionExecution

	query := r.db.Get().WithContext(ctx).
		Where(
			"tenant_id = ?",
			tenantID,
		)

	if workOrderID != nil {
		query = query.Where(
			"work_order_id = ?", *workOrderID,
		)
	}

	if status != nil {
		query = query.Where(
			"status = ?", *status,
		)
	}

	err := query.
		Order("created_at ASC").
		Find(&entities).
		Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *executionRepository) UpdateExecution(ctx context.Context, entity *model.ProductionExecution) error {
	return r.db.Get().WithContext(ctx).
		Save(entity).
		Error
}

func (r *executionRepository) GetOperation(ctx context.Context, executionID, operationID uint) (*model.ExecutionOperation, error) {
	var entity model.ExecutionOperation

	err := r.db.Get().WithContext(ctx).
		Where(
			"execution_id = ? AND id = ?",
			executionID,
			operationID,
		).
		First(&entity).
		Error

	// if errors.Is(
	// 	err,
	// 	gorm.ErrRecordNotFound,
	// ) {
	// 	return nil, execution.ErrOperationNotFound
	// }

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *executionRepository) ListOperations(ctx context.Context, executionID uint) ([]*model.ExecutionOperation, error) {
	var entities []*model.ExecutionOperation

	err := r.db.Get().WithContext(ctx).
		Where(
			"execution_id = ?",
			executionID,
		).
		Order("sequence ASC").
		Find(&entities).
		Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *executionRepository) UpdateOperation(ctx context.Context, entity *model.ExecutionOperation) error {
	return r.db.Get().WithContext(ctx).
		Save(entity).
		Error
}

func (r *executionRepository) CountExecutions(ctx context.Context, tenantID uuid.UUID, workOrderID uint) (int64, error) {
	var count int64

	err := r.db.Get().WithContext(ctx).
		Model(
			&model.ProductionExecution{},
		).
		Where(
			"tenant_id = ? AND work_order_id = ?",
			tenantID,
			workOrderID,
		).
		Count(&count).
		Error

	return count, err
}

func (r *executionRepository) GetCurrentOperation(ctx context.Context, executionID uint) (*model.ExecutionOperation, error) {
	var entity model.ExecutionOperation

	err := r.db.Get().WithContext(ctx).
		Where(
			"execution_id = ? AND status = ?",
			executionID,
			model.ExecutionOperationStatusInProgress,
		).
		Order("sequence ASC").
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Create implements execution.Repository.
// It creates the ProductionExecution and its associated ExecutionOperations in a single transaction.
func (r *executionRepository) Create(ctx context.Context, tx *gorm.DB, exec *model.ProductionExecution) error {
	// Use the provided transaction 'tx' to ensure atomicity.
	db := tx.WithContext(ctx)

	// GORM's Create will automatically handle the main object (ProductionExecution)
	// and its associated objects (the slice of ExecutionOperation) because of the
	// model struct tags defining the relationship.
	if err := db.Create(exec).Error; err != nil {
		return err
	}

	return nil
}
func (r *executionRepository) GetOperationForUpdateTx(
	ctx context.Context,
	executionID uint,
	operationID uint,
) (*model.ExecutionOperation, error) {
	var entity model.ExecutionOperation

	err := dbFromContext(ctx, r.db.Get()).
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where(
			"execution_id = ? AND id = ?",
			executionID,
			operationID,
		).
		First(&entity).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}
