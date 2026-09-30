package postgres

import (
	"context"
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/nexus/database"
	"github.com/google/uuid"
)

type userRepository struct {
	db *database.DBProvider
}

func NewUserRepository(
	db *database.DBProvider,
) *userRepository {
	return &userRepository{
		db: db,
	}
}

// =====================
// Tenant
// =====================

func (r *userRepository) GetTenantByID(
	ctx context.Context,
	id uint,
) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.Get().
		WithContext(ctx).
		Where("id = ?", id).
		First(&tenant).
		Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *userRepository) GetTenantByCode(
	ctx context.Context,
	code string,
) (*model.Tenant, error) {
	var tenant model.Tenant

	err := r.db.Get().
		WithContext(ctx).
		Where("code = ?", code).
		First(&tenant).
		Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// =====================
// User
// =====================

func (r *userRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(user).
		Error
}

func (r *userRepository) GetUserByID(
	ctx context.Context,
	tenantID uint,
	userID uuid.UUID,
) (*model.User, error) {
	var user model.User

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"uuid = ? AND tenant_id = ?",
			userID,
			tenantID,
		).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

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
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(
	ctx context.Context,
	user *model.User,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(user).
		Error
}

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

	q := r.db.Get().
		WithContext(ctx).
		Where("tenant_id = ?", tenantID)

	if status != "" {
		q = q.Where("status = ?", status)
	}

	if userType != "" {
		q = q.Where("user_type = ?", userType)
	}

	if keyword != "" {
		keywordPattern := "%" + keyword + "%"

		q = q.Where(
			"(name ILIKE ? OR email ILIKE ?)",
			keywordPattern,
			keywordPattern,
		)
	}

	err := q.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).
		Error

	return users, err
}

// =====================
// Role
// =====================

// GetRole returns a role belonging to the specified tenant.
//
// Role names are unique within a tenant.
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
		First(&role).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetRoleByID returns a role only when it belongs to the specified tenant.
//
// tenantID is intentionally part of the query to prevent a role from
// another tenant from being used.
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
		First(&role).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// CheckRolePermission checks whether a role belonging to the specified
// tenant has the requested permission.
//
// The permission relationship is:
//
// roles
//   -> role_permissions
//       -> permissions
//
// Both tenantID and roleID are required so a role ID from another tenant
// cannot be used to obtain permissions.
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
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// =====================
// Principal
// =====================

func (r *userRepository) CreatePrincipal(
	ctx context.Context,
	principal *model.Principal,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(principal).
		Error
}

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
		First(&principal).
		Error

	if err != nil {
		return nil, err
	}

	return &principal, nil
}

func (r *userRepository) UpdatePrincipal(
	ctx context.Context,
	principal *model.Principal,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(principal).
		Error
}

// =====================
// Invitation
// =====================

func (r *userRepository) CreateInvitation(
	ctx context.Context,
	invitation *model.UserInvitation,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(invitation).
		Error
}

func (r *userRepository) GetInvitationByToken(
	ctx context.Context,
	tokenHash string,
) (*model.UserInvitation, error) {
	var invitation model.UserInvitation

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"token_hash = ?",
			tokenHash,
		).
		First(&invitation).
		Error

	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (r *userRepository) GetInvitationByID(
	ctx context.Context,
	tenantID uint,
	id uuid.UUID,
) (*model.UserInvitation, error) {
	var invitation model.UserInvitation

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"uuid = ? AND tenant_id = ?",
			id,
			tenantID,
		).
		First(&invitation).
		Error

	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (r *userRepository) UpdateInvitation(
	ctx context.Context,
	invitation *model.UserInvitation,
) error {
	return r.db.Get().
		WithContext(ctx).
		Save(invitation).
		Error
}

// =====================
// Permission
// =====================

// CreatePermission is kept for system/bootstrap usage.
//
// Permissions are global and are not tenant-specific.
func (r *userRepository) CreatePermission(
	ctx context.Context,
	permission *model.Permission,
) error {
	return r.db.Get().
		WithContext(ctx).
		Create(permission).
		Error
}

func (r *userRepository) GetPermission(
	ctx context.Context,
	resourceKey string,
	action string,
) (*model.Permission, error) {
	var permission model.Permission

	permissionName := fmt.Sprintf(
		"%s:%s",
		resourceKey,
		action,
	)

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"name = ?",
			permissionName,
		).
		First(&permission).
		Error

	if err != nil {
		return nil, err
	}

	return &permission, nil
}