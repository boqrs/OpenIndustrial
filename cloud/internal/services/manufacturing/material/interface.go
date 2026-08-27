package material

import (
	"context"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	//"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
)

// MaterialRepository defines the persistence interface for Material entities.
type Repository interface {
	Create(ctx context.Context, material *model.Material) error
	GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.Material, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*model.Material, error)
	List(ctx context.Context, tenantID uuid.UUID, offset int, limit int) ([]*model.Material, int64, error)
	Update(ctx context.Context, material *model.Material) error
	Delete(ctx context.Context, tenantID uuid.UUID, id uint) error
}

// Service 定义了物料服务的接口。
// 它封装了与物料相关的核心业务逻辑。
type Service interface {
	// Create 创建一个新的物料。
	// tenantID 用于确保数据隔离。
	Create(ctx context.Context, tenantID uuid.UUID, material *model.Material) error

	// GetByID 根据 ID 获取物料信息。
	GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.Material, error)

	// GetByCode 根据物料编码获取物料信息。
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*model.Material, error)

	// List 获取物料列表，支持分页。
	List(ctx context.Context, tenantID uuid.UUID, offset int, limit int) ([]*model.Material, int64, error)

	// Update 更新一个已存在的物料。
	Update(ctx context.Context, tenantID uuid.UUID, material *model.Material) error

	// Delete 删除一个物料。
	Delete(ctx context.Context, tenantID uuid.UUID, id uint) error
}
