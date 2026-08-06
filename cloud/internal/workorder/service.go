package workorder

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service encapsulates the business logic for the work order domain.
type Service struct {
	repo Repository
}

// NewService creates a new work order service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateWorkOrder handles the business logic of creating a new work order.
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

// StartWorkOrder marks a work order as "In_Progress".
func (s *Service) StartWorkOrder(ctx context.Context, orgID, workOrderID uuid.UUID) (*WorkOrder, error) {
	wo, err := s.repo.FindByID(ctx, orgID, workOrderID)
	if err != nil {
		return nil, err
	}

	// Add business logic here, e.g., check if status is "Pending"
	wo.Status = "In_Progress"
	wo.StartedAt = time.Now().UTC()
	wo.UpdatedAt = wo.StartedAt

	if err := s.repo.Update(ctx, wo); err != nil {
		return nil, err
	}
	return wo, nil
}