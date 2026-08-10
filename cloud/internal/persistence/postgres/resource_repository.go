package postgres

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// --- ResourceRepository Implementation ---

type PgResourceRepository struct {
	db *sqlx.DB
}

func NewPgResourceRepository(db *sqlx.DB) *PgResourceRepository {
	return &PgResourceRepository{db: db}
}

func (r *PgResourceRepository) CreateResource(ctx context.Context, res *resource.Resource) error {
	query := `INSERT INTO resources (id, tenant_id, type, name, properties, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, res.ID, res.TenantID, res.Type, res.Name, res.Properties, res.CreatedAt, res.UpdatedAt)
	return err
}

func (r *PgResourceRepository) GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*resource.Resource, error) {
	var res resource.Resource
	query := `SELECT * FROM resources WHERE tenant_id = $1 AND id = $2`
	err := r.db.GetContext(ctx, &res, query, tenantID, resourceID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *PgResourceRepository) UpdateResource(ctx context.Context, res *resource.Resource) error {
	query := `UPDATE resources SET name = $1, properties = $2, updated_at = NOW()
              WHERE id = $3 AND tenant_id = $4`
	_, err := r.db.ExecContext(ctx, query, res.Name, res.Properties, res.ID, res.TenantID)
	return err
}

func (r *PgResourceRepository) DeleteResource(ctx context.Context, tenantID, resourceID uuid.UUID) error {
	query := `DELETE FROM resources WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, tenantID, resourceID)
	return err
}

func (r *PgResourceRepository) ListResourcesByType(ctx context.Context, tenantID uuid.UUID, resourceType string, offset, limit int) ([]resource.Resource, error) {
	var resources []resource.Resource
	query := `SELECT * FROM resources WHERE tenant_id = $1 AND type = $2 ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	err := r.db.SelectContext(ctx, &resources, query, tenantID, resourceType, offset, limit)
	return resources, err
}

func (r *PgResourceRepository) CreateRelation(ctx context.Context, rel *resource.ResourceRelation) error {
	query := `INSERT INTO resource_relations (id, tenant_id, source_resource_id, target_resource_id, relation_type, properties, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, rel.ID, rel.TenantID, rel.SourceResourceID, rel.TargetResourceID, rel.RelationType, rel.Properties, rel.CreatedAt, rel.UpdatedAt)
	return err
}

func (r *PgResourceRepository) GetRelationsBySource(ctx context.Context, tenantID, sourceID uuid.UUID) ([]resource.ResourceRelation, error) {
	var relations []resource.ResourceRelation
	query := `SELECT * FROM resource_relations WHERE tenant_id = $1 AND source_resource_id = $2`
	err := r.db.SelectContext(ctx, &relations, query, tenantID, sourceID)
	return relations, err
}

func (r *PgResourceRepository) GetRelationsByTarget(ctx context.Context, tenantID, targetID uuid.UUID) ([]resource.ResourceRelation, error) {
	var relations []resource.ResourceRelation
	query := `SELECT * FROM resource_relations WHERE tenant_id = $1 AND target_resource_id = $2`
	err := r.db.SelectContext(ctx, &relations, query, tenantID, targetID)
	return relations, err
}

func (r *PgResourceRepository) DeleteRelation(ctx context.Context, tenantID, relationID uuid.UUID) error {
	query := `DELETE FROM resource_relations WHERE tenant_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, tenantID, relationID)
	return err
}


// --- GroupRepository & AuthorizationRepository Implementation ---

type PgGroupRepository struct {
	db *sqlx.DB
}

func NewPgGroupRepository(db *sqlx.DB) *PgGroupRepository {
	return &PgGroupRepository{db: db}
}

func (r *PgGroupRepository) CreateGroup(ctx context.Context, group *resource.Group) error {
	query := `INSERT INTO groups (id, tenant_id, name, description, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, group.ID, group.TenantID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt)
	return err
}

func (r *PgGroupRepository) GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*resource.Group, error) {
	var group resource.Group
	query := `SELECT * FROM groups WHERE tenant_id = $1 AND id = $2`
	err := r.db.GetContext(ctx, &group, query, tenantID, groupID)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *PgGroupRepository) AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	query := `INSERT INTO user_groups (user_id, group_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, groupID, tenantID)
	return err
}

func (r *PgGroupRepository) RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	query := `DELETE FROM user_groups WHERE user_id = $1 AND group_id = $2 AND tenant_id = $3`
	_, err := r.db.ExecContext(ctx, query, userID, groupID, tenantID)
	return err
}

func (r *PgGroupRepository) AddResourceToGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error {
	query := `INSERT INTO resource_groups (resource_id, group_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, resourceID, groupID, tenantID)
	return err
}

func (r *PgGroupRepository) RemoveResourceFromGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error {
	query := `DELETE FROM resource_groups WHERE resource_id = $1 AND group_id = $2 AND tenant_id = $3`
	_, err := r.db.ExecContext(ctx, query, resourceID, groupID, tenantID)
	return err
}

// CheckUserInSameGroupAsResource implements the AuthorizationRepository interface.
func (r *PgGroupRepository) CheckUserInSameGroupAsResource(ctx context.Context, tenantID, userID, resourceID uuid.UUID) (bool, error) {
	var hasAccess bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM user_groups ug
            JOIN resource_groups rg ON ug.group_id = rg.group_id
            WHERE ug.tenant_id = $1
              AND ug.user_id = $2
              AND rg.resource_id = $3
        )`
	err := r.db.GetContext(ctx, &hasAccess, query, tenantID, userID, resourceID)
	if err != nil {
		return false, err
	}
	return hasAccess, nil
}