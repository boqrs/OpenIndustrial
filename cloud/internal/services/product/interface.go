package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"


)

type Repository interface {
	Create(ctx context.Context,entity *model.ProductModel) error
	GetByID(ctx context.Context,id uint) (*model.ProductModel, error)
	GetByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.ProductModel, error)
	GetByCodeAndVersion(ctx context.Context,code string,version string) (*model.ProductModel, error)
	List(ctx context.Context,req ListProductModelsRequest) ([]*model.ProductModel, int64, error)
	Update(ctx context.Context,entity *model.ProductModel) error
	Delete(ctx context.Context,id uint) error
}

type Service interface {
	// =========================================================
	// Product Model
	// =========================================================

	CreateProductModel(ctx context.Context,req *CreateProductModelRequest,) (*CreateProductModelResponse, error)
	GetProductModel(ctx context.Context,id uint) (*ProductDetailResponse, error)
	ListProductModels(ctx context.Context,req *ListProductModelsRequest) (*ProductModelListResponse, error)
	UpdateProductModel(ctx context.Context,id uint,req *UpdateProductModelRequest) (*UpdateProductModelResponse, error)
	//TODO： 改为resource定义的状态类型
	UpdateProductModelStatus(ctx context.Context,id uint,status string) error

	// =========================================================
	// Product Model Attribute Definition
	// =========================================================

	GetAttributeDefinitions(ctx context.Context, productModelID uint) ([]AttributeDefinitionResponse, error)

	UpdateAttributeDefinitions(ctx context.Context,productModelID uint,req *UpdateAttributeDefinitionsRequest) error
}
