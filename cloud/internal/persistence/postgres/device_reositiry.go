package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/OpenIndustrial/cloud/internal/device"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)

// --- DeviceTypeRepository Implementation ---

type deviceTypeRepository struct {
	db *gorm.DB
}

// NewDeviceTypeRepository creates a new GORM-based repository for device types.
func NewDeviceTypeRepository(db *gorm.DB) device.DeviceTypeRepository {
	return &deviceTypeRepository{db: db}
}

func (r *deviceTypeRepository) Create(ctx context.Context, entity *model.DeviceType) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *deviceTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.DeviceType, error) {
	var dt model.DeviceType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceTypeNotFound
		}
		return nil, err
	}
	return &dt, nil
}

func (r *deviceTypeRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.DeviceType, error) {
	var dt model.DeviceType
	err := r.db.WithContext(ctx).Where("resource_id = ?", resourceID).First(&dt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceTypeNotFound
		}
		return nil, err
	}
	return &dt, nil
}

func (r *deviceTypeRepository) GetByCode(ctx context.Context, code string) (*model.DeviceType, error) {
	var dt model.DeviceType
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&dt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceTypeNotFound
		}
		return nil, err
	}
	return &dt, nil
}

func (r *deviceTypeRepository) List(ctx context.Context) ([]*model.DeviceType, error) {
	var dts []*model.DeviceType
	err := r.db.WithContext(ctx).Find(&dts).Error
	return dts, err
}

func (r *deviceTypeRepository) Update(ctx context.Context, entity *model.DeviceType) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// --- DeviceRepository Implementation ---

type deviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository creates a new GORM-based repository for devices.
func NewDeviceRepository(db *gorm.DB) device.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, entity *model.Device) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *deviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	var d model.Device
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&d).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *deviceRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID) (*model.Device, error) {
	var d model.Device
	err := r.db.WithContext(ctx).Where("resource_id = ?", resourceID).First(&d).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, device.ErrDeviceNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *deviceRepository) List(ctx context.Context, deviceTypeID *uuid.UUID) ([]*model.Device, error) {
	var devices []*model.Device
	query := r.db.WithContext(ctx)
	if deviceTypeID != nil {
		query = query.Where("device_type_id = ?", *deviceTypeID)
	}
	err := query.Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) Update(ctx context.Context, entity *model.Device) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *deviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Device{}).Error
}