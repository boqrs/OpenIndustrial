package factory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	zlog "github.com/boqrs/nexus/log"
)

var (
	ErrFactoryNotFound                 = errors.New("factory not found")
	ErrFactoryCodeExists               = errors.New("factory code already exists")
	ErrResourceNotFound                = errors.New("resource not found")
	ErrResourceTypeMismatch            = errors.New("resource type mismatch")
	ErrInvalidTopologyType             = errors.New("invalid topology resource type")
	ErrInvalidParent                   = errors.New("invalid topology parent")
	ErrTopologyCycle                   = errors.New("topology operation would create a cycle")
	ErrCannotDeleteFactoryWithChildren = errors.New("cannot delete factory with children")
	ErrCannotDeleteNodeWithChildren    = errors.New("cannot delete topology node with children")
	ErrNodeNotFound                    = errors.New("topology node not found")
)

var topologyTypes = map[resource.ResourceType]struct{}{
	resource.ResourceTypeProductionLine: {},
	resource.ResourceTypeWorkStation:    {},
}

type serviceImpl struct {
	resourceSvc resource.Service
	repository  Repository
	l           *zlog.Provider
}

func NewService(
	resourceSvc resource.Service,
	repository Repository,
	l *zlog.Provider,
) Service {
	return &serviceImpl{
		resourceSvc: resourceSvc,
		repository:  repository,
		l:           l,
	}
}

// getFactoryForTenant ensures that a factory belongs to the current tenant.
func (s *serviceImpl) getFactoryForTenant(
	ctx context.Context,
	factoryID uint,
) (*model.Factory, uint, error) {
	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil, 0, errors.New("tenant id not found")
	}

	entity, err := s.repository.GetByID(ctx, factoryID)
	if err != nil {
		if errors.Is(err, ErrFactoryNotFound) {
			return nil, 0, ErrFactoryNotFound
		}
		return nil, 0, fmt.Errorf("get factory: %w", err)
	}

	if entity.TenantID != tenantID {
		return nil, 0, ErrFactoryNotFound
	}

	return entity, tenantID, nil
}

func (s *serviceImpl) getFactoryResource(
	ctx context.Context,
	tenantID uint,
	entity *model.Factory,
) (*model.Resource, error) {
	res, err := s.resourceSvc.GetResourceByID(
		ctx,
		tenantID,
		entity.ResourceID,
	)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if resource.ResourceType(res.ResourceType) != resource.ResourceTypeFactory {
		return nil, ErrResourceTypeMismatch
	}

	return res, nil
}

func (s *serviceImpl) CreateFactory(
	ctx context.Context,
	req *CreateFactoryRequest,
) (*FactoryResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil,
			errors.New("tenant id not found")
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

	existing, err := s.repository.GetByCode(ctx, tenantID, code)
	if err == nil && existing != nil {
		return nil, ErrFactoryCodeExists
	}
	if err != nil && !errors.Is(err, ErrFactoryNotFound) {
		return nil, fmt.Errorf("check factory code: %w", err)
	}

	resourceEntity, err := s.resourceSvc.CreateResource(
		ctx,
		&resource.CreateResource{
			TenantID: tenantID,
			Type:     string(resource.ResourceTypeFactory),
			Name:     name,
			Status:   model.StatusActive,
			Metadata: []byte(`{}`),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create factory resource: %w", err)
	}

	entity := &model.Factory{
		TenantID:   tenantID,
		ResourceID: resourceEntity.ID,
		Name:       req.Name,
		Code:       code,
		Address:    strings.TrimSpace(req.Address),
		Timezone:   timezone,
	}

	if err := s.repository.Create(ctx, entity); err != nil {
		// Best-effort cleanup if the Factory row cannot be created.
		_ = s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.ID)
		return nil, fmt.Errorf("create factory record: %w", err)
	}

	return s.factoryResponse(resourceEntity, entity), nil
}

func (s *serviceImpl) GetFactory(
	ctx context.Context,
	factoryID uint,
) (*FactoryResponse, error) {
	entity, tenantID, err := s.getFactoryForTenant(ctx, factoryID)
	if err != nil {
		return nil, err
	}

	res, err := s.getFactoryResource(ctx, tenantID, entity)
	if err != nil {
		return nil, err
	}

	return s.factoryResponse(res, entity), nil
}

func (s *serviceImpl) UpdateFactory(
	ctx context.Context,
	factoryID uint,
	req *UpdateFactoryRequest,
) (*FactoryResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	entity, tenantID, err := s.getFactoryForTenant(ctx, factoryID)
	if err != nil {
		return nil, err
	}

	res, err := s.getFactoryResource(ctx, tenantID, entity)
	if err != nil {
		return nil, err
	}

	resourceNeedsUpdate := false
	factoryNeedsUpdate := false

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		if res.ResourceName != name {
			res.ResourceName = name
			resourceNeedsUpdate = true
		}
	}

	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, errors.New("code cannot be empty")
		}

		if code != entity.Code {
			existing, lookupErr := s.repository.GetByCode(
				ctx,
				tenantID,
				code,
			)
			if lookupErr == nil && existing.ID != entity.ID {
				return nil, ErrFactoryCodeExists
			}
			if lookupErr != nil && !errors.Is(lookupErr, ErrFactoryNotFound) {
				return nil, fmt.Errorf("check factory code: %w", lookupErr)
			}

			entity.Code = code
			factoryNeedsUpdate = true
		}
	}

	if req.Address != nil {
		address := strings.TrimSpace(*req.Address)
		if entity.Address != address {
			entity.Address = address
			factoryNeedsUpdate = true
		}
	}

	if req.Timezone != nil {
		timezone := strings.TrimSpace(*req.Timezone)
		if timezone == "" {
			return nil, errors.New("timezone cannot be empty")
		}
		if entity.Timezone != timezone {
			entity.Timezone = timezone
			factoryNeedsUpdate = true
		}
	}

	if resourceNeedsUpdate {
		updateReq := &resource.UpdateResource{
			Name:     res.ResourceName,
			Code:     res.Code,
			Status:   res.ResourceStatus,
			Metadata: res.Metadata,
			Version:  res.Version,
			ParentID: res.ParentID,
		}
		updated, updateErr := s.resourceSvc.UpdateResource(
			ctx,
			res.ID,
			updateReq,
		)
		if updateErr != nil {
			return nil, fmt.Errorf("update factory resource: %w", updateErr)
		}
		if updated != nil {
			res = updated
		}
	}

	if factoryNeedsUpdate {
		if err := s.repository.Update(ctx, entity); err != nil {
			return nil, fmt.Errorf("update factory record: %w", err)
		}
	}

	return s.factoryResponse(res, entity), nil
}

func (s *serviceImpl) DeleteFactory(
	ctx context.Context,
	factoryID uint,
) error {
	entity, tenantID, err := s.getFactoryForTenant(ctx, factoryID)
	if err != nil {
		return err
	}

	if _, err := s.getFactoryResource(ctx, tenantID, entity); err != nil {
		return err
	}

	children, err := s.resourceSvc.GetChildren(ctx, entity.ResourceID)
	if err != nil {
		return fmt.Errorf("get factory children: %w", err)
	}
	if len(children) > 0 {
		return ErrCannotDeleteFactoryWithChildren
	}

	// Delete the Resource first to ensure tenant-scoped deletion succeeds
	// before removing the Factory record.
	if err := s.resourceSvc.DeleteResource(
		ctx,
		tenantID,
		entity.ResourceID,
	); err != nil {
		return fmt.Errorf("delete factory resource: %w", err)
	}

	if err := s.repository.Delete(ctx, entity.ID); err != nil {
		return fmt.Errorf("delete factory record: %w", err)
	}

	return nil
}

func (s *serviceImpl) CreateTopologyNode(
	ctx context.Context,
	req *CreateTopologyNodeRequest,
) (*TopologyNodeResponse, error) {
	if req == nil {
		s.l.Get().Errorf("request is nil")
		return nil, errors.New("request is nil")
	}
	if req.FactoryID == 0 {
		s.l.Get().Errorf("factory id is empty")
		return nil, errors.New("factory_id is required")
	}

	var nodeType resource.ResourceType

	switch strings.ToLower(strings.TrimSpace(req.Type)) {
	case "production_line":
		nodeType = resource.ResourceTypeProductionLine
	case "work_station":
		nodeType = resource.ResourceTypeWorkStation
	default:
		s.l.Get().Errorf(
			"invalid topology resource type: %q",
			req.Type,
		)
		return nil, ErrInvalidTopologyType
	}

	if _, ok := topologyTypes[nodeType]; !ok {
		s.l.Get().Errorf(
			"topology resource type is not allowed: %q",
			nodeType,
		)
		return nil, ErrInvalidTopologyType
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		s.l.Get().Errorf("name is empty")
		return nil, errors.New("name is required")
	}

	entity, tenantID, err := s.getFactoryForTenant(ctx, req.FactoryID)
	if err != nil {
		s.l.Get().Errorf("factory not found by tenant")
		return nil, err
	}
	if _, err := s.getFactoryResource(ctx, tenantID, entity); err != nil {
		s.l.Get().Errorf("resource not found")
		return nil, err
	}

	parentResourceID := entity.ResourceID
	if req.ParentResourceID != nil {
		parentResourceID = *req.ParentResourceID
	}

	if err := s.validateTopologyParent(
		ctx,
		tenantID,
		entity.ResourceID,
		parentResourceID,
		nodeType,
	); err != nil {
		s.l.Get().Errorf("invalid topology parent resource, error: %s", err.Error())
		return nil, err
	}

	metadata := []byte(`{}`)
	if len(req.Metadata) > 0 {
		metadata, err = json.Marshal(req.Metadata)
		if err != nil {
			s.l.Get().Errorf("json error: %s", err.Error())
			return nil, fmt.Errorf("marshal metadata: %w", err)
		}
	}

	res, err := s.resourceSvc.CreateResource(
		ctx,
		&resource.CreateResource{
			TenantID: tenantID,
			ParentID: &parentResourceID,
			Type:     string(nodeType),
			Name:     name,
			Code:     optionalString(req.Code),
			Status:   model.StatusActive,
			Metadata: metadata,
		},
	)
	if err != nil {
		s.l.Get().Errorf("failed to create resource error: %s", err.Error())
		return nil, fmt.Errorf("create topology node: %w", err)
	}

	return s.topologyNodeResponse(res), nil
}

func (s *serviceImpl) UpdateTopologyNode(
	ctx context.Context,
	resourceID uint,
	req *UpdateTopologyNodeRequest,
) (*TopologyNodeResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil, errors.New("tenant id not found")
	}

	entity, err := s.resourceSvc.GetResourceByID(ctx, tenantID, resourceID)
	if err != nil {
		return nil, ErrNodeNotFound
	}
	if resource.ResourceType(entity.ResourceType) == resource.ResourceTypeFactory {
		return nil, errors.New("factory resource must be updated through factory service")
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
		if (entity.Code == nil && code != nil) ||
			(entity.Code != nil && code == nil) ||
			(entity.Code != nil && code != nil && *entity.Code != *code) {
			entity.Code = code
			needsUpdate = true
		}
	}

	if req.Metadata != nil {
		metadata, marshalErr := json.Marshal(req.Metadata)
		if marshalErr != nil {
			return nil, fmt.Errorf("marshal metadata: %w", marshalErr)
		}
		entity.Metadata = metadata
		needsUpdate = true
	}

	if needsUpdate {
		updateReq := &resource.UpdateResource{
			Name:     entity.ResourceName,
			Code:     entity.Code,
			Status:   entity.ResourceStatus,
			Metadata: entity.Metadata,
			Version:  entity.Version,
			ParentID: entity.ParentID,
		}
		updated, updateErr := s.resourceSvc.UpdateResource(
			ctx,
			entity.ID,
			updateReq,
		)
		if updateErr != nil {
			return nil, fmt.Errorf("update topology node: %w", updateErr)
		}
		if updated != nil {
			entity = updated
		}
	}

	return s.topologyNodeResponse(entity), nil
}

func (s *serviceImpl) MoveTopologyNode(
	ctx context.Context,
	req *MoveTopologyNodeRequest,
) error {
	if req == nil {
		return errors.New("request is nil")
	}
	if req.ResourceID == 0 {
		return errors.New("resource_id is required")
	}
	if req.ParentResourceID == nil || *req.ParentResourceID == 0 {
		return ErrInvalidParent
	}

	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return errors.New("tenant id not found")
	}

	node, err := s.resourceSvc.GetResourceByID(
		ctx,
		tenantID,
		req.ResourceID,
	)
	if err != nil {
		return ErrNodeNotFound
	}
	if resource.ResourceType(node.ResourceType) == resource.ResourceTypeFactory {
		return errors.New("factory resource cannot be moved")
	}

	parentID := *req.ParentResourceID
	if parentID == req.ResourceID ||
		s.createsCycle(ctx, tenantID, req.ResourceID, parentID) {
		return ErrTopologyCycle
	}

	factoryResourceID, err := s.findFactoryResourceID(
		ctx,
		tenantID,
		req.ResourceID,
	)
	if err != nil {
		return err
	}

	if err := s.validateTopologyParent(
		ctx,
		tenantID,
		factoryResourceID,
		parentID,
		resource.ResourceType(node.ResourceType),
	); err != nil {
		return err
	}

	return s.resourceSvc.UpdateParent(
		ctx,
		tenantID,
		req.ResourceID,
		parentID,
	)
}

func (s *serviceImpl) DeleteTopologyNode(
	ctx context.Context,
	resourceID uint,
) error {
	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return errors.New("tenant id not found")
	}

	entity, err := s.resourceSvc.GetResourceByID(ctx, tenantID, resourceID)
	if err != nil {
		return ErrNodeNotFound
	}
	if resource.ResourceType(entity.ResourceType) == resource.ResourceTypeFactory {
		return errors.New("factory must be deleted through factory service")
	}

	children, err := s.resourceSvc.GetChildren(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("get topology node children: %w", err)
	}
	if len(children) > 0 {
		return ErrCannotDeleteNodeWithChildren
	}

	if err := s.resourceSvc.DeleteResource(ctx, tenantID, resourceID); err != nil {
		return fmt.Errorf("delete topology node: %w", err)
	}

	return nil
}

func (s *serviceImpl) GetTopology(
	ctx context.Context,
	factoryID uint,
) (*FactoryTopologyResponse, error) {
	entity, tenantID, err := s.getFactoryForTenant(ctx, factoryID)
	if err != nil {
		return nil, err
	}

	res, err := s.getFactoryResource(ctx, tenantID, entity)
	if err != nil {
		return nil, err
	}

	nodes, err := s.collectTopology(ctx, tenantID, entity.ResourceID)
	if err != nil {
		return nil, err
	}

	return &FactoryTopologyResponse{
		Factory: *s.factoryResponse(res, entity),
		Nodes:   nodes,
	}, nil
}

// validateTopologyParent enforces:
// Factory -> ProductionLine -> WorkStation.
func (s *serviceImpl) validateTopologyParent(
	ctx context.Context,
	tenantID uint,
	factoryResourceID uint,
	parentResourceID uint,
	childType resource.ResourceType,
) error {
	parent, err := s.resourceSvc.GetResourceByID(
		ctx,
		tenantID,
		parentResourceID,
	)
	if err != nil {
		return ErrInvalidParent
	}

	isUnder, err := s.isResourceUnderFactory(
		ctx,
		tenantID,
		parent.ID,
		factoryResourceID,
	)
	if err != nil {
		return fmt.Errorf("check parent ancestry: %w", err)
	}
	if !isUnder {
		return ErrInvalidParent
	}

	switch childType {
	case resource.ResourceTypeProductionLine:
		if resource.ResourceType(parent.ResourceType) != resource.ResourceTypeFactory {
			return ErrInvalidParent
		}
	case resource.ResourceTypeWorkStation:
		if resource.ResourceType(parent.ResourceType) != resource.ResourceTypeProductionLine {
			return ErrInvalidParent
		}
	default:
		return ErrInvalidTopologyType
	}

	return nil
}

func (s *serviceImpl) isResourceUnderFactory(
	ctx context.Context,
	tenantID uint,
	resourceID uint,
	factoryResourceID uint,
) (bool, error) {
	currentID := resourceID

	for {
		if currentID == factoryResourceID {
			return true, nil
		}

		current, err := s.resourceSvc.GetResourceByID(
			ctx,
			tenantID,
			currentID,
		)
		if err != nil {
			return false, err
		}

		if current.ParentID == 0 {
			return false, nil
		}

		currentID = current.ParentID
	}
}

func (s *serviceImpl) findFactoryResourceID(
	ctx context.Context,
	tenantID uint,
	resourceID uint,
) (uint, error) {
	currentID := resourceID

	for {
		current, err := s.resourceSvc.GetResourceByID(
			ctx,
			tenantID,
			currentID,
		)
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

func (s *serviceImpl) createsCycle(
	ctx context.Context,
	tenantID uint,
	resourceID uint,
	newParentID uint,
) bool {
	currentID := newParentID

	for {
		if currentID == resourceID {
			return true
		}

		current, err := s.resourceSvc.GetResourceByID(
			ctx,
			tenantID,
			currentID,
		)
		if err != nil {
			return true
		}

		if current.ParentID == 0 {
			return false
		}

		currentID = current.ParentID
	}
}

func (s *serviceImpl) collectTopology(
	ctx context.Context,
	tenantID uint,
	parentID uint,
) ([]TopologyNodeResponse, error) {
	children, err := s.resourceSvc.GetChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("get resource children: %w", err)
	}

	result := make([]TopologyNodeResponse, 0, len(children))

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

func (s *serviceImpl) factoryResponse(
	res *model.Resource,
	entity *model.Factory,
) *FactoryResponse {
	return &FactoryResponse{
		ID:         entity.ID,
		ResourceID: entity.ResourceID,
		Name:       res.ResourceName,
		Code:       entity.Code,
		Address:    entity.Address,
		Timezone:   entity.Timezone,
		Status:     string(res.ResourceStatus),
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}

func (s *serviceImpl) topologyNodeResponse(
	entity *model.Resource,
) *TopologyNodeResponse {
	return &TopologyNodeResponse{
		ResourceID:       entity.ID,
		Type:             string(entity.ResourceType),
		Name:             entity.ResourceName,
		Status:           string(entity.ResourceStatus),
		ParentResourceID: &entity.ParentID,
		Metadata:         decodeMetadata(entity.Metadata),
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
		return nil
	}

	return result
}

func (s *serviceImpl) ListFactories(
	ctx context.Context,
) ([]FactoryResponse, error) {
	tenantID, has := pkg.TenantIDFromContext(ctx)
	if !has {
		return nil, errors.New("tenant id not found")
	}
	entities, err := s.repository.ListByTenant(ctx, tenantID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("list factories: %w", err)
	}
	result := make([]FactoryResponse, 0, len(entities))
	for _, entity := range entities {
		res, err := s.getFactoryResource(ctx, tenantID, entity)
		if err != nil {
			return nil, fmt.Errorf(
				"get resource for factory %d: %w",
				entity.ID,
				err,
			)
		}

		result = append(result, *s.factoryResponse(res, entity))
	}

	return result, nil
}
