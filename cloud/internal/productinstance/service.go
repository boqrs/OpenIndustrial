package productinstance

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req CreateRequest,
) (*ProductInstance, error) {
	instance := &ProductInstance{
		ID:        uuid.New().String(),
		SN:        req.SN,
		ProductID: req.ProductID,
		OrgID:     req.OrgID,
		State:     "CREATED",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(
		ctx,
		instance,
	)
	if err != nil {
		return nil, err
	}

	return instance, nil
}