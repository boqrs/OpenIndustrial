package resource

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Service provides business logic for resources, including attribute management.
type Service struct {
	resourceRepo ResourceRepository
	attrDefRepo  AttributeDefinitionRepository
	resAttrRepo  ResourceAttributeRepository
}

// NewService creates a new resource service.
func NewService(
	resourceRepo ResourceRepository,
	attrDefRepo AttributeDefinitionRepository,
	resAttrRepo ResourceAttributeRepository,
) *Service {
	return &Service{
		resourceRepo: resourceRepo,
		attrDefRepo:  attrDefRepo,
		resAttrRepo:  resAttrRepo,
	}
}

// CreateProductParams defines the parameters for creating a product.
// This remains as a clear Data Transfer Object (DTO).
type CreateProductParams struct {
	Name         string
	Description  string
	Type         string
	SerialNumber string
	OwnerGroupID uuid.UUID
}

// CreateProduct creates a new product resource, sets its owner, and saves its specific attributes.
// This function is PRESERVED and UPGRADED to the new architecture.
func (s *Service) CreateProduct(ctx context.Context, tenantID uuid.UUID, params CreateProductParams) (*Resource, error) {
	// 1. Create the core resource object.
	// Note: Description and SerialNumber are NOT part of the core resource anymore.
	// We now correctly set the OwnerGroupID directly on the resource.
	resource := &Resource{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         params.Name,
		Type:         params.Type,
		OwnerGroupID: &params.OwnerGroupID, // This is the NEW way to set ownership.
		RecordVersion: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 2. Save the core resource.
	if err := s.resourceRepo.CreateResource(ctx, resource); err != nil {
		return nil, err
	}

	// 3. Handle Description and SerialNumber as DYNAMIC ATTRIBUTES.
	// This is the "add new interface implementation" part.
	attributesToSet := []*ResourceAttribute{}

	// Handle Description attribute
	if params.Description != "" {
		def, err := s.attrDefRepo.GetAttributeDefinitionByKey(ctx, tenantID, "description")
		if err != nil {
			// You might want to create the definition if it doesn't exist, or return an error.
			// For now, we'll assume it exists and return an error if not.
			return nil, errors.New("attribute definition for 'description' not found")
		}
		attributesToSet = append(attributesToSet, &ResourceAttribute{
			ResourceID:  resource.ID,
			AttributeID: def.ID,
			ValueString: &params.Description,
		})
	}

	// Handle SerialNumber attribute
	if params.SerialNumber != "" {
		def, err := s.attrDefRepo.GetAttributeDefinitionByKey(ctx, tenantID, "serial_number")
		if err != nil {
			return nil, errors.New("attribute definition for 'serial_number' not found")
		}
		attributesToSet = append(attributesToSet, &ResourceAttribute{
			ResourceID:  resource.ID,
			AttributeID: def.ID,
			ValueString: &params.SerialNumber,
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

// --- Generic CRUD Methods ---

// CreateResourceParams defines the parameters for creating a generic resource.
type CreateResourceParams struct {
	TenantID     uuid.UUID
	Type         string
	Name         string
	Code         *string
	Status       string
	Metadata     []byte
	ParentID     *uuid.UUID
	OwnerGroupID *uuid.UUID
}

// CreateResource creates a new, generic resource.
func (s *Service) CreateResource(ctx context.Context, params CreateResourceParams) (*Resource, error) {
	resource := &Resource{
		ID:            uuid.New(),
		TenantID:      params.TenantID,
		Type:          params.Type,
		Name:          params.Name,
		Code:          params.Code,
		Status:        params.Status,
		Metadata:      params.Metadata,
		ParentID:      params.ParentID,
		OwnerGroupID:  params.OwnerGroupID,
		RecordVersion: 1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.resourceRepo.CreateResource(ctx, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateResourceParams defines the parameters for updating a resource.
type UpdateResourceParams struct {
	Name     string
	Code     *string
	Status   string
	Metadata []byte
	// Note: TenantID, Type, ParentID, OwnerGroupID are generally not updated via a simple PATCH.
}

// UpdateResource updates an existing resource. It uses optimistic locking.
func (s *Service) UpdateResource(ctx context.Context, tenantID, resourceID uuid.UUID, version int, params UpdateResourceParams) (*Resource, error) {
	// 1. Get the existing resource to ensure it exists and to get its current state.
	res, err := s.resourceRepo.GetResourceByID(ctx, tenantID, resourceID)
	if err != nil {
		return nil, err // Handles not found, etc.
	}

	// 2. Optimistic locking check.
	if res.RecordVersion != version {
		return nil, errors.New("update conflict: resource has been modified by another process")
	}

	// 3. Apply changes from params.
	res.Name = params.Name
	res.Code = params.Code
	res.Status = params.Status
	res.Metadata = params.Metadata
	res.UpdatedAt = time.Now()
	// The repository will handle incrementing the version.

	// 4. Persist the changes.
	if err := s.resourceRepo.UpdateResource(ctx, res); err != nil {
		return nil, err
	}

	// The resource object 'res' now has an old version number.
	// We should return the resource with the *new* version number.
	res.RecordVersion++

	return res, nil
}

// DeleteResource performs a soft delete on a resource.
func (s *Service) DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error {
	// We could add a check here to ensure the resource exists before deleting.
	// For now, we delegate this to the repository.
	return s.resourceRepo.DeleteResource(ctx, tenantID, resourceID)
}

// GetResource retrieves a single resource by its ID.
// This can be enhanced later to also fetch and compose its attributes.
func (s *Service) GetResource(ctx context.Context, tenantID, resourceID uuid.UUID) (*Resource, error) {
	// TODO: Add authorization check here.
	return s.resourceRepo.GetResourceByID(ctx, tenantID, resourceID)
}

// ListResources retrieves a list of resources for a tenant, with filtering and pagination.
// This is the corrected implementation.
func (s *Service) ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*Resource, error) {
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
func (s *Service) SetAttribute(ctx context.Context, tenantID, resourceID uuid.UUID, attrKey string, attrValue interface{}) error {
	def, err := s.attrDefRepo.GetAttributeDefinitionByKey(ctx, tenantID, attrKey)
	if err != nil {
		return err
	}

	attr := &ResourceAttribute{
		ResourceID:  resourceID,
		AttributeID: def.ID,
	}

	// This is a simplified type switch. A real implementation would be more robust.
	switch v := attrValue.(type) {
	case string:
		attr.ValueString = &v
	case int:
		val := int64(v)
		attr.ValueInteger = &val
	case float64:
		attr.ValueFloat = &v
	case bool:
		attr.ValueBoolean = &v
	case time.Time:
		attr.ValueDateTime = &v
	default:
		return errors.New("unsupported attribute type")
	}

	return s.resAttrRepo.SetAttribute(ctx, attr)
}

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