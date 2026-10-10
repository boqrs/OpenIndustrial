import request from "../../api/request";

import type {
  AttributeDefinition,
  CreateProductModelRequest,
  ProductDetail,
  ProductListQuery,
  ProductListResponse,
  UpdateProductModelRequest,
} from "./types";

function normalizeProductList(
  response: Record<string, unknown>,
): ProductListResponse {
  const items = Array.isArray(response.items) ? response.items : [];

  return {
    items: items as ProductListResponse["items"],
    total: Number(response.total ?? 0),
    next: Boolean(response.next ?? false),
  };
}

/** 查询产品型号列表 */
export async function listProductModels(
  query: ProductListQuery = {},
): Promise<ProductListResponse> {
  const response = await request.get("/product-models", { params: query });
  return normalizeProductList(
    (response ?? {}) as unknown as Record<string, unknown>,
  );
}

/** 创建产品型号 */
export async function createProductModel(
  data: CreateProductModelRequest,
): Promise<ProductDetail> {
  return request.post("/product-models", data);
}

/** 查询产品型号详情 */
export async function getProductModel(id: number): Promise<ProductDetail> {
  return request.get(`/product-models/${id}`);
}

/** 更新产品型号基础信息 */
export async function updateProductModel(
  id: number,
  data: UpdateProductModelRequest,
): Promise<unknown> {
  return request.patch(`/product-models/${id}`, data);
}

/** 获取产品型号的属性定义 */
export async function getAttributeDefinitions(
  id: number,
): Promise<AttributeDefinition[]> {
  const response = await request.get(`/product-models/${id}/attributes`);

  return Array.isArray(response) ? (response as AttributeDefinition[]) : [];
}

/** 替换产品型号的属性定义 */
export async function updateAttributeDefinitions(
  id: number,
  attributes: Record<string, import("./types").AttributeDefinitionInput>,
): Promise<void> {
  await request.put(`/product-models/${id}/attributes`, {
    attributes,
  });
}

/** 激活产品型号 */
export async function activateProductModel(id: number): Promise<void> {
  await request.post(`/product-models/${id}/activate`);
}

/** 停用产品型号 */
export async function deactivateProductModel(id: number): Promise<void> {
  await request.post(`/product-models/${id}/deactivate`);
}

/** 归档产品型号 */
export async function archiveProductModel(id: number): Promise<void> {
  await request.post(`/product-models/${id}/archive`);
}
