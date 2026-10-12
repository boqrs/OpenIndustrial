import request from "../../api/request";

import type {
  CreateMaterialRequest,
  Material,
  MaterialListQuery,
  UpdateMaterialRequest,
} from "./types";

/** 查询物料列表。后端当前返回数组，不返回 total。 */
export async function listMaterials(
  query: MaterialListQuery = {},
): Promise<Material[]> {
  const response = await request.get("/materials", {
    params: {
      offset: query.offset ?? 0,
      limit: query.limit ?? 20,
    },
  });

  return Array.isArray(response) ? (response as Material[]) : [];
}

/** 查询物料详情 */
export async function getMaterial(id: number): Promise<Material> {
  return request.get(`/materials/${id}`);
}

/** 创建物料 */
export async function createMaterial(
  data: CreateMaterialRequest,
): Promise<unknown> {
  return request.post("/materials", data);
}

/** 更新物料 */
export async function updateMaterial(
  id: number,
  data: UpdateMaterialRequest,
): Promise<Material> {
  return request.put(`/materials/${id}`, data);
}

/** 删除物料。后端采用 GORM 软删除。 */
export async function deleteMaterial(id: number): Promise<void> {
  await request.delete(`/materials/${id}`);
}
