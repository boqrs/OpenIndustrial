package postgres

import (
	"context"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// GroupRepository implements the identity.GroupRepository interface for PostgreSQL.
type GroupRepository struct {
	db *sqlx.DB
}

// NewGroupRepository creates a new GroupRepository.
func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// CreateGroup creates a new group in the database.
func (r *GroupRepository) CreateGroup(ctx context.Context, group *identity.Group) error {
	query := `
		INSERT INTO groups (id, tenant_id, name, description, created_at, updated_at)
		VALUES (:id, :tenant_id, :name, :description, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, group)
	return err
}

// GetGroupByID retrieves a group by its ID.
func (r *GroupRepository) GetGroupByID(ctx context.Context, tenantID, groupID uuid.UUID) (*identity.Group, error) {
	var group identity.Group
	query := `SELECT * FROM groups WHERE id = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &group, query, groupID, tenantID)
	return &group, err
}

// AddUserToGroup adds a user to a group in the group_members table.
func (r *GroupRepository) AddUserToGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	query := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

// RemoveUserFromGroup removes a user from a group.
func (r *GroupRepository) RemoveUserFromGroup(ctx context.Context, tenantID, userID, groupID uuid.UUID) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

// ListGroupsByUserID retrieves all groups a specific user is a member of.
// This is the core implementation for our API endpoint.
func (r *GroupRepository) ListGroupsByUserID(ctx context.Context, tenantID, userID uuid.UUID) ([]*identity.Group, error) {
	var groups []*identity.Group
	query := `
		SELECT g.*
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE g.tenant_id = $1 AND gm.user_id = $2
		ORDER BY g.name`
	err := r.db.SelectContext(ctx, &groups, query, tenantID, userID)
	return groups, err
}