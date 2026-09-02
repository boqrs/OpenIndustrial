package factory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"

)

var (
	ErrFactoryNotFound = errors.New("factory not found")
	ErrFactoryCodeExists = errors.New("factory code already exists")
	ErrResourceNotFound = errors.New("resource not found")
	ErrResourceTypeMismatch = errors.New("resource type mismatch")
	ErrInvalidTopologyType = errors.New("invalid topology resource type")
	ErrInvalidParent = errors.New("invalid topology parent")
	ErrTopologyCycle = errors.New("topology operation would create a cycle")
	ErrCannotDeleteFactoryWithChildren = errors.New("cannot delete factory with children")
	ErrCannotDeleteNodeWithChildren = errors.New("cannot delete topology node with children")
	ErrNodeNotFound = errors.New("topology node not found")
)

// topologyTypes defines the allowed resource types for topology nodes in this first phase.
var topologyTypes = map[resource.ResourceType]struct{}{
	resource.ResourceTypeProductionLine: {},
	resource.ResourceTypeWorkStation:     {},
}

// serviceImpl implements the Service interface.
type serviceImpl struct {
	resourceSvc resource.Service
	repository  Repository
}

// NewService creates a new factory service.
func NewService(resourceSvc resource.Service,repository Repository) Service {
	return &serviceImpl{
		resourceSvc: resourceSvc,
		repository:  repository,
	}
}

func (s *serviceImpl) CreateFactory(ctx context.Context, req *CreateFactoryRequest) (*FactoryResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	timezone := strings.TrimSpace(req.Timezone)

	if name == "" {
		return nil, errors.New("name is required")
	}
	if code == "" {
		return nil, errors.New("code is required")
	}
	if timezone == "" {
		timezone = "UTC"
	}

	existing, err := s.repository.GetByCode(ctx, code)
	if err == nil && existing != nil {
		return nil, ErrFactoryCodeExists
	}
	if err != nil && !errors.Is(err, ErrFactoryNotFound) {
		return nil, fmt.Errorf("check factory code: %w", err)
	}

	tenantID := pkg.TenantIDFromContext(ctx) 
	resourceEntity, err := s.resourceSvc.CreateResource(ctx, &resource.CreateResource{
		TenantID: tenantID,
		Type:     string(resource.ResourceTypeFactory),
		Name:     name,
		Status:   model.StatusActive,
		Metadata: []byte(`{}`),
	})
	if err != nil {
		return nil, fmt.Errorf("create factory resource: %w", err)
	}

	entity := &model.Factory{
		ResourceID: resourceEntity.ID,
		Code:       code,
		Address:    req.Address,
		Timezone:   timezone,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		// Rollback resource creation
		_ = s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.ID)
		return nil, fmt.Errorf("create factory in repository: %w", err)
	}

	return s.factoryResponse(resourceEntity, entity), nil
}

func (s *serviceImpl) GetFactory(ctx context.Context, factoryID uint) (*FactoryResponse, error) {
	factory, err := s.repository.GetByID(ctx, factoryID)
	if err != nil {
		return nil, ErrFactoryNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx) // Placeholder

	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantID, factory.ResourceID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if resource.ResourceType(resourceEntity.ResourceType) != resource.ResourceTypeFactory {
		return nil, ErrResourceTypeMismatch
	}

	return s.factoryResponse(resourceEntity, factory), nil
}

func (s *serviceImpl) UpdateFactory(ctx context.Context, factoryID uint, req *UpdateFactoryRequest) (*FactoryResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	factory, err := s.repository.GetByID(ctx, factoryID)
	if err != nil {
		return nil, ErrFactoryNotFound
	}

	// TODO: tenantID should be properly extracted from context
	tenantID := pkg.TenantIDFromContext(ctx) // Placeholder

	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantID, factory.ResourceID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	resourceNeedsUpdate := false
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		if resourceEntity.ResourceName != name {
			resourceEntity.ResourceName = name
			resourceNeedsUpdate = true
		}
	}

	factoryNeedsUpdate := false
	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, errors.New("code cannot be empty")
		}
		if code != factory.Code {
			existing, lookupErr := s.repository.GetByCode(ctx, code)
			if lookupErr == nil && existing != nil && existing.ID != factory.ID {
				return nil, ErrFactoryCodeExists
			}
			factory.Code = code
			factoryNeedsUpdate = true
		}
	}

	if req.Address != nil && factory.Address != *req.Address {
		factory.Address = *req.Address
		factoryNeedsUpdate = true
	}

	if req.Timezone != nil {
		timezone := strings.TrimSpace(*req.Timezone)
		if timezone == "" {
			return nil, errors.New("timezone cannot be empty")
		}
		if factory.Timezone != timezone {
			factory.Timezone = timezone
			factoryNeedsUpdate = true
		}
	}

	if resourceNeedsUpdate {
		req := &resource.UpdateResource{
			Name: resourceEntity.ResourceName,
			Code: resourceEntity.Code,
			Status: resourceEntity.ResourceStatus,
			Metadata:resourceEntity.Metadata,
			Version: resourceEntity.Version,
			ParentID: resourceEntity.ParentID,
		}
		if _, err := s.resourceSvc.UpdateResource(ctx, resourceEntity.ID, req); err != nil {
			return nil, fmt.Errorf("update factory resource: %w", err)
		}
	}

	if factoryNeedsUpdate {
		if err := s.repository.Update(ctx, factory); err != nil {
			return nil, fmt.Errorf("update factory in repository: %w", err)
		}
	}

	return s.factoryResponse(resourceEntity, factory), nil
}

func (s *serviceImpl) DeleteFactory(ctx context.Context, factoryID uint) error {
	factory, err := s.repository.GetByID(ctx, factoryID)
	if err != nil {
		return ErrFactoryNotFound
	}

	children, err := s.resourceSvc.GetChildren(ctx, factory.ResourceID)
	if err != nil {
		return fmt.Errorf("get factory children: %w", err)
	}
	if len(children) > 0 {
		return ErrCannotDeleteFactoryWithChildren
	}

	if err := s.repository.Delete(ctx, factory.ID); err != nil {
		return fmt.Errorf("delete factory from repository: %w", err)
	}

	// TODO: tenantID should be properly extracted from context
	tenantID := pkg.TenantIDFromContext(ctx)
	if err := s.resourceSvc.DeleteResource(ctx, tenantID, factory.ResourceID); err != nil {
		// Note: The factory entry is already deleted, this could lead to orphaned resources.
		// A transaction would be ideal here.
		return fmt.Errorf("delete factory resource: %w", err)
	}

	return nil
}

func (s *serviceImpl) CreateTopologyNode(ctx context.Context, req *CreateTopologyNodeRequest) (*TopologyNodeResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.FactoryID == 0 {
		return nil, errors.New("factory_id is required")
	}
	if _, ok := topologyTypes[resource.ResourceType(req.Type)]; !ok {
		return nil, ErrInvalidTopologyType
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	factory, err := s.repository.GetByID(ctx, req.FactoryID)
	if err != nil {
		return nil, ErrFactoryNotFound
	}

	// TODO: tenantID should be properly extracted from context
	tenantID := pkg.TenantIDFromContext(ctx)
	parentResourceID := factory.ResourceID // Default parent is the factory itself
	if req.ParentResourceID != nil {
		parentResourceID = *req.ParentResourceID
		if err := s.validateTopologyParent(ctx, tenantID, factory.ResourceID, parentResourceID, resource.ResourceType(req.Type)); err != nil {
			return nil, err
		}
	}

	var metadata []byte
	if len(req.Metadata) > 0 {
		metadata, err = json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("marshal metadata: %w", err)
		}
	} else {
		metadata = []byte(`{}`)
	}

	resourceEntity, err := s.resourceSvc.CreateResource(ctx, &resource.CreateResource{
		TenantID:  tenantID,
		ParentID:  &parentResourceID,
		Type:      req.Type,
		Name:      name,
		Code:      optionalString(req.Code),
		Status:    model.StatusActive,
		Metadata:  metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("create topology node resource: %w", err)
	}

	return s.topologyNodeResponse(resourceEntity), nil
}

func (s *serviceImpl) UpdateTopologyNode(ctx context.Context, resourceID uint, req *UpdateTopologyNodeRequest) (*TopologyNodeResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	// TODO: tenantID should be properly extracted from context
	tenantID := pkg.TenantIDFromContext(ctx)
	entity, err := s.resourceSvc.GetResourceByID(ctx, tenantID, resourceID)
	if err != nil {
		return nil, ErrNodeNotFound
	}

	needsUpdate := false
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		if entity.ResourceName != name {
			entity.ResourceName = name
			needsUpdate = true
		}
	}

	if req.Code != nil {
		code := optionalString(*req.Code)
		if (entity.Code == nil && code != nil) || (entity.Code != nil && code == nil) || (entity.Code != nil && code != nil && *entity.Code != *code) {
			entity.Code = code
			needsUpdate = true
		}
	}

	if req.Metadata != nil {
		metadata, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("marshal metadata: %w", err)
		}
		entity.Metadata = metadata
		needsUpdate = true
	}

	if needsUpdate {
		req := &resource.UpdateResource{
			Name: entity.ResourceName,
			Code: entity.Code,
			Status: entity.ResourceStatus,
			Metadata:entity.Metadata,
			Version: entity.Version,
			ParentID: entity.ParentID,
		}
		if _, err := s.resourceSvc.UpdateResource(ctx, entity.ID, req); err != nil {
			return nil, fmt.Errorf("update topology node resource: %w", err)
		}
	}

	return s.topologyNodeResponse(entity), nil
}

func (s *serviceImpl) MoveTopologyNode(ctx context.Context, req *MoveTopologyNodeRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	tenantID := pkg.TenantIDFromContext(ctx) 
	node, err := s.resourceSvc.GetResourceByID(ctx, tenantID, req.ResourceID)
	if err != nil {
		return ErrNodeNotFound
	}
	if resource.ResourceType(node.ResourceType) == resource.ResourceTypeFactory {
		return errors.New("factory resource cannot be moved")
	}

	if req.ParentResourceID != nil {
		if *req.ParentResourceID == req.ResourceID {
			return ErrTopologyCycle
		}
		if s.createsCycle(ctx, tenantID, req.ResourceID, *req.ParentResourceID) {
			return ErrTopologyCycle
		}

		factoryResourceID, err := s.findFactoryResourceID(ctx, tenantID, req.ResourceID)
		if err != nil {
			return err
		}

		if err := s.validateTopologyParent(ctx, tenantID, factoryResourceID, *req.ParentResourceID, resource.ResourceType(node.ResourceType)); err != nil {
			return err
		}
	}

	return s.resourceSvc.UpdateParent(ctx, tenantID, req.ResourceID, *req.ParentResourceID)
}

func (s *serviceImpl) DeleteTopologyNode(ctx context.Context, resourceID uint) error {
	tenantID := pkg.TenantIDFromContext(ctx) 

	entity, err := s.resourceSvc.GetResourceByID(ctx, tenantID, resourceID)
	if err != nil {
		return ErrNodeNotFound
	}
	if resource.ResourceType(entity.ResourceType) == resource.ResourceTypeFactory {
		return errors.New("factory must be deleted through factory service")
	}

	children, err := s.resourceSvc.GetChildren(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("get children: %w", err)
	}
	if len(children) > 0 {
		return ErrCannotDeleteNodeWithChildren
	}

	if err := s.resourceSvc.DeleteResource(ctx, tenantID, resourceID); err != nil {
		return fmt.Errorf("delete topology node resource: %w", err)
	}

	return nil
}

func (s *serviceImpl) GetTopology(ctx context.Context, factoryID uint) (*FactoryTopologyResponse, error) {
	factory, err := s.repository.GetByID(ctx, factoryID)
	if err != nil {
		return nil, ErrFactoryNotFound
	}

	tenantID := pkg.TenantIDFromContext(ctx)

	factoryResource, err := s.resourceSvc.GetResourceByID(ctx, tenantID, factory.ResourceID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	nodes, err := s.collectTopology(ctx, tenantID, factory.ResourceID)
	if err != nil {
		return nil, err
	}

	return &FactoryTopologyResponse{
		Factory: *s.factoryResponse(factoryResource, factory),
		Nodes:   nodes,
	}, nil
}

// --- Helper Functions ---

func (s *serviceImpl) validateTopologyParent(ctx context.Context, tenantID uuid.UUID, factoryResourceID, parentResourceID uint, childType resource.ResourceType) error {
	parent, err := s.resourceSvc.GetResourceByID(ctx, tenantID, parentResourceID)
	if err != nil {
		return ErrInvalidParent
	}

	isUnder, err := s.isResourceUnderFactory(ctx, tenantID, parent.ID, factoryResourceID)
	if err != nil {
		return fmt.Errorf("checking parent ancestry: %w", err)
	}
	if !isUnder {
		return ErrInvalidParent
	}

	// Parent validation rules
	switch childType {
	case resource.ResourceTypeWorkStation:
		if resource.ResourceType(parent.ResourceType) != resource.ResourceTypeProductionLine {
			return ErrInvalidParent
		}
	case resource.ResourceTypeProductionLine:
		if resource.ResourceType(parent.ResourceType) != resource.ResourceTypeWorkStation {
			return ErrInvalidParent
		}
	}


	return nil
}

func (s *serviceImpl) isResourceUnderFactory(ctx context.Context, tenantID uuid.UUID, resourceID, factoryResourceID uint) (bool, error) {
	currentID := resourceID
	for {
		if currentID == factoryResourceID {
			return true, nil
		}
		current, err := s.resourceSvc.GetResourceByID(ctx, tenantID, currentID)
		if err != nil {
			return false, err
		}
		if current.ParentID == 0 {
			return false, nil
		}
		currentID = current.ParentID
	}
}

func (s *serviceImpl) findFactoryResourceID(ctx context.Context, tenantID uuid.UUID, resourceID uint) (uint, error) {
	currentID := resourceID
	for {
		current, err := s.resourceSvc.GetResourceByID(ctx, tenantID, currentID)
		if err != nil {
			return 0, err
		}
		if resource.ResourceType(current.ResourceType) == resource.ResourceTypeFactory {
			return current.ID, nil
		}
		if current.ParentID == 0 {
			return 0, ErrFactoryNotFound
		}
		currentID = current.ParentID
	}
}

func (s *serviceImpl) createsCycle(ctx context.Context, tenantID uuid.UUID, resourceID, newParentID uint) bool {
	currentID := newParentID
	for {
		if currentID == resourceID {
			return true
		}
		current, err := s.resourceSvc.GetResourceByID(ctx, tenantID, currentID)
		if err != nil {
			return true // Fail safe
		}
		if current.ParentID == 0 {
			return false
		}
		currentID = current.ParentID
	}
}

func (s *serviceImpl) collectTopology(ctx context.Context, tenantID uuid.UUID, parentID uint) ([]TopologyNodeResponse, error) {
	children, err := s.resourceSvc.GetChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("get resource children: %w", err)
	}

	var result []TopologyNodeResponse
	for _, child := range children {
		result = append(result, *s.topologyNodeResponse(child))
		descendants, err := s.collectTopology(ctx, tenantID, child.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, descendants...)
	}
	return result, nil
}

func (s *serviceImpl) factoryResponse(resourceEntity *model.Resource, factory *model.Factory) *FactoryResponse {
	return &FactoryResponse{
		ID:         factory.ID,
		ResourceID: factory.ResourceID,
		Name:       resourceEntity.ResourceName,
		Code:       factory.Code,
		Address:    factory.Address,
		Timezone:   factory.Timezone,
		Status:     string(resourceEntity.ResourceStatus),
		CreatedAt:  factory.CreatedAt,
		UpdatedAt:  factory.UpdatedAt,
	}
}

func (s *serviceImpl) topologyNodeResponse(resourceEntity *model.Resource) *TopologyNodeResponse {
	return &TopologyNodeResponse{
		ResourceID:       resourceEntity.ID,
		Type:             string(resourceEntity.ResourceType),
		Name:             resourceEntity.ResourceName,
		Status:           string(resourceEntity.ResourceStatus),
		ParentResourceID: &resourceEntity.ParentID,
		Metadata:         decodeMetadata(resourceEntity.Metadata),
	}
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func decodeMetadata(value []byte) map[string]interface{} {
	if len(value) == 0 {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(value, &result); err != nil {
		return nil // Should probably log this error
	}
	return result
}
