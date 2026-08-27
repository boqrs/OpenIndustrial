package bom

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context,bom *model.BOM) error
	GetByID(ctx context.Context,tenantID uuid.UUID,id uint) (*model.BOM, error)
	GetByNoVersion(ctx context.Context,tenantID uuid.UUID,bomNo string,version int) (*model.BOM, error)
	List(ctx context.Context,tenantID uuid.UUID,productID uuid.UUID,offset int,limit int) ([]*model.BOM, int64, error)
	Update(ctx context.Context,bom *model.BOM) error
	CreateItems(ctx context.Context,items []*model.BOMItem) error
	GetItems(ctx context.Context,tenantID uuid.UUID,bomID uint) ([]*model.BOMItem, error)
	DeleteItems(ctx context.Context,tenantID uuid.UUID,bomID uint) error
}

type Service interface {
	Create(ctx context.Context,tenantID uuid.UUID,req *CreateRequest) (*Response, error)
	GetByID(ctx context.Context,tenantID uuid.UUID,id uint) (*Response, error)
	List(ctx context.Context,tenantID uuid.UUID,productID uuid.UUID,offset int,limit int) ([]*Response, int64, error)
	Update(ctx context.Context,tenantID uuid.UUID,id uint,req *UpdateRequest) (*Response, error)
	Release(ctx context.Context,tenantID uuid.UUID,id uint) (*Response, error)
	Obsolete(ctx context.Context,tenantID uuid.UUID,id uint) error
}