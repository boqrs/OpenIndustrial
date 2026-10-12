export type MaterialType =
  | "raw"
  | "component"
  | "consumable"
  | "packaging"
  | "";

export interface Material {
  ID: number;
  TenantID: number;
  Code: string;
  Name: string;
  MaterialType: MaterialType | string;
  Unit: string;
  Description?: string;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt?: string | null;
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
