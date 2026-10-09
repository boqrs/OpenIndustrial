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

/**
 * 数据库存储的业务岗位名称。
 *
 * 注意：这些名称来自 Cloud 的 roles 表，
 * 不等于 users.user_type。
 */
export const BusinessRoleNames = {
  EMPLOYEE: "Employee",
  PROCESS_ENGINEER: "Process Engineer",
  PRODUCTION_PLANNER: "Production Planner",
  OPERATOR: "Operator",
  QUALITY_INSPECTOR: "Quality Inspector",
} as const;

const businessRoleLabels: Record<string, string> = {
  Admin: "工厂管理员",
  Employee: "普通员工",
  "Process Engineer": "工艺工程师",
  "Production Planner": "生产计划员",
  Operator: "生产操作员",
  "Quality Inspector": "质量检验员",
};

export function getBusinessRoleLabel(roleName?: string | null): string {
  const name = roleName?.trim();

  if (!name) {
    return "未分配岗位";
  }

  return businessRoleLabels[name] ?? name;
}
