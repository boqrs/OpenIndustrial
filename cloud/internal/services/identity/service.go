package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/notification"
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

	ErrAdminRequired         = errors.New("administrator privileges required")
	ErrAdminCannotBeDisabled = errors.New("administrator cannot be disabled")
	ErrAdminCannotBeModified = errors.New("administrator cannot be modified")

	ErrEmailAlreadyExists = errors.New("email already exists")
)

type service struct {
	repo          UserRepository
	notification  notification.Service
	jwtSecret     string
	accessExpire  time.Duration
	refreshExpire time.Duration
}

func NewService(
	repo UserRepository,
	notificationService notification.Service,
	jwtSecret string,
) Service {
	return &service{
		repo:          repo,
		notification:  notificationService,
		jwtSecret:     jwtSecret,
		accessExpire:  time.Hour * 2,
		refreshExpire: time.Hour * 24 * 7,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// =====================================================
// Authentication
// =====================================================

func (s *service) Login(
	ctx context.Context,
	req LoginRequest,
) (*LoginResponse, error) {
	req.Email = normalizeEmail(req.Email)
	req.TenantCode = strings.TrimSpace(req.TenantCode)

	tenant, err := s.repo.GetTenantByCode(ctx, req.TenantCode)
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

	accessToken, err := GenerateAccessToken(
		s.jwtSecret,
		user,
		s.accessExpire,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(
		s.jwtSecret,
		user,
		s.refreshExpire,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpire.Seconds()),
		TokenType:    "Bearer",
		User:         NewUserResponse(user),
	}, nil
}

func (s *service) Logout(
	ctx context.Context,
	req LogoutRequest,
) error {
	return nil
}

func (s *service) RefreshToken(
	ctx context.Context,
	req RefreshTokenRequest,
) (*LoginResponse, error) {
	claims, err := ParseRefreshToken(
		s.jwtSecret,
		req.RefreshToken,
	)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	user, err := s.repo.GetUserByID(
		ctx,
		claims.TenantID,
		claims.UserID,
	)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	if user.Status == model.UserStatusDisabled {
		return nil, ErrUserDisabled
	}

	if user.Status != model.UserStatusActive {
		return nil, ErrUserNotActive
	}

	accessToken, err := GenerateAccessToken(
		s.jwtSecret,
		user,
		s.accessExpire,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(
		s.jwtSecret,
		user,
		s.refreshExpire,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpire.Seconds()),
		TokenType:    "Bearer",
		User:         NewUserResponse(user),
	}, nil
}

// =====================================================
// Access Request
// =====================================================

func (s *service) RequestAccess(
	ctx context.Context,
	req RequestAccessRequest,
) error {
	req.TenantCode = strings.TrimSpace(req.TenantCode)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = normalizeEmail(req.Email)

	tenant, err := s.repo.GetTenantByCode(
		ctx,
		req.TenantCode,
	)
	if err != nil {
		return err
	}

	admin, err := s.repo.GetAdminByTenantID(
		ctx,
		tenant.ID,
	)
	if err != nil {
		return err
	}

	if s.notification == nil {
		return errors.New("notification service is not configured")
	}

	return s.notification.SendAccessRequest(
		ctx,
		notification.AccessRequestEmail{
			AdminEmail: admin.Email,
			Name:       req.Name,
			Email:      req.Email,
			TenantName: tenant.Name,
			TenantCode: tenant.Code,
		},
	)
}

// =====================================================
// Invitation
// =====================================================

func (s *service) InviteUser(
	ctx context.Context,
	tenantID uint,
	operatorID uuid.UUID,
	req InviteUserRequest,
) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = normalizeEmail(req.Email)

	operator, err := s.repo.GetUserByID(
		ctx,
		tenantID,
		operatorID,
	)
	if err != nil {
		return err
	}

	if operator.Status != model.UserStatusActive ||
		operator.UserType != model.UserTypeAdmin {
		return ErrAdminRequired
	}

	if req.RoleID == 0 {
		return ErrRoleInvalid
	}

	_, err = s.repo.GetRoleByID(
		ctx,
		tenantID,
		req.RoleID,
	)
	if err != nil {
		return ErrRoleInvalid
	}

	if existing, err := s.repo.GetUserByEmail(
		ctx,
		tenantID,
		req.Email,
	); err == nil && existing != nil {
		return ErrEmailAlreadyExists
	}

	tenant, err := s.repo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return err
	}

	user := &model.User{
		UUID:     uuid.New(),
		TenantID: tenantID,
		RoleID:   req.RoleID,
		Name:     req.Name,
		Email:    req.Email,
		UserType: model.UserTypeEmployee,
		Status:   model.UserStatusInvited,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	token, hash, err := GenerateInvitationToken()
	if err != nil {
		return err
	}

	invitation := &model.UserInvitation{
		UUID:      uuid.New(),
		TenantID:  tenantID,
		UserID:    user.ID,
		Email:     req.Email,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedBy: operator.ID,
	}

	if err := s.repo.CreateInvitation(
		ctx,
		invitation,
	); err != nil {
		return err
	}

	if s.notification == nil {
		return errors.New("notification service is not configured")
	}

	return s.notification.SendInvitation(
		ctx,
		notification.InvitationEmail{
			Email:      user.Email,
			Name:       user.Name,
			TenantName: tenant.Name,
			Token:      token,
		},
	)
}

func (s *service) AcceptInvitation(
	ctx context.Context,
	req AcceptInvitationRequest,
) error {
	tokenHash, err := HashToken(req.Token)
	if err != nil {
		return ErrInvitationInvalid
	}

	invitation, err := s.repo.GetInvitationByToken(
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
		normalizeEmail(invitation.Email),
	)
	if err != nil {
		return ErrInvitationInvalid
	}

	if user.ID != invitation.UserID {
		return ErrInvitationInvalid
	}

	if user.Status != model.UserStatusInvited {
		return ErrInvitationInvalid
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	principal := &model.Principal{
		UUID:       uuid.New(),
		TenantID:   user.TenantID,
		UserID:     user.ID,
		Provider:   model.PrincipalProviderPassword,
		Identifier: user.Email,
		SecretHash: hash,
		Status:     model.UserStatusActive,
	}

	if err := s.repo.CreatePrincipal(
		ctx,
		principal,
	); err != nil {
		return err
	}

	user.Status = model.UserStatusActive

	if req.Name != "" {
		user.Name = strings.TrimSpace(req.Name)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
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
) (*model.UserInvitation, error) {
	return s.repo.GetInvitationByID(ctx, tenantID, id)
}

func (s *service) CreateUser(
	ctx context.Context,
	req CreateUserRequest,
) (*model.User, error) {
	if req.TenantID == 0 || req.RoleID == 0 {
		return nil, ErrRoleInvalid
	}

	if _, err := s.repo.GetRoleByID(
		ctx,
		req.TenantID,
		req.RoleID,
	); err != nil {
		return nil, ErrRoleInvalid
	}

	req.Email = normalizeEmail(req.Email)

	user := &model.User{
		UUID:     req.UUID,
		TenantID: req.TenantID,
		RoleID:   req.RoleID,
		Name:     strings.TrimSpace(req.Name),
		Email:    req.Email,
		UserType: model.UserTypeEmployee,
		Status:   model.UserStatusInvited,
	}

	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateUser(
	ctx context.Context,
	req UpdateUserRequest,
) error {
	user, err := s.repo.GetUserByID(
		ctx,
		req.TenantID,
		req.UserID,
	)
	if err != nil {
		return err
	}

	if user.UserType == model.UserTypeAdmin {
		return ErrAdminCannotBeModified
	}

	if req.Name != "" {
		user.Name = strings.TrimSpace(req.Name)
	}

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) DisableUser(
	ctx context.Context,
	req DisableUserRequest,
) error {
	user, err := s.repo.GetUserByID(
		ctx,
		req.TenantID,
		req.UserID,
	)
	if err != nil {
		return err
	}

	if user.UserType == model.UserTypeAdmin {
		return ErrAdminCannotBeDisabled
	}

	if user.Status == model.UserStatusInvited {
		return ErrUserNotActive
	}

	user.Status = model.UserStatusDisabled

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) EnableUser(
	ctx context.Context,
	req EnableUserRequest,
) error {
	user, err := s.repo.GetUserByID(
		ctx,
		req.TenantID,
		req.UserID,
	)
	if err != nil {
		return err
	}

	if user.UserType == model.UserTypeAdmin {
		return nil
	}

	if user.Status != model.UserStatusDisabled {
		return ErrUserNotActive
	}

	user.Status = model.UserStatusActive

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) GetUser(
	ctx context.Context,
	tenantID uint,
	userID uuid.UUID,
) (*model.User, error) {
	return s.repo.GetUserByID(ctx, tenantID, userID)
}

func (s *service) GetUserByEmail(
	ctx context.Context,
	tenantID uint,
	email string,
) (*model.User, error) {
	return s.repo.GetUserByEmail(
		ctx,
		tenantID,
		normalizeEmail(email),
	)
}

func (s *service) ListUsers(
	ctx context.Context,
	req ListUsersRequest,
) ([]*model.User, error) {
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

func (s *service) UpdatePassword(
	ctx context.Context,
	req UpdatePasswordRequest,
) error {
	user, err := s.repo.GetUserByID(
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

	principal, err := s.repo.GetPrincipal(
		ctx,
		req.TenantID,
		model.PrincipalProviderPassword,
		user.Email,
	)
	if err != nil {
		return err
	}

	if !VerifyPassword(
		principal.SecretHash,
		req.OldPassword,
	) {
		return errors.New("old password incorrect")
	}

	hash, err := HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	principal.SecretHash = hash

	return s.repo.UpdatePrincipal(ctx, principal)
}

func (s *service) ResetPassword(
	ctx context.Context,
	req ResetPasswordRequest,
) error {
	operator, err := s.repo.GetUserByID(
		ctx,
		req.TenantID,
		req.OperatorID,
	)
	if err != nil {
		return err
	}

	if operator.Status != model.UserStatusActive ||
		operator.UserType != model.UserTypeAdmin {
		return ErrAdminRequired
	}

	user, err := s.repo.GetUserByID(
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

	principal, err := s.repo.GetPrincipal(
		ctx,
		req.TenantID,
		model.PrincipalProviderPassword,
		user.Email,
	)
	if err != nil {
		return err
	}

	hash, err := HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	principal.SecretHash = hash

	return s.repo.UpdatePrincipal(ctx, principal)
}

func (s *service) GetTenant(
	ctx context.Context,
	tenantID uint,
) (*model.Tenant, error) {
	return s.repo.GetTenantByID(ctx, tenantID)
}

func (s *service) GetTenantByCode(
	ctx context.Context,
	code string,
) (*model.Tenant, error) {
	return s.repo.GetTenantByCode(
		ctx,
		strings.TrimSpace(code),
	)
}
