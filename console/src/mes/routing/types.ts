export type RoutingStatus = "draft" | "active" | "inactive" | "archived";

export interface Routing {
  id: number;
  resourceId: number;
  productId: number;
  code: string;
  name: string;
  version: number;
  description: string;
  status: RoutingStatus | string;
  createdAt: string;
  updatedAt: string;
}

export interface RoutingOperation {
  id: number;
  routingId: number;
  code: string;
  name: string;
  description: string;
  workstationId?: number | null;
  sequence: number;
  standardDurationSeconds: number;
  required: boolean;
  parameters?: Record<string, unknown> | unknown[] | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRoutingRequest {
  productId: number;
  name: string;
  description?: string;
}

export interface UpdateRoutingRequest {
  name?: string;
  description?: string;
}

export interface CreateOperationRequest {
  code: string;
  name: string;
  description?: string;
  workstationId?: number;
  sequence: number;
  standardDurationSeconds: number;
  required: boolean;
  parameters?: Record<string, unknown>;
}

export interface UpdateOperationRequest {
  code?: string;
  name?: string;
  description?: string;
  workstationId?: number;
  sequence?: number;
  standardDurationSeconds?: number;
  required?: boolean;
  parameters?: Record<string, unknown>;
}
