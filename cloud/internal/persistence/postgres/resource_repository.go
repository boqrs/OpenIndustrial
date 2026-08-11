package postgres

import (
	"context"
	"fmt"
	"errors"

	"github.com/OpenIndustrial/cloud/internal/kernel/resource"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ResourceRepository implements the resource.ResourceRepository interface for PostgreSQL.
type ResourceRepository struct {
	db *sqlx.DB
}

// NewResourceRepository creates a new ResourceRepository.
func NewResourceRepository(db *sqlx.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// CreateResource creates a new resource in the database.
func (r *ResourceRepository) CreateResource(ctx context.Context, res *resource.Resource) error {
	query := `
		INSERT INTO resources (id, tenant_id, type, name, code, status, metadata, created_at, updated_at, version, parent_id, owner_group_id)
		VALUES (:id, :tenant_id, :type, :name, :code, :status, :metadata, :created_at, :updated_at, :version, :parent_id, :owner_group_id)`
	_, err := r.db.NamedExecContext(ctx, query, res)
	return err
}

// GetResourceByID retrieves a resource by its ID.
func (r *ResourceRepository) GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*resource.Resource, error) {
	var res resource.Resource
	query := `SELECT * FROM resources WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &res, query, resourceID, tenantID)
	return &res, err
}

// UpdateResource updates an existing resource.
func (r *ResourceRepository) UpdateResource(ctx context.Context, res *resource.Resource) error {
	// Optimistic locking check
	query := `
		UPDATE resources SET
			name = :name,
			code = :code,
			status = :status,
			metadata = :metadata,
			updated_at = NOW(),
			version = version + 1
		WHERE id = :id AND tenant_id = :tenant_id AND version = :version`

	result, err := r.db.NamedExecContext(ctx, query, res)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("update failed: resource not found or version mismatch")
	}
	return nil
}

// DeleteResource performs a soft delete on a resource.
func (r *ResourceRepository) DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error {
	query := `UPDATE resources SET deleted_at = NOW() WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, resourceID, tenantID)
	return err
}

// ListResources retrieves a list of resources with pagination.
func (r *ResourceRepository) ListResources(ctx context.Context, tenantID uuid.UUID, resourceType string, limit, offset int) ([]*resource.Resource, error) {
	var resources []*resource.Resource
	query := `
		SELECT * FROM resources
		WHERE tenant_id = $1 AND type = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`
	err := r.db.SelectContext(ctx, &resources, query, tenantID, resourceType, limit, offset)
	return resources, err
}

// CheckUserInSameGroupAsResource is a placeholder. A real implementation would be more complex.
func (r *ResourceRepository) CheckUserInSameGroupAsResource(ctx context.Context, userID, resourceID uuid.UUID) (bool, error) {
	// This is a complex query that likely involves joining resources, groups, and group_members.
	// For now, we'll return true to allow development to proceed.
	return true, nil // Placeholder
}

// CreateResourceRelation creates a new relationship between two resources.
func (r *ResourceRepository) CreateResourceRelation(ctx context.Context, relation *resource.ResourceRelation) error {
	query := `INSERT INTO resource_relations (from_id, to_id, relation_type) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, relation.FromID, relation.ToID, relation.RelationType)
	return err
}

// ListResourceRelations lists relationships for a given resource.
func (r *ResourceRepository) ListResourceRelations(ctx context.Context, resourceID uuid.UUID) ([]*resource.ResourceRelation, error) {
	var relations []*resource.ResourceRelation
	query := `SELECT * FROM resource_relations WHERE from_id = $1 OR to_id = $1`
	err := r.db.SelectContext(ctx, &relations, query, resourceID)
	return relations, err
}

// AttributeDefinitionRepository implements the resource.AttributeDefinitionRepository interface for PostgreSQL.
type AttributeDefinitionRepository struct {
	db *sqlx.DB
}

// NewAttributeDefinitionRepository creates a new AttributeDefinitionRepository.
func NewAttributeDefinitionRepository(db *sqlx.DB) *AttributeDefinitionRepository {
	return &AttributeDefinitionRepository{db: db}
}

func (r *AttributeDefinitionRepository) CreateAttributeDefinition(ctx context.Context, def *resource.AttributeDefinition) error {
	query := `
		INSERT INTO attribute_definitions (id, tenant_id, key, name, description, value_type)
		VALUES (:id, :tenant_id, :key, :name, :description, :value_type)`
	_, err := r.db.NamedExecContext(ctx, query, def)
	return err
}

func (r *AttributeDefinitionRepository) GetAttributeDefinitionByID(ctx context.Context, tenantID, defID uuid.UUID) (*resource.AttributeDefinition, error) {
	var def resource.AttributeDefinition
	query := `SELECT * FROM attribute_definitions WHERE id = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &def, query, defID, tenantID)
	return &def, err
}

func (r *AttributeDefinitionRepository) GetAttributeDefinitionByKey(ctx context.Context, tenantID uuid.UUID, key string) (*resource.AttributeDefinition, error) {
	var def resource.AttributeDefinition
	query := `SELECT * FROM attribute_definitions WHERE key = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &def, query, key, tenantID)
	return &def, err
}

func (r *AttributeDefinitionRepository) ListAttributeDefinitions(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*resource.AttributeDefinition, error) {
	var defs []*resource.AttributeDefinition
	query := `SELECT * FROM attribute_definitions WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &defs, query, tenantID, limit, offset)
	return defs, err
}

func (r *AttributeDefinitionRepository) UpdateAttributeDefinition(ctx context.Context, def *resource.AttributeDefinition) error {
	query := `
		UPDATE attribute_definitions SET
			name = :name,
			description = :description,
			updated_at = NOW()
		WHERE id = :id AND tenant_id = :tenant_id`
	_, err := r.db.NamedExecContext(ctx, query, def)
	return err
}

func (r *AttributeDefinitionRepository) DeleteAttributeDefinition(ctx context.Context, tenantID, defID uuid.UUID) error {
	query := `DELETE FROM attribute_definitions WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, defID, tenantID)
	return err
}


// ResourceAttributeRepository implements the resource.ResourceAttributeRepository interface for PostgreSQL.
type ResourceAttributeRepository struct {
	db *sqlx.DB
}

// NewResourceAttributeRepository creates a new ResourceAttributeRepository.
func NewResourceAttributeRepository(db *sqlx.DB) *ResourceAttributeRepository {
	return &ResourceAttributeRepository{db: db}
}

// SetAttribute uses an UPSERT operation to create or update a resource attribute.
func (r *ResourceAttributeRepository) SetAttribute(ctx context.Context, attr *resource.ResourceAttribute) error {
	query := `
		INSERT INTO resource_attributes (
			resource_id, attribute_id, value_string, value_text, value_integer,
			value_float, value_boolean, value_datetime, value_json
		) VALUES (
			:resource_id, :attribute_id, :value_string, :value_text, :value_integer,
			:value_float, :value_boolean, :value_datetime, :value_json
		)
		ON CONFLICT (resource_id, attribute_id) DO UPDATE SET
			value_string = EXCLUDED.value_string,
			value_text = EXCLUDED.value_text,
			value_integer = EXCLUDED.value_integer,
			value_float = EXCLUDED.value_float,
			value_boolean = EXCLUDED.value_boolean,
			value_datetime = EXCLUDED.value_datetime,
			value_json = EXCLUDED.value_json,
			updated_at = NOW()`
	_, err := r.db.NamedExecContext(ctx, query, attr)
	return err
}

// SetAttributes sets multiple attributes in a single transaction.
func (r *ResourceAttributeRepository) SetAttributes(ctx context.Context, attrs []*resource.ResourceAttribute) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback is a no-op if the transaction is committed.

	for _, attr := range attrs {
		if err := r.setAttributeInTx(tx, attr); err != nil {
			return fmt.Errorf("failed to set attribute %s for resource %s: %w", attr.AttributeID, attr.ResourceID, err)
		}
	}

	return tx.Commit()
}

func (r *ResourceAttributeRepository) setAttributeInTx(tx *sqlx.Tx, attr *resource.ResourceAttribute) error {
	query := `
		INSERT INTO resource_attributes (
			resource_id, attribute_id, value_string, value_text, value_integer,
			value_float, value_boolean, value_datetime, value_json
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (resource_id, attribute_id) DO UPDATE SET
			value_string = EXCLUDED.value_string,
			value_text = EXCLUDED.value_text,
			value_integer = EXCLUDED.value_integer,
			value_float = EXCLUDED.value_float,
			value_boolean = EXCLUDED.value_boolean,
			value_datetime = EXCLUDED.value_datetime,
			value_json = EXCLUDED.value_json,
			updated_at = NOW()`
	_, err := tx.Exec(query,
		attr.ResourceID, attr.AttributeID, attr.ValueString, attr.ValueText, attr.ValueInteger,
		attr.ValueFloat, attr.ValueBoolean, attr.ValueDateTime, attr.ValueJSON,
	)
	return err
}


func (r *ResourceAttributeRepository) GetAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) (*resource.ResourceAttribute, error) {
	var attr resource.ResourceAttribute
	query := `SELECT * FROM resource_attributes WHERE resource_id = $1 AND attribute_id = $2`
	err := r.db.GetContext(ctx, &attr, query, resourceID, attributeID)
	return &attr, err
}

func (r *ResourceAttributeRepository) GetAttributesByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*resource.ResourceAttribute, error) {
	var attrs []*resource.ResourceAttribute
	query := `SELECT * FROM resource_attributes WHERE resource_id = $1`
	err := r.db.SelectContext(ctx, &attrs, query, resourceID)
	return attrs, err
}

func (r *ResourceAttributeRepository) DeleteAttribute(ctx context.Context, resourceID, attributeID uuid.UUID) error {
	query := `DELETE FROM resource_attributes WHERE resource_id = $1 AND attribute_id = $2`
	_, err := r.db.ExecContext(ctx, query, resourceID, attributeID)
	return err
}