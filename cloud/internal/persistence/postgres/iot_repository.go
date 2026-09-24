package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
)

type IoTRepository struct {
	db *database.DBProvider
}

func NewIoTRepository(
	db *database.DBProvider,
) *IoTRepository {

	return &IoTRepository{
		db: db,
	}
}

func (r *IoTRepository) GetDeviceByID(
	ctx context.Context,
	deviceID uint,
) (*model.Device, error) {

	if deviceID == 0 {
		return nil, errors.New(
			"device is invalid",
		)
	}

	var device model.Device

	err := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		First(
			&device,
			deviceID,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *IoTRepository) GetDeviceByResourceID(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {

	if resourceID == 0 {
		return nil, errors.New(
			"resource is invalid",
		)
	}

	var device model.Device

	err := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Where(
			"resource_id = ?",
			resourceID,
		).
		First(
			&device,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *IoTRepository) SetOnline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {

	if resourceID == 0 {
		return nil, errors.New(
			"resource is invalid",
		)
	}

	now := time.Now().UTC()

	result := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Model(&model.Device{}).
		Where(
			"resource_id = ?",
			resourceID,
		).
		Updates(map[string]interface{}{
			"connection_status": model.ConnectionStatusConnected,

			"last_online_at": now,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	// if result.RowsAffected == 0 {
	// 	return nil, ErrDeviceNotFound
	// }

	return r.GetDeviceByResourceID(
		ctx,
		resourceID,
	)
}

func (r *IoTRepository) SetOffline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {

	if resourceID == 0 {
		return nil, errors.New(
			"resource is invalid",
		)
	}

	result := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Model(&model.Device{}).
		Where(
			"resource_id = ?",
			resourceID,
		).
		Update(
			"connection_status",
			model.ConnectionStatusDisconnected,
		)

	if result.Error != nil {
		return nil, result.Error
	}

	// if result.RowsAffected == 0 {
	// 	return nil, ErrDeviceNotFound
	// }

	return r.GetDeviceByResourceID(
		ctx,
		resourceID,
	)
}

func (r *IoTRepository) UpdateLastOnline(
	ctx context.Context,
	resourceID uint,
) (*model.Device, error) {

	if resourceID == 0 {
		return nil, errors.New(
			"resource is invalid",
		)
	}

	now := time.Now().UTC()

	result := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Model(&model.Device{}).
		Where(
			"resource_id = ?",
			resourceID,
		).
		Updates(map[string]interface{}{
			"connection_status": model.ConnectionStatusConnected,

			"last_online_at": now,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	// if result.RowsAffected == 0 {
	// 	return nil, ErrDeviceNotFound
	// }

	return r.GetDeviceByResourceID(
		ctx,
		resourceID,
	)
}

func (r *IoTRepository) CreateCommand(
	ctx context.Context,
	cmd *model.DeviceCommand,
) error {

	if cmd == nil {
		return errors.New(
			"command is nil",
		)
	}

	return dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Create(cmd).
		Error
}

func (r *IoTRepository) GetCommand(
	ctx context.Context,
	id uint,
) (*model.DeviceCommand, error) {

	if id == 0 {
		return nil, errors.New(
			"command is invalid",
		)
	}

	var cmd model.DeviceCommand

	err := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
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

	if cmd == nil {
		return errors.New(
			"command is nil",
		)
	}

	return dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Save(cmd).
		Error
}

func (r *IoTRepository) ListDeviceCommands(
	ctx context.Context,
	deviceID uint,
) ([]*model.DeviceCommand, error) {

	if deviceID == 0 {
		return nil, errors.New(
			"device is invalid",
		)
	}

	var result []*model.DeviceCommand

	err := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Where(
			"device_id = ?",
			deviceID,
		).
		Order(
			"id desc",
		).
		Find(
			&result,
		).
		Error

	return result, err
}
