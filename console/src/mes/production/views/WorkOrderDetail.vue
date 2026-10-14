<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  CircleCheck,
  CircleClose,
  EditPen,
  Refresh,
} from "@element-plus/icons-vue";

import {
  cancelWorkOrder,
  getWorkOrder,
  releaseWorkOrder,
  startWorkOrder,
  updateWorkOrder,
} from "../workorder-api";
import type { UpdateWorkOrderRequest, WorkOrder } from "../types";

const props = defineProps<{ id: string | number }>();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const actionLoading = ref(false);
const editing = ref(false);
const loadError = ref("");
const workOrder = ref<WorkOrder | null>(null);

const form = reactive({
  code: "",
  plannedQuantity: 1,
  priority: 0,
  dueDate: "",
});

const workOrderId = computed(() => Number(props.id));
const canEdit = computed(() => workOrder.value?.status === "draft");
const canRelease = computed(() => workOrder.value?.status === "draft");
const canStart = computed(() => workOrder.value?.status === "released");
const canCancel = computed(
  () =>
    workOrder.value?.status === "draft" ||
    workOrder.value?.status === "released",
);

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
    case "completed":
      return "success";
    case "in_progress":
      return "primary";
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

function toDatePickerValue(value?: string | null): string {
  if (!value) return "";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

function toRFC3339(value: string): string | undefined {
  if (!value) return undefined;
  const date = new Date(value.replace(" ", "T"));
  return Number.isNaN(date.getTime()) ? undefined : date.toISOString();
}

function syncForm(item: WorkOrder) {
  form.code = item.code;
  form.plannedQuantity = item.planned_quantity;
  form.priority = item.priority;
  form.dueDate = toDatePickerValue(item.due_date);
}

async function loadWorkOrder() {
  if (!Number.isInteger(workOrderId.value) || workOrderId.value <= 0) {
    loadError.value = "无效的工单 ID";
    return;
  }

  loading.value = true;
  loadError.value = "";
  try {
    workOrder.value = await getWorkOrder(workOrderId.value);
    syncForm(workOrder.value);
  } catch (error: unknown) {
    loadError.value =
      error instanceof Error ? error.message : "工单详情加载失败";
  } finally {
    loading.value = false;
  }
}

function startEdit() {
  if (!workOrder.value || !canEdit.value) return;
  syncForm(workOrder.value);
  editing.value = true;
}

function discardEdit() {
  if (workOrder.value) syncForm(workOrder.value);
  editing.value = false;
}

async function saveChanges() {
  if (!workOrder.value || saving.value) return;

  if (!form.code.trim()) {
    ElMessage.warning("工单编号不能为空");
    return;
  }
  if (!Number.isInteger(form.plannedQuantity) || form.plannedQuantity <= 0) {
    ElMessage.warning("计划数量必须是大于 0 的整数");
    return;
  }

  const payload: UpdateWorkOrderRequest = {
    code: form.code.trim(),
    planned_quantity: form.plannedQuantity,
    priority: form.priority,
    due_date: toRFC3339(form.dueDate),
  };

  saving.value = true;
  try {
    await updateWorkOrder(workOrderId.value, payload);
    ElMessage.success("工单已更新");
    editing.value = false;
    await loadWorkOrder();
  } catch (error: unknown) {
    ElMessage.error(error instanceof Error ? error.message : "更新工单失败");
  } finally {
    saving.value = false;
  }
}

async function release() {
  if (!workOrder.value) return;
  try {
    await ElMessageBox.confirm(
      `确认发布工单「${workOrder.value.code}」吗？发布后即可启动生产。`,
      "发布工单",
      {
        type: "warning",
        confirmButtonText: "确认发布",
        cancelButtonText: "返回",
      },
    );
    actionLoading.value = true;
    await releaseWorkOrder(workOrderId.value);
    ElMessage.success("工单已发布");
    await loadWorkOrder();
  } catch (error: unknown) {
    if (error instanceof Error) ElMessage.error(error.message);
  } finally {
    actionLoading.value = false;
  }
}

async function start() {
  if (!workOrder.value) return;
  try {
    await ElMessageBox.confirm(
      `确认启动工单「${workOrder.value.code}」吗？`,
      "启动生产",
      {
        type: "warning",
        confirmButtonText: "启动生产",
        cancelButtonText: "返回",
      },
    );
    actionLoading.value = true;
    await startWorkOrder(workOrderId.value);
    ElMessage.success("工单已启动");
    await loadWorkOrder();
  } catch (error: unknown) {
    if (error instanceof Error) ElMessage.error(error.message);
  } finally {
    actionLoading.value = false;
  }
}

async function cancel() {
  if (!workOrder.value) return;
  try {
    await ElMessageBox.confirm(
      `确认取消工单「${workOrder.value.code}」吗？`,
      "取消工单",
      {
        type: "warning",
        confirmButtonText: "确认取消",
        cancelButtonText: "返回",
      },
    );
    actionLoading.value = true;
    await cancelWorkOrder(workOrderId.value);
    ElMessage.success("工单已取消");
    await loadWorkOrder();
  } catch (error: unknown) {
    if (error instanceof Error) ElMessage.error(error.message);
  } finally {
    actionLoading.value = false;
  }
}

onMounted(() => {
  void loadWorkOrder();
});
</script>

<template>
  <div class="detail-page" v-loading="loading">
    <button
      class="back-link"
      type="button"
      @click="router.push('/mes/work-orders')"
    >
      <el-icon><ArrowLeft /></el-icon>
      返回工单列表
    </button>

    <header class="page-header">
      <div>
        <div class="eyebrow">MES / WORK ORDER DETAIL</div>
        <div class="title-row">
          <h1>{{ workOrder?.code || "工单详情" }}</h1>
          <el-tag
            v-if="workOrder"
            :type="statusType(workOrder.status)"
            effect="light"
          >
            {{ statusLabel(workOrder.status) }}
          </el-tag>
        </div>
        <p class="description">查看工单关联关系、计划数量、交期与状态流转。</p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadWorkOrder">
          刷新
        </el-button>
        <el-button
          v-if="canEdit && !editing"
          type="primary"
          :icon="EditPen"
          @click="startEdit"
        >
          编辑草稿
        </el-button>
        <el-button
          v-if="canRelease && !editing"
          type="success"
          :icon="CircleCheck"
          :loading="actionLoading"
          @click="release"
        >
          发布工单
        </el-button>
        <el-button
          v-if="canStart && !editing"
          type="primary"
          :loading="actionLoading"
          @click="start"
        >
          启动生产
        </el-button>
        <el-button
          v-if="canCancel && !editing"
          type="danger"
          plain
          :icon="CircleClose"
          :loading="actionLoading"
          @click="cancel"
        >
          取消工单
        </el-button>
      </div>
    </header>

    <el-alert
      v-if="loadError"
      :title="loadError"
      type="error"
      show-icon
      :closable="false"
    />

    <template v-if="workOrder">
      <section class="info-grid">
        <article class="info-card">
          <span>计划数量</span>
          <strong>{{ workOrder.planned_quantity }}</strong>
        </article>
        <article class="info-card">
          <span>生产计划 ID</span>
          <strong>{{ workOrder.production_plan_id }}</strong>
        </article>
        <article class="info-card">
          <span>产品 ID</span>
          <strong>{{ workOrder.product_id }}</strong>
        </article>
        <article class="info-card">
          <span>生产线 ID</span>
          <strong>{{ workOrder.production_line_id }}</strong>
        </article>
      </section>

      <section class="panel">
        <div class="section-heading">
          <div>
            <h2>工单信息</h2>
            <p>工单基本字段与关联的制造资源</p>
          </div>
          <span class="record-id">ID {{ workOrder.id }}</span>
        </div>

        <el-form label-position="top">
          <div class="form-grid">
            <el-form-item label="工单编号" required>
              <el-input v-if="editing" v-model="form.code" maxlength="100" />
              <el-input v-else :model-value="workOrder.code" disabled />
            </el-form-item>
            <el-form-item label="状态">
              <el-input :model-value="statusLabel(workOrder.status)" disabled />
            </el-form-item>
            <el-form-item label="计划数量" required>
              <el-input-number
                v-if="editing"
                v-model="form.plannedQuantity"
                :min="1"
                :precision="0"
                controls-position="right"
                style="width: 100%"
              />
              <el-input
                v-else
                :model-value="String(workOrder.planned_quantity)"
                disabled
              />
            </el-form-item>
            <el-form-item label="优先级">
              <el-input-number
                v-if="editing"
                v-model="form.priority"
                :min="0"
                :precision="0"
                controls-position="right"
                style="width: 100%"
              />
              <el-input
                v-else
                :model-value="String(workOrder.priority)"
                disabled
              />
            </el-form-item>
            <el-form-item label="交付日期">
              <el-date-picker
                v-if="editing"
                v-model="form.dueDate"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                clearable
                placeholder="选择交期"
                style="width: 100%"
              />
              <el-input
                v-else
                :model-value="formatDate(workOrder.due_date)"
                disabled
              />
            </el-form-item>
            <el-form-item label="生产计划 ID">
              <el-input
                :model-value="String(workOrder.production_plan_id)"
                disabled
              />
            </el-form-item>
            <el-form-item label="产品 ID">
              <el-input :model-value="String(workOrder.product_id)" disabled />
            </el-form-item>
            <el-form-item label="BOM ID">
              <el-input :model-value="String(workOrder.bom_id)" disabled />
            </el-form-item>
            <el-form-item label="工艺路线 ID">
              <el-input :model-value="String(workOrder.routing_id)" disabled />
            </el-form-item>
            <el-form-item label="工厂 ID">
              <el-input :model-value="String(workOrder.factory_id)" disabled />
            </el-form-item>
            <el-form-item label="生产线 ID">
              <el-input
                :model-value="String(workOrder.production_line_id)"
                disabled
              />
            </el-form-item>
          </div>
        </el-form>

        <div v-if="editing" class="edit-actions">
          <el-button @click="discardEdit">放弃修改</el-button>
          <el-button type="primary" :loading="saving" @click="saveChanges">
            保存修改
          </el-button>
        </div>
      </section>

      <section class="panel">
        <div class="section-heading">
          <div>
            <h2>生产时间</h2>
            <p>由后端状态流转记录的实际时间</p>
          </div>
        </div>
        <div class="metadata-grid">
          <div>
            <span>开始时间</span
            ><strong>{{ formatDate(workOrder.started_at) }}</strong>
          </div>
          <div>
            <span>完成时间</span
            ><strong>{{ formatDate(workOrder.completed_at) }}</strong>
          </div>
          <div>
            <span>创建时间</span
            ><strong>{{ formatDate(workOrder.created_at) }}</strong>
          </div>
          <div>
            <span>更新时间</span
            ><strong>{{ formatDate(workOrder.updated_at) }}</strong>
          </div>
        </div>
      </section>

      <el-alert
        title="工单启动后，工序开始、完成、失败及生产结果确认由生产执行模块负责。"
        type="info"
        show-icon
        :closable="false"
      />
    </template>
  </div>
</template>

<style scoped>
.detail-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 28px 0 56px;
  color: #17283b;
}
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 22px;
  padding: 0;
  border: 0;
  color: #64748b;
  background: transparent;
  font-size: 12px;
  cursor: pointer;
}
.back-link:hover {
  color: #d28b25;
}
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
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
.title-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.title-row h1 {
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
.header-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-top: 22px;
}
.info-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 22px;
}
.info-card,
.panel {
  border: 1px solid #e6ebf1;
  border-radius: 10px;
  background: #fff;
  box-shadow: 0 3px 12px rgb(14 31 51 / 3%);
}
.info-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
}
.info-card span {
  color: #718096;
  font-size: 12px;
}
.info-card strong {
  color: #0e1f33;
  font-size: 20px;
  overflow-wrap: anywhere;
}
.panel {
  margin-bottom: 22px;
  padding: 24px;
}
.section-heading {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
}
.section-heading h2 {
  margin: 0;
  color: #0e1f33;
  font-size: 17px;
}
.section-heading p {
  margin: 8px 0 0;
  color: #8a97a8;
  font-size: 12px;
}
.record-id {
  color: #9aa6b5;
  font-size: 12px;
}
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}
.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 18px;
  border-top: 1px solid #edf0f4;
}
.metadata-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}
.metadata-grid div {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.metadata-grid span {
  color: #8a97a8;
  font-size: 12px;
}
.metadata-grid strong {
  color: #344256;
  font-size: 13px;
  font-weight: 500;
}
@media (max-width: 900px) {
  .page-header {
    flex-direction: column;
  }
  .header-actions {
    padding-top: 0;
  }
  .info-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .detail-page {
    width: calc(100% - 32px);
  }
}
@media (max-width: 600px) {
  .form-grid,
  .metadata-grid {
    grid-template-columns: 1fr;
  }
  .panel {
    padding: 16px;
  }
}
</style>
