package factory

import (
	"context"

	"github.com/google/uuid"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context,entity *model.Factory) error
	GetByUUID(ctx context.Context,id uuid.UUID) (*model.Factory, error)
	GetByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.Factory, error)
	GetByCode(ctx context.Context,code string) (*model.Factory, error)
	Update(	ctx context.Context,entity *model.Factory) error
	Delete(ctx context.Context,id uuid.UUID) error
}