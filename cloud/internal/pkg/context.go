package pkg

import (
	"context"

	"github.com/google/uuid"
)

// 每种 Context Key 使用独立类型，避免不同身份字段互相冲突。
type userIDContextKey struct{}
type tenantIDContextKey struct{}
type roleIDContextKey struct{}
type claimsContextKey struct{}

// WithUserID 将用户 UUID 写入 Context。
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// UserIDFromContext 从 Context 中读取用户 UUID。
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey{}).(uuid.UUID)
	return userID, ok
}

// WithTenantID 将租户 ID 写入 Context。
func WithTenantID(ctx context.Context, tenantID uint) context.Context {
	return context.WithValue(ctx, tenantIDContextKey{}, tenantID)
}

// TenantIDFromContext 从 Context 中读取租户 ID。
func TenantIDFromContext(ctx context.Context) (uint, bool) {
	tenantID, ok := ctx.Value(tenantIDContextKey{}).(uint)
	return tenantID, ok
}

// WithRoleID 将角色 ID 写入 Context。
func WithRoleID(ctx context.Context, roleID uint) context.Context {
	return context.WithValue(ctx, roleIDContextKey{}, roleID)
}

// RoleIDFromContext 从 Context 中读取角色 ID。
func RoleIDFromContext(ctx context.Context) (uint, bool) {
	roleID, ok := ctx.Value(roleIDContextKey{}).(uint)
	return roleID, ok
}

// WithClaims 将认证 Claims 写入 Context。
func WithClaims(ctx context.Context, claims any) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

// ClaimsFromContext 从 Context 中读取认证 Claims。
func ClaimsFromContext(ctx context.Context) (any, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(any)
	return claims, ok
}
