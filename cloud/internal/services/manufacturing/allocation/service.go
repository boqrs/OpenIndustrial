package allocation

import (
	"context"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/manufacturing/planning"
	salesorder "github.com/boqrs/OpenIndustrial/cloud/internal/services/salesorder"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidAllocation          = errors.New("invalid production plan allocation")
	ErrAllocationNotFound         = errors.New("production plan allocation not found")
	ErrAllocationQuantityExceeded = errors.New("allocation quantity exceeds available quantity")
	ErrAllocationProductMismatch  = errors.New("sales order item product does not match production plan product")
	ErrProductionPlanNotFound     = errors.New("production plan not found")
	ErrSalesOrderItemNotFound     = errors.New("sales order item not found")
)

type service struct {
	uow         postgres.UnitOfWork
	repository  Repository
	planning    planning.Repository
	salesOrders salesorder.Repository
}

func NewService(
	uow postgres.UnitOfWork,
	repository Repository,
	planningRepo planning.Repository,
	salesOrderRepo salesorder.Repository,
) Service {
	return &service{
		uow:         uow,
		repository:  repository,
		planning:    planningRepo,
		salesOrders: salesOrderRepo,
	}
}

func (s *service) Create(
	ctx context.Context,
	tenantID uuid.UUID,
	req *CreateRequest,
) (*Response, error) {
	if req == nil {
		return nil, ErrInvalidAllocation
	}

	if tenantID == uuid.Nil ||
		req.ProductionPlanID == 0 ||
		req.SalesOrderItemID == 0 ||
		req.AllocatedQuantity <= 0 {
		return nil, ErrInvalidAllocation
	}

	var result *model.ProductionPlanAllocation

	err := s.uow.Execute(ctx, func(txCtx context.Context) error {
		// Always lock SalesOrderItem first.
		orderItem, err := s.salesOrders.GetItemByIDForUpdateTx(
			txCtx,
			req.SalesOrderItemID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSalesOrderItemNotFound
			}
			return err
		}

		// Then lock ProductionPlan.
		plan, err := s.planning.GetByIDForUpdateTx(
			txCtx,
			tenantID,
			req.ProductionPlanID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProductionPlanNotFound
			}
			return err
		}

		if orderItem.ProductID != plan.ProductID {
			return ErrAllocationProductMismatch
		}

		itemAllocated, err := s.repository.SumBySalesOrderItemID(
			txCtx,
			orderItem.ID,
		)
		if err != nil {
			return fmt.Errorf(
				"sum sales order item allocation: %w",
				err,
			)
		}

		planAllocated, err := s.repository.SumByProductionPlanID(
			txCtx,
			plan.ID,
		)
		if err != nil {
			return fmt.Errorf(
				"sum production plan allocation: %w",
				err,
			)
		}

		if itemAllocated+req.AllocatedQuantity >
			orderItem.OrderedQuantity {
			return ErrAllocationQuantityExceeded
		}

		if planAllocated+req.AllocatedQuantity >
			plan.PlannedQuantity {
			return ErrAllocationQuantityExceeded
		}

		entity := &model.ProductionPlanAllocation{
			ProductionPlanID:  plan.ID,
			SalesOrderItemID:  orderItem.ID,
			AllocatedQuantity: req.AllocatedQuantity,
		}

		if err := s.repository.CreateTx(txCtx, entity); err != nil {
			return fmt.Errorf(
				"create production plan allocation: %w",
				err,
			)
		}

		result = entity

		return nil
	})

	if err != nil {
		return nil, err
	}

	return toResponse(result), nil
}

func (s *service) GetByID(
	ctx context.Context,
	tenantID uuid.UUID,
	id uint,
) (*Response, error) {
	if tenantID == uuid.Nil || id == 0 {
		return nil, ErrInvalidAllocation
	}

	entity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAllocationNotFound
		}

		return nil, err
	}

	// Validate that the production plan belongs to this tenant.
	if _, err := s.planning.GetByID(
		ctx,
		tenantID,
		entity.ProductionPlanID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAllocationNotFound
		}

		return nil, err
	}

	return toResponse(entity), nil
}

func (s *service) ListBySalesOrderItemID(
	ctx context.Context,
	tenantID uuid.UUID,
	salesOrderItemID uint,
) ([]*Response, error) {
	if tenantID == uuid.Nil || salesOrderItemID == 0 {
		return nil, ErrInvalidAllocation
	}

	entities, err := s.repository.ListBySalesOrderItemID(
		ctx,
		salesOrderItemID,
	)
	if err != nil {
		return nil, err
	}

	responses := make([]*Response, 0, len(entities))

	for _, entity := range entities {
		if _, err := s.planning.GetByID(
			ctx,
			tenantID,
			entity.ProductionPlanID,
		); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			return nil, err
		}

		responses = append(responses, toResponse(entity))
	}

	return responses, nil
}

func (s *service) ListByProductionPlanID(
	ctx context.Context,
	tenantID uuid.UUID,
	productionPlanID uint,
) ([]*Response, error) {
	if tenantID == uuid.Nil || productionPlanID == 0 {
		return nil, ErrInvalidAllocation
	}

	if _, err := s.planning.GetByID(
		ctx,
		tenantID,
		productionPlanID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductionPlanNotFound
		}

		return nil, err
	}

	entities, err := s.repository.ListByProductionPlanID(
		ctx,
		productionPlanID,
	)
	if err != nil {
		return nil, err
	}

	responses := make([]*Response, len(entities))

	for i, entity := range entities {
		responses[i] = toResponse(entity)
	}

	return responses, nil
}

func toResponse(
	entity *model.ProductionPlanAllocation,
) *Response {
	if entity == nil {
		return nil
	}

	return &Response{
		ID:                entity.ID,
		ProductionPlanID:  entity.ProductionPlanID,
		SalesOrderItemID:  entity.SalesOrderItemID,
		AllocatedQuantity: entity.AllocatedQuantity,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}

func tenantIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value("tenant_id").(uuid.UUID); ok {
		return id
	}

	return uuid.Nil
}

func (s *service) GetAllocatedQuantityByProductionPlanID(
	ctx context.Context,
	tenantID uuid.UUID,
	productionPlanID uint,
) (int64, error) {
	if tenantID == uuid.Nil || productionPlanID == 0 {
		return 0, ErrInvalidAllocation
	}

	if _, err := s.planning.GetByID(
		ctx,
		tenantID,
		productionPlanID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrProductionPlanNotFound
		}

		return 0, err
	}

	return s.repository.SumByProductionPlanID(
		ctx,
		productionPlanID,
	)
}
