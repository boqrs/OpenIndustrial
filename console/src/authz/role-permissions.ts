import { Permissions, type Permission } from "./permissions";
import { Roles } from "./roles";

export const rolePermissions: Record<string, readonly Permission[]> = {
  [Roles.ADMIN]: [
    Permissions.DASHBOARD_VIEW,
    Permissions.USER_LIST,
    Permissions.USER_INVITE,
    Permissions.USER_UPDATE,
    Permissions.USER_DISABLE,
  ],

  [Roles.EMPLOYEE]: [Permissions.DASHBOARD_VIEW],

  [Roles.OPERATOR]: [Permissions.DASHBOARD_VIEW],

  [Roles.VIEWER]: [Permissions.DASHBOARD_VIEW],
};
