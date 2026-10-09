export const Roles = {
  ADMIN: "admin",
  EMPLOYEE: "employee",
  OPERATOR: "operator",
  VIEWER: "viewer",
} as const;

export type BuiltInRole = (typeof Roles)[keyof typeof Roles];

export interface RoleDefinition {
  label: string;
  dashboard: string;
}

export const roleDefinitions: Record<BuiltInRole, RoleDefinition> = {
  [Roles.ADMIN]: {
    label: "工厂管理员",
    dashboard: "admin",
  },
  [Roles.EMPLOYEE]: {
    label: "工厂员工",
    dashboard: "employee",
  },
  [Roles.OPERATOR]: {
    label: "操作员",
    dashboard: "employee",
  },
  [Roles.VIEWER]: {
    label: "只读用户",
    dashboard: "employee",
  },
};

export function normalizeRole(value?: string | null): string {
  return value?.trim().toLowerCase() ?? "";
}

export function getRoleLabel(value?: string | null): string {
  const role = normalizeRole(value);

  if (!role) {
    return "未知角色";
  }

  return role in roleDefinitions
    ? roleDefinitions[role as BuiltInRole].label
    : role;
}

export function getDefaultDashboard(value?: string | null): string {
  const role = normalizeRole(value);

  if (role in roleDefinitions) {
    return roleDefinitions[role as BuiltInRole].dashboard;
  }

  return "unavailable";
}
