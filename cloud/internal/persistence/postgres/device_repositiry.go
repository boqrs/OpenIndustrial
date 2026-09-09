package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/nexus/database"
)

type DeviceRepository struct {
	db *database.DBProvider
}

func NewDeviceRepository(db *database.DBProvider) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) CreateTx(
	ctx context.Context,
	entity *model.Device,
) error {
	return dbFromContext(ctx, r.db.Get()). //TODO： 设备创建是在生产阶段的事物中完成的
						WithContext(ctx).
						Create(entity).
						Error
}

func (r *DeviceRepository) Create(
	ctx context.Context,
	entity *model.Device,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(entity).
		Error
}

func (r *DeviceRepository) CreateBatchTx(
	ctx context.Context,
	devices []*model.Device,
) error {

	if len(devices) == 0 {
		return nil
	}

	return dbFromContext(ctx, r.db.Get()).CreateInBatches(&devices, len(devices)).
		Error
}

func (r *DeviceRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Device, error) {
	var d model.Device

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}

		return nil, err
	}

	return &d, nil
}

func (r *DeviceRepository) GetByResourceID(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {
	var d model.Device

	err := r.db.Get().
		WithContext(ctx).
		Where("resource_id = ?", resourceID).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}

		return nil, err
	}

	return &d, nil
}

func (r *DeviceRepository) GetBySerialNumber(
	ctx context.Context,
	serialNumber string,
) (*model.Device, error) {
	var d model.Device

	err := r.db.Get().
		WithContext(ctx).
		Where("serial_number = ?", serialNumber).
		First(&d).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}

		return nil, err
	}

	return &d, nil
}

func (r *DeviceRepository) List(
	ctx context.Context,
	req *device.ListDevicesRequest,
) ([]*model.Device, int64, error) {

	var items []*model.Device
	var total int64

	query := r.db.Get().
		WithContext(ctx).
		Model(&model.Device{})

	if req.ProductID != nil {
		query = query.Where(
			"product_id = ?",
			*req.ProductID,
		)
	}

	if req.Status != nil {
		query = query.Where(
			"status = ?",
			*req.Status,
		)
	}

	if req.ParentID != nil {
		query = query.
			Joins(
				"JOIN resources ON resources.id = devices.resource_id",
			).
			Where(
				"resources.parent_id = ?",
				*req.ParentID,
			)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.CurrentPage - 1) * req.PageSize

	if err := query.
		Offset(offset).
		Limit(req.PageSize).
		Order("devices.created_at DESC").
		Find(&items).
		Error; err != nil {

		return nil, 0, err
	}

	return items, total, nil
}

func (r *DeviceRepository) Update(
	ctx context.Context,
	entity *model.Device,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(entity).
		Error
}

func (r *DeviceRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	return r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Device{}).
		Error
}
