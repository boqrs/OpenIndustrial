import request from "../../api/request";
import type {
  ExecutionOperation,
  ProductionExecution,
  OperationResultRequest,
} from "./types";

const BASE_URL = "/executions";

/** 查询生产执行列表，可按工单和状态过滤。 */
export async function listExecutions(params?: {
  workOrderId?: number;
  status?: string;
}): Promise<ProductionExecution[]> {
  const response = await request.get(BASE_URL, {
    params: {
      ...(params?.workOrderId ? { work_order_id: params.workOrderId } : {}),
      ...(params?.status ? { status: params.status } : {}),
    },
  });

  if (Array.isArray(response)) {
    return response as ProductionExecution[];
  }

  throw new Error("生产执行列表接口返回的数据格式不正确");
}

/** 查询执行详情。 */
export async function getExecution(id: number): Promise<ProductionExecution> {
  return request.get(`${BASE_URL}/${id}`);
}

/** 查询执行下的全部工序。 */
export async function listExecutionOperations(
  executionId: number,
): Promise<ExecutionOperation[]> {
  const response = await request.get(`${BASE_URL}/${executionId}/operations`);

  if (Array.isArray(response)) {
    return response as ExecutionOperation[];
  }

  throw new Error("工序列表接口返回的数据格式不正确");
}

/** 为工单创建生产执行。 */
export async function createExecution(
  workOrderId: number,
): Promise<ProductionExecution> {
  return request.post(`${BASE_URL}/execution/${workOrderId}`);
}

/** 启动生产执行。 */
export async function startExecution(id: number): Promise<void> {
  await request.post(`${BASE_URL}/start/${id}`);
}

/** 开始工序。 */
export async function startOperation(
  executionId: number,
  operationId: number,
): Promise<void> {
  await request.post(
    `${BASE_URL}/${executionId}/operations/start/${operationId}`,
  );
}

/** 完成工序并提交工序结果。 */
export async function completeOperation(
  executionId: number,
  operationId: number,
  data: OperationResultRequest,
): Promise<void> {
  await request.post(
    `${BASE_URL}/${executionId}/operations/complete/${operationId}`,
    data,
  );
}

/** 标记工序失败并提交失败信息。 */
export async function failOperation(
  executionId: number,
  operationId: number,
  data: OperationResultRequest,
): Promise<void> {
  await request.post(
    `${BASE_URL}/${executionId}/operations/fail/${operationId}`,
    data,
  );
}
