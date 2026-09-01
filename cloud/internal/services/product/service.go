package product

import (
	"context"
	"errors"
	"strings"
	"fmt"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/resource"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
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

func (s *serviceImpl) CreateProductModel(ctx context.Context, req *CreateProductModelRequest) (*CreateProductModelResponse, error) {
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
		//ID:          uuid.New(),
		ResourceID:  resourceEntity.ID,
		Code:        code,
		Version:     version,
		Category:    category,
		Description: req.Description,
	}
	if err := s.repository.Create(ctx, entity); err != nil {
		_ = s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.ID) // Rollback
		return nil, fmt.Errorf("create product model: %w", err)
	}

	// 5. Create attribute definitions
	definitions := make([]*model.AttributeDefinition, 0, len(req.Attributes))
	if len(req.Attributes) > 0 {
		for attrName, attribute := range req.Attributes {
			definitions = append(definitions, &model.AttributeDefinition{
				ResourceID:  resourceEntity.ID,
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
			_ = s.resourceSvc.DeleteResource(ctx, tenantID, resourceEntity.ID)
			return nil, fmt.Errorf("create product model attribute definitions: %w", err)
		}
	}

	// 6. Build and return response
	return s.buildCreateProductModelResponse(resourceEntity, entity, definitions), nil
}

func (s *serviceImpl) buildCreateProductModelResponse(resourceEntity *model.Resource, entity *model.ProductModel, ad []*model.AttributeDefinition) *CreateProductModelResponse {
	
	atts := make([]AttributeDefinitionResponse, len(ad))	
	for _, def := range ad {
		atts = append(atts, AttributeDefinitionResponse{
			Name:        def.Name,
			Label:       def.Label,
			Description: def.Description,
			DataType:    string(def.DataType),
			Unit:        def.Unit})
	}	
	
	return &CreateProductModelResponse{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		Name:        resourceEntity.ResourceName,
		Code:        entity.Code,
		Version:     entity.Version,
		Category:    entity.Category,
		Description: entity.Description,
		Status:      resourceEntity.ResourceStatus,
		Attributes:  atts,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// GetProductModel retrieves the full detail of a product model, including all its
// attribute definitions and any attribute values set at the model level.
func (s *serviceImpl) GetProductModel(ctx context.Context, id uint) (*ProductDetailResponse, error) {
	// 1. Get ProductModel
	entity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductModelNotFound
	}

	// 2. Get associated Resource
	resourceEntity, err := s.resourceSvc.GetResourceByID(ctx, tenantIDFromContext(ctx), entity.ResourceID)
	if err != nil {
		return nil, ErrProductModelNotFound
	}

	if resource.ResourceType(resourceEntity.ResourceType) != resource.ResourceTypeProductModel {
		return nil, ErrProductModelNotFound
	}

	// 3. Get all Attribute Definitions for the model
	definitions, err := s.resourceSvc.FindAttributeDefinitionByResourceID(ctx, entity.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute definitions for resource %d: %w", entity.ResourceID, err)
	}

	// 4. Get all Attribute Values for the model
	attributes, err := s.resourceSvc.GetAttributesByResourceID(ctx, entity.ResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attributes for resource %d: %w", entity.ResourceID, err)
	}

	// 5. Construct the detailed response
	response := &ProductDetailResponse{
		ID:                  entity.ID,
		ResourceID:          entity.ResourceID,
		Name:                resourceEntity.ResourceName,
		Code:                entity.Code,
		Version:             entity.Version,
		Category:            entity.Category,
		Description:         entity.Description,
		Status:              string(resourceEntity.ResourceStatus),
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
		Attribute:           attributes,
		AttributeDefinition: definitions,
	}

	return response, nil
}

//TODO: 这里是获取单个产品的详细信息包括产品 资源和属性



func (s *serviceImpl) ListProductModels(ctx context.Context, req *ListProductModelsRequest) (*ProductModelListResponse, error) {
	if req == nil {
		req = &ListProductModelsRequest{}
	}

	if req.CurrentPage <= 0 {
		req.CurrentPage	 = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	var offset = (req.CurrentPage - 1) * req.PageSize
	products, total, err := s.repository.List(ctx, *req)
	if err != nil {
		return nil, fmt.Errorf("list product models: %w", err)
	}

	if len(products) == 0 {
		return &ProductModelListResponse{
			Items:    []*ProductModelResponse{},
			PageBaseResp: pkg.PageBaseResp{
				Total:       total,
				Next: false,
},
		}, nil
	}

	// Collect resource IDs
	resourceIDs := make([]uint, len(products))
	for i, item := range products {
		resourceIDs[i] = item.ResourceID
	}

	// Batch fetch resources
	resources, err := s.resourceSvc.GetResourcesAndAttributesByIDs(ctx, tenantIDFromContext(ctx), resourceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources by ids: %w", err)
	}
	 resourceMap := make(map[uint][]model.ResourceAttribute, len(resources))
	for _, r := range products {
		resourceMap[r.ID] = make([]model.ResourceAttribute, 0)
		for _, rs := range resources {
			if rs.ResourceID == r.ResourceID{
				resourceMap[r.ID] = append(resourceMap[r.ID], rs)
			}
		}
	}
	
	// Build responses
	result := &ProductModelListResponse{
		Items:    []*ProductModelResponse{},
		PageBaseResp: pkg.PageBaseResp{
			Total:       total,
			Next: false,
		},
	}

	for _, product := range products {
		response := s.buildProductModelResponse(resourceMap, product)
		result.Items = append(result.Items, response)
	}

	if int64(offset+req.PageSize) < total {
		result.PageBaseResp.Next = true
	}

	return result, nil
}
func (s *serviceImpl) UpdateProductModel(ctx context.Context, id uint, req *UpdateProductModelRequest) (*UpdateProductModelResponse, error) {
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
		if _, err := s.resourceSvc.UpdateResource(ctx, resourceEntity.ID, req); err != nil {
			return nil, fmt.Errorf("update product model resource: %w", err)
		}
	}

	if productNeedsUpdate {
		if err := s.repository.Update(ctx, entity); err != nil {
			return nil, fmt.Errorf("update product model: %w", err)
		}
	}

	return &UpdateProductModelResponse{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		Name:        resourceEntity.ResourceName,
		Code:        *resourceEntity.Code,
		Description: entity.Description,
		Status:      resourceEntity.ResourceStatus,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func (s *serviceImpl) UpdateProductModelStatus(ctx context.Context, id uint, status string) error {
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

	if _, err := s.resourceSvc.UpdateResource(ctx, resourceEntity.ID, req); err != nil {
		return fmt.Errorf("update product model status: %w", err)
	}

	return nil
}

func (s *serviceImpl) GetAttributeDefinitions(ctx context.Context, productModelID uint) ([]AttributeDefinitionResponse, error) {
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
			ID:          definition.ID,
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

func (s *serviceImpl) UpdateAttributeDefinitions(ctx context.Context, productModelID uint, req *UpdateAttributeDefinitionsRequest) error {
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
			ResourceID:  entity.ResourceID,
			Name:        strings.TrimSpace(name),
			Label:       attribute.Label,
			Description: attribute.Description,
			DataType:    model.AttributeValueType(attribute.DataType),
			Unit:        attribute.Unit,
		})
	}

	return s.resourceSvc.ReplaceAttributeDefinitions(ctx, entity.ResourceID, definitions)
}

// --- Helper Functions ---

func (s *serviceImpl) buildProductModelResponse(dataM map[uint][]model.ResourceAttribute, entity *model.ProductModel) *ProductModelResponse {
	prs := &ProductModelResponse{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		Code:        entity.Code,
		Version:     entity.Version,
		Category:    entity.Category,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
	
	if  v, has:= dataM[prs.ID]; has{
		prs.Attributes = make([]AttributeResponse, 0, len(v))
		for _, attr := range v {
			prs.Attributes = append(prs.Attributes, AttributeResponse{
				ID:          attr.ID,
				Name:        attr.AttributeDefinition.Name,
				Label:       attr.AttributeDefinition.Label,
				Description: attr.AttributeDefinition.Description,
				DataType:    string(attr.AttributeDefinition.DataType),
				Unit:        attr.AttributeDefinition.Unit,
			})
			prs.Name = attr.Resource.ResourceName
			prs.Status = string(attr.Resource.ResourceStatus)
		}
	}

	return prs
}

type UpdateProductResponse struct {

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