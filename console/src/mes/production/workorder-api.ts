import request from "../../api/request";
import type {
  CreateWorkOrderRequest,
  UpdateWorkOrderRequest,
  WorkOrder,
  WorkOrderListResponse,
} from "./types";

const BASE_URL = "/work-orders";

export async function listWorkOrders(
  query: {
    current_page?: number;
    page_size?: number;
    product_id?: number;
  } = {},
): Promise<WorkOrderListResponse> {
  const response = await request.get(BASE_URL, {
    params: query,
  });

  if (Array.isArray(response)) {
    return { detail: response as WorkOrder[] };
  }

  if (response && typeof response === "object") {
    const body = response as Partial<WorkOrderListResponse>;

    return {
      detail: Array.isArray(body.detail) ? body.detail : [],
      total: body.total,
      next: body.next,
    };
  }

  throw new Error("工单列表接口返回的数据格式不正确");
}

export async function getWorkOrder(id: number): Promise<WorkOrder> {
  return request.get(`${BASE_URL}/${id}`);
}

export async function createWorkOrder(
  data: CreateWorkOrderRequest,
): Promise<WorkOrder> {
  return request.post(BASE_URL, data);
}

export async function updateWorkOrder(
  id: number,
  data: UpdateWorkOrderRequest,
): Promise<WorkOrder> {
  return request.put(`${BASE_URL}/${id}`, data);
}

export async function releaseWorkOrder(id: number): Promise<void> {
  await request.post(`${BASE_URL}/${id}/release`);
}

export async function startWorkOrder(id: number): Promise<void> {
  await request.post(`${BASE_URL}/${id}/start`);
}

export async function cancelWorkOrder(id: number): Promise<void> {
  await request.post(`${BASE_URL}/${id}/cancel`);
}
