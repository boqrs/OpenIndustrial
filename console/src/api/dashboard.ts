import request from "./request";

export interface UserStats {
  total: number;
  init: number;
  invited: number;
  active: number;
  disabled: number;
}

export interface DashboardOverview {
  users: UserStats;
}

export async function getDashboardOverview(): Promise<DashboardOverview> {
  return request.get("/dashboard/overview");
}
