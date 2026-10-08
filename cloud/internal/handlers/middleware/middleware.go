package middleware

import (
	"errors"
	//"fmt"
	"context"
	"net/http"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Permission interface {
	CheckRolePermission(
		ctx context.Context,
		tenantID uint,
		roleID uint,
		permissionName string,
	) (bool, error)

	GetUserByID(
		ctx context.Context,
		tenantID uint,
		userID uuid.UUID,
	) (*model.User, error)
}

const (
	contextUserID   = "user_id"
	contextTenantID = "tenant_id"
	contextRoleID   = "role_id"
	contextClaims   = "auth_claims"
)

type service struct {
	jwtSecret *string
	perm      Permission
}

type Service interface {
	Authenticate() gin.HandlerFunc
	RequirePermission(permissionKey string) gin.HandlerFunc
	RequireAdmin() gin.HandlerFunc
}

func NewAuthService(
	jwtSecret *string,
	repo Permission,
) Service {
	return &service{
		jwtSecret: jwtSecret,
		perm:      repo,
	}
}

// Authenticate validates the JWT access token and stores the
// authenticated identity information in the Gin context.
//
// Context values:
//
//	user_id   -> uuid.UUID
//	tenant_id -> uint
//	role_id   -> uint
//	auth_claims -> jwt.MapClaims
func (s *service) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.jwtSecret == nil || *s.jwtSecret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "jwt secret is not configured",
			})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimSpace(
			strings.TrimPrefix(authHeader, bearerPrefix),
		)

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token is required",
			})
			return
		}

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				// Only accept HMAC SHA-256.
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}

				return []byte(*s.jwtSecret), nil
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		// --------------------------------------------------------
		// token_type
		// --------------------------------------------------------

		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token type",
			})
			return
		}

		if tokenType != identity.TokenTypeAccess {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "access token required",
			})
			return
		}

		// --------------------------------------------------------
		// user_id
		// --------------------------------------------------------

		userIDRaw, ok := claims["user_id"].(string)
		if !ok || userIDRaw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id",
			})
			return
		}

		userID, err := uuid.Parse(userIDRaw)
		if err != nil || userID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id",
			})
			return
		}

		// --------------------------------------------------------
		// tenant_id
		// --------------------------------------------------------

		tenantIDRaw, ok := claims["tenant_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid tenant_id",
			})
			return
		}

		if tenantIDRaw <= 0 ||
			tenantIDRaw != float64(uint(tenantIDRaw)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid tenant_id",
			})
			return
		}

		tenantID := uint(tenantIDRaw)

		// --------------------------------------------------------
		// role_id
		// --------------------------------------------------------

		roleIDRaw, ok := claims["role_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid role_id",
			})
			return
		}

		if roleIDRaw <= 0 ||
			roleIDRaw != float64(uint(roleIDRaw)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid role_id",
			})
			return
		}

		roleID := uint(roleIDRaw)

		// --------------------------------------------------------
		// Store authenticated identity in context.
		// --------------------------------------------------------

		c.Set(contextUserID, userID)
		c.Set(contextTenantID, tenantID)
		c.Set(contextRoleID, roleID)
		c.Set(contextClaims, claims)

		c.Next()
	}
}

// RequirePermission checks whether the authenticated user's role
// has the specified permission.
//
// Authentication must be applied before this middleware.
func (s *service) RequirePermission(
	permissionKey string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if permissionKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "permission key is empty",
			})
			return
		}

		tenantID, err := GetTenantIDFromContextV2(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		roleID, err := GetRoleIDFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if s.perm == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "permission repository is not configured",
			})
			return
		}

		hasPermission, err := s.perm.CheckRolePermission(
			c.Request.Context(),
			tenantID,
			roleID,
			permissionKey,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check permission",
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
			})
			return
		}

		c.Next()
	}
}

func (s *service) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, err := GetTenantIDFromContextV2(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if s.perm == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "identity repository is not configured",
			})
			return
		}

		user, err := s.perm.GetUserByID(
			c.Request.Context(),
			tenantID,
			userID,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "administrator privileges required",
			})
			return
		}

		if user.UserType != model.UserTypeAdmin ||
			user.Status != model.UserStatusActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "administrator privileges required",
			})
			return
		}

		c.Next()
	}
}

// GetRoleIDFromContext returns the authenticated role ID.
func GetRoleIDFromContext(c *gin.Context) (uint, error) {
	value, exists := c.Get(contextRoleID)
	if !exists {
		return 0, errors.New("role_id not found in context")
	}

	roleID, ok := value.(uint)
	if !ok {
		return 0, errors.New("invalid role_id in context")
	}

	if roleID == 0 {
		return 0, errors.New("invalid role_id in context")
	}

	return roleID, nil
}

// GetTenantIDFromContext returns the authenticated tenant ID.
func GetTenantIDFromContextV2(c *gin.Context) (uint, error) {
	value, exists := c.Get("tenant_id")
	if !exists {
		return 0, errors.New("tenant_id not found in context")
	}

	tenantID, ok := value.(uint)
	if !ok {
		return 0, errors.New("invalid tenant_id in context")
	}

	if tenantID == 0 {
		return 0, errors.New("invalid tenant_id in context")
	}

	return tenantID, nil
}

// GetUserIDFromContext returns the authenticated user UUID.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	value, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id not found in context")
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user_id in context")
	}

	if userID == uuid.Nil {
		return uuid.Nil, errors.New("invalid user_id in context")
	}

	return userID, nil
}

// GetClaimsFromContext returns the raw JWT claims.
func GetClaimsFromContext(c *gin.Context) (jwt.MapClaims, error) {
	value, exists := c.Get("auth_claims")
	if !exists {
		return nil, errors.New("auth claims not found in context")
	}

	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid auth claims in context")
	}

	return claims, nil
}

func GetTenantIDFromContext(c *gin.Context) (uuid.UUID, error) {
	tenantIDStr, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil, errors.New("tenant_id not found in context")
	}
	tenantID, err := uuid.Parse(tenantIDStr.(string))
	if err != nil {
		return uuid.Nil, errors.New("invalid tenant_id format in context")
	}
	return tenantID, nil
}
