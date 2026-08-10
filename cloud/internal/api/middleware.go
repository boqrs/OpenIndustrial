package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/OpenIndustrial/cloud/internal/identity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationPayloadKey = "authorization_payload"
)

// jwtSecretKey is the secret key for signing JWTs.
// 在真实的应用中，这应该从安全配置中加载！
var jwtSecretKey = []byte("my-super-secret-key")

// AuthMiddleware 创建一个用于 JWT 认证的 gin 中间件。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		accessToken := fields[1]
		claims := &identity.Claims{}

		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			// 确保 token 的签名算法是我们期望的
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 将 claims 存入 context，供后续 handler 使用
		c.Set(AuthorizationPayloadKey, claims)
		c.Next()
	}
}

// GetAuthPayload 从 context 中获取认证用户的 payload。
func GetAuthPayload(c *gin.Context) (*identity.Claims, error) {
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization payload not found in context")
	}

	claims, ok := payload.(*identity.Claims)
	if !ok {
		return nil, errors.New("invalid payload type in context")
	}

	return claims, nil
}

// PermissionMiddleware 创建一个用于检查用户权限的 gin 中间件。
// 它依赖于 AuthMiddleware 已经运行并成功将 user payload 存入 context。
func (h *IdentityHandler) PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 context 获取用户信息
		payload, err := GetAuthPayload(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied: user not authenticated"})
			return
		}

		// 2. 检查权限
		hasPermission, err := h.permRepo.CheckPermissionForUser(c.Request.Context(), payload.UserID, requiredPermission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		// 3. 权限检查通过，继续处理请求
		c.Next()
	}
}