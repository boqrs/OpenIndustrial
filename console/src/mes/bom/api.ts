import request from "../../api/request";
import type { BOM, CreateBOMRequest, UpdateBOMRequest } from "./types";

import type { Material } from "../material/types";

/** 创建 BOM 草稿 */
export async function createBOM(data: CreateBOMRequest): Promise<BOM> {
  return request.post("/boms", data);
}

/** 查询 BOM 详情 */
export async function getBOM(id: number): Promise<BOM> {
  return request.get(`/boms/${id}`);
}

/** 查询指定产品的 BOM 列表 */
export async function listBOMsByProduct(
  productId: number,
  page = 1,
  pageSize = 20,
): Promise<BOM[]> {
  const response = await request.get(`/boms/product/${productId}`, {
    params: {
      page,
      page_size: pageSize,
    },
  });

  return Array.isArray(response) ? (response as BOM[]) : [];
}

/** 更新 BOM 草稿 */
export async function updateBOM(
  id: number,
  data: UpdateBOMRequest,
): Promise<BOM> {
  return request.put(`/boms/${id}`, data);
}

/** 发布 BOM */
export async function releaseBOM(id: number): Promise<BOM> {
  return request.post(`/boms/${id}/release`);
}

/** 作废 BOM */
export async function obsoleteBOM(id: number): Promise<void> {
  await request.post(`/boms/${id}/obsolete`);
}

/** 查询物料列表 */
export async function listMaterials(
  offset = 0,
  limit = 100,
): Promise<Material[]> {
  const response = await request.get("/materials", {
    params: { offset, limit },
  });

  return Array.isArray(response) ? (response as Material[]) : [];
}
