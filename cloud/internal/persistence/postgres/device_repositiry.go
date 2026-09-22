package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/device"
	"github.com/boqrs/nexus/database"
)

type DeviceRepository struct {
	db *database.DBProvider
}

func NewDeviceRepository(
	db *database.DBProvider,
) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) CreateTx(
	ctx context.Context,
	entity *model.Device,
) error {
	return dbFromContext(
		ctx,
		r.db.Get(),
	).
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

func (r *DeviceRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Device, error) {
	var d model.Device

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"id = ?",
			id,
		).
		First(&d).
		Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
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
		Where(
			"resource_id = ?",
			resourceID,
		).
		First(&d).
		Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
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

	err := dbFromContext(
		ctx,
		r.db.Get(),
	).
		WithContext(ctx).
		Where(
			"serial_number = ?",
			serialNumber,
		).
		First(&d).
		Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, device.ErrDeviceNotFound
		}

		return nil, err
	}

	return &d, nil
}

// ActivateBySerialNumber atomically binds a device to an end user.
//
// The update succeeds only when the device has not been activated.
//
// This prevents concurrent activation requests from assigning the
// same device to two different users.
func (r *DeviceRepository) ActivateBySerialNumber(
	ctx context.Context,
	serialNumber string,
	customerUUID uuid.UUID,
) (*model.Device, error) {
	if serialNumber == "" {
		return nil, device.ErrInvalidActivationRequest
	}

	if customerUUID == uuid.Nil {
		return nil, device.ErrUserNotAuthenticated
	}

	now := time.Now().UTC()

	db := dbFromContext(
		ctx,
		r.db.Get(),
	).WithContext(ctx)

	result := db.
		Model(&model.Device{}).
		Where(
			"serial_number = ?",
			serialNumber,
		).
		Where(
			"customer_uuid IS NULL",
		).
		Where(
			"activated_at IS NULL",
		).
		Updates(map[string]interface{}{
			"customer_uuid": customerUUID,
			"activated_at":  now,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		// Re-read the device to distinguish between:
		//
		// 1. device not found
		// 2. device already activated
		deviceEntity, err := r.GetBySerialNumber(
			ctx,
			serialNumber,
		)
		if err != nil {
			return nil, err
		}

		if deviceEntity == nil {
			return nil, device.ErrDeviceNotFound
		}

		if deviceEntity.CustomerUUID != nil ||
			deviceEntity.ActivatedAt != nil {
			return nil, device.ErrDeviceAlreadyActivated
		}

		// This should be extremely unusual. Treat it as an
		// activation failure rather than pretending activation
		// succeeded.
		return nil, device.ErrDeviceAlreadyActivated
	}

	return r.GetBySerialNumber(
		ctx,
		serialNumber,
	)
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

	if err := query.Count(
		&total,
	).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.CurrentPage - 1) * req.PageSize

	if err := query.
		Offset(offset).
		Limit(req.PageSize).
		Order(
			"devices.created_at DESC",
		).
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
