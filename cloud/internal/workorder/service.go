package workorder

import (
	"context"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/event"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/product" // Import product package
)

const TaskFinished = "finished"

// Service defines the business logic for the workorder domain.
type Service interface {
	FinishTask(ctx context.Context, taskID string, result string) error
	FindByID(ctx context.Context, orgID, workOrderID string) (*WorkOrder, error)
	ListByOrg(ctx context.Context, orgID string) ([]*WorkOrder, error)
	CreateWorkOrder(ctx context.Context, orgID, productID string, quantity int, scheduledAt time.Time) (*WorkOrder, error)
	StartWorkOrder(ctx context.Context, orgID, workOrderID string) error
}

type service struct {
	repo     Repository
	eventBus event.Bus
}

// NewService creates a new workorder service.
func NewService(repo Repository, bus event.Bus) Service {
	return &service{
		repo:     repo,
		eventBus: bus,
	}
}

// CreateWorkOrder creates a new work order.
func (s *service) CreateWorkOrder(ctx context.Context, orgID, productID string, quantity int, scheduledAt time.Time) (*WorkOrder, error) {
	wo, err := NewWorkOrder(orgID, productID, quantity, scheduledAt)
	if err != nil {
		return nil, err
	}
	return wo, s.repo.Create(ctx, wo)
}

// StartWorkOrder starts a work order.
func (s *service) StartWorkOrder(ctx context.Context, orgID, workOrderID string) error {
	wo, err := s.repo.FindByID(ctx, orgID, workOrderID)
	if err != nil {
		return err
	}
	wo.Status = "IN_PROGRESS"
	wo.StartedAt = time.Now()
	return s.repo.Update(ctx, wo)
}

// FindByID retrieves a work order by its ID.
func (s *service) FindByID(ctx context.Context, orgID, workOrderID string) (*WorkOrder, error) {
	return s.repo.FindByID(ctx, orgID, workOrderID)
}

// ListByOrg lists all work orders for an organization.
func (s *service) ListByOrg(ctx context.Context, orgID string) ([]*WorkOrder, error) {
	return s.repo.ListByOrg(ctx, orgID)
}

// FinishTask marks a task as finished and publishes a production finished event.
func (s *service) FinishTask(ctx context.Context, taskID string, result string) error {
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	task.Status = TaskFinished

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	// Publish the domain event using the concrete type from the product package
	e := product.ProductionFinishedEvent{
		WorkOrderID:       task.WorkOrderID,
		ProductInstanceID: task.ProductInstanceID,
		SN:                task.SN,
		StationID:         task.StationID,
		Result:            result,
		FinishedAt:        time.Now(),
	}

	return s.eventBus.Publish(e.ToDomainEvent())
}