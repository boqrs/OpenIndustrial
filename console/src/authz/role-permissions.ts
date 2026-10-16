import { Permissions, type Permission } from "./permissions";
import { Roles, BusinessRoleNames } from "./roles";

const P = Permissions;

const adminPermissions: Permission[] = [
  P.DASHBOARD_VIEW,

  P.USER_LIST,
  P.USER_INVITE,
  P.USER_UPDATE,
  P.USER_DISABLE,

  P.MES_VIEW,

  P.PRODUCT_LIST,
  P.PRODUCT_CREATE,
  P.PRODUCT_UPDATE,
  P.PRODUCT_STATUS_UPDATE,

  P.MATERIAL_VIEW,
  P.MATERIAL_CREATE,
  P.MATERIAL_UPDATE,

  P.BOM_VIEW,
  P.BOM_CREATE,
  P.BOM_UPDATE,
  P.BOM_RELEASE,
  P.BOM_OBSOLETE,

  P.ROUTING_VIEW,
  P.ROUTING_CREATE,
  P.ROUTING_UPDATE,
  P.ROUTING_RELEASE,

  P.FACTORY_VIEW,
  P.FACTORY_CONFIGURE,

  P.PRODUCTION_PLAN_VIEW,
  P.PRODUCTION_PLAN_CREATE,
  P.PRODUCTION_PLAN_UPDATE,

  P.WORK_ORDER_VIEW,
  P.WORK_ORDER_CREATE,
  P.WORK_ORDER_UPDATE,

  P.EXECUTION_VIEW,
  P.EXECUTION_CREATE,
  P.EXECUTION_START,
  P.OPERATION_START,
  P.OPERATION_COMPLETE,
  P.OPERATION_FAIL,

  P.EXECUTION_RESULT_CREATE,
  P.EXECUTION_RESULT_CONFIRM,
  P.EXECUTION_RESULT_CANCEL,

  P.QUALITY_VIEW,
  P.QUALITY_INSPECT,
  P.QUALITY_CONFIRM,

  P.DEVICE_LIST,
  P.DEVICE_VIEW,
  P.WMS_VIEW,
  P.IOT_VIEW,
];

const processEngineerPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,

  P.PRODUCT_LIST,
  P.PRODUCT_CREATE,
  P.PRODUCT_UPDATE,
  P.PRODUCT_STATUS_UPDATE,

  P.MATERIAL_VIEW,
  P.MATERIAL_CREATE,
  P.MATERIAL_UPDATE,

  P.BOM_VIEW,
  P.BOM_CREATE,
  P.BOM_UPDATE,
  P.BOM_RELEASE,
  P.BOM_OBSOLETE,

  P.ROUTING_VIEW,
  P.ROUTING_CREATE,
  P.ROUTING_UPDATE,
  P.ROUTING_RELEASE,

  P.FACTORY_VIEW,
  P.FACTORY_CONFIGURE,

  P.PRODUCTION_PLAN_VIEW,
  P.WORK_ORDER_VIEW,
  P.EXECUTION_VIEW,
];

const productionPlannerPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,

  P.PRODUCT_LIST,
  P.MATERIAL_VIEW,
  P.BOM_VIEW,
  P.ROUTING_VIEW,
  P.FACTORY_VIEW,

  P.PRODUCTION_PLAN_VIEW,
  P.PRODUCTION_PLAN_CREATE,
  P.PRODUCTION_PLAN_UPDATE,

  P.WORK_ORDER_VIEW,
  P.WORK_ORDER_CREATE,
  P.WORK_ORDER_UPDATE,

  P.EXECUTION_VIEW,
];

const operatorPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,

  P.PRODUCT_LIST,
  P.BOM_VIEW,
  P.ROUTING_VIEW,
  P.FACTORY_VIEW,

  P.PRODUCTION_PLAN_VIEW,
  P.WORK_ORDER_VIEW,

  P.EXECUTION_VIEW,
  P.EXECUTION_START,
  P.OPERATION_START,
  P.OPERATION_COMPLETE,
  P.OPERATION_FAIL,
  P.EXECUTION_RESULT_CREATE,
];

const qualityInspectorPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,

  P.PRODUCT_LIST,
  P.BOM_VIEW,
  P.ROUTING_VIEW,
  P.WORK_ORDER_VIEW,
  P.EXECUTION_VIEW,

  P.QUALITY_VIEW,
  P.QUALITY_INSPECT,
  P.QUALITY_CONFIRM,
];

const employeePermissions: Permission[] = [P.DASHBOARD_VIEW];

const viewerPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,
  P.PRODUCT_LIST,
  P.MATERIAL_VIEW,
  P.BOM_VIEW,
  P.ROUTING_VIEW,
  P.FACTORY_VIEW,
  P.PRODUCTION_PLAN_VIEW,
  P.WORK_ORDER_VIEW,
  P.EXECUTION_VIEW,
  P.QUALITY_VIEW,
];

function normalizeBusinessRole(value: string): string {
  return value.trim().toLowerCase();
}

/**
 * 内置账户类型权限。
 * 业务岗位权限单独维护，不能将 employee 自动提升为 MES 配置人员。
 */
export const rolePermissions: Record<string, readonly Permission[]> = {
  [Roles.ADMIN]: adminPermissions,
  [Roles.EMPLOYEE]: employeePermissions,
  [Roles.OPERATOR]: operatorPermissions,
  [Roles.VIEWER]: viewerPermissions,

  [normalizeBusinessRole(BusinessRoleNames.EMPLOYEE)]: employeePermissions,
  [normalizeBusinessRole(BusinessRoleNames.PROCESS_ENGINEER)]:
    processEngineerPermissions,
  [normalizeBusinessRole(BusinessRoleNames.PRODUCTION_PLANNER)]:
    productionPlannerPermissions,
  [normalizeBusinessRole(BusinessRoleNames.OPERATOR)]: operatorPermissions,
  [normalizeBusinessRole(BusinessRoleNames.QUALITY_INSPECTOR)]:
    qualityInspectorPermissions,
};
