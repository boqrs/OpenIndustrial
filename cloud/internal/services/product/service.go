package product

import (
	"context"
	"errors"
	"strings"
	"fmt"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
)

var (
	ErrProductModelNotFound        = errors.New("product model not found")
	ErrProductModelCodeVersionExists = errors.New("product model code and version already exists")
	ErrProductModelImmutable       = errors.New("product model code and version are immutable")
	ErrInvalidProductModel         = errors.New("invalid product model")
	ErrInvalidAttributeDefinition  = errors.New("invalid attribute definition")
	ErrAttributeDefinitionNotFound = errors.New("attribute definition not found")
	ErrProductModelAlreadyActive   = errors.New("product model is already active")
	ErrProductModelCannotModify    = errors.New("active product model cannot modify schema")
)


// serviceImpl implements the Service interface.
type serviceImpl struct {
	resourceSvc resource.Service
	repository  Repository
}

// NewService creates a new product service.
func NewService(resourceSvc resource.Service,repository Repository) Service {
	return &serviceImpl{
		resourceSvc: resourceSvc,
		repository:  repository,
	}
}

// --- tenantIDFromContext is a placeholder ---
// In a real application, this would extract the tenant ID from the context,
// likely from a JWT or other authentication middleware.
func tenantIDFromContext(ctx context.Context) uuid.UUID {
	// TODO: Implement actual tenant ID extraction from context.
	return uuid.Nil
}

func (s *serviceImpl) CreateProductModel(ctx context.Context, req *CreateProductModelRequest) (*ProductModelResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, errors.New("request is nil")
	}
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	version := strings.TrimSpace(req.Version)
	category := strings.TrimSpace(req.Category)

	if name == "" || code == "" || version == "" || category == "" {
		return nil, ErrInvalidProductModel
	}
	if err := validateAttributeDefinitions(req.Attributes); err != nil {
		return nil, err
	}

	// 2. Check for uniqueness
	existing, err := s.repository.GetByCodeAndVersion(ctx, code, version)
	if err == nil && existing != nil {
		return nil, ErrProductModelCodeVersionExists
	}
	if err != nil && !errors.Is(err, ErrProductModelNotFound) {
		return nil, fmt.Errorf("failed to check for existing product model: %w", err)
	}

	// 3. Create core resource
	tenantID := tenantIDFromContext(ctx)
	resourceEntity, err := s.resourceSvc.CreateResource(ctx, &resource.CreateResource{
		TenantID: tenantID,
		Type:     string(resource.ResourceTypeProductModel),
		Name:     name,
		Status:   model.StatusPending, // Always start as pending
	})
	if err != nil {
		return nil, fmt.Errorf("create product model resource: %w", err)
	}

	// 4. Create product model domain entity
	entity := &model.ProductModel{
		ID:          uuid.New(),
		ResourceID:  resourceEntity.UUID,
		Code:        code,
		Version:     version,
		Category:    category,
		Description: req.Description,
	}
	if err := s.repository.Create(ctx, entity); err != nil {
		_ = s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.UUID) // Rollback
		return nil, fmt.Errorf("create product model: %w", err)
	}

	// 5. Create attribute definitions
	if len(req.Attributes) > 0 {
		definitions := make([]*model.AttributeDefinition, 0, len(req.Attributes))
		for attrName, attribute := range req.Attributes {
			definitions = append(definitions, &model.AttributeDefinition{
				UUID:        uuid.New(),
				ResourceID:  resourceEntity.UUID,
				Name:        strings.TrimSpace(attrName),
				Label:       attribute.Label,
				Description: attribute.Description,
				DataType:    model.AttributeValueType(attribute.DataType),
				Unit:        attribute.Unit,
				//Required:    attribute.Required,
			})
		}
		// Assuming BatchCreateAttributeDefinition exists on resourceSvc
		if err := s.resourceSvc.BatchCreateAttributeDefinition(ctx, definitions); err != nil {
			// Full rollback
			_ = s.repository.Delete(ctx, entity.ID)
			_ = s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.UUID)
			return nil, fmt.Errorf("create product model attribute definitions: %w", err)
		}
	}

	// 6. Build and return response
	return s.buildProductModelResponse(ctx, resourceEntity, entity)
}

func (s *serviceImpl) GetProductModel(ctx context.Context, id uuid.UUID) (*ProductModelResponse, error) {
	entity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductModelNotFound // GORM's not found is already handled in repo, this is for service layer consistency
	}

	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), entity.ResourceID)
	if err != nil {
		return nil, ErrProductModelNotFound
	}

	if resource.ResourceType(resourceEntity.ResourceType) != resource.ResourceTypeProductModel {
		return nil, ErrProductModelNotFound
	}

	return s.buildProductModelResponse(ctx, resourceEntity, entity)
}

func (s *serviceImpl) ListProductModels(ctx context.Context, req *ListProductModelsRequest) (*ProductModelListResponse, error) {
	if req == nil {
		req = &ListProductModelsRequest{}
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	items, total, err := s.repository.List(ctx, *req)
	if err != nil {
		return nil, fmt.Errorf("list product models: %w", err)
	}

	result := &ProductModelListResponse{
		Items:    make([]*ProductModelResponse, 0, len(items)),
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	for _, item := range items {
		resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), item.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("get product model resource %s: %w", item.ID, err)
		}
		response, err := s.buildProductModelResponse(ctx, resourceEntity, item)
		if err != nil {
			return nil, fmt.Errorf("build response for model %s: %w", item.ID, err)
		}
		result.Items = append(result.Items, response)
	}

	return result, nil
}

func (s *serviceImpl) UpdateProductModel(ctx context.Context, id uuid.UUID, req *UpdateProductModelRequest) (*ProductModelResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	entity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductModelNotFound
	}

	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), entity.ResourceID)
	if err != nil {
		return nil, ErrProductModelNotFound
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

	productNeedsUpdate := false
	if req.Category != nil {
		category := strings.TrimSpace(*req.Category)
		if category == "" {
			return nil, errors.New("category cannot be empty")
		}
		if entity.Category != category {
			entity.Category = category
			productNeedsUpdate = true
		}
	}
	if req.Description != nil {
		if entity.Description != *req.Description {
			entity.Description = *req.Description
			productNeedsUpdate = true
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
		if _, err := s.resourceSvc.UpdateResource(ctx, resourceEntity.UUID, req); err != nil {
			return nil, fmt.Errorf("update product model resource: %w", err)
		}
	}

	if productNeedsUpdate {
		if err := s.repository.Update(ctx, entity); err != nil {
			return nil, fmt.Errorf("update product model: %w", err)
		}
	}

	return s.buildProductModelResponse(ctx, resourceEntity, entity)
}

func (s *serviceImpl) UpdateProductModelStatus(ctx context.Context, id uuid.UUID, status string) error {
	entity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return ErrProductModelNotFound
	}

	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), entity.ResourceID)
	if err != nil {
		return ErrProductModelNotFound
	}

	// As per your design, a mapping or direct assignment is needed.
	// Assuming resource.ResourceStatus maps directly to model.Status for now.
	var targetStatus string
	switch status {
	case model.StatusPending:
		targetStatus = model.StatusPending
	case model.StatusActive:
		targetStatus = model.StatusActive
	case model.StatusInactive:
		targetStatus = model.StatusInactive
	case model.StatusArchived:
		targetStatus = model.StatusArchived
	default:
		return errors.New("invalid product model status")
	}

	if resourceEntity.ResourceStatus == model.StatusArchived && targetStatus == model.StatusActive {
		return errors.New("archived product model cannot be activated")
	}

	resourceEntity.ResourceStatus = targetStatus
	req := &resource.UpdateResource{
			Name: resourceEntity.ResourceName,
			Code: resourceEntity.Code,
			Status: resourceEntity.ResourceStatus,
			Metadata:resourceEntity.Metadata,
			Version: resourceEntity.Version,
			ParentID: resourceEntity.ParentID,
	}

	if _, err := s.resourceSvc.UpdateResource(ctx, resourceEntity.UUID, req); err != nil {
		return fmt.Errorf("update product model status: %w", err)
	}

	return nil
}

func (s *serviceImpl) GetAttributeDefinitions(ctx context.Context, productModelID uuid.UUID) ([]AttributeDefinitionResponse, error) {
	entity, err := s.repository.GetByID(ctx, productModelID)
	if err != nil {
		return nil, ErrProductModelNotFound
	}

	definitions, err := s.resourceSvc.FindAttributeDefinitionByResourceID(ctx, entity.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("get attribute definitions: %w", err)
	}

	result := make([]AttributeDefinitionResponse, 0, len(definitions))
	for _, definition := range definitions {
		result = append(result, AttributeDefinitionResponse{
			ID:          definition.UUID,
			Name:        definition.Name,
			Label:       definition.Label,
			Description: definition.Description,
			DataType:    string(definition.DataType),
			Unit:        definition.Unit,
			//Required:    definition.Required,
		})
	}

	return result, nil
}

func (s *serviceImpl) UpdateAttributeDefinitions(ctx context.Context, productModelID uuid.UUID, req *UpdateAttributeDefinitionsRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	entity, err := s.repository.GetByID(ctx, productModelID)
	if err != nil {
		return ErrProductModelNotFound
	}

	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), entity.ResourceID)
	if err != nil {
		return ErrProductModelNotFound
	}

	if resourceEntity.ResourceStatus != model.StatusPending {
		return ErrProductModelCannotModify
	}

	if err := validateAttributeDefinitions(req.Attributes); err != nil {
		return err
	}

	definitions := make([]*model.AttributeDefinition, 0, len(req.Attributes))
	for name, attribute := range req.Attributes {
		definitions = append(definitions, &model.AttributeDefinition{
			UUID:        uuid.New(), // New UUIDs for replacement
			ResourceID:  entity.ResourceID,
			Name:        strings.TrimSpace(name),
			Label:       attribute.Label,
			Description: attribute.Description,
			DataType:    model.AttributeValueType(attribute.DataType),
			Unit:        attribute.Unit,
			//Required:    attribute.Required,
		})
	}

	// Per your design, this relies on a Replace/Sync capability in the Resource Kernel.
	return s.resourceSvc.ReplaceAttributeDefinitions(ctx, entity.ResourceID, definitions)
}

// --- Helper Functions ---

func (s *serviceImpl) buildProductModelResponse(ctx context.Context, resourceEntity *model.Resource, entity *model.ProductModel) (*ProductModelResponse, error) {
	definitions, err := s.resourceSvc.FindAttributeDefinitionByResourceID(ctx, entity.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to find attribute definitions: %w", err)
	}

	attributes := make([]AttributeDefinitionResponse, 0, len(definitions))
	for _, definition := range definitions {
		attributes = append(attributes, AttributeDefinitionResponse{
			ID:          definition.UUID,
			Name:        definition.Name,
			Label:       definition.Label,
			Description: definition.Description,
			DataType:    string(definition.DataType),
			Unit:        definition.Unit,
			//Required:    definition.Required,
		})
	}

	return &ProductModelResponse{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		Name:        resourceEntity.ResourceName,
		Code:        entity.Code,
		Version:     entity.Version,
		Category:    entity.Category,
		Description: entity.Description,
		Status:      string(resourceEntity.ResourceStatus),
		Attributes:  attributes,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func validateAttributeDefinitions(attributes map[string]AttributeDefinitionRequest) error {
	seen := make(map[string]struct{}, len(attributes))
	for name, definition := range attributes {
		name = strings.TrimSpace(name)
		if name == "" {
			return errors.New("attribute name cannot be empty")
		}
		if _, exists := seen[name]; exists {
			return fmt.Errorf("duplicate attribute: %s", name)
		}
		seen[name] = struct{}{}

		if !isSupportedAttributeType(definition.DataType) {
			return fmt.Errorf("unsupported attribute data type: %s", definition.DataType)
		}
	}
	return nil
}

func isSupportedAttributeType(value string) bool {
	switch model.AttributeValueType(value) {
	case model.AttributeValueTypeString,
		model.AttributeValueTypeText,
		model.AttributeValueTypeInteger,
		model.AttributeValueTypeFloat,
		model.AttributeValueTypeBoolean,
		model.AttributeValueTypeDateTime,
		model.AttributeValueTypeJSON:
		//model.AttributeValueTypeDecimal,
		//model.AttributeValueTypeResourceReference,
		//model.AttributeValueTypeResourceReferenceList:
		return true
	default:
		return false
	}
}