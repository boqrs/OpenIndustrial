
<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Plus, Refresh, Search, View } from "@element-plus/icons-vue";

import { listProductionPlans } from "../api";
import {
  cancelWorkOrder,
  createWorkOrder,
  listWorkOrders,
} from "../workorder-api";
import type {
  CreateWorkOrderRequest,
  WorkOrder,
  WorkOrderStatus,
} from "../types";
import type { ProductionPlan } from "../types";
import { listBOMsByProduct } from "../../bom/api";
import type { BOM } from "../../bom/types";
import { listRoutings } from "../../routing/api";
import type { Routing } from "../../routing/types";

const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const plansLoading = ref(false);
const referencesLoading = ref(false);

const workOrders = ref<WorkOrder[]>([]);
const plans = ref<ProductionPlan[]>([]);
const boms = ref<BOM[]>([]);
const routings = ref<Routing[]>([]);

const keyword = ref("");
const statusFilter = ref<WorkOrderStatus | "">("");

const dialogVisible = ref(false);
const formError = ref("");

const form = reactive({
  productionPlanId: undefined as number | undefined,
  productionLineId: undefined as number | undefined,
  bomId: undefined as number | undefined,
  routingId: undefined as number | undefined,
  code: "",
  plannedQuantity: 1,
  priority: 0,
  dueDate: "",
});

const selectedPlan = computed(() =>
  plans.value.find((item) => item.id === form.productionPlanId),
);

const filteredWorkOrders = computed(() => {
  const query = keyword.value.trim().toLowerCase();

  return workOrders.value.filter((item) => {
    const matchesKeyword =
      !query ||
      item.code.toLowerCase().includes(query) ||
      String(item.id).includes(query) ||
      String(item.production_plan_id).includes(query) ||
      String(item.product_id).includes(query);

    const matchesStatus =
      !statusFilter.value || item.status === statusFilter.value;

    return matchesKeyword && matchesStatus;
  });
});

const summary = computed(() => ({
  total: workOrders.value.length,
  draft: workOrders.value.filter((x) => x.status === "draft").length,
  released: workOrders.value.filter((x) => x.status === "released").length,
  inProgress: workOrders.value.filter((x) => x.status === "in_progress").length,
  completed: workOrders.value.filter((x) => x.status === "completed").length,
}));

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    draft: "草稿",
    released: "已发布",
    in_progress: "执行中",
    completed: "已完成",
    cancelled: "已取消",
  };
  return labels[status] ?? status;
}

function statusType(
  status: string,
): "info" | "success" | "warning" | "danger" | "primary" {
  switch (status) {
    case "released":
      return "success";
    case "in_progress":
      return "primary";
    case "completed":
      return "success";
    case "cancelled":
      return "danger";
    default:
      return "info";
  }
}

function formatDate(value?: string | null): string {
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

function toRFC3339(value: string): string | undefined {
  if (!value) return undefined;
  const date = new Date(value.replace(" ", "T"));
  return Number.isNaN(date.getTime()) ? undefined : date.toISOString();
}

async function loadWorkOrders() {
  loading.value = true;
  try {
    const result = await listWorkOrders({
      current_page: 1,
      page_size: 100,
    });
    workOrders.value = result.detail;
  } catch (error: unknown) {
    workOrders.value = [];
    ElMessage.error(
      error instanceof Error ? error.message : "工单列表加载失败",
    );
  } finally {
    loading.value = false;
  }
}

async function loadPlans() {
  plansLoading.value = true;
  try {
    const [released, inProgress] = await Promise.all([
      listProductionPlans("released"),
      listProductionPlans("in_progress"),
    ]);

    const unique = new Map<number, ProductionPlan>();
    [...released, ...inProgress].forEach((plan) => unique.set(plan.id, plan));
    plans.value = [...unique.values()];
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "生产计划加载失败",
    );
  } finally {
    plansLoading.value = false;
  }
}

watch(
  () => form.productionPlanId,
  async (planId) => {
    boms.value = [];
    routings.value = [];
    form.bomId = undefined;
    form.routingId = undefined;

    const plan = plans.value.find((item) => item.id === planId);
    if (!plan) return;

    referencesLoading.value = true;
    try {
      const [bomList, routingList] = await Promise.all([
        listBOMsByProduct(plan.productId, 1, 100),
        listRoutings({ productId: plan.productId, status: "active" }),
      ]);

      boms.value = bomList.filter((bom) => bom.status === "released");
      routings.value = routingList.filter(
        (routing) => routing.status === "active",
      );

      form.bomId = boms.value[0]?.id;
      form.routingId = routings.value[0]?.id;
    } catch (error: unknown) {
      ElMessage.error(
        error instanceof Error ? error.message : "加载 BOM 或工艺路线失败",
      );
    } finally {
      referencesLoading.value = false;
    }
  },
);

function openCreate() {
  Object.assign(form, {
    productionPlanId: plans.value[0]?.id,
    productionLineId: undefined,
    bomId: undefined,
    routingId: undefined,
    code: "",
    plannedQuantity: 1,
    priority: 0,
    dueDate: "",
  });
  formError.value = "";
  dialogVisible.value = true;
}

async function submitCreate() {
  if (submitting.value) return;

  const plan = selectedPlan.value;
  if (!plan) {
    formError.value = "请选择有效的生产计划";
    return;
  }
  if (!form.productionLineId || !Number.isInteger(form.productionLineId)) {
    formError.value = "请输入有效的生产线 ID";
    return;
  }
  if (!form.bomId) {
    formError.value = "请选择已发布的 BOM";
    return;
  }
  if (!form.routingId) {
    formError.value = "请选择有效的工艺路线";
    return;
  }
  if (!form.code.trim()) {
    formError.value = "请输入工单编号";
    return;
  }
  if (!Number.isInteger(form.plannedQuantity) || form.plannedQuantity <= 0) {
    formError.value = "工单数量必须是大于 0 的整数";
    return;
  }

  const payload: CreateWorkOrderRequest = {
    production_plan_id: plan.id,
    production_line_id: form.productionLineId,
    product_id: plan.productId,
    bom_id: form.bomId,
    routing_id: form.routingId,
    code: form.code.trim(),
    planned_quantity: form.plannedQuantity,
    priority: form.priority,
    due_date: toRFC3339(form.dueDate),
  };

  submitting.value = true;
  formError.value = "";

  try {
    const created = await createWorkOrder(payload);
    ElMessage.success("工单草稿创建成功");
    dialogVisible.value = false;
    await loadWorkOrders();
    if (created?.id) {
      await router.push(`/mes/work-orders/${created.id}`);
    }
  } catch (error: unknown) {
    formError.value =
      error instanceof Error ? error.message : "创建工单失败";
  } finally {
    submitting.value = false;
  }
}

async function cancel(row: WorkOrder) {
  try {
    await ElMessageBox.confirm(
      `确定取消工单「${row.code}」吗？`,
      "取消工单",
      { type: "warning", confirmButtonText: "确认取消", cancelButtonText: "返回" },
    );
    await cancelWorkOrder(row.id);
    ElMessage.success("工单已取消");
    await loadWorkOrders();
  } catch (error: unknown) {
    if (error instanceof Error) {
      ElMessage.error(error.message);
    }
  }
}

import { ElMessageBox } from "element-plus";

onMounted(() => {
  void loadWorkOrders();
  void loadPlans();
});
</script>

<template>
  <div class="work-order-page">
    <header class="page-header">
      <div>
        <div class="eyebrow">MES / WORK ORDER MANAGEMENT</div>
        <h1>生产工单</h1>
        <p class="description">
          将已发布的生产计划拆分为具体制造任务，并指定 BOM、工艺路线及生产线。
        </p>
      </div>
      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadWorkOrders">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">
          新建工单
        </el-button>
      </div>
    </header>

    <section class="summary-grid">
      <article class="summary-card">
        <span>工单总数</span><strong>{{ summary.total }}</strong>
      </article>
      <article class="summary-card">
        <span>草稿</span><strong>{{ summary.draft }}</strong>
      </article>
      <article class="summary-card">
        <span>已发布</span><strong>{{ summary.released }}</strong>
      </article>
      <article class="summary-card">
        <span>执行中</span><strong>{{ summary.inProgress }}</strong>
      </article>
      <article class="summary-card">
        <span>已完成</span><strong>{{ summary.completed }}</strong>
      </article>
    </section>

    <section class="panel">
      <div class="filter-row">
        <el-input
          v-model="keyword"
          clearable
          :prefix-icon="Search"
          placeholder="搜索工单编号、工单 ID、计划 ID 或产品 ID"
          style="width: 330px"
        />
        <el-select
          v-model="statusFilter"
          clearable
          placeholder="全部状态"
          style="width: 150px"
        >
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="released" />
          <el-option label="执行中" value="in_progress" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
        <span class="result-count">共 {{ filteredWorkOrders.length }} 条记录</span>
      </div>

      <el-table
        v-loading="loading"
        :data="filteredWorkOrders"
        row-key="id"
        empty-text="暂无工单，请先发布生产计划"
        style="width: 100%"
      >
        <el-table-column label="工单编号" min-width="180">
          <template #default="{ row }">
            <button
              class="record-link"
              type="button"
              @click="router.push(`/mes/work-orders/${row.id}`)"
            >
              {{ row.code }}
            </button>
            <div class="secondary-text">ID {{ row.id }}</div>
          </template>
        </el-table-column>
        <el-table-column label="生产计划" prop="production_plan_id" width="110" />
        <el-table-column label="产品 ID" prop="product_id" width="100" />
        <el-table-column label="生产线 ID" prop="production_line_id" width="110" />
        <el-table-column label="计划数量" prop="planned_quantity" width="100" align="right" />
        <el-table-column label="优先级" prop="priority" width="90" align="right" />
        <el-table-column label="交期" min-width="160">
          <template #default="{ row }">{{ formatDate(row.due_date) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :icon="View"
              @click="router.push(`/mes/work-orders/${row.id}`)"
            >详情</el-button>
            <el-button
              v-if="row.status === 'draft' || row.status === 'released'"
              link
              type="danger"
              @click="cancel(row)"
            >取消</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>

    <el-dialog v-model="dialogVisible" title="新建生产工单" width="680px" destroy-on-close>
      <el-alert
        title="工单必须关联已发布或执行中的生产计划。BOM 与工艺路线必须属于计划对应的产品。"
        type="info"
        show-icon
        :closable="false"
        style="margin-bottom: 18px"
      />

      <el-form label-position="top">
        <el-form-item label="生产计划" required>
          <el-select
            v-model="form.productionPlanId"
            filterable
            :loading="plansLoading"
            placeholder="选择已发布的生产计划"
            style="width: 100%"
          >
            <el-option
              v-for="plan in plans"
              :key="plan.id"
              :label="`${plan.planNo}（产品 ID ${plan.productId}，计划数量 ${plan.plannedQuantity}）`"
              :value="plan.id"
            />
          </el-select>
        </el-form-item>

        <div class="form-grid">
          <el-form-item label="产品 ID">
            <el-input :model-value="selectedPlan ? String(selectedPlan.productId) : ''" disabled />
          </el-form-item>
          <el-form-item label="生产线 ID" required>
            <el-input-number
              v-model="form.productionLineId"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 100%"
              placeholder="输入生产线 ID"
            />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="BOM（已发布）" required>
            <el-select
              v-model="form.bomId"
              filterable
              :loading="referencesLoading"
              placeholder="选择 BOM"
              style="width: 100%"
            >
              <el-option
                v-for="bom in boms"
                :key="bom.id"
                :label="`${bom.bom_no}（版本 ${bom.version}）`"
                :value="bom.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="工艺路线（有效）" required>
            <el-select
              v-model="form.routingId"
              filterable
              :loading="referencesLoading"
              placeholder="选择工艺路线"
              style="width: 100%"
            >
              <el-option
                v-for="routing in routings"
                :key="routing.id"
                :label="`${routing.code} - ${routing.name}`"
                :value="routing.id"
              />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="工单编号" required>
          <el-input v-model="form.code" maxlength="100" placeholder="例如：WO-2026-001" />
        </el-form-item>

        <div class="form-grid">
          <el-form-item label="计划数量" required>
            <el-input-number
              v-model="form.plannedQuantity"
              :min="1"
              :precision="0"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="优先级">
            <el-input-number
              v-model="form.priority"
              :min="0"
              :precision="0"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
        </div>

        <el-form-item label="交付日期">
          <el-date-picker
            v-model="form.dueDate"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="选择交付日期（选填）"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <el-alert
        v-if="formError"
        :title="formError"
        type="error"
        :closable="false"
        show-icon
      />

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate">
          创建工单草稿
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.work-order-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 34px 0 56px;
  color: #17283b;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
  margin-bottom: 26px;
}
.eyebrow {
  margin-bottom: 10px;
  color: #a16a19;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.5px;
}
h1 {
  margin: 0;
  color: #0e1f33;
  font-size: 28px;
}
.description {
  margin-top: 12px;
  color: #718096;
  font-size: 13px;
  line-height: 1.8;
}
.header-actions { display: flex; gap: 10px; padding-top: 24px; }
.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 22px;
}
.summary-card, .panel {
  border: 1px solid #e6ebf1;
  border-radius: 10px;
  background: #fff;
  box-shadow: 0 3px 12px rgb(14 31 51 / 3%);
}
.summary-card { display: flex; flex-direction: column; padding: 20px; }
.summary-card span { color: #718096; font-size: 12px; }
.summary-card strong {
  margin-top: 12px;
  color: #0e1f33;
  font-size: 27px;
  font-variant-numeric: tabular-nums;
}
.panel { padding: 22px; }
.filter-row { display: flex; align-items: center; flex-wrap: wrap; gap: 12px; margin-bottom: 20px; }
.result-count { margin-left: auto; color: #8a97a8; font-size: 12px; }
.record-link { padding: 0; border: 0; color: #245b91; background: transparent; font: inherit; font-weight: 650; cursor: pointer; }
.record-link:hover { color: #d28b25; }
.secondary-text { margin-top: 5px; color: #9aa6b5; font-size: 11px; }
.form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 0 18px; }
@media (max-width: 1000px) {
  .summary-grid { grid-template-columns: repeat(3, minmax(0, 1fr)); }
  .work-order-page { width: calc(100% - 32px); }
}
@media (max-width: 640px) {
  .page-header { flex-direction: column; }
  .header-actions { padding-top: 0; }
  .summary-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .form-grid { grid-template-columns: 1fr; }
  .panel { padding: 14px; }
}
</style>