package postgres

import (
	"context"
	"errors"

	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PgResourceRepository implements the resource.ResourceRepository interface using PostgreSQL.
type PgResourceRepository struct {
	db *sqlx.DB
}

func NewPgResourceRepository(db *sqlx.DB) *PgResourceRepository {
	return &PgResourceRepository{db: db}
}

func (r *PgResourceRepository) CreateResource(ctx context.Context, res *resource.Resource) error {
	query := `INSERT INTO resources (id, tenant_id, name, description, type, serial_number, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, res.ID, res.TenantID, res.Name, res.Description, res.Type, res.SerialNumber, res.CreatedAt, res.UpdatedAt)
	return err
}

func (r *PgResourceRepository) CreateResourceRelation(ctx context.Context, rel *resource.ResourceRelation) error {
	query := `INSERT INTO resource_relations (id, tenant_id, source_resource_id, target_resource_id, relation_type, created_at)
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, rel.ID, rel.TenantID, rel.SourceResourceID, rel.TargetResourceID, rel.RelationType, rel.CreatedAt)
	return err
}

// PgGroupRepository implements the resource.GroupRepository and resource.AuthorizationRepository interfaces.
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
	query := `SELECT id, tenant_id, name, description, created_at, updated_at FROM groups WHERE id = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &group, query, groupID, tenantID)
	return &group, err
}

func (r *PgGroupRepository) AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	query := `INSERT INTO group_members (tenant_id, group_id, user_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, tenantID, groupID, userID)
	return err
}

// ListGroupsByUserID retrieves all groups a specific user belongs to.
func (r *PgGroupRepository) ListGroupsByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*resource.Group, error) {
	groups := make([]*resource.Group, 0)
	query := `
		SELECT g.id, g.tenant_id, g.name, g.description, g.created_at, g.updated_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE g.tenant_id = $1 AND gm.user_id = $2
		ORDER BY g.name ASC`
	err := r.db.SelectContext(ctx, &groups, query, tenantID, userID)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// --- STUB IMPLEMENTATIONS ---
// We will implement these in future steps when they are needed.

func (r *PgGroupRepository) RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	// TODO: Implement in a future step
	return errors.New("not implemented")
}

func (r *PgGroupRepository) AddResourceToGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error {
	// TODO: Implement in a future step
	return errors.New("not implemented")
}

func (r *PgGroupRepository) RemoveResourceFromGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error {
	// TODO: Implement in a future step
	return errors.New("not implemented")
}

func (r *PgGroupRepository) CheckUserInSameGroupAsResource(ctx context.Context, tenantID, userID, resourceID uuid.UUID) (bool, error) {
	// TODO: Implement in a future step
	return false, errors.New("not implemented")
}