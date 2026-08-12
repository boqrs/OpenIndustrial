package device

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenIndustrial/cloud/internal/kernel/resource" // 正确且唯一的服务依赖
	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

const (
	ResourceTypeProductModel = "PRODUCT_MODEL"
	ResourceTypeIoTProduct   = "IOT_PRODUCT"
	ResourceTypeFactoryAsset = "FACTORY_ASSET"
)

type Service interface{
 CreateProductModel(ctx context.Context, req *param.CreateProductModelRequest) (*model.Resource, error)
 RegisterIoTProduct(ctx context.Context, req *param.RegisterDeviceRequest) (*param.ResourceResponse, error)
 RegisterFactoryAsset(ctx context.Context, req *param.RegisterDeviceRequest) (*param.ResourceResponse, error) 
}

// serviceImpl 是 device.Service 的具体实现
// 它现在只依赖于 resource.Service，这是整个系统数据操作的唯一入口
type serviceImpl struct {
	resourceSvc *resource.Service // 正确: 唯一的服务依赖
}

// NewService 创建一个新的 device 服务实例
func NewService(resourceSvc *resource.Service) Service {
	return &serviceImpl{
		resourceSvc: resourceSvc,
	}
}

// CreateProductModel 通过 resource 服务创建产品模型及其属性定义
func (s *serviceImpl) CreateProductModel(ctx context.Context, req *param.CreateProductModelRequest) (*model.Resource, error) {
	// 1. 通过 resource 服务创建产品模型的 Resource 记录
	productModel, err := s.resourceSvc.CreateResource(ctx, &param.CreateResource{
		TenantID: req.TenantID,
		Type:     ResourceTypeProductModel,
		Name:     req.Name,
		Status:   "active",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource for product model: %w", err)
	}

	// 2. 准备属性定义，并通过 resource 服务进行批量创建
	if len(req.Attributes) > 0 {
		definitions := make([]*model.AttributeDefinition, 0, len(req.Attributes))
		for name, attr := range req.Attributes {
			definitions = append(definitions, &model.AttributeDefinition{
				UUID: uuid.New(),
				ResourceID: productModel.UUID, //
				Name:           name,
				Label:          attr.Label,
				DataType:       model.AttributeValueType(attr.DataType),
				Unit:           attr.Unit,
			})
		}
		// 调用 resourceSvc 的方法来处理持久化
		if err := s.resourceSvc.BatchCreateAttributeDefinition(ctx, definitions); err != nil {
			return nil, fmt.Errorf("failed to batch create attribute definitions via resource service: %w", err)
		}
	}

	return productModel, nil
}

func (s *serviceImpl) RegisterIoTProduct(ctx context.Context, req *param.RegisterDeviceRequest) (*param.ResourceResponse, error) {
	return s.registerDevice(ctx, req, ResourceTypeIoTProduct)
}

func (s *serviceImpl) RegisterFactoryAsset(ctx context.Context, req *param.RegisterDeviceRequest) (*param.ResourceResponse, error) {
	return s.registerDevice(ctx, req, ResourceTypeFactoryAsset)
}

// registerDevice 是注册设备实例的通用逻辑
func (s *serviceImpl) registerDevice(ctx context.Context, req *param.RegisterDeviceRequest, resourceType string) (*param.ResourceResponse, error) {
	// 1. 通过 resource 服务验证产品模型是否存在
	_, err := s.resourceSvc.GetResourceByID(ctx, req.TenantID, req.ProductModelID)
	if err != nil {
		return nil, fmt.Errorf("product model with ID %s not found: %w", req.ProductModelID, err)
	}

	// 2. 通过 resource 服务获取产品模型允许的属性定义
	definitions, err := s.resourceSvc.FindAttributeDefinitionByResourceID(ctx, req.ProductModelID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attribute definitions for model %s: %w", req.ProductModelID, err)
	}


	// 4. 通过 resource 服务创建设备实例的 Resource 记录
	deviceInstance, err := s.resourceSvc.CreateResource(ctx, &param.CreateResource{
		TenantID: req.TenantID,
		Type:     resourceType,
		Name:     req.InstanceName,
		Status:   model.StatusProvisioned,
		ParentID: &req.ProductModelID,
		Code:     &req.SerialNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource for device instance: %w", err)
	}

		// 3. 校验传入的属性值，并准备批量创建
	attributesToCreate, err := s.validateAndPrepareAttributes(req.Attributes, definitions, deviceInstance.UUID)
	if err != nil {
		return nil, err
	}

	// 5. 批量创建设备实例的属性值
	if len(attributesToCreate) > 0 {
		for _, attr := range attributesToCreate {
			attr.ResourceID = deviceInstance.UUID // 关联到新创建的设备实例
		}
		// 调用 resourceSvc 的方法来处理持久化
		if err := s.resourceSvc.BatchCreateResourceAttributes(ctx, attributesToCreate); err != nil {
			return nil, fmt.Errorf("failed to batch create resource attributes via resource service: %w", err)
		}
	}

	return &param.ResourceResponse{
		Resource:   deviceInstance,
		Attributes: req.Attributes,
	}, nil
}

func (s *serviceImpl) validateAndPrepareAttributes(attrs map[string]interface{}, definitions []*model.AttributeDefinition, resourceID uuid.UUID) ([]*model.ResourceAttribute, error) {
	defMap := make(map[string]*model.AttributeDefinition, len(definitions))
	for _, def := range definitions {
		defMap[def.Name] = def
	}

	attributesToCreate := make([]*model.ResourceAttribute, 0, len(attrs))
	for key, value := range attrs {
		def, ok := defMap[key]
		if !ok {
			return nil, fmt.Errorf("attribute '%s' is not defined for this product model", key)
		}

		valueBytes, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize value for attribute '%s': %w", key, err)
		}

		attributesToCreate = append(attributesToCreate, &model.ResourceAttribute{
			ResourceID: resourceID,
			AttributeDefinitionID: def.UUID,
			Value:                 valueBytes,
		})
	}
	return attributesToCreate, nil
}