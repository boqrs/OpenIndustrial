package mes

import (
	"context"
	"fmt"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/product"
)

// Service provides business logic for the Manufacturing Execution System.
type Service struct {
	repo        Repository
	productRepo product.Repository
}

// NewService creates a new MES service.
func NewService(repo Repository, productRepo product.Repository) *Service {
	return &Service{
		repo:        repo,
		productRepo: productRepo,
	}
}

// GetCurrentTaskForProductSN retrieves the current manufacturing task for a given product serial number.
func (s *Service) GetCurrentTaskForProductSN(ctx context.Context, sn string) (*Task, error) {
	// First, find the product instance by its serial number.
	instance, err := s.productRepo.GetInstanceBySN(ctx, sn)
	if err != nil {
		return nil, fmt.Errorf("failed to find product instance by sn %s: %w", sn, err)
	}

	// Then, find the current task associated with that product instance ID.
	task, err := s.repo.GetCurrentTaskForProduct(ctx, instance.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get current task for product id %s: %w", instance.ID, err)
	}

	return task, nil
}