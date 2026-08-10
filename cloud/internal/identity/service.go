package identity

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/OpenIndustrial/cloud/internal/resource"
	"github.com/google/uuid"
)

// ServiceUseCase defines the interface for the identity service.
type ServiceUseCase interface {
	RegisterNewTenant(ctx context.Context, params RegisterNewTenantParams) (*RegisterNewTenantResult, error)
	Login(ctx context.Context, params LoginParams) (*LoginResult, error)
	GetCurrentUser(ctx context.Context, tenantID, userID uuid.UUID) (*GetCurrentUserResult, error)
	CreateUser(ctx context.Context, tenantID uuid.UUID, params CreateUserParams) (*CreateUserResult, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID, params ListUsersParams) ([]*User, error)
	UpdateUser(ctx context.Context, tenantID, userID uuid.UUID, params UpdateUserParams) error
	DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]*Role, error)
	AssignRoleToUser(ctx context.Context, tenantID, userID uuid.UUID, params AssignRoleToUserParams) error
}

// Service provides use cases for the identity domain.
type Service struct {
	tenantRepo TenantRepository
	userRepo   UserRepository
	roleRepo   RoleRepository
	groupRepo  resource.GroupRepository
	jwtSecret  string
}

// NewService creates a new identity service.
func NewService(tenantRepo TenantRepository, userRepo UserRepository, roleRepo RoleRepository, groupRepo resource.GroupRepository, jwtSecret string) *Service {
	return &Service{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		groupRepo:  groupRepo,
		jwtSecret:  jwtSecret,
	}
}

// RegisterNewTenantParams defines the parameters for registering a new tenant.
type RegisterNewTenantParams struct {
	TenantName    string
	AdminEmail    string
	AdminPassword string
}

// RegisterNewTenantResult defines the result of a tenant registration.
type RegisterNewTenantResult struct {
	TenantID    uuid.UUID
	AdminUserID uuid.UUID
}

// RegisterNewTenant handles the business logic of creating a new tenant.
func (s *Service) RegisterNewTenant(ctx context.Context, params RegisterNewTenantParams) (*RegisterNewTenantResult, error) {
	tenant := &Tenant{
		Name:   params.TenantName,
		Status: "active",
	}
	if err := s.tenantRepo.CreateTenant(ctx, tenant); err != nil {
		return nil, err
	}

	profileJSON, _ := json.Marshal(map[string]string{"name": params.AdminEmail})
	adminUser := &User{
		TenantID: tenant.ID,
		UserType: string(UserTypeAdmin),
		Profile:  profileJSON,
	}
	if err := s.userRepo.CreateUser(ctx, adminUser); err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(params.AdminPassword)
	if err != nil {
		return nil, err
	}

	adminPrincipal := &Principal{
		UserID:     adminUser.ID,
		TenantID:   tenant.ID,
		Provider:   "password",
		Identifier: params.AdminEmail,
		Credential: hashedPassword,
	}
	if err := s.userRepo.CreatePrincipal(ctx, adminPrincipal); err != nil {
		return nil, err
	}

	adminRole, err := s.roleRepo.GetRoleByName(ctx, uuid.Nil, "Admin")
	if err != nil {
		return nil, err
	}

	if err := s.roleRepo.AddUserToRole(ctx, adminUser.ID, adminRole.ID, tenant.ID); err != nil {
		return nil, err
	}

	adminGroup := &resource.Group{
		ID:          uuid.New(),
		TenantID:    tenant.ID,
		Name:        "Administrators",
		Description: "Default administrators group",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.groupRepo.CreateGroup(ctx, adminGroup); err != nil {
		return nil, err
	}

	if err := s.groupRepo.AddUserToGroup(ctx, tenant.ID, adminUser.ID, adminGroup.ID); err != nil {
		return nil, err
	}

	return &RegisterNewTenantResult{
		TenantID:    tenant.ID,
		AdminUserID: adminUser.ID,
	}, nil
}

// LoginParams defines the parameters for user login.
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id"`
}

// LoginResult defines the result of a successful login.
type LoginResult struct {
	Token string `json:"token"`
}

// Login handles the user authentication logic.
func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error) {
	tenantID, err := uuid.Parse(params.TenantID)
	if err != nil {
		return nil, errors.New("invalid tenant id format")
	}

	principal, err := s.userRepo.GetPrincipal(ctx, tenantID, "password", params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !CheckPasswordHash(params.Password, principal.Credential) {
		return nil, errors.New("invalid credentials")
	}

	token, err := GenerateToken(principal.UserID, principal.TenantID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token}, nil
}

// GetCurrentUserResult defines the result for getting the current user.
type GetCurrentUserResult struct {
	ID       uuid.UUID `json:"id"`
	UserType string    `json:"user_type"`
	Profile  any       `json:"profile"`
}

// GetCurrentUser retrieves the currently authenticated user's information.
func (s *Service) GetCurrentUser(ctx context.Context, tenantID, userID uuid.UUID) (*GetCurrentUserResult, error) {
	user, err := s.userRepo.GetUserByID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	return &GetCurrentUserResult{
		ID:       user.ID,
		UserType: user.UserType,
		Profile:  user.Profile,
	}, nil
}

// CreateUserParams defines the parameters for creating a new user.
type CreateUserParams struct {
	UserType string          `json:"user_type" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	RoleName string          `json:"role_name" binding:"required"`
	Profile  json.RawMessage `json:"profile"`
}

// CreateUserResult defines the result of creating a new user.
type CreateUserResult struct {
	ID uuid.UUID `json:"id"`
}

// CreateUser creates a new user within a tenant.
func (s *Service) CreateUser(ctx context.Context, tenantID uuid.UUID, params CreateUserParams) (*CreateUserResult, error) {
	return nil, errors.New("not implemented")
}

// ListUsersParams defines the parameters for listing users.
type ListUsersParams struct{}

// ListUsers retrieves a list of users for a tenant.
func (s *Service) ListUsers(ctx context.Context, tenantID uuid.UUID, params ListUsersParams) ([]*User, error) {
	return nil, errors.New("not implemented")
}

// UpdateUserParams defines the parameters for updating a user.
type UpdateUserParams struct {
	Profile  json.RawMessage `json:"profile"`
	RoleName string          `json:"role_name"`
}

// UpdateUser updates a user's information.
func (s *Service) UpdateUser(ctx context.Context, tenantID, userID uuid.UUID, params UpdateUserParams) error {
	return errors.New("not implemented")
}

// DeleteUser deletes a user.
func (s *Service) DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error {
	return errors.New("not implemented")
}

// ListRoles retrieves all available roles for a tenant.
func (s *Service) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]*Role, error) {
	return nil, errors.New("not implemented")
}

// AssignRoleToUserParams defines the parameters for assigning a role to a user.
type AssignRoleToUserParams struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// AssignRoleToUser assigns a role to a user.
func (s *Service) AssignRoleToUser(ctx context.Context, tenantID, userID uuid.UUID, params AssignRoleToUserParams) error {
	return errors.New("not implemented")
}