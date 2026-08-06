package workorder

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for data access operations for a WorkOrder.
type Repository interface {
	Create(ctx context.Context, wo *WorkOrder) error
	FindByID(ctx context.Context, orgID, workOrderID uuid.UUID) (*WorkOrder, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]*WorkOrder, error)
	Update(ctx context.Context, wo *WorkOrder) error
}