package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/OpenIndustrial/cloud/internal/param"


)

type Repository interface {
	Create(ctx context.Context,entity *model.ProductModel) error
	GetByID(ctx context.Context,id uuid.UUID) (*model.ProductModel, error)
	GetByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.ProductModel, error)
	GetByCodeAndVersion(ctx context.Context,code string,version string) (*model.ProductModel, error)
	List(ctx context.Context,req param.ListProductModelsRequest) ([]*model.ProductModel, int64, error)
	Update(ctx context.Context,entity *model.ProductModel) error
	Delete(ctx context.Context,id uuid.UUID) error
}