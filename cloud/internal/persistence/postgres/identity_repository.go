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
// User
// =====================
func (r *userRepository) GetTenantByID(
	ctx context.Context,
	id uint,
) (
	*model.Tenant,
	error,
) {
	var tenant model.Tenant

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"id=?",
			id,
		).
		First(&tenant).
		Error

	return &tenant, err
}
func (r *userRepository) GetTenantByCode(
	ctx context.Context,
	code string,
) (
	*model.Tenant,
	error,
) {
	var tenant model.Tenant

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"code=?",
			code,
		).
		First(&tenant).
		Error

	return &tenant, err
}

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
			"uuid=? AND tenant_id=?",
			userID,
			tenantID,
		).
		First(&user).
		Error

	return &user, err
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
			"tenant_id=? AND email=?",
			tenantID,
			email,
		).
		First(&user).
		Error

	return &user, err
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
) (
	[]*model.User,
	error,
) {
	var users []*model.User

	q := r.db.Get().
		WithContext(ctx).
		Where(
			"tenant_id=?",
			tenantID,
		)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if keyword != "" {
		q = q.Where("keyword = ?", keyword)
	}

	if userType != "" {
		q = q.Where("userType = ?", userType)
	}

	err := q.
		Order(
			"created_at DESC",
		).
		Limit(
			limit,
		).
		Offset(
			offset,
		).
		Find(&users).
		Error

	return users, err
}

// =====================
// Principal
// =====================

func (r *userRepository) CreatePrincipal(
	ctx context.Context,
	p *model.Principal,
) error {

	return r.db.Get().
		WithContext(ctx).
		Create(p).
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
			"tenant_id=? AND provider=? AND identifier=?",
			tenantID,
			provider,
			identifier,
		).
		First(&principal).
		Error

	return &principal, err
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
			"token_hash=?",
			tokenHash,
		).
		First(&invitation).
		Error

	return &invitation, err
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
func (r *userRepository) GetInvitationByID(
	ctx context.Context,
	tenantID uint,
	id uuid.UUID,
) (
	*model.UserInvitation,
	error,
) {
	var invitation model.UserInvitation

	err := r.db.Get().
		WithContext(ctx).
		Where(
			"uuid=? AND tenant_id=?",
			id,
			tenantID,
		).
		First(&invitation).
		Error

	return &invitation, err
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

type permissionRepository struct {
	db *database.DBProvider
}

func NewPermissionRepository(db *database.DBProvider) *permissionRepository {
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
