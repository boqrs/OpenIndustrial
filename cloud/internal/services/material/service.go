package material

import (
	"context"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// service 是 Service 接口的实现。
type service struct {
	repo Repository
}

// NewService 创建一个新的物料服务实例。
// 这是实现依赖注入的关键，它接收一个 Repository 接口作为参数。
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create 实现了创建物料的逻辑。
// 在调用 repository 之前，它会为物料对象设置正确的 tenantID。
func (s *service) Create(ctx context.Context, tenantID uuid.UUID, material *model.Material) error {
	material.TenantID = tenantID
	return s.repo.Create(ctx, material)
}

// GetByID 直接调用 repository 的同名方法来获取数据。
func (s *service) GetByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.Material, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetByCode 直接调用 repository 的同名方法来获取数据。
func (s *service) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*model.Material, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// List 直接调用 repository 的同名方法来获取数据列表。
func (s *service) List(ctx context.Context, tenantID uuid.UUID, offset int, limit int) ([]*model.Material, int64, error) {
	return s.repo.List(ctx, tenantID, offset, limit)
}

// Update 实现了更新物料的逻辑。
// 同样，它会先设置 tenantID 以确保操作的安全性。
func (s *service) Update(ctx context.Context, tenantID uuid.UUID, material *model.Material) error {
	material.TenantID = tenantID
	return s.repo.Update(ctx, material)
}

// Delete 直接调用 repository 的同名方法来删除数据。
func (s *service) Delete(ctx context.Context, tenantID uuid.UUID, id uint) error {
	return s.repo.Delete(ctx, tenantID, id)
}