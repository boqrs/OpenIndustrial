export type ProductionPlanStatus =
  | "draft"
  | "released"
  | "in_progress"
  | "completed"
  | "cancelled";

export interface ProductionPlan {
  id: number;
  tenantId: number;
  planNo: string;
  productId: number;
  factoryId: number;
  plannedQuantity: number;
  plannedStartAt: string;
  plannedEndAt: string;
  status: ProductionPlanStatus | string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateProductionPlanRequest {
  planNo: string;
  productID: number;
  factoryID: number;
  plannedQuantity: number;
  plannedStartAt: string;
  plannedEndAt: string;
  description: string;
}

export interface UpdateProductionPlanRequest {
  plannedQuantity?: number;
  plannedStartAt?: string;
  plannedEndAt?: string;
  description?: string;
}
