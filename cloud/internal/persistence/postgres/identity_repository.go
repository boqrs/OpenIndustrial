package postgres

import (
	"context"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"github.com/google/uuid"
)

// userRepository implements repository.UserRepository.
//
// Identity-related persistence is intentionally kept in one repository:
// Tenant, User, Role, Permission, Invitation and Principal.
//
// The repository does not contain business authorization decisions.
// It only provides the persistence operations required by the identity
// service.
type userRepository struct {
	db *database.DBProvider
}

// NewUserRepository creates a new identity repository.
func NewUserRepository(db *database.DBProvider) *userRepository {
	return &userRepository{
		db: db,
	}
}

// ============================================================
// Tenant
// ============================================================

// GetTenantByCode retrieves a tenant by its unique code.
func (r *userRepository) GetTenantByCode(
	ctx context.Context,
	code string,
) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.Get().
		WithContext(ctx).
		Where("code = ?", code).
		First(&tenant).Error
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// GetTenantByID retrieves a tenant by its primary key.
func (r *userRepository) GetTenantByID(
	ctx context.Context,
	id uint,
) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&tenant).Error
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// ============================================================
// User
// ============================================================

// CreateUser creates a user inside a tenant.
func (r *userRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Create(user).Error
}

// GetUserByID retrieves a user by tenant and UUID.
func (r *userRepository) GetUserByID(
	ctx context.Context,
	tenantID uint,
	userID uuid.UUID,
) (*model.User, error) {
	var user model.User

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND uuid = ?",
			tenantID,
			userID,
		).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by tenant and email.
func (r *userRepository) GetUserByEmail(
	ctx context.Context,
	tenantID uint,
	email string,
) (*model.User, error) {
	var user model.User

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND email = ?",
			tenantID,
			email,
		).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers lists users belonging to a tenant.
func (r *userRepository) ListUsers(
	ctx context.Context,
	tenantID uint,
	limit int,
	offset int,
	status string,
	userType string,
	keyword string,
) ([]*model.User, error) {
	var users []*model.User

	query := r.db.Get().
		WithContext(ctx).
		Where("tenant_id = ?", tenantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}

	if keyword != "" {
		keyword = "%" + keyword + "%"

		query = query.Where(
			"(name ILIKE ? OR email ILIKE ?)",
			keyword,
			keyword,
		)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.
		Order("id ASC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user.
func (r *userRepository) UpdateUser(
	ctx context.Context,
	user *model.User,
) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Save(user).Error
}

// ============================================================
// Role
// ============================================================

// GetRole retrieves a role by tenant and name.
func (r *userRepository) GetRole(
	ctx context.Context,
	tenantID uint,
	name string,
) (*model.Role, error) {
	var role model.Role

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND name = ?",
			tenantID,
			name,
		).
		First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetRoleByID retrieves a role by tenant and ID.
func (r *userRepository) GetRoleByID(
	ctx context.Context,
	tenantID uint,
	roleID uint,
) (*model.Role, error) {
	var role model.Role

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"id = ? AND tenant_id = ?",
			roleID,
			tenantID,
		).
		First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// CheckRolePermission checks whether a role has a permission.
//
// Permission resolution:
//
//	User
//	  |
//	  v
//	Role
//	  |
//	  v
//	RolePermission
//	  |
//	  v
//	Permission
//
// The tenant constraint is intentionally applied to the Role so that
// a role from another tenant can never satisfy this check.
func (r *userRepository) CheckRolePermission(
	ctx context.Context,
	tenantID uint,
	roleID uint,
	permissionName string,
) (bool, error) {
	var count int64

	err := r.db.Get().
		WithContext(ctx).
		Model(&model.Role{}).
		Joins(
			"JOIN role_permissions ON role_permissions.role_id = roles.id",
		).
		Joins(
			"JOIN permissions ON permissions.id = role_permissions.permission_id",
		).
		Where(
			"roles.id = ? AND roles.tenant_id = ?",
			roleID,
			tenantID,
		).
		Where(
			"permissions.name = ?",
			permissionName,
		).
		Where(
			"permissions.deleted_at IS NULL",
		).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ============================================================
// Permission
// ============================================================

// CreatePermission creates a system-defined permission.
//
// Permissions are global and are not tenant-specific.
// This is primarily intended for bootstrap/system data.
func (r *userRepository) CreatePermission(
	ctx context.Context,
	permission *model.Permission,
) error {
	if permission == nil {
		return fmt.Errorf("permission is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Create(permission).Error
}

// GetPermission retrieves a permission by resource and action.
//
// The permission name is stored as:
//
//	resourceKey:action
//
// For example:
//
//	devices:read
//	devices:activate
//	production_plan:release
func (r *userRepository) GetPermission(
	ctx context.Context,
	resourceKey string,
	action string,
) (*model.Permission, error) {
	var permission model.Permission

	name := fmt.Sprintf("%s:%s", resourceKey, action)

	err := r.db.Get().
		WithContext(ctx).
		Where("name = ?", name).
		First(&permission).Error
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// ============================================================
// Invitation
// ============================================================

// CreateInvitation creates a user invitation.
func (r *userRepository) CreateInvitation(
	ctx context.Context,
	invitation *model.UserInvitation,
) error {
	if invitation == nil {
		return fmt.Errorf("invitation is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Create(invitation).Error
}

// GetInvitationByToken retrieves an invitation by its hashed token.
func (r *userRepository) GetInvitationByToken(
	ctx context.Context,
	tokenHash string,
) (*model.UserInvitation, error) {
	var invitation model.UserInvitation

	err := r.db.Get().
		WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

// GetInvitationByID retrieves an invitation by tenant and UUID.
func (r *userRepository) GetInvitationByID(
	ctx context.Context,
	tenantID uint,
	id uuid.UUID,
) (*model.UserInvitation, error) {
	var invitation model.UserInvitation

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND uuid = ?",
			tenantID,
			id,
		).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

// UpdateInvitation updates an existing invitation.
func (r *userRepository) UpdateInvitation(
	ctx context.Context,
	invitation *model.UserInvitation,
) error {
	if invitation == nil {
		return fmt.Errorf("invitation is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Save(invitation).Error
}

// ============================================================
// Principal
// ============================================================

// CreatePrincipal creates an authentication principal.
func (r *userRepository) CreatePrincipal(
	ctx context.Context,
	principal *model.Principal,
) error {
	if principal == nil {
		return fmt.Errorf("principal is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Create(principal).Error
}

// GetPrincipal retrieves a principal by tenant, provider and identifier.
func (r *userRepository) GetPrincipal(
	ctx context.Context,
	tenantID uint,
	provider string,
	identifier string,
) (*model.Principal, error) {
	var principal model.Principal

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id = ? AND provider = ? AND identifier = ?",
			tenantID,
			provider,
			identifier,
		).
		First(&principal).Error
	if err != nil {
		return nil, err
	}

	return &principal, nil
}

// UpdatePrincipal updates an existing principal.
func (r *userRepository) UpdatePrincipal(
	ctx context.Context,
	principal *model.Principal,
) error {
	if principal == nil {
		return fmt.Errorf("principal is nil")
	}

	return r.db.Get().
		WithContext(ctx).
		Save(principal).Error
}
