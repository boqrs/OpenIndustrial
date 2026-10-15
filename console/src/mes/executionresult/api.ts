import request from "../../api/request";

export interface CreateExecutionResultRequest {
  execution_id: number;
  produced_quantity: number;
  qualified_quantity: number;
  rejected_quantity: number;
}

export interface ExecutionResult {
  id: number;
  tenant_id: number;
  execution_id: number;
  work_order_id: number;
  produced_quantity: number;
  qualified_quantity: number;
  rejected_quantity: number;
  status: "draft" | "confirmed" | "cancelled" | string;
  confirmed_at?: string | null;
  created_at: string;
  updated_at: string;
}

const BASE_URL = "/execution-results";

/** 上传生产执行结果，创建草稿。 */
export async function createExecutionResult(
  data: CreateExecutionResultRequest,
): Promise<ExecutionResult> {
  return request.post(BASE_URL, data);
}

/** 根据结果 ID 查询。 */
export async function getExecutionResult(id: number): Promise<ExecutionResult> {
  return request.get(`${BASE_URL}/${id}`);
}

/** 根据执行单 ID 查询结果。 */
export async function getExecutionResultByExecutionId(
  executionId: number,
): Promise<ExecutionResult> {
  return request.get(`/executions/${executionId}/result`);
}

/** 取消尚未确认的结果。 */
export async function cancelExecutionResult(id: number): Promise<void> {
  await request.post(`${BASE_URL}/${id}/cancel`);
}

/** 确认生产执行结果。 */
export async function confirmExecutionResult(id: number): Promise<void> {
  await request.post(`/executions/results/${id}/confirm`);
}
