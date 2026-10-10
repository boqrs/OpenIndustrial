package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg"
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

		ctx := c.Request.Context()
		ctx = pkg.WithUserID(ctx, userID)
		ctx = pkg.WithTenantID(ctx, tenantID)
		ctx = pkg.WithRoleID(ctx, roleID)
		ctx = pkg.WithClaims(ctx, claims)
		c.Request = c.Request.WithContext(ctx)

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

		tenantID, has := pkg.TenantIDFromContext(c.Request.Context())
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "tenant id is not set",
			})
			return
		}

		roleID, has := pkg.RoleIDFromContext(c.Request.Context())
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "role id is not set",
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
		tenantID, has := pkg.TenantIDFromContext(c.Request.Context())
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "tenant id is not set",
			})
			return
		}

		userID, has := pkg.UserIDFromContext(c.Request.Context())
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user id is not set",
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
