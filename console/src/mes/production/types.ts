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

export type WorkOrderStatus =
  | "draft"
  | "released"
  | "in_progress"
  | "completed"
  | "cancelled";

export interface WorkOrder {
  id: number;
  tenant_id: number;
  factory_id: number;
  production_line_id: number;
  production_plan_id: number;
  product_id: number;
  bom_id: number;
  routing_id: number;
  code: string;
  planned_quantity: number;
  priority: number;
  due_date?: string | null;
  status: WorkOrderStatus | string;
  started_at?: string | null;
  completed_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface WorkOrderListResponse {
  detail: WorkOrder[];
  total?: number;
  next?: boolean;
}

export interface CreateWorkOrderRequest {
  production_plan_id: number;
  production_line_id: number;
  product_id: number;
  bom_id: number;
  routing_id: number;
  code: string;
  planned_quantity: number;
  priority: number;
  due_date?: string;
}

export interface UpdateWorkOrderRequest {
  code: string;
  planned_quantity: number;
  priority: number;
  due_date?: string;
}
