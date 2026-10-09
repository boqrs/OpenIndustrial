export const Permissions = {
  DASHBOARD_VIEW: "dashboard:view",

  USER_LIST: "users:list",
  USER_INVITE: "users:invite",
  USER_UPDATE: "users:update",
  USER_DISABLE: "users:disable",

  DEVICE_LIST: "devices:list",
  DEVICE_VIEW: "devices:view",

  MES_VIEW: "mes:view",
  WMS_VIEW: "wms:view",
  IOT_VIEW: "iot:view",
} as const;

export type Permission = (typeof Permissions)[keyof typeof Permissions];
