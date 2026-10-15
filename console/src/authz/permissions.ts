export const Permissions = {
  DASHBOARD_VIEW: "dashboard:view",

  // 用户管理
  USER_LIST: "users:list",
  USER_INVITE: "users:invite",
  USER_UPDATE: "users:update",
  USER_DISABLE: "users:disable",

  // 产品管理
  PRODUCT_LIST: "products:list",
  PRODUCT_CREATE: "products:create",
  PRODUCT_UPDATE: "products:update",
  PRODUCT_STATUS_UPDATE: "products:status:update",

  // 设备管理
  DEVICE_LIST: "devices:list",
  DEVICE_VIEW: "devices:view",

  // MES 通用访问
  MES_VIEW: "mes:view",

  // 生产执行
  EXECUTION_VIEW: "execution:view",
  EXECUTION_CREATE: "execution:create",
  EXECUTION_START: "execution:start",

  // 工序作业
  OPERATION_START: "operation:start",
  OPERATION_COMPLETE: "operation:complete",
  OPERATION_FAIL: "operation:fail",

  // 生产结果
  EXECUTION_RESULT_CREATE: "execution-result:create",
  EXECUTION_RESULT_CONFIRM: "execution-result:confirm",
  EXECUTION_RESULT_CANCEL: "execution-result:cancel",

  WMS_VIEW: "wms:view",
  IOT_VIEW: "iot:view",
} as const;

export type Permission = (typeof Permissions)[keyof typeof Permissions];
