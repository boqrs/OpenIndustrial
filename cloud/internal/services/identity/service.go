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
		repo: repo,

		jwtSecret: jwtSecret,

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
		return nil, err
	}

	user, err := s.repo.GetUserByEmail(
		ctx,
		tenant.ID,
		req.Email,
	)

	if err != nil {
		return nil, ErrInvalidCredential
	}

	if user.Status != model.UserStatusActive {
		return nil, ErrUserDisabled
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
		return nil, err
	}

	user, err :=
		s.repo.GetUserByID(
			ctx,
			claims.TenantID,
			claims.UserID,
		)

	if err != nil {
		return nil, err
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

	user := &model.User{

		UUID: uuid.New(),

		TenantID: req.TenantID,

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

	_ = token

	return nil
}

func (s *service) AcceptInvitation(
	ctx context.Context,
	req AcceptInvitationRequest,
) error {

	invitation, err :=
		s.repo.GetInvitationByToken(
			ctx,
			req.Token,
		)

	if err != nil {
		return ErrInvitationInvalid
	}

	if invitation.UsedAt != nil ||
		time.Now().After(invitation.ExpiresAt) {

		return ErrInvitationInvalid
	}

	user, err :=
		s.repo.GetUserByID(
			ctx,
			invitation.TenantID,
			invitation.User.UUID,
		)

	if err != nil {
		return err
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

	return nil, errors.New(
		"repository method required",
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

	user := &model.User{

		UUID: req.UUID,

		TenantID: req.TenantID,

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

	// repository负责tenant隔离
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

	// 当前repo缺少UpdatePrincipal
	//
	// 后续补充

	return errors.New(
		"UpdatePrincipal repository method required",
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

	return errors.New(
		"UpdatePrincipal repository method required",
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
