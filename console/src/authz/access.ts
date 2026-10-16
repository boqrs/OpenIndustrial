import type { User } from "../api/auth";
import { rolePermissions } from "./role-permissions";
import type { Permission } from "./permissions";
import { getDefaultDashboard, normalizeRole } from "./roles";

export function getPermissionRole(user: User | null | undefined): string {
  if (!user) {
    return "";
  }

  const userType = normalizeRole(user.user_type);

  // 管理员身份由账户类型确定，不能仅靠岗位名称提升权限。
  if (userType === "admin") {
    return "admin";
  }

  const businessRole = normalizeRole(user.role_name);

  // 有业务岗位时，优先按明确的岗位权限判断。
  if (businessRole) {
    return businessRole;
  }

  // 兼容尚未返回 role_name 的旧账户，但只允许有限的内置权限。
  return userType;
}

export function hasPermission(
  user: User | null | undefined,
  permission: Permission,
): boolean {
  if (!user) {
    return false;
  }

  const role = getPermissionRole(user);
  const permissions = rolePermissions[role];

  return permissions?.includes(permission) ?? false;
}

export function getUserDashboard(user: User | null | undefined): string {
  if (!user) {
    return "unavailable";
  }

  return getDefaultDashboard(user.user_type);
}

export function isKnownRole(user: User | null | undefined): boolean {
  if (!user) {
    return false;
  }

  const role = getPermissionRole(user);

  return Object.prototype.hasOwnProperty.call(rolePermissions, role);
}
