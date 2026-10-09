import type { User } from "../api/auth";
import { rolePermissions } from "./role-permissions";
import type { Permission } from "./permissions";
import { getDefaultDashboard, normalizeRole } from "./roles";

export function hasPermission(
  user: User | null | undefined,
  permission: Permission,
): boolean {
  if (!user) {
    return false;
  }

  const role = normalizeRole(user.user_type);
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

  return Object.prototype.hasOwnProperty.call(
    rolePermissions,
    normalizeRole(user.user_type),
  );
}
