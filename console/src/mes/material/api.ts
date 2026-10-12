import request from "../../api/request";

import type {
  CreateMaterialRequest,
  Material,
  MaterialListQuery,
  UpdateMaterialRequest,
} from "./types";

/** 查询物料列表 */
export async function listMaterials(
  query: MaterialListQuery = {},
): Promise<Material[]> {
  const response = await request.get("/materials", {
    params: {
      offset: query.offset ?? 0,
      limit: query.limit ?? 20,
    },
  });

  if (Array.isArray(response)) {
    return response as Material[];
  }

  if (
    response &&
    typeof response === "object" &&
    "data" in response &&
    Array.isArray(response.data)
  ) {
    return response.data as Material[];
  }

  throw new Error("物料列表接口返回的数据格式不正确");
}

export async function getMaterial(id: number): Promise<Material> {
  return request.get(`/materials/${id}`);
}

export async function createMaterial(
  data: CreateMaterialRequest,
): Promise<unknown> {
  return request.post("/materials", data);
}

export async function updateMaterial(
  id: number,
  data: UpdateMaterialRequest,
): Promise<Material> {
  return request.put(`/materials/${id}`, data);
}

export async function deleteMaterial(id: number): Promise<void> {
  await request.delete(`/materials/${id}`);
}
