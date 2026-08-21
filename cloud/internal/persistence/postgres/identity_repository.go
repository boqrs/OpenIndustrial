package postgres

import (
	"context"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"github.com/boqrs/nexus/database"
)

// --- TenantRepository ---

type tenantRepository struct {
	db *database.DBProvider
}

func NewTenantRepository(db *database.DBProvider) identity.TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	return r.db.Get().WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Get().WithContext(ctx).Where("uuid = ?", id).First(&tenant).Error
	return &tenant, err
}

func (r *tenantRepository) UpdateTenant(ctx context.Context, tenant *model.Tenant) error {
	// GORM's Save method will update all fields of the struct,
	// or create a new record if it does not exist.
	// It uses the primary key value to find the record.
	return r.db.Get().WithContext(ctx).Save(tenant).Error
}
// --- UserRepository ---

type userRepository struct {
	db *database.DBProvider
}

func NewUserRepository(	db *database.DBProvider) identity.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.Get().WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Get().WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", userID, tenantID).
		First(&user).Error
	return &user, err
}

func (r *userRepository) CreatePrincipal(ctx context.Context, p *model.Principal) error {
	return r.db.Get().WithContext(ctx).Create(p).Error
}

func (r *userRepository) GetPrincipal(ctx context.Context, tenantID uuid.UUID, provider, identifier string) (*model.Principal, error) {
	var p model.Principal
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ? AND provider = ? AND identifier = ?", tenantID, provider, identifier).
		First(&p).Error
	return &p, err
}

func (r *userRepository) ListUsers(ctx context.Context, tenantID uuid.UUID, params identity.ListUsersRepoReq) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Get().WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&users).Error
	return users, err
}

// --- RoleRepository ---

type roleRepository struct {
	db *database.DBProvider
}

func NewRoleRepository(	db *database.DBProvider) identity.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.db.Get().WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetRoleByID(ctx context.Context, tenantID, roleID uuid.UUID) (*model.Role, error) {
	var role model.Role
	err := r.db.Get().WithContext(ctx).
		Where("uuid = ? AND tenant_id = ?", roleID, tenantID).
		First(&role).Error
	return &role, err
}

func (r *roleRepository) GetRoleByName(ctx context.Context, tenantID uuid.UUID, name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Get().WithContext(ctx).
		Where("name = ? AND (tenant_id = ? OR is_system = true)", name, tenantID).
		Order("is_system ASC"). // Prioritize non-system roles
		First(&role).Error
	return &role, err
}

func (r *roleRepository) AddUserToRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error {
	user := model.User{}
	if err := r.db.Get().WithContext(ctx).First(&user, "uuid = ?", userID).Error; err != nil {
		return err // User not found
	}

	role := model.Role{}
	if err := r.db.Get().WithContext(ctx).First(&role, "uuid = ?", roleID).Error; err != nil {
		return err // Role not found
	}

	return r.db.Get().WithContext(ctx).Model(&user).Association("Roles").Append(&role)
}

// --- PermissionRepository ---

type permissionRepository struct {
	db *database.DBProvider
}

func NewPermissionRepository(	db *database.DBProvider) identity.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) CheckPermissionForUser(ctx context.Context, userID uuid.UUID, permissionName string) (bool, error) {
	var count int64
	err := r.db.Get().WithContext(ctx).Model(&model.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.uuid").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("users.uuid = ? AND permissions.name = ?", userID, permissionName).
		Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) CreatePermission(ctx context.Context, p *model.Permission) error {
	return r.db.Get().WithContext(ctx).Create(p).Error
}

func (r *permissionRepository) GetPermission(ctx context.Context, resourceKey, action string) (*model.Permission, error) {
	var p model.Permission
	permissionName := fmt.Sprintf("%s:%s", resourceKey, action)
	err := r.db.Get().WithContext(ctx).Where("name = ?", permissionName).First(&p).Error
	return &p, err
}

func (r *permissionRepository) ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*model.Permission, error) {
	var role model.Role
	err := r.db.Get().WithContext(ctx).
		Preload("Permissions").
		Where("uuid = ?", roleID).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
