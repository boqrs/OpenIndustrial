package bom

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context, bom *model.BOM) error
	GetByID(ctx context.Context, tenantID uint, id uint) (*model.BOM, error)
	GetByNoVersion(ctx context.Context, tenantID uint, bomNo string, version int) (*model.BOM, error)
	List(ctx context.Context, tenantID uint, productID uint, offset int, limit int) ([]*model.BOM, int64, error)
	Update(ctx context.Context, bom *model.BOM) error
	CreateItems(ctx context.Context, items []*model.BOMItem) error
	GetItems(ctx context.Context, tenantID uint, bomID uint) ([]*model.BOMItem, error)
	DeleteItems(ctx context.Context, tenantID uint, bomID uint) error
}

type Service interface {
	Create(ctx context.Context, tenantID uint, req *CreateRequest) (*Response, error)
	GetByID(ctx context.Context, tenantID uint, id uint) (*Response, error)
	List(ctx context.Context, tenantID uint, productID uint, offset int, limit int) ([]*Response, int64, error)
	Update(ctx context.Context, tenantID uint, id uint, req *UpdateRequest) (*Response, error)
	Release(ctx context.Context, tenantID uint, id uint) (*Response, error)
	Obsolete(ctx context.Context, tenantID uint, id uint) error
}
