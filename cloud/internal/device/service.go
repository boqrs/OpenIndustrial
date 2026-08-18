package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/OpenIndustrial/cloud/internal/kernel/resource" // 正确且唯一的服务依赖
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrDeviceTypeNotFound    = errors.New("device type not found")
	ErrDeviceTypeCodeExists  = errors.New("device type code already exists")
	ErrDeviceNotFound        = errors.New("device not found")
	ErrResourceNotFound      = errors.New("resource not found")
	ErrInvalidDeviceResource = errors.New("resource is not a device")
	ErrInvalidParentResource = errors.New("invalid parent resource")
	ErrDeviceAlreadyAttached = errors.New("device already attached")
	ErrConnectionExists      = errors.New("resource connection already exists")
	ErrInvalidConnection     = errors.New("invalid resource connection")
)

type Service interface {
	// Device Type
	CreateDeviceType(ctx context.Context, req *param.CreateDeviceTypeRequest) (*param.DeviceTypeResponse, error)
	GetDeviceType(ctx context.Context, id uuid.UUID) (*param.DeviceTypeResponse, error)
	ListDeviceTypes(ctx context.Context) ([]param.DeviceTypeResponse, error)
	UpdateDeviceType(ctx context.Context, id uuid.UUID, req *param.UpdateDeviceTypeRequest) (*param.DeviceTypeResponse, error)

	// Device
	CreateDevice(ctx context.Context, req *param.CreateDeviceRequest) (*param.DeviceResponse, error)
	GetDevice(ctx context.Context, deviceID uuid.UUID) (*param.DeviceResponse, error)
	ListDevices(ctx context.Context, deviceTypeID *uuid.UUID) ([]*model.Device, error)
	UpdateDevice(ctx context.Context, deviceID uuid.UUID, req *param.UpdateDeviceRequest) (*param.DeviceResponse, error)
	DeleteDevice(ctx context.Context, deviceID uuid.UUID) error

	// Topology
	AttachDevice(ctx context.Context, deviceID uuid.UUID, req *param.AttachDeviceRequest) error
	DetachDevice(ctx context.Context, deviceID uuid.UUID) error
	ConnectDevice(ctx context.Context, deviceID uuid.UUID, req *param.ConnectDeviceRequest) error
	DisconnectDevice(ctx context.Context, deviceID uuid.UUID, connectionID uint) error
	GetTopology(ctx context.Context, deviceID uuid.UUID) (*param.DeviceTopologyResponse, error)
}

type serviceImpl struct {
	resourceSvc    resource.Service
	deviceTypeRepo DeviceTypeRepository
	deviceRepo     DeviceRepository
}

func NewService(
	resourceSvc resource.Service,
	deviceTypeRepo DeviceTypeRepository,
	deviceRepo DeviceRepository,
) Service {
	return &serviceImpl{
		resourceSvc:    resourceSvc,
		deviceTypeRepo: deviceTypeRepo,
		deviceRepo:     deviceRepo,
	}
}

// tenantIDFromContext is a placeholder for getting the tenant ID from the context.
// This should be replaced with the actual implementation from your auth/context logic.
func tenantIDFromContext(ctx context.Context) uuid.UUID {
	// In a real application, you would extract this from a JWT token or similar.
	val := ctx.Value("tenant_id")
	if val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id
		}
	}
	// Fallback for testing or unauthenticated contexts
	return uuid.Nil
}

// --- DeviceType Methods ---

func (s *serviceImpl) CreateDeviceType(ctx context.Context, req *param.CreateDeviceTypeRequest) (*param.DeviceTypeResponse, error) {
	if req.Code == "" || req.Name == "" {
		return nil, errors.New("device type name and code are required")
	}

	existing, err := s.deviceTypeRepo.GetByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check for existing device type: %w", err)
	}
	if existing != nil {
		return nil, ErrDeviceTypeCodeExists
	}

	// 1. Create the underlying Resource
	resReq := &param.CreateResource{
		TenantID: tenantIDFromContext(ctx),
		Name:     req.Name,
		Type:     string(resource.ResourceTypeDeviceType),
	}
	res, err := s.resourceSvc.CreateResource(ctx, resReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource for device type: %w", err)
	}

	// 2. Create the DeviceType extension
	deviceType := &model.DeviceType{
		ID:          uuid.New(),
		ResourceID:  res.UUID,
		Code:        req.Code,
		Category:    req.Category,
		Description: req.Description,
		Enabled:     true,
	}

	if err := s.deviceTypeRepo.Create(ctx, deviceType); err != nil {
		// Rollback: delete the created resource
		_ = s.resourceSvc.DeleteResource(ctx, res.TenantID, res.UUID)
		return nil, fmt.Errorf("failed to create device type extension: %w", err)
	}

	// 3. Create AttributeDefinitions
	if len(req.Attributes) > 0 {
		definitions := make([]*model.AttributeDefinition, 0, len(req.Attributes))
		for name, defReq := range req.Attributes {
			definitions = append(definitions, &model.AttributeDefinition{
				UUID:        uuid.New(),
				ResourceID:  res.UUID, // Link definition to the DeviceType's resource
				Name:        name,
				Label:       defReq.Label,
				Description: defReq.Description,
				DataType:    model.AttributeValueType(defReq.DataType),
				Unit:        defReq.Unit,
				//Required:    defReq.Required,
			})
		}
		if err := s.resourceSvc.BatchCreateAttributeDefinition(ctx, definitions); err != nil {
			// NOTE: In a production system, a more robust transaction/rollback mechanism would be needed here.
			return nil, fmt.Errorf("failed to create attribute definitions: %w", err)
		}
	}

	return s.buildDeviceTypeResponse(ctx, deviceType)
}

func (s *serviceImpl) GetDeviceType(ctx context.Context, id uuid.UUID) (*param.DeviceTypeResponse, error) {
	deviceType, err := s.deviceTypeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceTypeNotFound
		}
		return nil, fmt.Errorf("failed to get device type: %w", err)
	}
	return s.buildDeviceTypeResponse(ctx, deviceType)
}

func (s *serviceImpl) ListDeviceTypes(ctx context.Context) ([]param.DeviceTypeResponse, error) {
	items, err := s.deviceTypeRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list device types: %w", err)
	}

	responses := make([]param.DeviceTypeResponse, 0, len(items))
	for _, item := range items {
		resp, err := s.buildDeviceTypeResponse(ctx, item)
		if err != nil {
			continue
		}
		responses = append(responses, *resp)
	}
	return responses, nil
}

func (s *serviceImpl) UpdateDeviceType(ctx context.Context, id uuid.UUID, req *param.UpdateDeviceTypeRequest) (*param.DeviceTypeResponse, error) {
	deviceType, err := s.deviceTypeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDeviceTypeNotFound
	}

	if req.Name != nil && *req.Name != "" {
		updateReq := &param.UpdateResource{Name: *req.Name}
		if _, err := s.resourceSvc.UpdateResource(ctx, deviceType.ResourceID, updateReq); err != nil {
			return nil, fmt.Errorf("failed to update device type resource name: %w", err)
		}
	}

	if req.Description != nil {
		deviceType.Description = *req.Description
	}
	if req.Enabled != nil {
		deviceType.Enabled = *req.Enabled
	}
	if err := s.deviceTypeRepo.Update(ctx, deviceType); err != nil {
		return nil, fmt.Errorf("failed to update device type: %w", err)
	}

	return s.buildDeviceTypeResponse(ctx, deviceType)
}

// --- Device Methods ---

func (s *serviceImpl) CreateDevice(ctx context.Context, req *param.CreateDeviceRequest) (*param.DeviceResponse, error) {
	if req.DeviceTypeID == uuid.Nil || req.Name == "" {
		return nil, errors.New("device_type_id and name are required")
	}

	deviceType, err := s.deviceTypeRepo.GetByID(ctx, req.DeviceTypeID)
	if err != nil {
		return nil, ErrDeviceTypeNotFound
	}
	definitions, err := s.resourceSvc.FindAttributeDefinitionByResourceID(ctx, deviceType.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute definitions for device type: %w", err)
	}

	resReq := &param.CreateResource{
		TenantID: tenantIDFromContext(ctx),
		Name:     req.Name,
		Type:     string(resource.ResourceTypeDevice),
		ParentID: req.ParentResourceID,
		Code:     &req.Code,
	}
	res, err := s.resourceSvc.CreateResource(ctx, resReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource for device: %w", err)
	}

	device := &model.Device{
		ID:           uuid.New(),
		ResourceID:   res.UUID,
		DeviceTypeID: req.DeviceTypeID,
	}
	if err := s.deviceRepo.Create(ctx, device); err != nil {
		_ = s.resourceSvc.DeleteResource(ctx, res.TenantID,res.UUID)
		return nil, fmt.Errorf("failed to create device extension: %w", err)
	}

	if len(req.Attributes) > 0 {
		attributes, err := validateAndPrepareAttributes(req.Attributes, definitions, res.UUID)
		if err != nil {
			_ = s.resourceSvc.DeleteResource(ctx, res.TenantID,res.UUID)
			return nil, err
		}
		if err := s.resourceSvc.BatchCreateResourceAttributes(ctx, attributes); err != nil {
			_ = s.resourceSvc.DeleteResource(ctx, res.TenantID,res.UUID)
			return nil, fmt.Errorf("failed to create device attributes: %w", err)
		}
	}

	return s.buildDeviceResponse(ctx, device)
}

func (s *serviceImpl) GetDevice(ctx context.Context, deviceID uuid.UUID) (*param.DeviceResponse, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return s.buildDeviceResponse(ctx, device)
}

func (s *serviceImpl) ListDevices(ctx context.Context, deviceTypeID *uuid.UUID) ([]*model.Device, error) {
	devices, err := s.deviceRepo.List(ctx, deviceTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	return devices, nil
}

func (s *serviceImpl) UpdateDevice(ctx context.Context, deviceID uuid.UUID, req *param.UpdateDeviceRequest) (*param.DeviceResponse, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	updateReq := &param.UpdateResource{
		Name:     *req.Name,
		Code:     req.Code,
		ParentID: *req.ParentResourceID,
	}
	if _, err := s.resourceSvc.UpdateResource(ctx, device.ResourceID, updateReq); err != nil {
		return nil, fmt.Errorf("failed to update device resource: %w", err)
	}


	if len(req.Attributes) > 0 {
		if err := s.resourceSvc.UpsertAttributesForResource(ctx, device.ResourceID, req.Attributes); err != nil {
			return nil, fmt.Errorf("failed to upsert device attributes: %w", err)
		}
	}

	return s.buildDeviceResponse(ctx, device)
}

func (s *serviceImpl) DeleteDevice(ctx context.Context, deviceID uuid.UUID) error {
	tenantID := tenantIDFromContext(ctx)

	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	if err := s.deviceRepo.Delete(ctx, device.ID); err != nil {
		return fmt.Errorf("failed to delete device extension: %w", err)
	}

	if err := s.resourceSvc.DeleteResource(ctx, tenantID, device.ResourceID); err != nil {
		return fmt.Errorf("failed to delete device resource: %w", err)
	}

	return nil
}

// --- Topology Methods ---

func (s *serviceImpl) AttachDevice(ctx context.Context, deviceID uuid.UUID, req *param.AttachDeviceRequest) error {
	tenantID := tenantIDFromContext(ctx)
	
	if req == nil || req.ParentResourceID == uuid.Nil {
		return errors.New("parent_resource_id is required")
	}

	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	if _, err := s.resourceSvc.GetResource(ctx, tenantID, req.ParentResourceID); err != nil {
		return ErrInvalidParentResource
	}

	if device.ResourceID == req.ParentResourceID {
		return errors.New("device cannot attach to itself")
	}

	updateReq := &param.UpdateResource{ParentID: req.ParentResourceID}
	if _, err := s.resourceSvc.UpdateResource(ctx, device.ResourceID, updateReq); err != nil {
		return fmt.Errorf("failed to attach device: %w", err)
	}

	return nil
}

func (s *serviceImpl) DetachDevice(ctx context.Context, deviceID uuid.UUID) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	if err := s.resourceSvc.ClearParent(ctx, device.ResourceID); err != nil {
		return fmt.Errorf("failed to detach device: %w", err)
	}

	return nil
}

func (s *serviceImpl) ConnectDevice(ctx context.Context, deviceID uuid.UUID, req *param.ConnectDeviceRequest) error {
	tenantID := tenantIDFromContext(ctx)
	
	if req == nil || req.TargetResourceID == uuid.Nil || strings.TrimSpace(req.ConnectionType) == "" {
		return errors.New("target_resource_id and connection_type are required")
	}

	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	if _, err := s.resourceSvc.GetResource(ctx, tenantID, req.TargetResourceID); err != nil {
		return ErrResourceNotFound
	}

	if device.ResourceID == req.TargetResourceID {
		return ErrInvalidConnection
	}

	if err := s.resourceSvc.CreateConnection(ctx, device.ResourceID, req.TargetResourceID); err != nil {
		return fmt.Errorf("failed to create device connection: %w", err)
	}

	return nil
}

func (s *serviceImpl) DisconnectDevice(ctx context.Context, deviceID uuid.UUID, connectionID uint) error {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return ErrDeviceNotFound
	}

	conn, err := s.resourceSvc.GetConnection(ctx, connectionID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	if conn.SourceResourceID != device.ResourceID && conn.TargetResourceID != device.ResourceID {
		return ErrInvalidConnection
	}

	if err := s.resourceSvc.DeleteConnection(ctx, connectionID); err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}

	return nil
}

func (s *serviceImpl) GetTopology(ctx context.Context, deviceID uuid.UUID) (*param.DeviceTopologyResponse, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	deviceResp, err := s.buildDeviceResponse(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to build main device response for topology: %w", err)
	}

	childrenResources, err := s.resourceSvc.GetChildren(ctx, device.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topology children: %w", err)
	}

	children := make([]param.DeviceResponse, 0, len(childrenResources))
	for _, childRes := range childrenResources {
		childDevice, err := s.deviceRepo.GetByResourceID(ctx, childRes.UUID)
		if err == nil && childDevice != nil {
			childResp, err := s.buildDeviceResponse(ctx, childDevice)
			if err == nil && childResp != nil {
				children = append(children, *childResp)
			}
		}
	}

	connections, err := s.resourceSvc.ListConnections(ctx, device.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topology connections: %w", err)
	}
	connResponses := make([]param.DeviceConnectionResponse, 0, len(connections))
	for _, conn := range connections {
		var meta map[string]interface{}
		if len(conn.Metadata) > 0 {
			_ = json.Unmarshal(conn.Metadata, &meta)
		}
		connResponses = append(connResponses, param.DeviceConnectionResponse{
			ID:               conn.ID,
			SourceResourceID: conn.SourceResourceID,
			TargetResourceID: conn.TargetResourceID,
			ConnectionType:   string(conn.ConnectionType),
			Metadata:         meta,
		})
	}

	return &param.DeviceTopologyResponse{
		Device:      *deviceResp,
		Children:    children,
		Connections: connResponses,
	}, nil
}

// --- Helper Methods ---

func (s *serviceImpl) buildDeviceTypeResponse(ctx context.Context, deviceType *model.DeviceType) (*param.DeviceTypeResponse, error) {
	tenantID := tenantIDFromContext(ctx)

	
	res, err := s.resourceSvc.GetResource(ctx, tenantID, deviceType.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("resource not found for device type %s: %w", deviceType.ID, err)
	}

	return &param.DeviceTypeResponse{
		ID:          deviceType.ID,
		ResourceID:  deviceType.ResourceID,
		Name:        res.ResourceName,
		Code:        deviceType.Code,
		Category:    deviceType.Category,
		Description: deviceType.Description,
		Enabled:     deviceType.Enabled,
		CreatedAt:   deviceType.CreatedAt,
		UpdatedAt:   deviceType.UpdatedAt,
	}, nil
}

func (s *serviceImpl) buildDeviceResponse(ctx context.Context, device *model.Device) (*param.DeviceResponse, error) {
	tenantID := tenantIDFromContext(ctx)

	
	res, err := s.resourceSvc.GetResource(ctx, tenantID, device.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("resource not found for device %s: %w", device.ID, err)
	}

	deviceTypeResp, err := s.GetDeviceType(ctx, device.DeviceTypeID)
	if err != nil {
		return nil, fmt.Errorf("device type not found for device %s: %w", device.ID, err)
	}

	attrs, err := s.resourceSvc.GetAttributesForResource(ctx, device.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attributes for device %s: %w", device.ID, err)
	}

	return &param.DeviceResponse{
		ID:               device.ID,
		ResourceID:       device.ResourceID,
		DeviceType:       *deviceTypeResp,
		Name:             res.ResourceName,
		Code:             *res.Code,
		Status:           res.ResourceStatus,
		ParentResourceID: &res.ParentID,
		Attributes:       attrs,
		CreatedAt:        device.CreatedAt,
		UpdatedAt:        device.UpdatedAt,
	}, nil
}

func validateAndPrepareAttributes(attrs map[string]interface{},definitions []*model.AttributeDefinition,resourceID uuid.UUID) ([]*model.ResourceAttribute, error) {
	defMap := make(map[string]*model.AttributeDefinition, len(definitions))
	for _, d := range definitions {
		defMap[d.Name] = d
	}

	result := make([]*model.ResourceAttribute, 0, len(attrs))
	for name, value := range attrs {
		def, ok := defMap[name]
		if !ok {
			return nil, fmt.Errorf("attribute %q is not defined for this device type", name)
		}

		// TODO: Add strict data type validation based on def.DataType

		valueJSON, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal attribute %q: %w", name, err)
		}

		result = append(result, &model.ResourceAttribute{
			ResourceID:            resourceID,
			AttributeDefinitionID: def.UUID,
			Value:                 valueJSON,
		})
	}

	for _, def := range definitions {
		if _, exists := attrs[def.Name]; !exists {
			return nil, fmt.Errorf("missing required attribute %q", def.Name)
		}
	}

	return result, nil
}