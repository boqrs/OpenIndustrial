package postgres

import (
	"context"
	//"database/sql"
	"fmt"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	//"golang.org/x/crypto/bcrypt"

)

// TenantRepository implements the identity.TenantRepository interface for PostgreSQL.
type TenantRepository struct {
	db *sqlx.DB
}

func NewTenantRepository(db *sqlx.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) CreateTenant(ctx context.Context, tenant *identity.Tenant) error {
	query := `INSERT INTO tenants (name, status) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, tenant.Name, tenant.Status).StructScan(tenant)
}

// GetTenantByID is not implemented for brevity
func (r *TenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*identity.Tenant, error) {
	panic("not implemented")
}

// UpdateTenant is not implemented for brevity
func (r *TenantRepository) UpdateTenant(ctx context.Context, tenant *identity.Tenant) error {
	panic("not implemented")
}


// UserRepository implements the identity.UserRepository interface for PostgreSQL.
type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *identity.User) error {
	query := `INSERT INTO users (tenant_id, user_type, profile) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, user.TenantID, user.UserType, user.Profile).StructScan(user)
}

func (r *UserRepository) CreatePrincipal(ctx context.Context, p *identity.Principal) error {
	query := `INSERT INTO principals (user_id, tenant_id, provider, identifier, credential) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, p.UserID, p.TenantID, p.Provider, p.Identifier, p.Credential).StructScan(p)
}

func (r *UserRepository) GetPrincipal(ctx context.Context, tenantID uuid.UUID, provider, identifier string) (*identity.Principal, error) {
	var p identity.Principal
	query := `SELECT * FROM principals WHERE tenant_id = $1 AND provider = $2 AND identifier = $3`
	err := r.db.GetContext(ctx, &p, query, tenantID, provider, identifier)
	return &p, err
}


// RoleRepository implements the identity.RoleRepository interface for PostgreSQL.
type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// func (r *RoleRepository) GetRoleByName(ctx context.Context, tenantID uuid.UUID, name string) (*identity.Role, error) {
// 	var role identity.Role
//     // Handling NULL tenant_id for system roles
//     query := `SELECT * FROM roles WHERE name = $1 AND (tenant_id = $2 OR tenant_id IS NULL) ORDER BY is_system_role DESC LIMIT 1`
// 	err := r.db.GetContext(ctx, &role, query, name, tenantID)
// 	return &role, err
// }

func (r *RoleRepository) AddUserToRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error {
	query := `INSERT INTO user_roles (user_id, role_id, tenant_id) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, userID, roleID, tenantID)
	return err
}

// Other RoleRepository methods are not implemented for brevity
func (r *RoleRepository) CreateRole(ctx context.Context, role *identity.Role) error {
	panic("not implemented")
}
func (r *RoleRepository) GetRoleByID(ctx context.Context, tenantID, roleID uuid.UUID) (*identity.Role, error) {
	panic("not implemented")
}
func (r *RoleRepository) RemoveUserFromRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error {
	panic("not implemented")
}
func (r *RoleRepository) AddPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	panic("not implemented")
}

// Other UserRepository methods are not implemented for brevity
func (r *UserRepository) GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*identity.User, error) {
	var user identity.User
	query := `SELECT * FROM users WHERE id = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &user, query, userID, tenantID)
	return &user, err
}

func (r *RoleRepository) GetRoleByName(ctx context.Context, tenantID uuid.UUID, name string) (*identity.Role, error) {
	var role identity.Role
    // This query prefers tenant-specific roles over system roles if names conflict.
    query := `
        SELECT * FROM roles 
        WHERE name = $1 AND (tenant_id = $2 OR is_system_role = true)
        ORDER BY tenant_id IS NULL ASC, is_system_role ASC 
        LIMIT 1`
	err := r.db.GetContext(ctx, &role, query, name, tenantID)
	return &role, err
}

// ListUsersParams defines parameters for listing users.
type ListUsersParams struct {
	Limit  int
	Offset int
}

// ListUsers retrieves a paginated list of users for a tenant.
// 注意：这里的参数类型 identity.ListUsersRepoParams 与接口定义一致
func (r *UserRepository) ListUsers(ctx context.Context, tenantID uuid.UUID, params identity.ListUsersRepoParams) ([]identity.User, error) {
	var users []identity.User
	query := `SELECT * FROM users WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &users, query, tenantID, params.Limit, params.Offset)
	return users, err
}

// ListRoles retrieves all roles for a given tenant, including system roles.
func (r *RoleRepository) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]identity.Role, error) {
	var roles []identity.Role
	// 这个查询会获取租户自己的角色以及所有系统角色
	query := `SELECT * FROM roles WHERE tenant_id = $1 OR is_system_role = true ORDER BY name`
	err := r.db.SelectContext(ctx, &roles, query, tenantID)
	return roles, err
}

// UpdateUser updates a user's information in the database.
func (r *UserRepository) UpdateUser(ctx context.Context, user *identity.User) error {
	query := `UPDATE users SET user_type = $1, profile = $2, updated_at = NOW() WHERE id = $3 AND tenant_id = $4`
	_, err := r.db.ExecContext(ctx, query, user.UserType, user.Profile, user.ID, user.TenantID)
	return err
}

// DeleteUser removes a user from the database.
// Note: This is a hard delete. A soft delete (marking as inactive) is often preferred.
func (r *UserRepository) DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error {
	// We should also delete from the 'principals' and 'user_roles' tables in a transaction.
	// For simplicity, we only delete from the 'users' table here.
	query := `DELETE FROM users WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, tenantID)
	return err
}



// PermissionRepository implements the identity.PermissionRepository interface for PostgreSQL.
type PermissionRepository struct {
	db *sqlx.DB
}

func NewPermissionRepository(db *sqlx.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// CheckPermissionForUser checks if a user has a specific permission through their roles.
func (r *PermissionRepository) CheckPermissionForUser(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	var hasPermission bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM user_roles ur
            JOIN role_permissions rp ON ur.role_id = rp.role_id
            JOIN permissions p ON rp.permission_id = p.id
            WHERE ur.user_id = $1 AND p.name = $2
        )`
	err := r.db.GetContext(ctx, &hasPermission, query, userID, permissionName)
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}

// CreatePermission adds a new permission to the database.
func (r *PermissionRepository) CreatePermission(ctx context.Context, p *identity.Permission) error {
	query := `INSERT INTO permissions (id, name, description) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Description)
	return err
}

// GetPermission retrieves a permission by its key and action.
func (r *PermissionRepository) GetPermission(ctx context.Context, resourceKey, action string) (*identity.Permission, error) {
	var p identity.Permission
	permissionName := fmt.Sprintf("%s:%s", resourceKey, action)
	query := `SELECT * FROM permissions WHERE name = $1`
	err := r.db.GetContext(ctx, &p, query, permissionName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPermissionsByRole retrieves all permissions associated with a specific role.
func (r *PermissionRepository) ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*identity.Permission, error) {
	var permissions []*identity.Permission
	query := `
		SELECT p.*
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name`
	err := r.db.SelectContext(ctx, &permissions, query, roleID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}