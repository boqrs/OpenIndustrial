package resource

import (
	"context"
	"errors"
	"time"

	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"

	"github.com/google/uuid"
)

// Service provides business logic for resources, including attribute management.
type Service interface{
	CreateProduct(ctx context.Context, tenantID uuid.UUID, params *param.CreateProduct) (*model.Resource, error)
	CreateResource(ctx context.Context, params *param.CreateResource) (*model.Resource, error)
	UpdateResource(ctx context.Context, resourceID uuid.UUID, req *param.UpdateResource) (*model.Resource, error)
	DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error 
	GetResource(ctx context.Context, tenantID, resourceID uuid.UUID) (*model.Resource, error)
	ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*model.Resource, error)
	SetAttribute(ctx context.Context, tenantID, resourceID uuid.UUID, attrKey string, attrValue interface{}) error 
	BatchCreateResources(ctx context.Context, resources []*model.Resource) error 
	GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*model.Resource, error)
	FindResourceByNameAndType(ctx context.Context, tenantID uuid.UUID, name, resourceType string) (*model.Resource, error)
	BatchCreateAttributeDefinition(ctx context.Context, attrs []*model.AttributeDefinition)error
	GetAttributesForResource(ctx context.Context, resourceID uuid.UUID) (map[string]interface{}, error)
	FindAttributeDefinitionByResourceID(ctx context.Context, resourceID uuid.UUID)([]*model.AttributeDefinition, error)
	BatchCreateResourceAttributes(ctx context.Context, attr []*model.ResourceAttribute) error
	UpdateParent(ctx context.Context, tenantID, resourceID, newParentID uuid.UUID) error
	CreateConnection(ctx context.Context, sourceID, tragetID uuid.UUID) error
	UpsertAttributesForResource(ctx context.Context, resourceID uuid.UUID, attributes map[string]interface{}) error 
	ClearParent(ctx context.Context, resourceID uuid.UUID) error
	GetConnection(ctx context.Context, connectionID uint) (*model.ResourceConnection, error) 
	DeleteConnection(ctx context.Context, connectionID uint) error
	GetChildren(ctx context.Context, resourceID uuid.UUID) ([]*model.Resource, error)
	ListConnections(ctx context.Context, resourceID uuid.UUID) ([]*model.ResourceConnection, error)
}
type service struct {
	resourceRepo ResourceRepository
	attrDefRepo  AttributeDefinitionRepository
	resAttrRepo  ResourceAttributeRepository
	resourceConRepo ResourceConnectionsRepository
}

// NewService creates a new resource service.
func NewService(
	resourceRepo ResourceRepository,
	attrDefRepo AttributeDefinitionRepository,
	resAttrRepo ResourceAttributeRepository,
	resourceConRepo ResourceConnectionsRepository,
) Service {
	return &service{
		resourceRepo: resourceRepo,
		attrDefRepo:  attrDefRepo,
		resAttrRepo:  resAttrRepo,
		resourceConRepo: resourceConRepo,
	}
}

// tenantIDFromContext is a placeholder helper function.
// In a real application, this would be part of a shared context utility.
func tenantIDFromContext(ctx context.Context) uuid.UUID {
	// Dummy implementation for now
	// In a real application, you would extract this from a JWT token or similar.
	val := ctx.Value("tenant_id")
	if val != nil {
		if id, ok := val.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}


// CreateProduct creates a new product resource, sets its owner, and saves its specific attributes.
// This function is PRESERVED and UPGRADED to the new architecture.
func (s *service) CreateProduct(ctx context.Context, tenantID uuid.UUID, params *param.CreateProduct) (*model.Resource, error) {
	// 1. Create the core resource object.
	// Note: Description and SerialNumber are NOT part of the core resource anymore.
	// We now correctly set the OwnerGroupID directly on the resource.
	resource := &model.Resource{
		UUID:           uuid.New(),
		TenantID:     tenantID,
		ResourceName:         params.Name,
		ResourceType:         params.Type,
		OwnerGroupID: &params.OwnerGroupID, // This is the NEW way to set ownership.
		Version: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 2. Save the core resource.
	if err := s.resourceRepo.CreateResource(ctx, resource); err != nil {
		return nil, err
	}

	// 3. Handle Description and SerialNumber as DYNAMIC ATTRIBUTES.
	// This is the "add new interface implementation" part.
	attributesToSet := []*model.ResourceAttribute{}

	// Handle Description attribute
	if params.Description != "" {
		def, err := s.attrDefRepo.GetAttributeDefinitionByKey(ctx, tenantID, "description")
		if err != nil {
			// You might want to create the definition if it doesn't exist, or return an error.
			// For now, we'll assume it exists and return an error if not.
			return nil, errors.New("attribute definition for 'description' not found")
		}
		attributesToSet = append(attributesToSet, &model.ResourceAttribute{
			ResourceID:  resource.UUID,
			ID: def.ID,
			Value: []byte(params.Description), // Assuming Value is a []byte for JSONB storage.
		})
	}

	// Handle SerialNumber attribute
	if params.SerialNumber != "" {
		def, err := s.attrDefRepo.GetAttributeDefinitionByKey(ctx, tenantID, "serial_number")
		if err != nil {
			return nil, errors.New("attribute definition for 'serial_number' not found")
		}
		attributesToSet = append(attributesToSet, &model.ResourceAttribute{
			ResourceID:  resource.UUID,
			ID: def.ID,
			Value: []byte(params.SerialNumber), // Assuming Value is a []byte for JSONB storage.
		})
	}

	// 4. Batch-save the attributes.
	if len(attributesToSet) > 0 {
		if err := s.resAttrRepo.SetAttributes(ctx, attributesToSet); err != nil {
			// In a real transaction, you'd roll back the resource creation here.
			return nil, err
		}
	}

	return resource, nil
}

// CreateResource creates a new, generic resource.
func (s *service) CreateResource(ctx context.Context, params *param.CreateResource) (*model.Resource, error) {
	resource := &model.Resource{
		UUID:            uuid.New(),
		TenantID:      params.TenantID,
		ResourceType:          params.Type,
		ResourceName:          params.Name,
		Code:          params.Code,
		ResourceStatus:        params.Status,
		Metadata:      params.Metadata,
		ParentID:      *params.ParentID,
		OwnerGroupID:  params.OwnerGroupID,
		Version: 1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.resourceRepo.CreateResource(ctx, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateResource updates an existing resource. It uses optimistic locking.
func (s *service) UpdateResource(ctx context.Context, resourceID uuid.UUID, req *param.UpdateResource) (*model.Resource, error) {
	// 1. Get the existing resource to ensure it exists and to get its current state.
	res, err := s.resourceRepo.GetResourceByID(ctx, req.TenantID, resourceID)
	if err != nil {
		return nil, err // Handles not found, etc.
	}

	// 2. Optimistic locking check.
	if res.Version != req.Version {
		return nil, errors.New("update conflict: resource has been modified by another process")
	}

	// 3. Apply changes from params.
	res.ResourceName = req.Name
	res.Code = req.Code
	res.ResourceStatus = req.Status
	res.Metadata = req.Metadata
	res.UpdatedAt = time.Now()
	// The repository will handle incrementing the version.

	// 4. Persist the changes.
	if err := s.resourceRepo.UpdateResource(ctx, res); err != nil {
		return nil, err
	}

	// The resource object 'res' now has an old version number.
	// We should return the resource with the *new* version number.
	res.Version++

	return res, nil
}

// DeleteResource performs a soft delete on a resource.
func (s *service) DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error {
	// We could add a check here to ensure the resource exists before deleting.
	// For now, we delegate this to the repository.
	return s.resourceRepo.DeleteResource(ctx, tenantID, resourceID)
}

// GetResource retrieves a single resource by its ID.
// This can be enhanced later to also fetch and compose its attributes.
func (s *service) GetResource(ctx context.Context, tenantID, resourceID uuid.UUID) (*model.Resource, error) {
	// TODO: Add authorization check here.
	return s.resourceRepo.GetResourceByID(ctx, tenantID, resourceID)
}

// ListResources retrieves a list of resources for a tenant, with filtering and pagination.
// This is the corrected implementation.
func (s *service) ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*model.Resource, error) {
	// TODO: Add authorization filtering here.
	if limit <= 0 {
		limit = 100 // Default limit
	}
	if offset < 0 {
		offset = 0 // Default offset
	}
	return s.resourceRepo.ListResources(ctx, tenantID, resourceType, limit, offset)
}

// SetAttribute is a new function that properly uses the new interfaces.
func (s *service) SetAttribute(ctx context.Context, tenantID, resourceID uuid.UUID, attrKey string, attrValue interface{}) error {
	// 1. 编排步骤一：验证业务规则
	// 在更新属性之前，先确认这个属性的“定义”是否存在。
	// 这是一个核心的业务规则：不允许设置未定义的属性。
	// 我们调用 `attrDefRepo` 的接口来完成这个验证。
	_, err := s.attrDefRepo.FindByName(ctx, tenantID, attrKey)
	if err != nil {
		// 如果 FindByName 返回错误（比如 ErrNotFound），说明该属性定义不存在，
		// 直接将错误返回给上层，阻止后续操作。
		return err
	}

	// 2. 编排步骤二：准备数据
	// 因为 `resAttrRepo` 的 `UpsertForResource` 接口接收的是一个 map，
	// 所以我们在这里创建一个只包含一个键值对的 map。
	attributes := map[string]interface{}{
		attrKey: attrValue,
	}

	// 3. 编排步骤三：委托执行
	// 将准备好的数据（tenantID, resourceID, attributes map）传递给 `resAttrRepo`。
	// Service 层的任务到此结束。它不关心这个属性是如何被存入数据库的，
	// 不关心 JSON 转换，不关心是 INSERT 还是 UPDATE。
	// 这一切都由 `resAttrRepo` 的具体实现去负责。
	return s.resAttrRepo.UpsertForResource(ctx, tenantID, resourceID, attributes)
}

func (s *service) BatchCreateResources(ctx context.Context, resources []*model.Resource) error {
	return s.resourceRepo.BatchCreateResources(ctx, resources)
}

func (s *service) GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*model.Resource, error) {
	return s.resourceRepo.GetResourceByID(ctx, tenantID, resourceID)
}

func (s *service) FindResourceByNameAndType(ctx context.Context, tenantID uuid.UUID, name, resourceType string) (*model.Resource, error) {
	return s.resourceRepo.FindResourceByNameAndType(ctx, tenantID, name, resourceType)
}


func (s *service) BatchCreateAttributeDefinition(ctx context.Context, attrs []*model.AttributeDefinition)error{
		return s.attrDefRepo.BatchCreateAttributeDefinition(ctx, attrs)
}

func (s *service) GetAttributesForResource(ctx context.Context, resourceID uuid.UUID) (map[string]interface{}, error) {
	return s.resAttrRepo.GetAttributesForResource(ctx, tenantIDFromContext(ctx), resourceID)
}


func (s *service)	FindAttributeDefinitionByResourceID(ctx context.Context, resourceID uuid.UUID)([]*model.AttributeDefinition, error){
	return s.attrDefRepo.FindAttributeDefinitionByResourceID(ctx, resourceID)
} //TODO: 需要实现底层{}

func (s *service)	BatchCreateResourceAttributes(ctx context.Context, attr []*model.ResourceAttribute) error{
	return s.resAttrRepo.BatchCreateResourceAttributes(ctx, attr)
}


	// UpdateParent changes the hierarchical parent of a given resource.
func (s *service)UpdateParent(ctx context.Context, tenantID, resourceID, newParentID uuid.UUID) error{
	return s.resourceRepo.UpdateParent(ctx, tenantID, resourceID, newParentID)
}

func (s *service) CreateConnection(ctx context.Context, sourceID, tragetID uuid.UUID) error{
	res := &model.ResourceConnection{
		SourceResourceID: sourceID ,
		TargetResourceID: tragetID,
		ConnectionType: model.ConnectionTypeConnectedThrough, //TODO: 这里需要参数传入
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.resourceConRepo.CreateConnection(ctx, res)
}
// UpsertAttributesForResource updates or inserts a batch of attributes for a given resource.
func (s *service) UpsertAttributesForResource(ctx context.Context, resourceID uuid.UUID, attributes map[string]interface{}) error {
	// In a real implementation, you would first validate the attributes against their definitions.
	// For now, we delegate directly to the repository.
	return s.resAttrRepo.UpsertForResource(ctx, tenantIDFromContext(ctx), resourceID, attributes)
}

// ClearParent removes the parent from a resource, making it a root-level resource.
func (s *service) ClearParent(ctx context.Context, resourceID uuid.UUID) error {
	return s.resourceRepo.UpdateParent(ctx, tenantIDFromContext(ctx), resourceID, uuid.Nil)
}


func (s *service) GetConnection(ctx context.Context, connectionID uint) (*model.ResourceConnection, error) {
	return s.resourceConRepo.GetConnectionByID(ctx, connectionID)
}

func (s *service) DeleteConnection(ctx context.Context, connectionID uint) error {
	return s.resourceConRepo.DeleteConnection(ctx, connectionID)
}

func (s *service) GetChildren(ctx context.Context, resourceID uuid.UUID) ([]*model.Resource, error) {
	return s.resourceRepo.FindByParentID(ctx, tenantIDFromContext(ctx), resourceID)
}
func (s *service) ListConnections(ctx context.Context, resourceID uuid.UUID) ([]*model.ResourceConnection, error) {
	return s.resourceConRepo.ListConnectionsByResourceID(ctx, resourceID)
}
// CreateConnection establishes a new technical connection between two resources.

/*
NOTE ON ListUserGroups:

The function 'ListUserGroups' has been intentionally removed from this service.
Its responsibility is to answer "Which groups does a user belong to?". This is a core
question for the **Identity Kernel**.

Keeping it here would violate our architectural principle of separating the Resource Kernel
(the "what") from the Identity Kernel (the "who").

This function's logic should be moved to an 'identity.Service' which operates on the
'identity.GroupRepository'.
*/