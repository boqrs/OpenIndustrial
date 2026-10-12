export type BOMStatus = "draft" | "released" | "obsolete";

export interface BOMItem {
  id: number;
  material_id: number;
  quantity: number;
  unit: string;
  sequence: number;
  operation_code: string;
  is_optional: boolean;
  description: string;
}

export interface BOM {
  id: number;
  tenant_id: number;
  product_id: number;
  bom_no: string;
  version: number;
  status: BOMStatus | string;
  description: string;
  items: BOMItem[];
  created_at: string;
  updated_at: string;
}

export interface BOMItemInput {
  material_id: number;
  quantity: number;
  unit: string;
  sequence: number;
  operation_code?: string;
  is_optional: boolean;
  description?: string;
}

export interface CreateBOMRequest {
  product_id: number;
  bom_no: string;
  version: number;
  description?: string;
  items: BOMItemInput[];
}

export interface UpdateBOMRequest {
  description?: string;
  items: BOMItemInput[];
}

export interface Material {
  id: number;
  code: string;
  name: string;
  material_type: string;
  unit: string;
  description?: string;
}

export interface BOMListQuery {
  page?: number;
  page_size?: number;
}
