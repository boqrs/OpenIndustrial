package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
)

type DeviceRepository struct {
	db *database.DBProvider
}

// NewDeviceRepository creates a new GORM-based repository for devices.
func NewDeviceRepository(db *database.DBProvider) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(ctx context.Context, entity *model.Device) error {
	return r.db.Get().WithContext(ctx).Create(entity).Error
}

func (r *DeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	var d model.Device
	if err := r.db.Get().WithContext(ctx).Where("id = ?", id).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.Device, error) {
	var d model.Device
	if err := r.db.Get().WithContext(ctx).Where("resource_id = ?", resourceID).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) GetBySerialNumber(ctx context.Context, serialNumber string) (*model.Device, error) {
	var d model.Device
	if err := r.db.Get().WithContext(ctx).Where("serial_number = ?", serialNumber).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *DeviceRepository) List(ctx context.Context, req *device.ListDevicesRequest) ([]*model.Device, int64, error) {
	var items []*model.Device
	var total int64

	query := r.db.Get().WithContext(ctx).Model(&model.Device{})

	if req.ProductModelID != nil {
		query = query.Where("product_model_id = ?", *req.ProductModelID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	// Note: Filtering by parent requires a JOIN with the resources table.
	if req.ParentID != nil {
		query = query.Joins("JOIN resources ON resources.uuid = devices.resource_id").
			Where("resources.parent_id = ?", *req.ParentID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *DeviceRepository) Update(ctx context.Context, entity *model.Device) error {
	return r.db.Get().WithContext(ctx).Save(entity).Error
}

func (r *DeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Get().WithContext(ctx).Where("id = ?", id).Delete(&model.Device{}).Error
}