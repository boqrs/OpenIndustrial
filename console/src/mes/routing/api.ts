import request from "../../api/request";

import type {
  CreateOperationRequest,
  CreateRoutingRequest,
  Routing,
  RoutingOperation,
  RoutingStatus,
  UpdateOperationRequest,
  UpdateRoutingRequest,
} from "./types";

/** 查询工艺路线。后端返回数组，不是分页对象。 */
export async function listRoutings(
  query: {
    productId?: number;
    status?: RoutingStatus | string;
  } = {},
): Promise<Routing[]> {
  const response = await request.get("/routings", {
    params: {
      productID: query.productId,
      status: query.status,
    },
  });

  return Array.isArray(response) ? (response as Routing[]) : [];
}

/** 工艺路线详情 */
export async function getRouting(id: number): Promise<Routing> {
  return request.get(`/routings/${id}`);
}

/** 创建工艺路线 */
export async function createRouting(
  data: CreateRoutingRequest,
): Promise<Routing> {
  return request.post("/routings", data);
}

/** 更新工艺路线基本信息 */
export async function updateRouting(
  id: number,
  data: UpdateRoutingRequest,
): Promise<Routing> {
  return request.put(`/routings/${id}`, data);
}

/** 激活工艺路线 */
export async function activateRouting(id: number): Promise<void> {
  await request.post(`/routings/${id}/activate`);
}

/** 停用工艺路线 */
export async function deactivateRouting(id: number): Promise<void> {
  await request.post(`/routings/${id}/deactivate`);
}

/** 查询工艺路线工序 */
export async function listOperations(
  routingId: number,
): Promise<RoutingOperation[]> {
  const response = await request.get(`/routings/${routingId}/operations`);
  return Array.isArray(response) ? (response as RoutingOperation[]) : [];
}

/** 新增工序 */
export async function createOperation(
  routingId: number,
  data: CreateOperationRequest,
): Promise<RoutingOperation> {
  return request.post(`/routings/${routingId}/operations`, data);
}

/** 更新工序 */
export async function updateOperation(
  routingId: number,
  operationId: number,
  data: UpdateOperationRequest,
): Promise<RoutingOperation> {
  return request.put(`/routings/${routingId}/operations/${operationId}`, data);
}

/** 删除工序 */
export async function deleteOperation(
  routingId: number,
  operationId: number,
): Promise<void> {
  await request.delete(`/routings/${routingId}/operations/${operationId}`);
}
