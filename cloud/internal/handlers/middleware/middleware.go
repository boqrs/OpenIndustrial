package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/identity"
)

type service struct {
	jwtSecret string
	repo      identity.PermissionRepository
}

type Service interface {
	Authenticate() gin.HandlerFunc
	RequirePermission(permissionKey string) gin.HandlerFunc
}

func NewAuthService(
	jwtSecret string,
	repo identity.PermissionRepository,
) Service {
	return &service{
		jwtSecret: jwtSecret,
		repo:      repo,
	}
}

// Authenticate validates the JWT access token and stores the
// authenticated user information in the Gin context.
//
// Context values:
//
//	user_id   -> uuid.UUID
//	tenant_id -> uint
func (s *service) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.jwtSecret == "" {
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
				// Only accept HMAC signing methods.
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}

				return []byte(s.jwtSecret), nil
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

		// jwt/v5 validates registered claims such as exp when
		// Parse is used with MapClaims, but we still explicitly
		// check the token type below.
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

		// user_id is stored as UUID string in JWT.
		userIDRaw, ok := claims["user_id"].(string)
		if !ok || userIDRaw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id",
			})
			return
		}

		userID, err := uuid.Parse(userIDRaw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id",
			})
			return
		}

		// tenant_id is stored as a numeric JWT claim.
		//
		// JSON numbers are decoded into float64 by MapClaims,
		// so we convert it explicitly to uint.
		tenantIDRaw, ok := claims["tenant_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid tenant_id",
			})
			return
		}

		if tenantIDRaw <= 0 || tenantIDRaw != float64(uint(tenantIDRaw)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid tenant_id",
			})
			return
		}

		tenantID := uint(tenantIDRaw)

		c.Set("user_id", userID)
		c.Set("tenant_id", tenantID)
		c.Set("auth_claims", claims)

		c.Next()
	}
}

// RequirePermission checks whether the authenticated user
// has the specified permission.
//
// Authentication must be applied before this middleware.
func (s *service) RequirePermission(permissionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// tenantID, err := GetTenantIDFromContext(c)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		if s.repo == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "permission repository is not configured",
			})
			return
		}

		hasPermission, err := s.repo.CheckPermissionForUser(
			c.Request.Context(),
			//tenantID,
			userID,
			permissionKey,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
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
