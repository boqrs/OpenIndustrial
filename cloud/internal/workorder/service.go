package workorder

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/pkg/event"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/product"
	"github.com/google/uuid"
)

// Service provides business logic for work orders.
type Service struct {
	repo        Repository
	productRepo product.Repository
	eventBus    event.Bus
}

// NewService creates a new work order service.
func NewService(repo Repository, productRepo product.Repository, eventBus event.Bus) *Service {
	return &Service{
		repo:        repo,
		productRepo: productRepo,
		eventBus:    eventBus,
	}
}

// CreateWorkOrder creates a new work order.
func (s *Service) CreateWorkOrder(ctx context.Context, orgID, productID uuid.UUID, quantity int, scheduledAt time.Time) (*WorkOrder, error) {
	wo, err := NewWorkOrder(orgID, productID, quantity, scheduledAt)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, wo); err != nil {
		return nil, err
	}
	return wo, nil
}

// FindByID retrieves a work order by its ID.
func (s *Service) FindByID(ctx context.Context, orgID, workOrderID uuid.UUID) (*WorkOrder, error) {
	return s.repo.FindByID(ctx, orgID, workOrderID)
}

// ListByOrg lists all work orders for an organization.
func (s *Service) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*WorkOrder, error) {
	return s.repo.ListByOrg(ctx, orgID)
}

// StartWorkOrder changes the status of a work order to "InProgress".
func (s *Service) StartWorkOrder(ctx context.Context, orgID, workOrderID uuid.UUID) error {
	wo, err := s.repo.FindByID(ctx, orgID, workOrderID)
	if err != nil {
		return err
	}
	wo.Status = StatusInProgress
	wo.StartedAt = time.Now().UTC()
	wo.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, wo)
}

// CompleteTask marks a task as completed and may trigger subsequent actions.
func (s *Service) CompleteTask(ctx context.Context, taskID uuid.UUID, result string, sn string) error {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	task.Status = StatusCompleted
	task.CompletedAt = time.Now().UTC()
	task.Result = result
	task.SN = sn

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	isFinalTask := true // Placeholder
	if isFinalTask {
		instance, err := s.productRepo.GetInstanceByID(ctx, task.ProductInstanceID)
		if err != nil {
			return fmt.Errorf("could not fetch product instance %s: %w", task.ProductInstanceID, err)
		}

		// Use the event defined in the 'product' domain.
		evt := product.ProductionFinishedEvent{
			ProductInstanceID: instance.ID,
			SerialNumber:      instance.SerialNumber,
			ProductID:         instance.ProductID,
			OrgID:             instance.OrgID,
			FinishedAt:        time.Now().UTC(),
		}

		// Publish the event using the generic event bus.
		if err := s.eventBus.Publish(ctx, &evt); err != nil {
			// In a real app, we might want to log this or handle it more gracefully.
		}
	}

	return nil
}