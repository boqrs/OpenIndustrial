package identity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/OpenIndustrial/cloud/internal/param"
	"github.com/OpenIndustrial/cloud/internal/event"
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/google/uuid"
)

// Service defines the interface for the identity service.
// This is the contract that the rest of the application will use.
type Service interface {
	RegisterNewTenant(ctx context.Context, req *param.RegisterTenantRequest) (*param.RegisterTenantResponse, error)
	Login(ctx context.Context, req *param.LoginRequest) (*param.LoginResponse, error)
	GetCurrentUser(ctx context.Context, tenantID, userID uuid.UUID) (*param.GetCurrentUserResponse, error)
	CreateUser(ctx context.Context, tenantID uuid.UUID, req *param.CreateUserRequest) (*param.CreateUserResponse, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID, req *param.ListUsersRequest) ([]*param.UserResponse, error)
	UpdateUser(ctx context.Context, tenantID, userID uuid.UUID, req *param.UpdateUserRequest) error
	DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]*param.RoleResponse, error)
	AssignRoleToUser(ctx context.Context, tenantID, userID uuid.UUID, req *param.AssignRoleToUserRequest) error
	ListUserGroups(ctx context.Context, tenantID, userID uuid.UUID) ([]*param.GroupResponse, error)
}

// Service provides use cases for the identity domain.
type service struct {
	tenantRepo TenantRepository
	userRepo   UserRepository
	roleRepo   RoleRepository
	groupRepo  GroupRepository
	jwtSecret  string
	publisher   event.Publisher
}

// NewService creates a new identity service.
func NewService(tenantRepo TenantRepository, userRepo UserRepository, roleRepo RoleRepository, groupRepo GroupRepository, jwtSecret string, publisher event.Publisher) Service {
	return &service{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		groupRepo:  groupRepo,
		jwtSecret:  jwtSecret,
		publisher: publisher,
	}
}

// RegisterNewTenant handles the business logic of creating a new tenant.
func (s *service) RegisterNewTenant(ctx context.Context, req *param.RegisterTenantRequest) (*param.RegisterTenantResponse, error) {
	// This entire function should be wrapped in a transaction.
	tenant := &model.Tenant{
		Name:   req.TenantName,
		Status: "active",
	}
	if err := s.tenantRepo.CreateTenant(ctx, tenant); err != nil {
		return nil, err
	}

	profileJSON, _ := json.Marshal(map[string]string{"name": req.AdminEmail})
	adminUser := &model.User{
		TenantID: tenant.UUID,
		UserType: string(UserTypeAdmin),
		Profile:  profileJSON,
	}
	if err := s.userRepo.CreateUser(ctx, adminUser); err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(req.AdminPassword)
	if err != nil {
		return nil, err
	}

	adminPrincipal := &model.Principal{
		UserID:     adminUser.UUID,
		TenantID:   tenant.UUID,
		Provider:   "password",
		Identifier: req.AdminEmail,
		Credential: hashedPassword,
	}
	if err := s.userRepo.CreatePrincipal(ctx, adminPrincipal); err != nil {
		return nil, err
	}

	adminRole, err := s.roleRepo.GetRoleByName(ctx, uuid.Nil, "Admin")
	if err != nil {
		return nil, err
	}

	if err := s.roleRepo.AddUserToRole(ctx, adminUser.UUID, adminRole.UUID, tenant.UUID); err != nil {
		return nil, err
	}

	adminGroup := &model.Group{
		TenantID:    tenant.UUID,
		Name:        "Administrators",
		Description: "Default administrators group",
	}
	if err := s.groupRepo.CreateGroup(ctx, adminGroup); err != nil {
		return nil, err
	}

	if err := s.groupRepo.AddUserToGroup(ctx, tenant.UUID, adminUser.UUID, adminGroup.UUID); err != nil {
		return nil, err
	}

	return &param.RegisterTenantResponse{
		TenantID:    tenant.UUID,
		AdminUserID: adminUser.UUID,
	}, nil
}

// Login handles the user authentication logic.
func (s *service) Login(ctx context.Context, req *param.LoginRequest) (*param.LoginResponse, error) {
	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return nil, errors.New("invalid tenant id format")
	}

	principal, err := s.userRepo.GetPrincipal(ctx, tenantID, "password", req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !CheckPasswordHash(req.Password, principal.Credential) {
		return nil, errors.New("invalid credentials")
	}

	token, err := GenerateToken(principal.UserID, principal.TenantID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &param.LoginResponse{Token: token}, nil
}

// GetCurrentUser retrieves the currently authenticated user's information.
func (s *service) GetCurrentUser(ctx context.Context, tenantID, userID uuid.UUID) (*param.GetCurrentUserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	return &param.GetCurrentUserResponse{
		ID:       user.UUID,
		UserType: user.UserType,
		Profile:  user.Profile,
	}, nil
}

// CreateUser creates a new user within a tenant.
func (s *service) CreateUser(ctx context.Context, tenantID uuid.UUID, req *param.CreateUserRequest) (*param.CreateUserResponse, error) {
	// Note: In a real application, this entire function should run within a single database transaction.
	profile := req.Profile
	if profile == nil {
		profile = json.RawMessage("{}")
	}

	newUser := &model.User{
		TenantID: tenantID,
		UserType: req.UserType,
		Profile:  profile,
	}
	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user record: %w", err)
	}



	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newPrincipal := &model.Principal{
		UserID:     newUser.UUID,
		TenantID:   tenantID,
		Provider:   "password",
		Identifier: req.Email,
		Credential: hashedPassword,
	}
	if err := s.userRepo.CreatePrincipal(ctx, newPrincipal); err != nil {
		return nil, fmt.Errorf("failed to create user principal: %w", err)
	}

	roleToAssign, err := s.roleRepo.GetRoleByName(ctx, tenantID, req.RoleName)
	if err != nil {
		return nil, fmt.Errorf("role '%s' not found for this tenant: %w", req.RoleName, err)
	}

	if err := s.roleRepo.AddUserToRole(ctx, newUser.UUID, roleToAssign.UUID, tenantID); err != nil {
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	return &param.CreateUserResponse{
		ID: newUser.UUID,
	}, nil
}

// ListUsers retrieves a list of users for a tenant.
func (s *service) ListUsers(ctx context.Context, tenantID uuid.UUID, req *param.ListUsersRequest) ([]*param.UserResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	repoParams := param.ListUsersRepoReq{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	users, err := s.userRepo.ListUsers(ctx, tenantID, repoParams)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*param.UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, param.ToUserResponse(u))
	}
	if userResponses == nil {
		return []*param.UserResponse{}, nil
	}
	return userResponses, nil
}

// UpdateUser updates a user's information.
func (s *service) UpdateUser(ctx context.Context, tenantID, userID uuid.UUID, req *param.UpdateUserRequest) error {
	return errors.New("not implemented")
}

// DeleteUser deletes a user.
func (s *service) DeleteUser(ctx context.Context, tenantID, userID uuid.UUID) error {
	return errors.New("not implemented")
}

// ListRoles retrieves all available roles for a tenant.
func (s *service) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]*param.RoleResponse, error) {
	return nil, errors.New("not implemented")
}

// AssignRoleToUser assigns a role to a user.
func (s *service) AssignRoleToUser(ctx context.Context, tenantID, userID uuid.UUID, req *param.AssignRoleToUserRequest) error {
	return errors.New("not implemented")
}

// ListUserGroups lists all groups a user is a member of.
func (s *service) ListUserGroups(ctx context.Context, tenantID, userID uuid.UUID) ([]*param.GroupResponse, error) {
	groups, err := s.groupRepo.ListGroupsByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	groupResponses := make([]*param.GroupResponse, 0, len(groups))
	for _, g := range groups {
		groupResponses = append(groupResponses, param.ToGroupResponse(g))
	}
	if groupResponses == nil {
		return []*param.GroupResponse{}, nil
	}
	return groupResponses, nil
}

// --- Event Publishing Helper Methods ---

func (s *service) publishUserCreatedEvent(userID, tenantID uuid.UUID, email string) {
	payload := event.UserCreatedPayload{
		UserID: userID.String(),
		Email:  email,
	}
	userCreatedEvent, err := event.NewEnvelope(
		event.IdentityUserCreated,
		"user",
		userID.String(),
		tenantID.String(),
		payload,
	)
	if err != nil {
		log.Printf("ERROR: failed to create user.created event envelope: %v", err)
		return
	}
	s.publishWithRetry(userCreatedEvent)
}

func (s *service) publishTenantCreatedEvent(tenantID uuid.UUID, name, code string, adminUserID uuid.UUID) {
	payload := event.TenantCreatedPayload{
		TenantID:    tenantID.String(),
		Name:        name,
		Code:        code,
		AdminUserID: adminUserID.String(),
	}
	tenantCreatedEvent, err := event.NewEnvelope(
		event.IdentityTenantCreated,
		"tenant",
		tenantID.String(),
		tenantID.String(), // For tenant-level events, aggregate ID and tenant ID are the same
		payload,
	)
	if err != nil {
		log.Printf("ERROR: failed to create tenant.created event envelope: %v", err)
		return
	}
	s.publishWithRetry(tenantCreatedEvent)
}

func (s *service) publishWithRetry(evt *event.Envelope) {
	for i := 0; i < 3; i++ {
		err := s.publisher.Publish(context.Background(), "openindustrial:events", evt)
		if err == nil {
			log.Printf("INFO: successfully published event %s of type %s", evt.ID, evt.Type)
			return
		}
		log.Printf("WARN: failed to publish event %s, retrying... (%d/3): %v", evt.ID, i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	log.Printf("ERROR: failed to publish event %s after multiple retries", evt.ID)
}