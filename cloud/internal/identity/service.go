package identity

import (
	"context"
	"errors"


	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
)

// Service provides use cases for the identity domain.
type Service struct {
	tenantRepo TenantRepository
	userRepo   UserRepository
	roleRepo   RoleRepository
}

// NewService creates a new identity service.
func NewService(tenantRepo TenantRepository, userRepo UserRepository, roleRepo RoleRepository) *Service {
	return &Service{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
	}
}

// RegisterNewTenantParams defines the parameters for registering a new tenant.
type RegisterNewTenantParams struct {
	TenantName  string
	AdminEmail  string
	AdminPassword string
}

// RegisterNewTenantResult defines the result of a tenant registration.
type RegisterNewTenantResult struct {
	TenantID uuid.UUID
	AdminUserID uuid.UUID
}

// RegisterNewTenant handles the business logic of creating a new tenant,
// an admin user for that tenant, and assigning them the admin role.
func (s *Service) RegisterNewTenant(ctx context.Context, params RegisterNewTenantParams) (*RegisterNewTenantResult, error) {
	// Note: In a real application, this entire function should run within a single database transaction.
	// The transaction management logic is usually handled by a higher-level framework or a UoW (Unit of Work) pattern.

	// 1. Create the tenant
	tenant := &Tenant{
		Name:   params.TenantName,
		Status: "active",
	}
	if err := s.tenantRepo.CreateTenant(ctx, tenant); err != nil {
		return nil, err // wrap error
	}

	// 使用 AdminEmail 作为初始 Profile 中的名字
	profileJSON, _ := json.Marshal(map[string]string{"name": params.AdminEmail})
	adminUser := &User{
		TenantID: tenant.ID,
		UserType: string(UserTypeAdmin),
		Profile:  profileJSON,
	}
	if err := s.userRepo.CreateUser(ctx, adminUser); err != nil {
		return nil, err // wrap error
	}

	// 3. Create the principal for the admin user (for password auth)
	hashedPassword, err := HashPassword(params.AdminPassword)
	if err != nil {
		return nil, err // wrap error
	}

	adminPrincipal := &Principal{
		UserID:     adminUser.ID,
		TenantID:   tenant.ID,
		Provider:   "password",
		Identifier: params.AdminEmail,
		Credential: hashedPassword,
	}
	if err := s.userRepo.CreatePrincipal(ctx, adminPrincipal); err != nil {
		return nil, err // wrap error
	}

	// 4. Find the admin role for the tenant (or create it if it doesn't exist)
	// For now, let's assume a system 'Admin' role exists and we fetch it.
	adminRole, err := s.roleRepo.GetRoleByName(ctx, uuid.Nil, "Admin") // uuid.Nil for system roles
	if err != nil {
		return nil, err // wrap error
	}

	// 5. Assign the admin role to the user
	if err := s.roleRepo.AddUserToRole(ctx, adminUser.ID, adminRole.ID, tenant.ID); err != nil {
		return nil, err // wrap error
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
	TenantID string `json:"tenant_id"` // Can be optional if email is globally unique
}

// LoginResult defines the result of a successful login.
type LoginResult struct {
	Token string `json:"token"`
}

// Login handles the user authentication logic.
func (s *Service) Login(ctx context.Context, params LoginParams) (*LoginResult, error) {
	// 1. Find the principal by email (identifier)
	// Note: We need to parse params.TenantID to UUID
	tenantID, err := uuid.Parse(params.TenantID)
	if err != nil {
		return nil, err // Invalid tenant ID format
	}

	principal, err := s.userRepo.GetPrincipal(ctx, tenantID, "password", params.Email)
	if err != nil {
		return nil, err // User not found or other DB error
	}

	// 2. Check the password
	if !CheckPasswordHash(params.Password, principal.Credential) {
		return nil, errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token, err := GenerateToken(principal.UserID, principal.TenantID)
	if err != nil {
		return nil, err // Failed to generate token
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
		return nil, err // Could be sql.ErrNoRows, which should be handled as a 404 in the handler
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
// Note: This implementation is not yet fully transactional.
func (s *Service) CreateUser(ctx context.Context, tenantID uuid.UUID, params CreateUserParams) (*CreateUserResult, error) {
	// 1. Create the user object
	user := &User{
		TenantID: tenantID,
		UserType: params.UserType,
		Profile:  params.Profile,
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 2. Create the principal (login credentials)
	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		return nil, err
	}
	principal := &Principal{
		UserID:     user.ID,
		TenantID:   tenantID,
		Provider:   "password",
		Identifier: params.Email,
		Credential: hashedPassword,
	}
	if err := s.userRepo.CreatePrincipal(ctx, principal); err != nil {
		return nil, err
	}

	// 3. Find the role by name
	role, err := s.roleRepo.GetRoleByName(ctx, tenantID, params.RoleName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// 4. Assign role to user
	if err := s.roleRepo.AddUserToRole(ctx, user.ID, role.ID, tenantID); err != nil {
		return nil, err
	}

	return &CreateUserResult{ID: user.ID}, nil
}

// ListUsersParams defines parameters for listing users from the API.
type ListUsersParams struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=20"`
}

// ListUsersResult defines the result for listing users.
type ListUsersResult struct {
	Users []User `json:"users"`
}

// ListUsers retrieves a paginated list of users.
func (s *Service) ListUsers(ctx context.Context, tenantID uuid.UUID, params ListUsersParams) (*ListUsersResult, error) {
	// 使用在 identity 包中定义的 repo 参数，不再依赖 postgres 包
	repoParams := ListUsersRepoParams{
		Limit:  params.PageSize,
		Offset: (params.PageID - 1) * params.PageSize,
	}

	// 现在 s.userRepo 接口拥有 ListUsers 方法，可以被正确调用
	users, err := s.userRepo.ListUsers(ctx, tenantID, repoParams)
	if err != nil {
		return nil, err
	}

	// Avoid returning nil slice if there are no users
	if users == nil {
		users = []User{}
	}

	return &ListUsersResult{Users: users}, nil
}

// ListRolesResult defines the result for listing roles.
type ListRolesResult struct {
	Roles []Role `json:"roles"`
}

// ListRoles retrieves all available roles for a tenant.
func (s *Service) ListRoles(ctx context.Context, tenantID uuid.UUID) (*ListRolesResult, error) {
	roles, err := s.roleRepo.ListRoles(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if roles == nil {
		roles = []Role{}
	}

	return &ListRolesResult{Roles: roles}, nil
}

// AssignRoleToUserParams defines the parameters for assigning a role to a user.
type AssignRoleToUserParams struct {
	UserID uuid.UUID `json:"-"` // From URL path
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// AssignRoleToUser assigns a specific role to a user.
func (s *Service) AssignRoleToUser(ctx context.Context, tenantID uuid.UUID, params AssignRoleToUserParams) error {
	// 在生产代码中，这里应该有更多的检查：
	// 1. 检查 roleID 是否存在于该租户下。
	// 2. 检查 userID 是否存在于该租户下。
	// 3. 检查当前操作者是否有权限进行此操作。
	return s.roleRepo.AddUserToRole(ctx, params.UserID, params.RoleID, tenantID)
}

// UpdateUserParams defines the parameters for updating a user.
type UpdateUserParams struct {
	UserID   uuid.UUID       `json:"-"`
	UserType string          `json:"user_type"`
	Profile  json.RawMessage `json:"profile"`
}

// UpdateUser updates a user's details.
func (s *Service) UpdateUser(ctx context.Context, tenantID uuid.UUID, params UpdateUserParams) error {
	// 1. Get the existing user to ensure it exists
	user, err := s.userRepo.GetUserByID(ctx, tenantID, params.UserID)
	if err != nil {
		return err // Returns error if user not found
	}

	// 2. Update fields if they are provided in the request
	if params.UserType != "" {
		user.UserType = params.UserType
	}
	if params.Profile != nil {
		user.Profile = params.Profile
	}

	// 3. Save the updated user
	return s.userRepo.UpdateUser(ctx, user)
}

// DeleteUser deletes a user.
func (s *Service) DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error {
	// In a real app, you'd want to ensure the user exists before deleting.
	return s.userRepo.DeleteUser(ctx, tenantID, userID)
}