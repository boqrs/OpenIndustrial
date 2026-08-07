package workorder

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrWorkOrderNotFound is returned when a work order is not found.
	ErrWorkOrderNotFound = errors.New("work order not found")
	// ErrTaskNotFound is returned when a task is not found.
	ErrTaskNotFound = errors.New("task not found")
)

// Repository defines the persistence interface for workorder related entities.
type Repository interface {
	Create(ctx context.Context, wo *WorkOrder) error
	Update(ctx context.Context, wo *WorkOrder) error
	FindByID(ctx context.Context, orgID, workOrderID string) (*WorkOrder, error)
	ListByOrg(ctx context.Context, orgID string) ([]*WorkOrder, error)
	GetTask(ctx context.Context, taskID string) (*StationTask, error)
	UpdateTask(ctx context.Context, task *StationTask) error
}

// memoryRepository is an in-memory implementation of the Repository interface.
type memoryRepository struct {
	mu         sync.RWMutex
	workorders map[string]*WorkOrder
	tasks      map[string]*StationTask
}

// NewMemoryRepository creates a new in-memory workorder repository.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		workorders: make(map[string]*WorkOrder),
		tasks:      make(map[string]*StationTask),
	}
}

// Create saves a new work order.
func (r *memoryRepository) Create(ctx context.Context, wo *WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workorders[wo.ID] = wo
	return nil
}

// Update updates an existing work order.
func (r *memoryRepository) Update(ctx context.Context, wo *WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.workorders[wo.ID]; !ok {
		return ErrWorkOrderNotFound
	}
	r.workorders[wo.ID] = wo
	return nil
}

// FindByID retrieves a work order by its ID.
func (r *memoryRepository) FindByID(ctx context.Context, orgID, workOrderID string) (*WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if wo, ok := r.workorders[workOrderID]; ok && wo.OrgID == orgID {
		return wo, nil
	}
	return nil, ErrWorkOrderNotFound
}

// ListByOrg lists all work orders for an organization.
func (r *memoryRepository) ListByOrg(ctx context.Context, orgID string) ([]*WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var results []*WorkOrder
	for _, wo := range r.workorders {
		if wo.OrgID == orgID {
			results = append(results, wo)
		}
	}
	return results, nil
}

// GetTask retrieves a task by its ID.
func (r *memoryRepository) GetTask(ctx context.Context, taskID string) (*StationTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if task, ok := r.tasks[taskID]; ok {
		return task, nil
	}
	return nil, ErrTaskNotFound
}

// UpdateTask updates an existing task.
func (r *memoryRepository) UpdateTask(ctx context.Context, task *StationTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[task.ID]; !ok {
		return ErrTaskNotFound
	}
	r.tasks[task.ID] = task
	return nil
}