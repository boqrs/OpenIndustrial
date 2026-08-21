package device

import (
	"context"
	"errors"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource" // 正确且唯一的服务依赖
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/product"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
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
func (s *serviceImpl) CreateDevice(ctx context.Context, req *CreateDeviceRequest) (*BootstrapCredentialResponse, error) {
	if req == nil || req.Name == "" || req.ProductModelID == uuid.Nil {
		return nil, ErrInvalidCreateRequest
	}

	// 1. Validate ProductModel exists
	if _, err := s.productSvc.GetProductModel(ctx, req.ProductModelID); err != nil {
		if errors.Is(err, product.ErrProductModelNotFound) {
			return nil, ErrProductModelNotFound
		}
		return nil, err
	}

	// 2. Check for duplicate serial number if provided
	if req.SerialNumber != "" {
		_, err := s.repo.GetBySerialNumber(ctx, req.SerialNumber)
		if err == nil {
			return nil, ErrSerialNumberExists
		}
		if !errors.Is(err, ErrDeviceNotFound) {
			return nil, err // Handle unexpected repository errors
		}
	}

	var createdResource *model.Resource
	var createdDevice *model.Device

	// TODO: Use a real transaction manager to wrap these operations
	// tx, err := s.txManager.Begin(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// defer tx.Rollback()

	// 3. Create the core resource
	resourceReq := &resource.CreateResource{
		TenantID: tenantIDFromContext(ctx), // Assuming this helper exists
		Name:     req.Name,
		Type:     string(resource.ResourceTypeDevice),
		ParentID: req.ParentResourceID,
	}
	res, err := s.resourceSvc.CreateResource(ctx, resourceReq)
	if err != nil {
		return nil, err
	}
	createdResource = res

	// 4. Create the device entity
	newDevice := &model.Device{
		ID:             uuid.New(),
		ResourceID:     createdResource.UUID,
		ProductModelID: req.ProductModelID,
		SerialNumber:   req.SerialNumber,
		HardwareID:     req.HardwareID,
		Status:         model.DeviceStatusCreated,
	}
	if err := s.repo.Create(ctx, newDevice); err != nil {
		// Rollback resource creation would happen here
		return nil, err
	}
	createdDevice = newDevice

	var bootstrapCred *security.BootstrapCredentialResponse
	cbReq := security.CreateBootstrapCredentialRequest{
		ResourceID: createdResource.UUID,
	}
	bootstrapCred, err = s.securitySvc.CreateBootstrapCredential(ctx, cbReq)
	if err != nil {
		// Rollback would happen here
		return nil, err
	}
	return &BootstrapCredentialResponse{
		ResourceID:      createdDevice.ResourceID,
		CredentialID:    bootstrapCred.CredentialID,
		Token: bootstrapCred.Token,
		CreatedAt: bootstrapCred.CreatedAt,
	}, nil
}

func (s *serviceImpl) GetDevice(ctx context.Context, deviceID uuid.UUID) (*DeviceResponse, error) {
	d, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	res, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), d.ResourceID)
	if err != nil {
		// If resource is not found, it's an inconsistent state, but we can still return the device info
		return nil, err
	}

	return s.toDeviceResponse(d, res), nil
}

func (s *serviceImpl) ListDevices(ctx context.Context, req *ListDevicesRequest) (*ListDevicesResponse, error) {
	// Default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	items, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	responses := make([]*DeviceResponse, 0, len(items))
	for _, item := range items {
		// In a real-world scenario, fetching each resource individually is inefficient (N+1 problem).
		// A better approach would be to get all resource IDs and fetch them in a single batch call.
		res, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), item.ResourceID)
		if err != nil {
			// Log the error and skip this item, or handle as needed
			continue
		}
		responses = append(responses, s.toDeviceResponse(item, res))
	}

	return &ListDevicesResponse{
		Items:      responses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (s *serviceImpl) UpdateDevice(ctx context.Context, deviceID uuid.UUID, req *UpdateDeviceRequest) (*DeviceResponse, error) {
	d, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// Update resource if needed
	if req.Name != nil || req.ParentResourceID != nil {
		res, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), d.ResourceID)
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
			Name: res.ResourceName,
			Code: res.Code,
			Status: res.ResourceStatus,
			Metadata: res.Metadata,
			Version: res.Version,
			ParentID: res.ParentID,
		}

		if _, err := s.resourceSvc.UpdateResource(ctx, res.UUID, upReq); err != nil {
			return nil, err
		}
	}

	// Refetch the device and its resource to return the latest state
	return s.GetDevice(ctx, deviceID)
}

func (s *serviceImpl) DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	d, err := s.repo.GetByID(ctx, deviceID)
	if err != nil {
		return err
	}

	if d.Status == model.DeviceStatusOnline {
		return ErrCannotDeleteOnlineDevice
	}
	
	// TODO: Use a real transaction manager
	// Delete the resource, which should cascade or be handled appropriately
	if err := s.resourceSvc.DeleteResource(ctx, tenantIDFromContext(ctx), d.ResourceID); err != nil {
		return err
	}
	// Delete the device entity
	if err := s.repo.Delete(ctx, deviceID); err != nil {
		return err
	}
	// Invalidate/delete credentials
	if err := s.securitySvc.RevokeBootstrapCredential(ctx, d.ResourceID); err != nil {
		// Log this error, as the primary deletion succeeded
	}

	return nil
}

// toDeviceResponse is a helper to map the domain model to the DTO.
func (s *serviceImpl) toDeviceResponse(d *model.Device, r *model.Resource) *DeviceResponse {
	resp := &DeviceResponse{
		ID:               d.ID,
		ResourceID:       d.ResourceID,
		ProductModelID:   d.ProductModelID,
		Name:             r.ResourceName,
		SerialNumber:     d.SerialNumber,
		HardwareID:       d.HardwareID,
		Status:           d.Status,
		ParentResourceID: &r.ParentID,
		CreatedAt:        d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        d.UpdatedAt.Format(time.RFC3339),
	}
	if d.LastOnlineAt != nil {
		formatted := d.LastOnlineAt.Format(time.RFC3339)
		resp.LastOnlineAt = &formatted
	}
	return resp
}

// tenantIDFromContext is a placeholder for a helper function that extracts tenant ID from context.
func tenantIDFromContext(ctx context.Context) uuid.UUID {
	// In a real implementation, this would come from a JWT or other auth middleware.
	val := ctx.Value("tenant_id")
	if val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id
		}
	}
	// Return a default for now, but this should be handled properly.
	return uuid.Nil 
}
