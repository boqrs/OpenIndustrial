package identity

import (
	"context"
	"errors"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredential = errors.New("invalid credential")
	ErrUserDisabled      = errors.New("user disabled")

	ErrInvitationInvalid = errors.New("invalid invitation")
	ErrInvitationExpired = errors.New("invitation expired")
	ErrInvitationUsed    = errors.New("invitation already used")

	ErrRoleInvalid   = errors.New("invalid role")
	ErrUserNotActive = errors.New("user is not active")
)

type service struct {
	repo UserRepository

	jwtSecret string

	accessExpire  time.Duration
	refreshExpire time.Duration
}

func NewService(
	repo UserRepository,
	jwtSecret string,
) Service {

	return &service{
		repo:          repo,
		jwtSecret:     jwtSecret,
		accessExpire:  time.Hour * 2,
		refreshExpire: time.Hour * 24 * 7,
	}
}

// =====================================================
// Authentication
// =====================================================

func (s *service) Login(
	ctx context.Context,
	req LoginRequest,
) (
	*LoginResponse,
	error,
) {

	tenant, err := s.repo.GetTenantByCode(
		ctx,
		req.TenantCode,
	)

	if err != nil {
		return nil, ErrInvalidCredential
	}

	user, err := s.repo.GetUserByEmail(
		ctx,
		tenant.ID,
		req.Email,
	)

	if err != nil {
		return nil, ErrInvalidCredential
	}

	if user.Status == model.UserStatusDisabled {
		return nil, ErrUserDisabled
	}

	if user.Status != model.UserStatusActive {
		return nil, ErrInvalidCredential
	}

	principal, err := s.repo.GetPrincipal(
		ctx,
		tenant.ID,
		model.PrincipalProviderPassword,
		req.Email,
	)

	if err != nil {
		return nil, ErrInvalidCredential
	}

	if principal.Status != model.UserStatusActive {
		return nil, ErrInvalidCredential
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(principal.SecretHash),
		[]byte(req.Password),
	) != nil {

		return nil, ErrInvalidCredential
	}

	accessToken, err :=
		GenerateAccessToken(
			s.jwtSecret,
			user,
			s.accessExpire,
		)

	if err != nil {
		return nil, err
	}

	refreshToken, err :=
		GenerateRefreshToken(
			s.jwtSecret,
			user,
			s.refreshExpire,
		)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{

		AccessToken: accessToken,

		RefreshToken: refreshToken,

		User: user,
	}, nil
}

func (s *service) Logout(
	ctx context.Context,
	req LogoutRequest,
) error {

	// 当前JWT无状态
	//
	// 后续接入redis blacklist
	//
	// 当前仅完成接口闭环

	return nil
}

func (s *service) RefreshToken(
	ctx context.Context,
	req RefreshTokenRequest,
) (
	*LoginResponse,
	error,
) {

	claims, err :=
		ParseRefreshToken(
			s.jwtSecret,
			req.RefreshToken,
		)

	if err != nil {
		return nil, ErrInvalidCredential
	}

	user, err :=
		s.repo.GetUserByID(
			ctx,
			claims.TenantID,
			claims.UserID,
		)

	if err != nil {
		return nil, ErrInvalidCredential
	}

	if user.Status != model.UserStatusActive {
		if user.Status == model.UserStatusDisabled {
			return nil, ErrUserDisabled
		}

		return nil, ErrUserNotActive
	}

	accessToken, err :=
		GenerateAccessToken(
			s.jwtSecret,
			user,
			s.accessExpire,
		)

	if err != nil {
		return nil, err
	}

	refreshToken, err :=
		GenerateRefreshToken(
			s.jwtSecret,
			user,
			s.refreshExpire,
		)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{

		AccessToken: accessToken,

		RefreshToken: refreshToken,

		User: user,
	}, nil
}

// =====================================================
// Invitation
// =====================================================

func (s *service) InviteUser(
	ctx context.Context,
	req InviteUserRequest,
) error {

	if req.TenantID == 0 ||
		req.CreatedBy == 0 ||
		req.RoleID == 0 {

		return ErrInvalidCredential
	}

	// Role必须属于当前Tenant。
	_, err :=
		s.repo.GetRoleByID(
			ctx,
			req.TenantID,
			req.RoleID,
		)

	if err != nil {
		return ErrRoleInvalid
	}

	user := &model.User{

		UUID: uuid.New(),

		TenantID: req.TenantID,

		RoleID: req.RoleID,

		Name: req.Name,

		Email: req.Email,

		UserType: model.UserTypeEmployee,

		Status: model.UserStatusInvited,
	}

	if err :=
		s.repo.CreateUser(
			ctx,
			user,
		); err != nil {

		return err
	}

	token, hash, err :=
		GenerateInvitationToken()

	if err != nil {
		return err
	}

	invitation :=
		&model.UserInvitation{

			UUID: uuid.New(),

			TenantID: req.TenantID,

			UserID: user.ID,

			Email: req.Email,

			TokenHash: hash,

			ExpiresAt: time.Now().
				Add(24 * time.Hour),

			CreatedBy: req.CreatedBy,
		}

	if err :=
		s.repo.CreateInvitation(
			ctx,
			invitation,
		); err != nil {

		return err
	}

	// TODO:
	//
	// Email Service
	//
	// 发送:
	//
	// /register?token=<token>
	//
	// 当前阶段先保留token生成逻辑。
	// 后续接入Email Provider后，由通知层负责发送。
	_ = token

	return nil
}

func (s *service) AcceptInvitation(
	ctx context.Context,
	req AcceptInvitationRequest,
) error {

	tokenHash, err :=
		HashToken(
			req.Token,
		)

	if err != nil {
		return ErrInvitationInvalid
	}

	invitation, err :=
		s.repo.GetInvitationByToken(
			ctx,
			tokenHash,
		)

	if err != nil {
		return ErrInvitationInvalid
	}

	if invitation.UsedAt != nil {
		return ErrInvitationUsed
	}

	if !invitation.ExpiresAt.After(time.Now()) {
		return ErrInvitationExpired
	}

	user, err := s.repo.GetUserByEmail(
		ctx,
		invitation.TenantID,
		invitation.Email,
	)

	if err != nil {
		return ErrInvitationInvalid
	}

	if user.Status != model.UserStatusInvited {
		return ErrInvitationInvalid
	}

	// 防止用户邮箱已经被其他流程修改。
	if user.Email != invitation.Email {
		return ErrInvitationInvalid
	}

	hash, err :=
		HashPassword(
			req.Password,
		)

	if err != nil {
		return err
	}

	principal :=
		&model.Principal{

			UUID: uuid.New(),

			TenantID: user.TenantID,

			UserID: user.ID,

			Provider: model.PrincipalProviderPassword,

			Identifier: user.Email,

			SecretHash: hash,

			Status: model.UserStatusActive,
		}

	if err :=
		s.repo.CreatePrincipal(
			ctx,
			principal,
		); err != nil {

		return err
	}

	user.Status =
		model.UserStatusActive

	if req.Name != "" {
		user.Name = req.Name
	}

	if err :=
		s.repo.UpdateUser(
			ctx,
			user,
		); err != nil {

		return err
	}

	now := time.Now()

	invitation.UsedAt = &now

	return s.repo.UpdateInvitation(
		ctx,
		invitation,
	)
}

func (s *service) GetInvitation(
	ctx context.Context,
	tenantID uint,
	id uuid.UUID,
) (
	*model.UserInvitation,
	error,
) {

	return s.repo.GetInvitationByID(
		ctx,
		tenantID,
		id,
	)
}

// =====================================================
// User Management
// =====================================================

func (s *service) CreateUser(
	ctx context.Context,
	req CreateUserRequest,
) (
	*model.User,
	error,
) {

	if req.TenantID == 0 ||
		req.RoleID == 0 {

		return nil, ErrRoleInvalid
	}

	// Role必须属于当前Tenant。
	_, err :=
		s.repo.GetRoleByID(
			ctx,
			req.TenantID,
			req.RoleID,
		)

	if err != nil {
		return nil, ErrRoleInvalid
	}

	user := &model.User{

		UUID: req.UUID,

		TenantID: req.TenantID,

		RoleID: req.RoleID,

		Name: req.Name,

		Email: req.Email,

		UserType: req.UserType,

		Status: req.Status,
	}

	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}

	if user.UserType == "" {
		user.UserType =
			model.UserTypeEmployee
	}

	if user.Status == "" {
		user.Status =
			model.UserStatusInvited
	}

	if err :=
		s.repo.CreateUser(
			ctx,
			user,
		); err != nil {

		return nil, err
	}

	return user, nil
}

func (s *service) GetUser(
	ctx context.Context,
	tenantID uint,
	userID uuid.UUID,
) (
	*model.User,
	error,
) {

	return s.repo.GetUserByID(
		ctx,
		tenantID,
		userID,
	)
}

func (s *service) GetUserByEmail(
	ctx context.Context,
	tenantID uint,
	email string,
) (
	*model.User,
	error,
) {

	return s.repo.GetUserByEmail(
		ctx,
		tenantID,
		email,
	)
}

func (s *service) ListUsers(
	ctx context.Context,
	req ListUsersRequest,
) (
	[]*model.User,
	error,
) {

	return s.repo.ListUsers(
		ctx,
		req.TenantID,
		req.Limit,
		req.Offset,
		req.Status,
		req.UserType,
		req.Keyword,
	)
}

func (s *service) UpdateUser(
	ctx context.Context,
	req UpdateUserRequest,
) error {

	user, err :=
		s.repo.GetUserByID(
			ctx,
			req.TenantID,
			req.UserID,
		)

	if err != nil {
		return err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	return s.repo.UpdateUser(
		ctx,
		user,
	)
}

func (s *service) DisableUser(
	ctx context.Context,
	req DisableUserRequest,
) error {

	user, err :=
		s.repo.GetUserByID(
			ctx,
			req.TenantID,
			req.UserID,
		)

	if err != nil {
		return err
	}

	user.Status =
		model.UserStatusDisabled

	return s.repo.UpdateUser(
		ctx,
		user,
	)
}

func (s *service) EnableUser(
	ctx context.Context,
	req EnableUserRequest,
) error {

	user, err :=
		s.repo.GetUserByID(
			ctx,
			req.TenantID,
			req.UserID,
		)

	if err != nil {
		return err
	}

	user.Status =
		model.UserStatusActive

	return s.repo.UpdateUser(
		ctx,
		user,
	)
}

// =====================================================
// Principal
// =====================================================

func (s *service) UpdatePassword(
	ctx context.Context,
	req UpdatePasswordRequest,
) error {

	user, err :=
		s.repo.GetUserByID(
			ctx,
			req.TenantID,
			req.UserID,
		)

	if err != nil {
		return err
	}

	if user.Status != model.UserStatusActive {
		return ErrUserNotActive
	}

	principal, err :=
		s.repo.GetPrincipal(
			ctx,
			req.TenantID,
			model.PrincipalProviderPassword,
			user.Email,
		)

	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(principal.SecretHash),
		[]byte(req.OldPassword),
	) != nil {

		return errors.New(
			"old password incorrect",
		)
	}

	hash, err :=
		HashPassword(
			req.NewPassword,
		)

	if err != nil {
		return err
	}

	principal.SecretHash =
		hash

	return s.repo.UpdatePrincipal(
		ctx,
		principal,
	)
}

func (s *service) ResetPassword(
	ctx context.Context,
	req ResetPasswordRequest,
) error {

	user, err :=
		s.repo.GetUserByID(
			ctx,
			req.TenantID,
			req.UserID,
		)

	if err != nil {
		return err
	}

	if user.Status != model.UserStatusActive {
		return ErrUserNotActive
	}

	principal, err :=
		s.repo.GetPrincipal(
			ctx,
			req.TenantID,
			model.PrincipalProviderPassword,
			user.Email,
		)

	if err != nil {
		return err
	}

	hash, err :=
		HashPassword(
			req.NewPassword,
		)

	if err != nil {
		return err
	}

	principal.SecretHash =
		hash

	return s.repo.UpdatePrincipal(
		ctx,
		principal,
	)
}

// =====================================================
// Tenant
// =====================================================

func (s *service) GetTenant(
	ctx context.Context,
	tenantID uint,
) (
	*model.Tenant,
	error,
) {

	return s.repo.GetTenantByID(
		ctx,
		tenantID,
	)
}

func (s *service) GetTenantByCode(
	ctx context.Context,
	code string,
) (
	*model.Tenant,
	error,
) {

	return s.repo.GetTenantByCode(
		ctx,
		code,
	)
}
