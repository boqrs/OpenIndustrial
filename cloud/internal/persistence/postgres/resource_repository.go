package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ===================================================================
// ResourceRepository Implementation
// ===================================================================

// ResourceRepository implements the resource.ResourceRepository interface for PostgreSQL.
type ResourceRepository struct {
	db *gorm.DB
}

// NewResourceRepository creates a new ResourceRepository.
func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// CreateResource creates a new resource in the database using GORM.
func (r *ResourceRepository) CreateResource(ctx context.Context, res *model.Resource) error {
	return r.db.WithContext(ctx).Create(res).Error
}

// GetResourceByID retrieves a resource by its ID using GORM.
func (r *ResourceRepository) GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*model.Resource, error) {
	var res model.Resource
	err := r.db.WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", resourceID, tenantID).
		First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateResource updates an existing resource using GORM, with optimistic locking.
func (r *ResourceRepository) UpdateResource(ctx context.Context, res *model.Resource) error {
	result := r.db.WithContext(ctx).
		Model(res).
		Where("uuid = ? AND tenant_id = ? AND version = ?", res.UUID, res.TenantID, res.Version).
		Updates(map[string]interface{}{
			"resource_name":           res.ResourceName,
			"code":           res.Code,
			"resource_status":         res.ResourceStatus,
			"metadata":       res.Metadata,
			"version":        gorm.Expr("version + 1"),
			"parent_id":      res.ParentID,
			"owner_group_id": res.OwnerGroupID,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("update failed: resource not found or version mismatch")
	}
	res.Version++ // Manually increment version in the struct after successful update
	return nil
}

// DeleteResource performs a soft delete on a resource using GORM.
func (r *ResourceRepository) DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", resourceID, tenantID).
		Delete(&model.Resource{}).Error
}

// ListResources retrieves a list of resources with pagination using GORM.
func (r *ResourceRepository) ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*model.Resource, error) {
	var resources []*model.Resource
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&resources).Error
	return resources, err
}

// CheckUserInSameGroupAsResource is a placeholder.
func (r *ResourceRepository) CheckUserInSameGroupAsResource(ctx context.Context, userID, resourceID uuid.UUID) (bool, error) {
	// This is a complex query that likely involves joining resources, groups, and group_members.
	// For now, we'll return true to allow development to proceed.
	return true, nil // Placeholder
}

func (r *ResourceRepository) BatchCreateResources(ctx context.Context, resources []*model.Resource) error {
	if len(resources) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&resources).Error
}

func (r *ResourceRepository) FindResourceByNameAndType(ctx context.Context, tenantID uuid.UUID, name, resourceType string) (*model.Resource, error) {
	var resource model.Resource
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND resource_name = ? AND resource_type = ?", tenantID, name, resourceType).
		First(&resource).Error
	return &resource, err
}


// ===================================================================
// AttributeDefinitionRepository Implementation
// ===================================================================

// AttributeDefinitionRepository implements the resource.AttributeDefinitionRepository interface for PostgreSQL.
type AttributeDefinitionRepository struct {
	db *gorm.DB
}

// NewAttributeDefinitionRepository creates a new AttributeDefinitionRepository.
func NewAttributeDefinitionRepository(db *gorm.DB) *AttributeDefinitionRepository {
	return &AttributeDefinitionRepository{db: db}
}

func (r *AttributeDefinitionRepository) CreateAttributeDefinition(ctx context.Context, def *model.AttributeDefinition) error {
	return r.db.WithContext(ctx).Create(def).Error
}

func (r *AttributeDefinitionRepository) GetAttributeDefinitionByID(ctx context.Context, tenantID, defID uuid.UUID) (*model.AttributeDefinition, error) {
	var def model.AttributeDefinition
	err := r.db.WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", defID, tenantID).
		First(&def).Error
	return &def, err
}

func (r *AttributeDefinitionRepository) GetAttributeDefinitionByKey(ctx context.Context, tenantID uuid.UUID, key string) (*model.AttributeDefinition, error) {
	var def model.AttributeDefinition
	err := r.db.WithContext(ctx).
		Where("key = ? AND tenant_id = ?", key, tenantID).
		First(&def).Error
	return &def, err
}

func (r *AttributeDefinitionRepository) ListAttributeDefinitions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*model.AttributeDefinition, error) {
	var defs []*model.AttributeDefinition
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&defs).Error
	return defs, err
}

func (r *AttributeDefinitionRepository) UpdateAttributeDefinition(ctx context.Context, def *model.AttributeDefinition) error {
	return r.db.WithContext(ctx).Model(def).Updates(def).Error
}

func (r *AttributeDefinitionRepository) DeleteAttributeDefinition(ctx context.Context, tenantID, defID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", defID, tenantID).
		Delete(&model.AttributeDefinition{}).Error
}

// FindByName is a placeholder method to satisfy the interface.
func (r *AttributeDefinitionRepository) FindByName(ctx context.Context, tenantID uuid.UUID, name string) (*model.AttributeDefinition, error) {
	var def model.AttributeDefinition
	err := r.db.WithContext(ctx).
		Where("name = ? AND tenant_id = ?", name, tenantID).
		First(&def).Error
	return &def, err
}

// FindByIDs retrieves multiple definitions by their primary UUIDs.
func (r *AttributeDefinitionRepository) FindByIDs(ctx context.Context, tenantID uuid.UUID, ids []uuid.UUID) ([]*model.AttributeDefinition, error) {
	if len(ids) == 0 {
		return []*model.AttributeDefinition{}, nil
	}
	var definitions []*model.AttributeDefinition
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND uuid IN ?", tenantID, ids).Find(&definitions).Error
	return definitions, err
}

func (r *AttributeDefinitionRepository) BatchCreateDefinitions(ctx context.Context, defs []*model.AttributeDefinition) error {
	if len(defs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&defs).Error
}

func ( r *AttributeDefinitionRepository)FindAttributeDefinitionByResourceID(ctx context.Context, resourceID uuid.UUID)([]*model.AttributeDefinition, error){
	var results []*model.AttributeDefinition 
	
	if err := r.db.Model(&model.AttributeDefinition{}).Where("resource_id = ?", resourceID).Find(&results).Error; err != nil{
		return results, err
	}

	return results, nil
}

func (r * AttributeDefinitionRepository)BatchCreateAttributeDefinition(ctx context.Context, attrs []*model.AttributeDefinition)error{
	return r.db.Model(&model.AttributeDefinition{}).CreateInBatches(attrs, len(attrs)).Error
}


// ===================================================================
// ResourceAttributeRepository Implementation
// ===================================================================

// ResourceAttributeRepository implements the resource.ResourceAttributeRepository interface for PostgreSQL.
type ResourceAttributeRepository struct {
	db *gorm.DB
}

// NewResourceAttributeRepository creates a new ResourceAttributeRepository.
func NewResourceAttributeRepository(db *gorm.DB) *ResourceAttributeRepository {
	return &ResourceAttributeRepository{db: db}
}

// SetAttribute uses an UPSERT operation to create or update a resource attribute using GORM.
func (r *ResourceAttributeRepository) SetAttribute(ctx context.Context, attr *model.ResourceAttribute) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}, {Name: "attribute_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value_string", "value_text", "value_integer", "value_float", "value_boolean", "value_datetime", "value_json"}),
	}).Create(attr).Error
}

// SetAttributes sets multiple attributes in a single transaction using GORM.
func (r *ResourceAttributeRepository) SetAttributes(ctx context.Context, attrs []*model.ResourceAttribute) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, attr := range attrs {
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "resource_id"}, {Name: "attribute_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"value_string", "value_text", "value_integer", "value_float", "value_boolean", "value_datetime", "value_json"}),
			}).Create(attr).Error

			if err != nil {
				return fmt.Errorf("failed to set attribute for resource %s: %w", attr.ResourceID, err)
			}
		}
		return nil
	})
}

func (r *ResourceAttributeRepository) GetAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) (*model.ResourceAttribute, error) {
	var attr model.ResourceAttribute
	err := r.db.WithContext(ctx).
		Where("resource_id = ? AND attribute_id = ?", resourceID, attributeID).
		First(&attr).Error
	return &attr, err
}

func (r *ResourceAttributeRepository) GetAttributesByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*model.ResourceAttribute, error) {
	var attrs []*model.ResourceAttribute
	err := r.db.WithContext(ctx).
		Where("resource_id = ?", resourceID).
		Find(&attrs).Error
	return attrs, err
}

func (r *ResourceAttributeRepository) DeleteAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("resource_id = ? AND attribute_id = ?", resourceID, attributeID).
		Delete(&model.ResourceAttribute{}).Error
}

// GetForResource retrieves all attributes for a given resource and returns them as a map.
func (r *ResourceAttributeRepository) GetForResource(ctx context.Context, resourceID uuid.UUID) (map[string]interface{}, error) {
	var results []struct {
		Key           string `gorm:"column:key"`
		ValueString   *string
		ValueText     *string
		ValueInteger  *int64
		ValueFloat    *float64
		ValueBoolean  *bool
		ValueDateTime *time.Time
		ValueJSON     []byte
	}

	err := r.db.WithContext(ctx).
		Model(&model.ResourceAttribute{}).
		Select("attribute_definitions.key, resource_attributes.value_string, resource_attributes.value_text, resource_attributes.value_integer, resource_attributes.value_float, resource_attributes.value_boolean, resource_attributes.value_datetime, resource_attributes.value_json").
		Joins("JOIN attribute_definitions ON attribute_definitions.uuid = resource_attributes.attribute_id").
		Where("resource_attributes.resource_id = ?", resourceID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	attrMap := make(map[string]interface{})
	for _, res := range results {
		if res.ValueString != nil {
			attrMap[res.Key] = *res.ValueString
		} else if res.ValueText != nil {
			attrMap[res.Key] = *res.ValueText
		} else if res.ValueInteger != nil {
			attrMap[res.Key] = *res.ValueInteger
		} else if res.ValueFloat != nil {
			attrMap[res.Key] = *res.ValueFloat
		} else if res.ValueBoolean != nil {
			attrMap[res.Key] = *res.ValueBoolean
		} else if res.ValueDateTime != nil {
			attrMap[res.Key] = *res.ValueDateTime
		} else if res.ValueJSON != nil {
			var v interface{}
			if json.Unmarshal(res.ValueJSON, &v) == nil {
				attrMap[res.Key] = v
			} else {
				attrMap[res.Key] = string(res.ValueJSON)
			}
		}
	}

	return attrMap, nil
}

// UpsertForResource creates or updates a batch of attributes for a specific resource.
func (r *ResourceAttributeRepository) UpsertForResource(ctx context.Context, tenantID, resourceID uuid.UUID, attrs map[string]interface{}) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range attrs {
			// 1. Find the attribute definition to get its UUID.
			var def model.AttributeDefinition
			if err := tx.Where("key = ? AND tenant_id = ?", key, tenantID).First(&def).Error; err != nil {
				return fmt.Errorf("attribute definition for key '%s' not found: %w", key, err)
			}

			// 2. Build the ResourceAttribute model.
			resAttr := model.ResourceAttribute{
				ResourceID:  resourceID,
				ID: def.ID,
			}
			resAttr.SetValue(value) // Use the helper to set the correct value field.

			// 3. Upsert the attribute value.
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "resource_id"}, {Name: "attribute_id"}},
				DoUpdates: clause.AssignmentColumns(resAttr.GetValueColumns()),
			}).Create(&resAttr).Error

			if err != nil {
				return fmt.Errorf("failed to upsert attribute '%s': %w", key, err)
			}
		}
		return nil
	})
}


func( r *ResourceAttributeRepository)	BatchCreateResourceAttributes(ctx context.Context, attr []*model.ResourceAttribute) error{
	return r.db.Model(&model.ResourceAttribute{}).CreateInBatches(attr, len(attr)).Error
}
