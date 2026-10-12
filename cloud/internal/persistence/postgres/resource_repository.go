package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ===================================================================
// ResourceRepository Implementation
// ===================================================================

// ResourceRepository implements the resource.ResourceRepository interface for PostgreSQL.
type ResourceRepository struct {
	db *database.DBProvider
}

// NewResourceRepository creates a new ResourceRepository.
func NewResourceRepository(db *database.DBProvider) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// CreateResource creates a new resource in the database using GORM.
func (r *ResourceRepository) CreateResource(ctx context.Context, res *model.Resource) error {
	return r.db.Get().WithContext(ctx).Create(res).Error
}

func (r *ResourceRepository) CreateResourceTx(ctx context.Context, res *model.Resource) error {
	return dbFromContext(ctx, r.db.Get()).WithContext(ctx).Create(res).Error
}

func (r *ResourceRepository) CreateResourceBatchTx(ctx context.Context, res []*model.Resource) error {
	return dbFromContext(ctx, r.db.Get()).WithContext(ctx).CreateInBatches(res, len(res)).Error
}

func (r *ResourceRepository) FindByParentID(ctx context.Context, tenantID uint, parentID uint) ([]*model.Resource, error) {
	var resources []*model.Resource
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND parent_id = ?", tenantID, parentID).Find(&resources).Error
	return resources, err
}

// GetResourceByID retrieves a resource by its ID using GORM.
func (r *ResourceRepository) GetResourceByID(ctx context.Context, tenantID uint, resourceID uint) (*model.Resource, error) {
	var res model.Resource
	err := r.db.Get().WithContext(ctx).
		Where("id = ? AND tenant_id = ?", resourceID, tenantID).
		First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateResource updates an existing resource using GORM, with optimistic locking.

func (r *ResourceRepository) UpdateResource(ctx context.Context, res *model.Resource) error {
	result := r.db.Get().WithContext(ctx).
		Model(res).
		Where(
			"id = ? AND tenant_id = ? AND record_version = ?",
			res.ID,
			res.TenantID,
			res.Version,
		).
		Updates(map[string]interface{}{
			"resource_name":   res.ResourceName,
			"code":            res.Code,
			"resource_status": res.ResourceStatus,
			"metadata":        res.Metadata,
			"record_version":  gorm.Expr("record_version + 1"),
			"parent_id":       res.ParentID,
			"owner_group_id":  res.OwnerGroupID,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("update failed: resource not found or version mismatch")
	}

	res.Version++
	return nil
}

// DeleteResource performs a soft delete on a resource using GORM.
func (r *ResourceRepository) DeleteResource(ctx context.Context, tenantID uint, resourceID uint) error {
	return r.db.Get().WithContext(ctx).
		Where("id = ? AND tenant_id = ?", resourceID, tenantID).
		Delete(&model.Resource{}).Error
}

// ListResources retrieves a list of resources with pagination using GORM.
func (r *ResourceRepository) ListResources(ctx context.Context, tenantID uint, resourceType string, limit, offset int) ([]*model.Resource, error) {
	var resources []*model.Resource
	query := r.db.Get().WithContext(ctx).Where("tenant_id = ?", tenantID)
	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&resources).Error
	return resources, err
}

// CheckUserInSameGroupAsResource is a placeholder.
func (r *ResourceRepository) CheckUserInSameGroupAsResource(ctx context.Context, userID uuid.UUID, resourceID uint) (bool, error) {
	// This is a complex query that likely involves joining resources, groups, and group_members.
	// For now, we'll return true to allow development to proceed.
	return true, nil // Placeholder
}

func (r *ResourceRepository) BatchCreateResources(ctx context.Context, resources []*model.Resource) error {
	if len(resources) == 0 {
		return nil
	}
	return r.db.Get().WithContext(ctx).Create(&resources).Error
}

func (r *ResourceRepository) FindResourceByNameAndType(ctx context.Context, tenantID uint, name, resourceType string) (*model.Resource, error) {
	var resource model.Resource
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND resource_name = ? AND resource_type = ?", tenantID, name, resourceType).
		First(&resource).Error
	return &resource, err
}
func (r *ResourceRepository) UpdateParent(ctx context.Context, tenantID uint, resourceID, newParentID uint) error {
	var parentResource model.Resource
	if err := r.db.Get().WithContext(ctx).
		Where("id = ? AND tenant_id = ?", newParentID, tenantID).
		First(&parentResource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("parent resource with ID %s not found for tenant %s", newParentID, tenantID)
		}
		return fmt.Errorf("failed to verify parent resource: %w", err)
	}

	result := r.db.Get().WithContext(ctx).
		Model(&model.Resource{}).
		Where("id = ? AND tenant_id = ?", resourceID, tenantID).
		Update("parent_id", &newParentID)

	if result.Error != nil {
		return fmt.Errorf("failed to update parent for resource %s: %w", resourceID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("resource with ID %s not found for tenant %s", resourceID, tenantID)
	}
	return nil
}

func (r *ResourceRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.Get().WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ===================================================================
// AttributeDefinitionRepository Implementation
// ===================================================================

// AttributeDefinitionRepository implements the resource.AttributeDefinitionRepository interface for PostgreSQL.
type AttributeDefinitionRepository struct {
	db *database.DBProvider
}

// NewAttributeDefinitionRepository creates a new AttributeDefinitionRepository.
func NewAttributeDefinitionRepository(db *database.DBProvider) *AttributeDefinitionRepository {
	return &AttributeDefinitionRepository{db: db}
}

func (r *AttributeDefinitionRepository) CreateAttributeDefinition(ctx context.Context, def *model.AttributeDefinition) error {
	return r.db.Get().WithContext(ctx).Create(def).Error
}

func (r *AttributeDefinitionRepository) GetAttributeDefinitionByID(ctx context.Context, tenantID uint, defID uint) (*model.AttributeDefinition, error) {
	var def model.AttributeDefinition
	err := r.db.Get().WithContext(ctx).
		Where("id = ? AND tenant_id = ?", defID, tenantID).
		First(&def).Error
	return &def, err
}

func (r *AttributeDefinitionRepository) GetAttributeDefinitionByKey(ctx context.Context, tenantID uint, key string) (*model.AttributeDefinition, error) {
	var def model.AttributeDefinition
	err := r.db.Get().WithContext(ctx).
		Where("key = ? AND tenant_id = ?", key, tenantID).
		First(&def).Error
	return &def, err
}

func (r *AttributeDefinitionRepository) ListAttributeDefinitions(ctx context.Context, tenantID uint, limit, offset int) ([]*model.AttributeDefinition, error) {
	var defs []*model.AttributeDefinition
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&defs).Error
	return defs, err
}

func (r *AttributeDefinitionRepository) UpdateAttributeDefinition(ctx context.Context, def *model.AttributeDefinition) error {
	return r.db.Get().WithContext(ctx).Model(def).Updates(def).Error
}

func (r *AttributeDefinitionRepository) DeleteAttributeDefinition(ctx context.Context, tenantID uint, defID uint) error {
	return r.db.Get().WithContext(ctx).
		Where("id = ? AND tenant_id = ?", defID, tenantID).
		Delete(&model.AttributeDefinition{}).Error
}

// FindByName is a placeholder method to satisfy the interface.
func (r *AttributeDefinitionRepository) FindByName(ctx context.Context, tenantID uint, name string) (*model.AttributeDefinition, error) {
	var def model.AttributeDefinition
	err := r.db.Get().WithContext(ctx).
		Where("name = ? AND tenant_id = ?", name, tenantID).
		First(&def).Error
	return &def, err
}

// FindByIDs retrieves multiple definitions by their primary UUIDs.
func (r *AttributeDefinitionRepository) FindByIDs(ctx context.Context, tenantID uint, ids []uint) ([]*model.AttributeDefinition, error) {
	if len(ids) == 0 {
		return []*model.AttributeDefinition{}, nil
	}
	var definitions []*model.AttributeDefinition
	err := r.db.Get().WithContext(ctx).Where("tenant_id = ? AND id IN ?", tenantID, ids).Find(&definitions).Error
	return definitions, err
}

func (r *AttributeDefinitionRepository) BatchCreateDefinitions(ctx context.Context, defs []*model.AttributeDefinition) error {
	if len(defs) == 0 {
		return nil
	}
	return r.db.Get().WithContext(ctx).Create(&defs).Error
}

func (r *AttributeDefinitionRepository) FindAttributeDefinitionByResourceID(ctx context.Context, resourceID uint) ([]*model.AttributeDefinition, error) {
	var results []*model.AttributeDefinition

	if err := r.db.Get().Model(&model.AttributeDefinition{}).Where("resource_id = ?", resourceID).Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

func (r *AttributeDefinitionRepository) BatchCreateAttributeDefinition(ctx context.Context, attrs []*model.AttributeDefinition) error {
	return r.db.Get().Model(&model.AttributeDefinition{}).CreateInBatches(attrs, len(attrs)).Error
}

func (r *AttributeDefinitionRepository) ReplaceAttributeDefinitions(ctx context.Context, resourceID uint, definitions []*model.AttributeDefinition) error {
	return r.db.Get().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Delete all existing definitions associated with the resource ID.
		if err := tx.Where("resource_id = ?", resourceID).Delete(&model.AttributeDefinition{}).Error; err != nil {
			return fmt.Errorf("failed to delete old attribute definitions: %w", err)
		}

		// 2. If new definitions are provided, create them in a batch.
		// GORM's Create function can handle a slice of structs for batch insertion.
		if len(definitions) > 0 {
			if err := tx.Create(&definitions).Error; err != nil {
				return fmt.Errorf("failed to create new attribute definitions: %w", err)
			}
		}

		// If both operations succeed, the transaction is committed.
		return nil
	})
}

// ===================================================================
// ResourceAttributeRepository Implementation
// ===================================================================

// ResourceAttributeRepository implements the resource.ResourceAttributeRepository interface for PostgreSQL.
type ResourceAttributeRepository struct {
	db *database.DBProvider
}

// NewResourceAttributeRepository creates a new ResourceAttributeRepository.
func NewResourceAttributeRepository(db *database.DBProvider) *ResourceAttributeRepository {
	return &ResourceAttributeRepository{db: db}
}

// SetAttribute uses an UPSERT operation to create or update a resource attribute using GORM.
func (r *ResourceAttributeRepository) SetAttribute(
	ctx context.Context,
	attr *model.ResourceAttribute,
) error {
	if attr == nil {
		return errors.New("resource attribute is nil")
	}

	return r.db.Get().WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "resource_id"},
				{Name: "attribute_definition_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"value",
				"updated_at",
			}),
		}).
		Create(attr).Error
}

// SetAttributes sets multiple attributes in a single transaction using GORM.
func (r *ResourceAttributeRepository) SetAttributes(
	ctx context.Context,
	attrs []*model.ResourceAttribute,
) error {
	if len(attrs) == 0 {
		return nil
	}

	return r.db.Get().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, attr := range attrs {
			if attr == nil {
				return errors.New("resource attribute is nil")
			}

			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "resource_id"},
					{Name: "attribute_definition_id"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"value",
					"updated_at",
				}),
			}).Create(attr).Error; err != nil {
				return fmt.Errorf(
					"set attribute for resource %d, definition %d: %w",
					attr.ResourceID,
					attr.AttributeDefinitionID,
					err,
				)
			}
		}

		return nil
	})
}

func (r *ResourceAttributeRepository) GetAttribute(
	ctx context.Context,
	resourceID uint,
	attributeID uint,
) (*model.ResourceAttribute, error) {
	var attr model.ResourceAttribute

	err := r.db.Get().WithContext(ctx).
		Where(
			"resource_id = ? AND attribute_definition_id = ?",
			resourceID,
			attributeID,
		).
		First(&attr).Error
	if err != nil {
		return nil, err
	}

	return &attr, nil
}

func (r *ResourceAttributeRepository) GetAttributesByResourceID(
	ctx context.Context,
	resourceID uint,
) ([]*model.ResourceAttribute, error) {
	attrs := make([]*model.ResourceAttribute, 0)

	err := r.db.Get().WithContext(ctx).
		Preload("AttributeDefinition").
		Where("resource_id = ?", resourceID).
		Order("id ASC").
		Find(&attrs).Error
	if err != nil {
		return nil, err
	}

	return attrs, nil
}

func (r *ResourceAttributeRepository) DeleteAttribute(
	ctx context.Context,
	resourceID uint,
	attributeID uint,
) error {
	result := r.db.Get().WithContext(ctx).
		Where(
			"resource_id = ? AND attribute_definition_id = ?",
			resourceID,
			attributeID,
		).
		Delete(&model.ResourceAttribute{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetForResource retrieves all attributes for a given resource and returns them as a map.
func (r *ResourceAttributeRepository) GetForResource(
	ctx context.Context,
	resourceID uint,
) (map[string]interface{}, error) {
	attrs, err := r.GetAttributesByResourceID(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			continue
		}

		key := attr.AttributeDefinition.Name
		if key == "" {
			continue
		}

		var value interface{}
		if err := json.Unmarshal(attr.Value, &value); err != nil {
			return nil, fmt.Errorf(
				"decode attribute %q: %w",
				key,
				err,
			)
		}

		result[key] = value
	}

	return result, nil
}

// UpsertForResource creates or updates a batch of attributes for a specific resource.
func (r *ResourceAttributeRepository) UpsertForResource(
	ctx context.Context,
	tenantID uint,
	resourceID uint,
	attributes map[string]interface{},
) error {
	if len(attributes) == 0 {
		return nil
	}

	return r.db.Get().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Ensure the target resource belongs to the tenant.
		var resourceEntity model.Resource
		if err := tx.
			Where("id = ? AND tenant_id = ?", resourceID, tenantID).
			First(&resourceEntity).Error; err != nil {
			return fmt.Errorf("get target resource: %w", err)
		}

		for name, value := range attributes {
			// Product-level values must use definitions attached to this
			// product Resource, rather than unrelated tenant-wide definitions.
			var definition model.AttributeDefinition
			if err := tx.
				Where("resource_id = ? AND name = ?", resourceID, name).
				First(&definition).Error; err != nil {
				return fmt.Errorf(
					"attribute definition %q not found for resource %d: %w",
					name,
					resourceID,
					err,
				)
			}

			raw, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("encode attribute %q: %w", name, err)
			}

			attr := model.ResourceAttribute{
				ResourceID:            resourceID,
				AttributeDefinitionID: definition.ID,
				Value:                 raw,
			}

			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "resource_id"},
					{Name: "attribute_definition_id"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"value",
					"updated_at",
				}),
			}).Create(&attr).Error; err != nil {
				return fmt.Errorf("save attribute %q: %w", name, err)
			}
		}

		return nil
	})
}



func (r *ResourceAttributeRepository) BatchCreateResourceAttributes(
	ctx context.Context,
	attrs []*model.ResourceAttribute,
) error {
	if len(attrs) == 0 {
		return nil
	}

	return r.db.Get().WithContext(ctx).
		CreateInBatches(attrs, len(attrs)).Error
}

func (r *ResourceAttributeRepository) GetAttributesForResource(
	ctx context.Context,
	tenantID uint,
	resourceID uint,
) (map[string]interface{}, error) {
	var resourceEntity model.Resource
	if err := r.db.Get().WithContext(ctx).
		Where("id = ? AND tenant_id = ?", resourceID, tenantID).
		First(&resourceEntity).Error; err != nil {
		return nil, err
	}

	return r.GetForResource(ctx, resourceID)
}

func (r *ResourceAttributeRepository) GetResourcesAndAttributesByIDs(ctx context.Context, tenantID uint, resourceIDs []uint) ([]model.ResourceAttribute, error) {
	var attributes []model.ResourceAttribute
	if len(resourceIDs) == 0 {
		return attributes, nil
	}

	// We join with resources to filter by tenant_id, and Preload to efficiently fetch both associated Resource and AttributeDefinition objects.
	err := r.db.Get().WithContext(ctx).
		Preload("Resource").
		Preload("AttributeDefinition").
		Joins("JOIN resources ON resources.id = resource_attributes.resource_id").
		Where("resource_attributes.resource_id IN ?", resourceIDs).
		Where("resources.tenant_id = ?", tenantID).
		Find(&attributes).Error

	if err != nil {
		return nil, err
	}
	return attributes, nil
}


// ===================================================================
// ResourceAttributeRepository Implementation
// ===================================================================
type ResourceConnectRepository struct {
	db *database.DBProvider
}

// NewResourceAttributeRepository creates a new ResourceAttributeRepository.
func NewResourceConnectionsRepository(db *database.DBProvider) *ResourceConnectRepository {
	return &ResourceConnectRepository{db: db}
}

func (r *ResourceConnectRepository) CreateConnection(ctx context.Context, conn *model.ResourceConnection) error {
	return r.db.Get().WithContext(ctx).Create(conn).Error
}

func (r *ResourceConnectRepository) GetConnectionByID(ctx context.Context, connectionID uint) (*model.ResourceConnection, error) {
	var conn model.ResourceConnection
	err := r.db.Get().WithContext(ctx).First(&conn, connectionID).Error
	return &conn, err
}

func (r *ResourceConnectRepository) DeleteConnection(ctx context.Context, connectionID uint) error {
	return r.db.Get().WithContext(ctx).Delete(&model.ResourceConnection{}, connectionID).Error
}

func (r *ResourceConnectRepository) ListConnectionsByResourceID(ctx context.Context, resourceID uint) ([]*model.ResourceConnection, error) {
	var conns []*model.ResourceConnection
	err := r.db.Get().WithContext(ctx).Where("source_resource_id = ? OR target_resource_id = ?", resourceID, resourceID).Find(&conns).Error
	return conns, err
}

func (r *ResourceConnectRepository) GetAttributesForResource(ctx context.Context, tenantID uint, resourceID uint) (map[string]interface{}, error) {
	var results []struct {
		Name  string
		Value json.RawMessage
	}
	err := r.db.Get().WithContext(ctx).Model(&model.ResourceAttribute{}).
		Select("ad.name, ra.value").
		Joins("JOIN attribute_definitions ad ON ad.uuid = ra.attribute_definition_id").
		Where("ra.resource_id = ?", resourceID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	attrs := make(map[string]interface{})
	for _, r := range results {
		var val interface{}
		if err := json.Unmarshal(r.Value, &val); err != nil {
			// Decide how to handle malformed JSON in the DB
			continue
		}
		attrs[r.Name] = val
	}

	return attrs, nil
}