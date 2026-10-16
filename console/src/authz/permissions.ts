export const Permissions = {
  DASHBOARD_VIEW: "dashboard:view",

  // 用户管理
  USER_LIST: "users:list",
  USER_INVITE: "users:invite",
  USER_UPDATE: "users:update",
  USER_DISABLE: "users:disable",

  // 产品
  PRODUCT_LIST: "products:list",
  PRODUCT_CREATE: "products:create",
  PRODUCT_UPDATE: "products:update",
  PRODUCT_STATUS_UPDATE: "products:status:update",

  // 物料
  MATERIAL_VIEW: "materials:view",
  MATERIAL_CREATE: "materials:create",
  MATERIAL_UPDATE: "materials:update",

  // BOM
  BOM_VIEW: "boms:view",
  BOM_CREATE: "boms:create",
  BOM_UPDATE: "boms:update",
  BOM_RELEASE: "boms:release",
  BOM_OBSOLETE: "boms:obsolete",

  // 工艺路线
  ROUTING_VIEW: "routings:view",
  ROUTING_CREATE: "routings:create",
  ROUTING_UPDATE: "routings:update",
  ROUTING_RELEASE: "routings:release",

  // 工厂、生产线、工位
  FACTORY_VIEW: "factories:view",
  FACTORY_CONFIGURE: "factories:configure",

  // 生产计划、工单
  PRODUCTION_PLAN_VIEW: "production-plans:view",
  PRODUCTION_PLAN_CREATE: "production-plans:create",
  PRODUCTION_PLAN_UPDATE: "production-plans:update",
  WORK_ORDER_VIEW: "work-orders:view",
  WORK_ORDER_CREATE: "work-orders:create",
  WORK_ORDER_UPDATE: "work-orders:update",

  // MES 模块入口（兼容已有代码）
  MES_VIEW: "mes:view",

  // 生产执行
  EXECUTION_VIEW: "execution:view",
  EXECUTION_CREATE: "execution:create",
  EXECUTION_START: "execution:start",

  // 工序操作
  OPERATION_START: "operation:start",
  OPERATION_COMPLETE: "operation:complete",
  OPERATION_FAIL: "operation:fail",

  // 生产结果
  EXECUTION_RESULT_CREATE: "execution-result:create",
  EXECUTION_RESULT_CONFIRM: "execution-result:confirm",
  EXECUTION_RESULT_CANCEL: "execution-result:cancel",

  // 质量检验
  QUALITY_VIEW: "quality:view",
  QUALITY_INSPECT: "quality:inspect",
  QUALITY_CONFIRM: "quality:confirm",

  // 其他业务模块
  DEVICE_LIST: "devices:list",
  DEVICE_VIEW: "devices:view",
  WMS_VIEW: "wms:view",
  IOT_VIEW: "iot:view",
} as const;

export type Permission = (typeof Permissions)[keyof typeof Permissions];
