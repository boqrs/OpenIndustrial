package planning

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/factory"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	"github.com/google/uuid"
)

var (
	ErrProductionPlanNotFound = errors.New(
		"production plan not found",
	)

	ErrProductionPlanNoExists = errors.New(
		"production plan with this number already exists",
	)

	ErrInvalidProductionPlan = errors.New(
		"invalid production plan data, please check constraints",
	)

	ErrProductionPlanState = errors.New(
		"operation not allowed in the current production plan state",
	)

	ErrReferencedProductNotFound = errors.New(
		"referenced product not found",
	)

	ErrReferencedFactoryNotFound = errors.New(
		"referenced factory not found",
	)
)

type serviceImpl struct {
	repository Repository

	productSvc product.Service
	factorySvc factory.Service

	uow postgres.UnitOfWork
}

// NewService creates a new planning service.
func NewService(
	repo Repository,
	productSvc product.Service,
	factorySvc factory.Service,
	uow postgres.UnitOfWork,
) Service {
	return &serviceImpl{
		repository: repo,
		productSvc: productSvc,
		factorySvc: factorySvc,
		uow:        uow,
	}
}

func (s *serviceImpl) CreateProductionPlan(
	ctx context.Context,
	req *CreateProductionPlanRequest,
) (*ProductionPlanResponse, error) {
	if req == nil {
		return nil, ErrInvalidProductionPlan
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrInvalidProductionPlan
	}

	planNo := strings.TrimSpace(req.PlanNo)
	if planNo == "" ||
		req.ProductID == 0 ||
		req.FactoryID == 0 ||
		req.PlannedQuantity <= 0 ||
		!req.PlannedEndAt.After(req.PlannedStartAt) {
		return nil, ErrInvalidProductionPlan
	}

	// ---------------------------------------------------------------------
	// Plan number uniqueness
	// ---------------------------------------------------------------------

	_, err := s.repository.GetByPlanNo(
		ctx,
		tenantID,
		planNo,
	)

	if err == nil {
		return nil, ErrProductionPlanNoExists
	}

	if !errors.Is(err, ErrProductionPlanNotFound) {
		return nil, err
	}

	// ---------------------------------------------------------------------
	// Validate Product
	// ---------------------------------------------------------------------

	if _, err := s.productSvc.GetProductModel(
		ctx,
		req.ProductID,
	); err != nil {
		return nil, ErrReferencedProductNotFound
	}

	// ---------------------------------------------------------------------
	// Validate Factory
	//
	// factory service resolves the factory Resource using the authenticated
	// tenant from context, so a factory belonging to another tenant cannot
	// be used here.
	// ---------------------------------------------------------------------

	if _, err := s.factorySvc.GetFactory(
		ctx,
		req.FactoryID,
	); err != nil {
		return nil, ErrReferencedFactoryNotFound
	}

	// ---------------------------------------------------------------------
	// Create Resource first.
	//
	// ProductionPlan is a Resource-backed production-domain entity.
	// ---------------------------------------------------------------------
	entity := &model.ProductionPlan{
		TenantID:        tenantID,
		PlanNo:          planNo,
		ProductID:       req.ProductID,
		FactoryID:       req.FactoryID,
		PlannedQuantity: req.PlannedQuantity,
		PlannedStartAt:  req.PlannedStartAt,
		PlannedEndAt:    req.PlannedEndAt,
		Description:     strings.TrimSpace(req.Description),
		Status:          model.ProductionPlanStatusDraft,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		// ProductionPlan creation failed after Resource creation.
		// Remove the orphaned Resource.

		return nil, fmt.Errorf(
			"create production plan: %w",
			err,
		)
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) GetProductionPlanByID(
	ctx context.Context,
	id uint,
) (*ProductionPlanResponse, error) {
	if id == 0 {
		return nil, ErrInvalidProductionPlan
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrInvalidProductionPlan
	}

	entity, err := s.repository.GetByID(
		ctx,
		tenantID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return toResponse(entity), nil
}

func (s *serviceImpl) ListProductionPlans(
	ctx context.Context,
	status *model.ProductionPlanStatus,
) ([]*ProductionPlanResponse, error) {
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrInvalidProductionPlan
	}

	entities, err := s.repository.List(
		ctx,
		tenantID,
		status,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*ProductionPlanResponse,
		len(entities),
	)

	for i, entity := range entities {
		responses[i] = toResponse(entity)
	}

	return responses, nil
}

func (s *serviceImpl) UpdateProductionPlan(
	ctx context.Context,
	id uint,
	req *UpdateProductionPlanRequest,
) (*ProductionPlanResponse, error) {
	if id == 0 || req == nil {
		return nil, ErrInvalidProductionPlan
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, ErrInvalidProductionPlan
	}

	var result *model.ProductionPlan

	err := s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			entity, err := s.repository.GetByIDForUpdateTx(
				txCtx,
				tenantID,
				id,
			)
			if err != nil {
				return err
			}

			if entity.Status != model.ProductionPlanStatusDraft {
				return ErrProductionPlanState
			}

			if req.PlannedQuantity != nil {
				if *req.PlannedQuantity <= 0 {
					return ErrInvalidProductionPlan
				}

				entity.PlannedQuantity = *req.PlannedQuantity
			}

			if req.PlannedStartAt != nil {
				entity.PlannedStartAt = *req.PlannedStartAt
			}

			if req.PlannedEndAt != nil {
				entity.PlannedEndAt = *req.PlannedEndAt
			}

			if !entity.PlannedEndAt.After(
				entity.PlannedStartAt,
			) {
				return ErrInvalidProductionPlan
			}

			if req.Description != nil {
				entity.Description = strings.TrimSpace(
					*req.Description,
				)
			}

			if err := s.repository.UpdateTx(
				txCtx,
				entity,
			); err != nil {
				return fmt.Errorf(
					"update production plan: %w",
					err,
				)
			}

			result = entity

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return toResponse(result), nil
}

func (s *serviceImpl) ReleaseProductionPlan(
	ctx context.Context,
	id uint,
) error {
	if id == 0 {
		return ErrInvalidProductionPlan
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ErrInvalidProductionPlan
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			entity, err := s.repository.GetByIDForUpdateTx(
				txCtx,
				tenantID,
				id,
			)
			if err != nil {
				return err
			}

			if entity.Status != model.ProductionPlanStatusDraft {
				return ErrProductionPlanState
			}

			entity.Status = model.ProductionPlanStatusReleased

			if err := s.repository.UpdateTx(
				txCtx,
				entity,
			); err != nil {
				return fmt.Errorf(
					"release production plan: %w",
					err,
				)
			}

			return nil
		},
	)
}

func (s *serviceImpl) CancelProductionPlan(
	ctx context.Context,
	id uint,
) error {
	if id == 0 {
		return ErrInvalidProductionPlan
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return ErrInvalidProductionPlan
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			entity, err := s.repository.GetByIDForUpdateTx(
				txCtx,
				tenantID,
				id,
			)
			if err != nil {
				return err
			}

			switch entity.Status {
			case model.ProductionPlanStatusDraft,
				model.ProductionPlanStatusReleased:

				entity.Status =
					model.ProductionPlanStatusCancelled

			default:
				return ErrProductionPlanState
			}

			if err := s.repository.UpdateTx(
				txCtx,
				entity,
			); err != nil {
				return fmt.Errorf(
					"cancel production plan: %w",
					err,
				)
			}

			return nil
		},
	)
}

func toResponse(
	entity *model.ProductionPlan,
) *ProductionPlanResponse {
	if entity == nil {
		return nil
	}

	return &ProductionPlanResponse{
		ID: entity.ID,
		//	ResourceID:      entity.ResourceID,
		TenantID:        entity.TenantID,
		PlanNo:          entity.PlanNo,
		ProductID:       entity.ProductID,
		FactoryID:       entity.FactoryID,
		PlannedQuantity: entity.PlannedQuantity,
		PlannedStartAt:  entity.PlannedStartAt,
		PlannedEndAt:    entity.PlannedEndAt,
		Status:          entity.Status,
		Description:     entity.Description,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}
