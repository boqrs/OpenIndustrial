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

    // MES 通用访问
    Permissions.MES_VIEW,

    // 生产执行管理
    Permissions.EXECUTION_VIEW,
    Permissions.EXECUTION_CREATE,
    Permissions.EXECUTION_START,

    // 工序作业
    Permissions.OPERATION_START,
    Permissions.OPERATION_COMPLETE,
    Permissions.OPERATION_FAIL,

    // 生产结果
    Permissions.EXECUTION_RESULT_CREATE,
    Permissions.EXECUTION_RESULT_CONFIRM,
    Permissions.EXECUTION_RESULT_CANCEL,
  ],

  // 普通员工默认没有生产执行操作权限。
  [Roles.EMPLOYEE]: [Permissions.DASHBOARD_VIEW],

  // 操作员可以查看执行记录并进行现场生产操作，
  // 不允许创建生产执行，也不允许确认或取消生产结果。
  [Roles.OPERATOR]: [
    Permissions.DASHBOARD_VIEW,
    Permissions.EXECUTION_VIEW,
    Permissions.EXECUTION_START,
    Permissions.OPERATION_START,
    Permissions.OPERATION_COMPLETE,
    Permissions.OPERATION_FAIL,
    Permissions.EXECUTION_RESULT_CREATE,
  ],

  // 只读用户只能查看生产执行及其详情。
  [Roles.VIEWER]: [Permissions.DASHBOARD_VIEW, Permissions.EXECUTION_VIEW],
};
