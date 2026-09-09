package device

import (
	"context"
	"errors"
	"time"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource" // 正确且唯一的服务依赖
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	"github.com/google/uuid"
)

var (
	ErrDeviceNotFound         = errors.New("device not found")
	ErrProductModelNotFound   = errors.New("associated product model not found")
	ErrSerialNumberExists     = errors.New("a device with this serial number already exists")
	ErrInvalidCreateRequest   = errors.New("invalid create device request")
	ErrInvalidUpdateRequest   = errors.New("invalid update device request")
	ErrCannotDeleteOnlineDevice = errors.New("cannot delete a device that is currently online")
)

type serviceImpl struct {
	repo         Repository
	resourceSvc  resource.Service
	productSvc   product.Service
	securitySvc  security.Service
	// txManager transaction.Manager // Assuming a transaction manager exists
}

// NewService creates a new device service implementation.
func NewService(repo Repository,resourceSvc resource.Service,productSvc product.Service,securitySvc security.Service) Service {
	return &serviceImpl{
		repo:         repo,
		resourceSvc:  resourceSvc,
		productSvc:   productSvc,
		securitySvc:  securitySvc,
	}
}

// CreateDevice orchestrates the creation of a new device.
func (s *serviceImpl) CreateFromExecutionResultTx(
	ctx context.Context,
	req *CreateDeviceFromExecutionResultRequest,
) (*DeviceResponse, error) {

	if req == nil ||
		req.ProductID == 0 ||
		req.WorkOrderID == 0 ||
		req.ExecutionID == 0 ||
		req.ExecutionResultID == 0 ||
		req.SerialNumber == "" {
		return nil, ErrInvalidCreateRequest
	}

	// 1. Product must exist.
	if _, err := s.productSvc.GetProductModel(ctx, req.ProductID); err != nil {
		if errors.Is(err, product.ErrProductModelNotFound) {
			return nil, ErrProductModelNotFound
		}

		return nil, err
	}

	// 2. Serial number must be unique.
	existing, err := s.repo.GetBySerialNumber(
		ctx,
		req.SerialNumber,
	)

	if err == nil && existing != nil {
		return nil, ErrSerialNumberExists
	}

	if err != nil && !errors.Is(err, ErrDeviceNotFound) {
		return nil, err
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID not found in context")
	}

	// 3. Create Resource for the physical device.
	resourceReq := &resource.CreateResource{
		TenantID: tenantID,
		Name:     req.SerialNumber,
		Type:     string(resource.ResourceTypeDevice),
		ParentID: req.ParentResourceID,
	}

	res, err := s.resourceSvc.CreateResourceTx(ctx, resourceReq)
	if err != nil {
		return nil, err
	}

	// 4. Create Device.
	entity := &model.Device{
		ResourceID:        res.ID,
		ProductID:         req.ProductID,
		WorkOrderID:       req.WorkOrderID,
		ExecutionID:       req.ExecutionID,
		ExecutionResultID: req.ExecutionResultID,
		SerialNumber:      req.SerialNumber,
		HardwareID:        req.HardwareID,
		Status:            model.DeviceStatusCreated,
	}

	if err := s.repo.CreateTx(ctx, entity); err != nil {
		return nil, err
	}

	return s.toDeviceResponse(entity, res), nil
}


func validateCreateRequest(req *CreateDeviceFromExecutionResultRequest) error {
	if req == nil {
		return errors.New("invalid request")
	}
	if req.ProductID == 0 {
		return errors.New("invalid product ID")
	}
	if req.WorkOrderID == 0 {
		return errors.New("invalid work order ID")
	}
	if req.ExecutionID == 0 {
		return errors.New("invalid execution ID")
	}
	if req.ExecutionResultID == 0 {
		return errors.New("invalid execution result ID")
	}
	if req.SerialNumber == "" {
		return errors.New("invalid serial number")
	}
	return nil
}

func (s *serviceImpl) CreateFromExecutionResultBatchTx(
    ctx context.Context,
    reqs []*CreateDeviceFromExecutionResultRequest,
) ([]*DeviceResponse, error) {

    if len(reqs) == 0 {
        return nil, nil
    }

    resourceParams := make(
        []*resource.CreateResource,
        0,
        len(reqs),
    )

    for _, req := range reqs {
        if err := validateCreateRequest(req); err != nil {
            return nil, err
        }
        resourceParams = append(
            resourceParams,
            &resource.CreateResource{
                Type:     string(resource.ResourceTypeDevice),
                Name:     req.SerialNumber,
                ParentID: req.ParentResourceID,
            },
        )
    }

    resources, err := s.resourceSvc.CreateResourceBatchTx(ctx, resourceParams)
    if err != nil {
        return nil, fmt.Errorf(
            "create device resources: %w",
            err,
        )
    }

    devices := make(
        []*model.Device,
        0,
        len(reqs),
    )

    for i, req := range reqs {

        devices = append(devices, &model.Device{
            ResourceID:        resources[i].ID,
            ProductID:         req.ProductID,
            WorkOrderID:       req.WorkOrderID,
            ExecutionID:       req.ExecutionID,
            ExecutionResultID: req.ExecutionResultID,
            SerialNumber:      req.SerialNumber,
            HardwareID:        req.HardwareID,
            Status:            model.StatusInactive,
        })
    }

    if err := s.repo.CreateBatchTx(
        ctx,
        devices,
    ); err != nil {
        return nil, fmt.Errorf(
            "create devices: %w",
            err,
        )
    }

    responses := make(
        []*DeviceResponse,
        0,
        len(devices),
    )

    for i, d := range devices {
        responses = append(
            responses,
            s.toDeviceResponse(d, resources[i]),
        )
    }

    return responses, nil
}

func (s *serviceImpl) GetDevice(
	ctx context.Context,
	deviceID uint,
) (*DeviceResponse, error) {

	d, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if d == nil {
		return nil, ErrDeviceNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID not found in context")
	}

	res, err := s.resourceSvc.GetResourceByID(
		ctx,
		tenantID,
		d.ResourceID,
	)
	if err != nil {
		return nil, err
	}

	return s.toDeviceResponse(d, res), nil
}

func (s *serviceImpl) UpdateDevice(
	ctx context.Context,
	deviceID uint,
	req *UpdateDeviceRequest,
) (*DeviceResponse, error) {

	if req == nil {
		return nil, ErrInvalidUpdateRequest
	}

	d, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if d == nil {
		return nil, ErrDeviceNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID not found in context")
	}
	if req.Name != nil || req.ParentResourceID != nil {

		res, err := s.resourceSvc.GetResourceByID(
			ctx,
			tenantID,
			d.ResourceID,
		)
		if err != nil {
			return nil, err
		}

		if req.Name != nil {
			res.ResourceName = *req.Name
		}

		if req.ParentResourceID != nil {
			res.ParentID = *req.ParentResourceID
		}

		upReq := &resource.UpdateResource{
			Name:     res.ResourceName,
			Code:     res.Code,
			Status:   res.ResourceStatus,
			Metadata: res.Metadata,
			Version:  res.Version,
			ParentID: res.ParentID,
		}

		if _, err := s.resourceSvc.UpdateResource(
			ctx,
			res.ID,
			upReq,
		); err != nil {
			return nil, err
		}
	}

	return s.GetDevice(ctx, deviceID)
}

func (s *serviceImpl) DeleteDevice(
	ctx context.Context,
	deviceID uint,
) error {

	d, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return err
	}

	if d == nil {
		return ErrDeviceNotFound
	}

	if d.Status == model.DeviceStatusOnline {
		return ErrCannotDeleteOnlineDevice
	}

	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return errors.New("tenant ID not found in context")
	}

	if err := s.resourceSvc.DeleteResource(
		ctx,
		tenantID,
		d.ResourceID,
	); err != nil {
		return err
	}

	return s.repo.Delete(ctx, deviceID)
}

func (s *serviceImpl) toDeviceResponse(
	d *model.Device,
	r *model.Resource,
) *DeviceResponse {

	resp := &DeviceResponse{
		ID:                d.ID,
		ResourceID:        d.ResourceID,
		ProductID:         d.ProductID,
		Name:              r.ResourceName,
		SerialNumber:      d.SerialNumber,
		HardwareID:        d.HardwareID,
		WorkOrderID:       d.WorkOrderID,
		ExecutionID:       d.ExecutionID,
		ExecutionResultID: d.ExecutionResultID,
		Status:            d.Status,
		ParentResourceID:  &r.ParentID,
		CreatedAt:         d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         d.UpdatedAt.Format(time.RFC3339),
	}

	if d.LastOnlineAt != nil {
		formatted := d.LastOnlineAt.Format(time.RFC3339)
		resp.LastOnlineAt = &formatted
	}

	return resp
}

func (s *serviceImpl) ListDevices(
	ctx context.Context,
	req *ListDevicesRequest,
) (*ListDevicesResponse, error) {

	if req == nil {
		req = &ListDevicesRequest{}
	}

	if req.CurrentPage <= 0 {
		req.CurrentPage = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	items, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	responses := make([]*DeviceResponse, 0, len(items))
	tenantID := pkg.TenantIDFromContext(ctx)
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID not found in context")
	}
	for _, item := range items {
		res, err := s.resourceSvc.GetResourceByID(
			ctx,
			tenantID,
			item.ResourceID,
		)
		if err != nil {
			return nil, err
		}

		responses = append(
			responses,
			s.toDeviceResponse(item, res),
		)
	}

	return &ListDevicesResponse{
		Items:    responses,
		Total:    total,
		Page:     req.CurrentPage,
		PageSize: req.PageSize,
	}, nil
}