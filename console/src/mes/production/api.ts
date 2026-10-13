import request from "../../api/request";
import type {
  ProductionPlan,
  CreateProductionPlanRequest,
  UpdateProductionPlanRequest,
} from "./types";

const BASE_URL = "/production-plans";

/**
 * 获取生产计划列表
 */
export async function listProductionPlans(
  status?: string,
): Promise<ProductionPlan[]> {
  const response = await request.get(BASE_URL, {
    params: status ? { status } : {},
  });

  if (Array.isArray(response)) {
    return response as ProductionPlan[];
  }

  throw new Error("生产计划列表接口返回的数据格式不正确");
}

/**
 * 获取生产计划详情
 */
export async function getProductionPlan(id: number): Promise<ProductionPlan> {
  return request.get(`${BASE_URL}/${id}`);
}

/**
 * 创建生产计划
 */
export async function createProductionPlan(
  data: CreateProductionPlanRequest,
): Promise<ProductionPlan> {
  return request.post(BASE_URL, data);
}

/**
 * 修改生产计划
 */
export async function updateProductionPlan(
  id: number,
  data: UpdateProductionPlanRequest,
): Promise<ProductionPlan> {
  return request.put(`${BASE_URL}/${id}`, data);
}

/**
 * 发布生产计划
 */
export async function releaseProductionPlan(
  id: number,
): Promise<ProductionPlan> {
  return request.post(`${BASE_URL}/${id}/release`);
}

/**
 * 取消生产计划
 */
export async function cancelProductionPlan(
  id: number,
): Promise<ProductionPlan> {
  return request.post(`${BASE_URL}/${id}/cancel`);
}
