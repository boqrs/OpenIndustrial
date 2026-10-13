<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Plus,
  Refresh,
  Edit,
  Delete,
  Check,
  VideoPlay,
  VideoPause,
} from "@element-plus/icons-vue";

import {
  activateRouting,
  createOperation,
  deactivateRouting,
  deleteOperation,
  getRouting,
  listOperations,
  updateOperation,
  updateRouting,
} from "../api";
import type {
  CreateOperationRequest,
  Routing,
  RoutingOperation,
  UpdateOperationRequest,
} from "../types";
import { listProductModels } from "../../product/api";
import type { ProductModel } from "../../product/types";
import { getFactoryTopology, listFactories } from "../../factory/api";
import type {
  Factory,
  FactoryTopology,
  TopologyNode,
} from "../../factory/types";

const route = useRoute();
const router = useRouter();
const routingId = Number(route.params.id);

const loading = ref(false);
const saving = ref(false);
const operationSaving = ref(false);
const statusSaving = ref(false);
const routing = ref<Routing | null>(null);
const operations = ref<RoutingOperation[]>([]);
const products = ref<ProductModel[]>([]);
const factories = ref<Factory[]>([]);
const workstations = ref<TopologyNode[]>([]);

const routeDialogVisible = ref(false);
const operationDialogVisible = ref(false);
const editingOperation = ref<RoutingOperation | null>(null);
const routeError = ref("");
const operationError = ref("");

const routeForm = ref({ name: "", description: "" });
const operationForm = ref({
  code: "",
  name: "",
  description: "",
  workstationId: undefined as number | undefined,
  sequence: 1,
  standardDurationSeconds: 0,
  required: true,
  parametersText: "{}",
});

const product = computed(
  () =>
    products.value.find((item) => item.id === routing.value?.productId) ?? null,
);

const sortedOperations = computed(() =>
  [...operations.value].sort((a, b) => a.sequence - b.sequence),
);

const totalDuration = computed(() =>
  sortedOperations.value.reduce(
    (total, item) => total + (item.standardDurationSeconds || 0),
    0,
  ),
);

const editable = computed(
  () =>
    !!routing.value &&
    routing.value.status !== "active" &&
    routing.value.status !== "archived",
);

function statusLabel(status?: string) {
  const labels: Record<string, string> = {
    draft: "草稿",
    active: "已激活",
    inactive: "已停用",
    archived: "已归档",
  };
  return status ? (labels[status] ?? status) : "—";
}

function statusType(
  status?: string,
): "info" | "success" | "warning" | "danger" {
  if (status === "active") return "success";
  if (status === "inactive") return "warning";
  if (status === "archived") return "danger";
  return "info";
}

function formatDuration(seconds: number) {
  if (!seconds) return "0 秒";
  if (seconds < 60) return `${seconds} 秒`;
  const minutes = Math.floor(seconds / 60);
  const remaining = seconds % 60;
  if (minutes < 60) {
    return remaining ? `${minutes} 分 ${remaining} 秒` : `${minutes} 分钟`;
  }
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return remainingMinutes
    ? `${hours} 小时 ${remainingMinutes} 分`
    : `${hours} 小时`;
}

function formatDate(value?: string) {
  if (!value) return "—";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function workstationName(id?: number | null) {
  if (!id) return "未指定";
  return (
    workstations.value.find((item) => item.resource_id === id)?.name ??
    `工作站 #${id}`
  );
}

function flattenNodes(nodes: TopologyNode[], result: TopologyNode[] = []) {
  for (const node of nodes) {
    const normalizedType = node.type.toLowerCase();
    if (normalizedType === "work_station" || normalizedType === "workstation") {
      result.push(node);
    }
    if (node.children?.length) flattenNodes(node.children, result);
  }
  return result;
}

async function loadWorkstations() {
  try {
    factories.value = await listFactories();
    const topologies: FactoryTopology[] = await Promise.all(
      factories.value.map((factory) => getFactoryTopology(factory.id)),
    );

    const allNodes = topologies.flatMap((topology) =>
      flattenNodes(topology.nodes ?? []),
    );
    const unique = new Map<number, TopologyNode>();
    for (const node of allNodes) unique.set(node.resource_id, node);
    workstations.value = [...unique.values()];
  } catch {
    // 工作站列表加载失败不阻塞路线和工序查看。
    workstations.value = [];
  }
}

async function loadAll() {
  if (!Number.isInteger(routingId) || routingId <= 0) {
    ElMessage.error("工艺路线 ID 无效");
    void router.replace("/mes/routings");
    return;
  }

  loading.value = true;
  try {
    const [routingData, operationData, productData] = await Promise.all([
      getRouting(routingId),
      listOperations(routingId),
      listProductModels({ current_page: 1, page_size: 100 }),
    ]);
    routing.value = routingData;
    operations.value = operationData;
    products.value = productData.items ?? [];
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "工艺路线详情加载失败",
    );
  } finally {
    loading.value = false;
  }
  await loadWorkstations();
}

function openRouteEdit() {
  if (!routing.value || !editable.value) return;
  routeForm.value = {
    name: routing.value.name,
    description: routing.value.description ?? "",
  };
  routeError.value = "";
  routeDialogVisible.value = true;
}

async function submitRouteEdit() {
  if (!routing.value || saving.value) return;
  const name = routeForm.value.name.trim();
  if (!name) {
    routeError.value = "路线名称不能为空";
    return;
  }

  saving.value = true;
  routeError.value = "";
  try {
    routing.value = await updateRouting(routingId, {
      name,
      description: routeForm.value.description.trim(),
    });
    routeDialogVisible.value = false;
    ElMessage.success("路线信息已更新");
  } catch (error: unknown) {
    routeError.value = error instanceof Error ? error.message : "更新路线失败";
  } finally {
    saving.value = false;
  }
}

function openCreateOperation() {
  if (!editable.value) return;
  editingOperation.value = null;
  operationForm.value = {
    code: "",
    name: "",
    description: "",
    workstationId: undefined,
    sequence: sortedOperations.value.length
      ? Math.max(...sortedOperations.value.map((item) => item.sequence)) + 10
      : 10,
    standardDurationSeconds: 0,
    required: true,
    parametersText: "{}",
  };
  operationError.value = "";
  operationDialogVisible.value = true;
}

function openEditOperation(item: RoutingOperation) {
  if (!editable.value) return;
  editingOperation.value = item;
  operationForm.value = {
    code: item.code,
    name: item.name,
    description: item.description ?? "",
    workstationId: item.workstationId ?? undefined,
    sequence: item.sequence,
    standardDurationSeconds: item.standardDurationSeconds ?? 0,
    required: item.required,
    parametersText: JSON.stringify(item.parameters ?? {}, null, 2),
  };
  operationError.value = "";
  operationDialogVisible.value = true;
}

function parseParameters(): Record<string, unknown> | null {
  try {
    const parsed: unknown = JSON.parse(
      operationForm.value.parametersText || "{}",
    );
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      operationError.value = '工序参数必须是 JSON 对象，例如 {"torque": 12}';
      return null;
    }
    return parsed as Record<string, unknown>;
  } catch {
    operationError.value = "工序参数不是有效的 JSON，请检查格式";
    return null;
  }
}

async function submitOperation() {
  if (operationSaving.value || !routing.value || !editable.value) return;

  const code = operationForm.value.code.trim();
  const name = operationForm.value.name.trim();
  if (!code || !name) {
    operationError.value = "请填写工序编码和工序名称";
    return;
  }
  if (
    !Number.isInteger(operationForm.value.sequence) ||
    operationForm.value.sequence <= 0
  ) {
    operationError.value = "工序顺序必须是大于 0 的整数";
    return;
  }
  if (
    !Number.isFinite(operationForm.value.standardDurationSeconds) ||
    operationForm.value.standardDurationSeconds < 0
  ) {
    operationError.value = "标准工时不能为负数";
    return;
  }

  const parameters = parseParameters();
  if (parameters === null) return;

  const common = {
    code,
    name,
    description: operationForm.value.description.trim(),
    sequence: operationForm.value.sequence,
    standardDurationSeconds: operationForm.value.standardDurationSeconds,
    required: operationForm.value.required,
    parameters,
    ...(operationForm.value.workstationId
      ? { workstationId: operationForm.value.workstationId }
      : {}),
  };

  operationSaving.value = true;
  operationError.value = "";
  try {
    if (editingOperation.value) {
      const payload: UpdateOperationRequest = common;
      await updateOperation(routingId, editingOperation.value.id, payload);
      ElMessage.success("工序已更新");
    } else {
      const payload: CreateOperationRequest = common;
      await createOperation(routingId, payload);
      ElMessage.success("工序已添加");
    }

    operationDialogVisible.value = false;
    operations.value = await listOperations(routingId);
  } catch (error: unknown) {
    operationError.value =
      error instanceof Error ? error.message : "保存工序失败";
  } finally {
    operationSaving.value = false;
  }
}

async function removeOperation(item: RoutingOperation) {
  if (!editable.value) return;

  try {
    await ElMessageBox.confirm(
      `确定删除工序「${item.sequence}. ${item.name}」吗？`,
      "删除工序",
      {
        type: "warning",
        confirmButtonText: "确认删除",
        cancelButtonText: "取消",
      },
    );
    await deleteOperation(routingId, item.id);
    operations.value = await listOperations(routingId);
    ElMessage.success("工序已删除");
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;
    ElMessage.error(error instanceof Error ? error.message : "删除工序失败");
  }
}

async function changeRoutingStatus(action: "activate" | "deactivate") {
  if (!routing.value || statusSaving.value) return;
  const activating = action === "activate";
  const actionText = activating ? "激活" : "停用";

  try {
    await ElMessageBox.confirm(
      activating
        ? "激活后，此产品的其他激活路线可能会被切换。请确认工序配置已完成。"
        : "如果此路线存在未完成工单，后端将拒绝停用。",
      `${actionText}工艺路线`,
      {
        type: "warning",
        confirmButtonText: `确认${actionText}`,
        cancelButtonText: "取消",
      },
    );

    statusSaving.value = true;
    if (activating) await activateRouting(routingId);
    else await deactivateRouting(routingId);

    ElMessage.success(`工艺路线${actionText}成功`);
    await loadAll();
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;
    ElMessage.error(
      error instanceof Error ? error.message : `${actionText}失败`,
    );
  } finally {
    statusSaving.value = false;
  }
}

onMounted(() => {
  void loadAll();
});
</script>

<template>
  <div class="routing-detail-page" v-loading="loading">
    <div class="back-row">
      <el-button text :icon="ArrowLeft" @click="router.push('/mes/routings')">
        返回工艺路线
      </el-button>
      <el-button :icon="Refresh" :loading="loading" @click="loadAll"
        >刷新</el-button
      >
    </div>

    ```
    <template v-if="routing">
      <header class="detail-header">
        <div class="header-main">
          <div class="eyebrow">MES / PROCESS ROUTING / DETAIL</div>
          <div class="title-row">
            <h1>{{ routing.name }}</h1>
            <el-tag :type="statusType(routing.status)" effect="light">
              {{ statusLabel(routing.status) }}
            </el-tag>
          </div>
          <div class="route-meta">
            <span>{{ routing.code }}</span>
            <span class="meta-separator">·</span>
            <span>V{{ routing.version }}</span>
            <span class="meta-separator">·</span>
            <span>{{
              product
                ? `${product.code} / ${product.name}`
                : `产品 #${routing.productId}`
            }}</span>
          </div>
          <p class="route-description">
            {{ routing.description || "暂无工艺路线说明" }}
          </p>
        </div>

        <div class="header-actions">
          <el-button v-if="editable" :icon="Edit" @click="openRouteEdit">
            编辑路线
          </el-button>
          <el-button
            v-if="routing.status !== 'active' && routing.status !== 'archived'"
            type="primary"
            :icon="VideoPlay"
            :loading="statusSaving"
            @click="changeRoutingStatus('activate')"
          >
            激活路线
          </el-button>
          <el-button
            v-if="routing.status === 'active'"
            type="warning"
            :icon="VideoPause"
            :loading="statusSaving"
            @click="changeRoutingStatus('deactivate')"
          >
            停用路线
          </el-button>
        </div>
      </header>

      <section class="stats-grid">
        <article class="stat-card">
          <span>工序数量</span>
          <strong>{{ operations.length }}</strong>
          <small>按工序顺序执行</small>
        </article>
        <article class="stat-card">
          <span>标准总工时</span>
          <strong>{{ formatDuration(totalDuration) }}</strong>
          <small>各工序标准工时之和</small>
        </article>
        <article class="stat-card">
          <span>必需工序</span>
          <strong>{{
            operations.filter((item) => item.required).length
          }}</strong>
          <small>不可跳过的工序</small>
        </article>
        <article class="stat-card">
          <span>更新时间</span>
          <strong class="date-stat">{{ formatDate(routing.updatedAt) }}</strong>
          <small>路线最后修改时间</small>
        </article>
      </section>

      <section class="content-card">
        <div class="section-header">
          <div>
            <div class="section-eyebrow">PROCESS SEQUENCE</div>
            <h2>工艺流程</h2>
            <p>工序按顺序号从小到大排列，保存后即为该路线的当前工序定义。</p>
          </div>
          <el-button
            v-if="editable"
            type="primary"
            :icon="Plus"
            @click="openCreateOperation"
          >
            添加工序
          </el-button>
        </div>

        <div v-if="sortedOperations.length" class="process-flow">
          <template v-for="(item, index) in sortedOperations" :key="item.id">
            <div class="process-node">
              <div class="node-number">{{ item.sequence }}</div>
              <div class="node-card">
                <div class="node-code">{{ item.code }}</div>
                <strong>{{ item.name }}</strong>
                <span>{{ workstationName(item.workstationId) }}</span>
                <small>{{
                  formatDuration(item.standardDurationSeconds)
                }}</small>
              </div>
            </div>
            <div
              v-if="index < sortedOperations.length - 1"
              class="process-connector"
              aria-hidden="true"
            >
              <span></span>
            </div>
          </template>
        </div>
        <el-empty v-else description="尚未配置工序，添加工序后即可激活路线">
          <el-button
            v-if="editable"
            type="primary"
            @click="openCreateOperation"
          >
            添加第一道工序
          </el-button>
        </el-empty>
      </section>

      <section class="content-card operation-section">
        <div class="section-header">
          <div>
            <div class="section-eyebrow">OPERATION CONFIGURATION</div>
            <h2>工序明细</h2>
            <p>
              工序编码和顺序在同一路线内必须唯一；激活路线不可直接修改工序。
            </p>
          </div>
        </div>

        <el-table :data="sortedOperations" row-key="id" stripe>
          <el-table-column label="顺序" width="85" align="center">
            <template #default="{ row }">
              <span class="sequence-badge">{{ row.sequence }}</span>
            </template>
          </el-table-column>
          <el-table-column label="工序" min-width="190">
            <template #default="{ row }">
              <strong>{{ row.name }}</strong>
              <div class="muted">{{ row.code }}</div>
            </template>
          </el-table-column>
          <el-table-column label="工作站" min-width="150">
            <template #default="{ row }">{{
              workstationName(row.workstationId)
            }}</template>
          </el-table-column>
          <el-table-column label="标准工时" min-width="120">
            <template #default="{ row }">{{
              formatDuration(row.standardDurationSeconds)
            }}</template>
          </el-table-column>
          <el-table-column label="工序属性" width="120">
            <template #default="{ row }">
              <el-tag :type="row.required ? 'success' : 'info'" effect="plain">
                {{ row.required ? "必需" : "可选" }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            label="工序说明"
            min-width="180"
            show-overflow-tooltip
          >
            <template #default="{ row }">{{ row.description || "—" }}</template>
          </el-table-column>
          <el-table-column
            v-if="editable"
            label="操作"
            width="135"
            fixed="right"
          >
            <template #default="{ row }">
              <el-button
                link
                type="primary"
                :icon="Edit"
                @click="openEditOperation(row)"
              >
                编辑
              </el-button>
              <el-button
                link
                type="danger"
                :icon="Delete"
                @click="removeOperation(row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="暂无工序配置" />
          </template>
        </el-table>
      </section>

      <div class="detail-footnote">
        <el-icon><Check /></el-icon>
        <span>
          路线创建后默认为草稿。只有工序配置完成后才能激活；激活和停用结果以服务端校验为准。
        </span>
      </div>
    </template>

    <el-dialog
      v-model="routeDialogVisible"
      title="编辑工艺路线"
      width="540px"
      destroy-on-close
    >
      <el-alert
        v-if="routeError"
        :title="routeError"
        type="error"
        :closable="false"
        show-icon
        class="form-alert"
      />
      <el-form label-position="top">
        <el-form-item label="路线名称" required>
          <el-input v-model="routeForm.name" maxlength="255" />
        </el-form-item>
        <el-form-item label="工艺说明">
          <el-input v-model="routeForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="routeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitRouteEdit"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <el-dialog
      v-model="operationDialogVisible"
      :title="editingOperation ? '编辑工序' : '添加工序'"
      width="680px"
      destroy-on-close
    >
      <el-alert
        v-if="operationError"
        :title="operationError"
        type="error"
        :closable="false"
        show-icon
        class="form-alert"
      />

      <el-form label-position="top">
        <div class="form-grid">
          <el-form-item label="工序编码" required>
            <el-input
              v-model="operationForm.code"
              maxlength="100"
              placeholder="例如：ASSEMBLY"
            />
          </el-form-item>
          <el-form-item label="工序名称" required>
            <el-input
              v-model="operationForm.name"
              maxlength="255"
              placeholder="例如：总装"
            />
          </el-form-item>
          <el-form-item label="工序顺序" required>
            <el-input-number
              v-model="operationForm.sequence"
              :min="1"
              :step="10"
              controls-position="right"
              class="full-width"
            />
          </el-form-item>
          <el-form-item label="标准工时（秒）" required>
            <el-input-number
              v-model="operationForm.standardDurationSeconds"
              :min="0"
              :step="60"
              controls-position="right"
              class="full-width"
            />
          </el-form-item>
          <el-form-item label="执行工作站">
            <el-select
              v-model="operationForm.workstationId"
              clearable
              filterable
              placeholder="不指定工作站"
              class="full-width"
            >
              <el-option
                v-for="station in workstations"
                :key="station.resource_id"
                :label="station.name"
                :value="station.resource_id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="工序要求">
            <el-switch
              v-model="operationForm.required"
              active-text="必需工序"
              inactive-text="可选工序"
            />
          </el-form-item>
        </div>

        <el-form-item label="工序说明">
          <el-input
            v-model="operationForm.description"
            type="textarea"
            :rows="2"
          />
        </el-form-item>
        <el-form-item label="工序参数（JSON 对象）">
          <el-input
            v-model="operationForm.parametersText"
            type="textarea"
            :rows="5"
            placeholder='例如：{"torque": 12, "unit": "N·m"}'
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="operationDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="operationSaving"
          @click="submitOperation"
        >
          保存工序
        </el-button>
      </template>
    </el-dialog>
    ```
  </div>
</template>

<style scoped>
.routing-detail-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 24px 0 48px;
  color: #17283b;
}
.back-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
  padding: 28px;
  border: 1px solid #e5eaf0;
  border-radius: 14px;
  background: #fff;
}
.eyebrow,
.section-eyebrow {
  margin-bottom: 9px;
  color: #8492a5;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1.5px;
}
.title-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
h1 {
  margin: 0;
  font-size: 26px;
  font-weight: 750;
}
.route-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 9px;
  margin-top: 13px;
  color: #687a90;
  font-size: 12px;
}
.meta-separator {
  color: #c0c9d3;
}
.route-description {
  margin: 14px 0 0;
  color: #718197;
  font-size: 13px;
  line-height: 1.7;
}
.header-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  flex-shrink: 0;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin: 18px 0;
}
.stat-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
  padding: 19px;
  border: 1px solid #e5eaf0;
  border-radius: 12px;
  background: #fff;
}
.stat-card > span {
  color: #718197;
  font-size: 12px;
}
.stat-card strong {
  overflow-wrap: anywhere;
  font-size: 25px;
  font-weight: 750;
}
.stat-card small {
  color: #9aa7b6;
  font-size: 11px;
}
.stat-card .date-stat {
  font-size: 13px;
}
.content-card {
  margin-top: 18px;
  overflow: hidden;
  border: 1px solid #e5eaf0;
  border-radius: 12px;
  background: #fff;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 22px 24px;
  border-bottom: 1px solid #edf0f4;
}
.section-header h2 {
  margin: 0;
  font-size: 17px;
  font-weight: 750;
}
.section-header p {
  margin: 8px 0 0;
  color: #8b9aac;
  font-size: 12px;
}
.process-flow {
  display: flex;
  align-items: stretch;
  overflow-x: auto;
  padding: 28px 24px;
}
.process-node {
  display: flex;
  flex: 0 0 190px;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.node-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid #cbdcf1;
  border-radius: 50%;
  color: #245fa8;
  background: #edf4fd;
  font-size: 12px;
  font-weight: 750;
}
.node-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  min-height: 140px;
  gap: 9px;
  padding: 16px 12px;
  border: 1px solid #dfe8f2;
  border-radius: 10px;
  background: #f9fbfe;
  text-align: center;
}
.node-code {
  color: #5b7ea7;
  font-size: 10px;
  letter-spacing: 0.6px;
}
.node-card strong {
  font-size: 13px;
  overflow-wrap: anywhere;
}
.node-card span {
  color: #718197;
  font-size: 11px;
}
.node-card small {
  color: #9aa7b6;
  font-size: 11px;
}
.process-connector {
  display: flex;
  flex: 0 0 48px;
  align-items: center;
  padding-bottom: 42px;
}
.process-connector span {
  position: relative;
  width: 100%;
  height: 2px;
  background: #d8e2ed;
}
.process-connector span::after {
  position: absolute;
  top: -4px;
  right: 0;
  width: 8px;
  height: 8px;
  border-top: 2px solid #9fb2c8;
  border-right: 2px solid #9fb2c8;
  content: "";
  transform: rotate(45deg);
}
.operation-section :deep(.el-table th.el-table__cell) {
  color: #66788d;
  background: #f8fafc;
  font-size: 12px;
}
.sequence-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 30px;
  height: 26px;
  padding: 0 7px;
  border-radius: 6px;
  color: #245fa8;
  background: #edf4fd;
  font-weight: 700;
}
.muted {
  margin-top: 5px;
  color: #91a0b1;
  font-size: 11px;
}
.detail-footnote {
  display: flex;
  align-items: flex-start;
  gap: 9px;
  margin-top: 16px;
  color: #718197;
  font-size: 12px;
  line-height: 1.7;
}
.detail-footnote :deep(.el-icon) {
  margin-top: 2px;
  color: #31845a;
}
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 18px;
}
.full-width {
  width: 100%;
}
.form-alert {
  margin-bottom: 16px;
}
@media (max-width: 900px) {
  .routing-detail-page {
    width: calc(100% - 32px);
  }
  .detail-header {
    flex-direction: column;
    padding: 22px;
  }
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
@media (max-width: 560px) {
  .stats-grid {
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }
  .stat-card {
    padding: 14px;
  }
  .form-grid {
    grid-template-columns: 1fr;
  }
  .section-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
