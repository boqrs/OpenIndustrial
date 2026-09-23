package postgres

import (
	"context"
	"errors"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"time"
)

type IoTRepository struct {
	db *database.DBProvider
}

func NewIoTRepository(db *database.DBProvider) *IoTRepository {
	return &IoTRepository{db: db}
}

func (r *IoTRepository) GetDeviceByResourceID(ctx context.Context, resourceID uint) (*model.Device, error) {
	if resourceID == 0 {
		return nil, errors.New("resource is invalid")
	}
	var device model.Device
	err := dbFromContext(ctx, r.db.Get()).WithContext(ctx).Where("resource_id = ?", resourceID).First(&device).Error
	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, iot.ErrDeviceNotFound
		// }
		return nil, err
	}
	return &device, nil
}

func (r *IoTRepository) SetOnline(ctx context.Context, resourceID uint) (*model.Device, error) {
	if resourceID == 0 {
		return nil, errors.New("resource is invalid")
	}
	now := time.Now().UTC()
	result := dbFromContext(ctx, r.db.Get()).WithContext(ctx).Model(&model.Device{}).Where("resource_id = ?", resourceID).Updates(map[string]interface{}{"status": model.DeviceStatusOnline, "last_online_at": now})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("device no found")
	}
	return r.GetDeviceByResourceID(ctx, resourceID)
}

func (r *IoTRepository) SetOffline(ctx context.Context, resourceID uint) (*model.Device, error) {
	if resourceID == 0 {
		return nil, errors.New("resource is invalid")
	}
	result := dbFromContext(ctx, r.db.Get()).WithContext(ctx).Model(&model.Device{}).Where("resource_id = ?", resourceID).Updates(map[string]interface{}{"status": model.DeviceStatusOffline})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("device no found")
	}
	return r.GetDeviceByResourceID(ctx, resourceID)
}

func (r *IoTRepository) UpdateLastOnline(ctx context.Context, resourceID uint) (*model.Device, error) {
	if resourceID == 0 {
		return nil, errors.New("resource is invalid")
	}
	now := time.Now().UTC()
	result := dbFromContext(ctx, r.db.Get()).WithContext(ctx).Model(&model.Device{}).Where("resource_id = ?", resourceID).Updates(map[string]interface{}{"status": model.DeviceStatusOnline, "last_online_at": now})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("device no found")
	}
	return r.GetDeviceByResourceID(ctx, resourceID)
}

func (r *IoTRepository) CreateCommand(
	ctx context.Context,
	cmd *model.DeviceCommand,
) error {

	return r.db.Get().WithContext(ctx).
		Create(cmd).
		Error

}

func (r *IoTRepository) GetCommand(
	ctx context.Context,
	id uint,
) (
	*model.DeviceCommand,
	error,
) {

	var cmd model.DeviceCommand

	err := r.db.Get().WithContext(ctx).
		First(
			&cmd,
			id,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &cmd, nil

}

func (r *IoTRepository) UpdateCommand(
	ctx context.Context,
	cmd *model.DeviceCommand,
) error {

	return r.db.Get().WithContext(ctx).
		Save(cmd).
		Error

}

func (r *IoTRepository) ListDeviceCommands(
	ctx context.Context,
	deviceID uint,
) (
	[]*model.DeviceCommand,
	error,
) {

	var result []*model.DeviceCommand

	err := r.db.Get().WithContext(ctx).
		Where(
			"device_id = ?",
			deviceID,
		).
		Order(
			"id desc",
		).
		Find(&result).
		Error

	return result, err

}
