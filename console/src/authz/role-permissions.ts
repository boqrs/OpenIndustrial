import { Permissions, type Permission } from "./permissions";
import { Roles } from "./roles";

export const rolePermissions: Record<string, readonly Permission[]> = {
  [Roles.ADMIN]: [
    Permissions.DASHBOARD_VIEW,

    // 用户管理
    Permissions.USER_LIST,
    Permissions.USER_INVITE,
    Permissions.USER_UPDATE,
    Permissions.USER_DISABLE,

    // 产品管理
    Permissions.PRODUCT_LIST,
    Permissions.PRODUCT_CREATE,
    Permissions.PRODUCT_UPDATE,
    Permissions.PRODUCT_STATUS_UPDATE,

    // MES 模块
    Permissions.MES_VIEW,
  ],

  [Roles.EMPLOYEE]: [Permissions.DASHBOARD_VIEW],
  [Roles.OPERATOR]: [Permissions.DASHBOARD_VIEW],
  [Roles.VIEWER]: [Permissions.DASHBOARD_VIEW],
};
