export type MaterialType = "raw" | "component" | "consumable" | "packaging";

export interface Material {
  id: number;
  tenant_id: number;
  code: string;
  name: string;
  material_type: MaterialType | string;
  unit: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateMaterialRequest {
  code: string;
  name: string;
  material_type: MaterialType;
  unit: string;
  description?: string;
}

export type UpdateMaterialRequest = CreateMaterialRequest;

export interface MaterialListQuery {
  offset?: number;
  limit?: number;
}
