export type ExecutionStatus =
  | "pending"
  | "in_progress"
  | "completed"
  | "failed"
  | "cancelled";

export type OperationStatus =
  | "pending"
  | "in_progress"
  | "completed"
  | "skipped"
  | "failed";

export interface ProductionExecution {
  id: number;
  resourceId: number;
  tenantId: number;
  workOrderId: number;
  status: ExecutionStatus | string;
  startedAt?: string | null;
  completedAt?: string | null;
  productId: number;
  routingId: number;
  routingVersion: number;
  createdAt: string;
  updatedAt: string;
}

export interface ExecutionOperation {
  id: number;
  executionId: number;
  sequence: number;
  status: OperationStatus | string;
  result?: Record<string, unknown>;
  startedAt?: string | null;
  completedAt?: string | null;
  code: string;
  name: string;
  description?: string;
  workstationId?: number;
  routingOperationId?: number;
  createdAt: string;
  updatedAt: string;
}

export interface OperationResultRequest {
  result: Record<string, unknown>;
}
