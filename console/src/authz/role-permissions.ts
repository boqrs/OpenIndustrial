import { Permissions, type Permission } from "./permissions";
import { Roles, BusinessRoleNames } from "./roles";

const P = Permissions;

/**
 * 工厂管理员：
 * 管理本工厂已定义的业务域。
 *
 * 注意：WMS、IoT 当前只有 view 权限点，
 * 因此目前只能查看，不能视为完整管理授权。
 */
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

/**
 * 工艺工程师：
 * 仅维护工艺路线。
 * MES_VIEW 只用于进入 MES 导航，不授予其他子域操作权限。
 */
const processEngineerPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,

  P.ROUTING_VIEW,
  P.ROUTING_CREATE,
  P.ROUTING_UPDATE,
  P.ROUTING_RELEASE,
];

/**
 * 生产计划员：
 * 读取生产所需的基础数据，负责生产计划和工单。
 * 不修改产品、物料、BOM、工艺路线和工厂拓扑。
 */
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
  P.QUALITY_VIEW,
];

/**
 * 生产操作员：
 * 查看生产所需资料，执行生产任务并提交生产结果。
 * 不负责生产计划、工单配置、质量确认。
 */
const operatorPermissions: Permission[] = [
  P.DASHBOARD_VIEW,
  P.MES_VIEW,

  P.PRODUCT_LIST,
  P.BOM_VIEW,
  P.ROUTING_VIEW,

  P.WORK_ORDER_VIEW,

  P.EXECUTION_VIEW,
  P.EXECUTION_START,

  P.OPERATION_START,
  P.OPERATION_COMPLETE,
  P.OPERATION_FAIL,

  P.EXECUTION_RESULT_CREATE,
];

/**
 * 质量检验员：
 * 查看必要的制造资料，执行并确认质量检验。
 * 不负责生产执行和生产结果修改。
 */
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

/** 未分配业务岗位的普通员工：默认仅可访问 Dashboard。 */
const employeePermissions: Permission[] = [P.DASHBOARD_VIEW];

/**
 * Viewer 暂时采用最小权限。
 * 后续如果需要跨域只读报表，再按明确的数据范围增加读取权限。
 */
const viewerPermissions: Permission[] = [P.DASHBOARD_VIEW];

function normalizeBusinessRole(value: string): string {
  return value.trim().toLowerCase();
}

export const rolePermissions: Record<string, readonly Permission[]> = {
  // 内置账户类型
  [Roles.ADMIN]: adminPermissions,
  [Roles.EMPLOYEE]: employeePermissions,
  [Roles.OPERATOR]: operatorPermissions,
  [Roles.VIEWER]: viewerPermissions,

  // 数据库业务岗位
  [normalizeBusinessRole(BusinessRoleNames.EMPLOYEE)]: employeePermissions,

  [normalizeBusinessRole(BusinessRoleNames.PROCESS_ENGINEER)]:
    processEngineerPermissions,

  [normalizeBusinessRole(BusinessRoleNames.PRODUCTION_PLANNER)]:
    productionPlannerPermissions,

  [normalizeBusinessRole(BusinessRoleNames.OPERATOR)]: operatorPermissions,

  [normalizeBusinessRole(BusinessRoleNames.QUALITY_INSPECTOR)]:
    qualityInspectorPermissions,
};
