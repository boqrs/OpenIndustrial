import request from "../../api/request";

import type {
  CreateFactoryRequest,
  CreateTopologyNodeRequest,
  Factory,
  FactoryTopology,
  TopologyNode,
  UpdateFactoryRequest,
  UpdateTopologyNodeRequest,
} from "./types";

/** 获取当前租户的工厂列表。
 * 注意：后端需要先注册 GET /factories。
 */

/** 获取当前租户的工厂列表 */
export async function listFactories(): Promise<Factory[]> {
  const response = await request.get("/factories");

  // 格式一：拦截器返回数组
  if (Array.isArray(response)) {
    return response as Factory[];
  }

  if (response && typeof response === "object") {
    const result = response as {
      data?: unknown;
      items?: unknown;
    };

    // 格式二：响应对象直接包含 data 数组
    if (Array.isArray(result.data)) {
      return result.data as Factory[];
    }

    // 格式三：响应对象包含 items 数组
    if (Array.isArray(result.items)) {
      return result.items as Factory[];
    }

    // 格式四：data 内部再包装一层 items
    if (
      result.data &&
      typeof result.data === "object" &&
      "items" in result.data
    ) {
      const data = result.data as { items?: unknown };

      if (Array.isArray(data.items)) {
        return data.items as Factory[];
      }
    }
  }

  console.error("工厂列表接口响应格式：", response);
  throw new Error("工厂列表接口返回的数据格式不正确，请检查接口响应");
}

/** 获取工厂详情 */
export async function getFactory(id: number): Promise<Factory> {
  return request.get(`/factories/${id}`);
}

/** 创建工厂 */
export async function createFactory(
  data: CreateFactoryRequest,
): Promise<Factory> {
  return request.post("/factories", data);
}

/** 更新工厂 */
export async function updateFactory(
  id: number,
  data: UpdateFactoryRequest,
): Promise<Factory> {
  return request.put(`/factories/${id}`, data);
}

/** 删除工厂 */
export async function deleteFactory(id: number): Promise<void> {
  await request.delete(`/factories/${id}`);
}

/** 查询工厂生产拓扑 */
export async function getFactoryTopology(
  factoryId: number,
): Promise<FactoryTopology> {
  return request.get(`/factories/${factoryId}/topology`);
}

/** 创建生产线或工作站 */
export async function createTopologyNode(
  data: CreateTopologyNodeRequest,
): Promise<TopologyNode> {
  return request.post("/factories/topology/nodes", data);
}

/** 更新生产线或工作站 */
export async function updateTopologyNode(
  resourceId: number,
  data: UpdateTopologyNodeRequest,
): Promise<TopologyNode> {
  return request.put(`/factories/topology/nodes/${resourceId}`, data);
}

/** 删除生产线或工作站 */
export async function deleteTopologyNode(resourceId: number): Promise<void> {
  await request.delete(`/factories/topology/nodes/${resourceId}`);
}

/** 移动生产拓扑节点 */
export async function moveTopologyNode(
  resourceId: number,
  parentResourceId: number,
): Promise<void> {
  await request.post("/factories/topology/nodes/move", {
    resource_id: resourceId,
    parent_resource_id: parentResourceId,
  });
}
