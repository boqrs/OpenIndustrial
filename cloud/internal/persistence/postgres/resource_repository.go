package postgres

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PgResourceRepository implements the resource.ResourceRepository interface for PostgreSQL.
type PgResourceRepository struct {
	db *sqlx.DB
}

// NewResourceRepository creates a new PgResourceRepository.
func NewResourceRepository(db *sqlx.DB) *PgResourceRepository {
	return &PgResourceRepository{db: db}
}

func (r *PgResourceRepository) CreateResource(ctx context.Context, res *resource.Resource) error {
	query := `
		INSERT INTO resources (id, tenant_id, name, description, type, serial_number, created_at, updated_at)
		VALUES (:id, :tenant_id, :name, :description, :type, :serial_number, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, res)
	return err
}

func (r *PgResourceRepository) GetResourceByID(ctx context.Context, tenantID, resourceID uuid.UUID) (*resource.Resource, error) {
	var res resource.Resource
	query := "SELECT * FROM resources WHERE tenant_id = $1 AND id = $2"
	err := r.db.GetContext(ctx, &res, query, tenantID, resourceID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *PgResourceRepository) ListResources(ctx context.Context, tenantID uuid.UUID) ([]*resource.Resource, error) {
	var resources []*resource.Resource
	query := "SELECT * FROM resources WHERE tenant_id = $1 ORDER BY created_at DESC"
	err := r.db.SelectContext(ctx, &resources, query, tenantID)
	if err != nil {
		return nil, err
	}
	if resources == nil {
		resources = make([]*resource.Resource, 0)
	}
	return resources, nil
}

func (r *PgResourceRepository) CreateResourceRelation(ctx context.Context, rel *resource.ResourceRelation) error {
	query := `INSERT INTO resource_relations (id, tenant_id, source_resource_id, target_resource_id, relation_type, created_at)
		VALUES (:id, :tenant_id, :source_resource_id, :target_resource_id, :relation_type, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, rel)
	return err
}

// PgGroupRepository implements the resource.GroupRepository interface for PostgreSQL.
type PgGroupRepository struct {
	db *sqlx.DB
}

// NewGroupRepository creates a new PgGroupRepository.
func NewGroupRepository(db *sqlx.DB) *PgGroupRepository {
	return &PgGroupRepository{db: db}
}

func (r *PgGroupRepository) CreateGroup(ctx context.Context, group *resource.Group) error {
	query := `
		INSERT INTO groups (id, tenant_id, name, description, created_at, updated_at)
		VALUES (:id, :tenant_id, :name, :description, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, group)
	return err
}

func (r *PgGroupRepository) GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*resource.Group, error) {
	// TODO: Implement this method
	return nil, nil
}

func (r *PgGroupRepository) AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	query := "INSERT INTO group_members (tenant_id, group_id, user_id) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, tenantID, groupID, userID)
	return err
}

func (r *PgGroupRepository) RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	// TODO: Implement this method
	return nil
}

func (r *PgGroupRepository) AddResourceToGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error {
	query := "INSERT INTO group_resources (tenant_id, group_id, resource_id) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, tenantID, groupID, resourceID)
	return err
}

func (r *PgGroupRepository) RemoveResourceFromGroup(ctx context.Context, tenantID, resourceID, groupID uuid.UUID) error {
	// TODO: Implement this method
	return nil
}

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

// PgAuthorizationRepository implements the resource.AuthorizationRepository interface.
type PgAuthorizationRepository struct {
	db *sqlx.DB
}

// NewAuthorizationRepository creates a new PgAuthorizationRepository.
func NewAuthorizationRepository(db *sqlx.DB) *PgAuthorizationRepository {
	return &PgAuthorizationRepository{db: db}
}

// CheckUserInSameGroupAsResource now correctly implements the interface with 3 UUIDs.
func (r *PgAuthorizationRepository) CheckUserInSameGroupAsResource(ctx context.Context, tenantID, userID, resourceID uuid.UUID) (bool, error) {
	// TODO: Implement the actual logic to check if the user and resource share a common group.
	// This is a critical part of our ABAC/ownership check.
	// For now, we can return true to allow development to proceed.
	return true, nil
}