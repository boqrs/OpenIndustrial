<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  EditPen,
  Refresh,
  CircleCheck,
  CircleClose,
} from "@element-plus/icons-vue";

import {
  cancelProductionPlan,
  getProductionPlan,
  releaseProductionPlan,
  updateProductionPlan,
} from "../api";
import type { ProductionPlan, UpdateProductionPlanRequest } from "../types";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const actionLoading = ref(false);
const editing = ref(false);
const loadError = ref("");
const plan = ref<ProductionPlan | null>(null);

const form = reactive({
  plannedQuantity: 1,
  plannedStartAt: "",
  plannedEndAt: "",
  description: "",
});

const planId = computed(() => Number(route.params.id));

const canEdit = computed(() => plan.value?.status === "draft");
const canRelease = computed(() => plan.value?.status === "draft");
const canCancel = computed(() =>
  ["draft", "released", "in_progress"].includes(plan.value?.status ?? ""),
);

function statusLabel(status?: string): string {
  const labels: Record<string, string> = {
    draft: "草稿",
    released: "已发布",
    in_progress: "执行中",
    completed: "已完成",
    cancelled: "已取消",
  };
  return labels[status ?? ""] ?? status ?? "未知";
}

function statusType(
  status?: string,
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

function formatDate(value?: string): string {
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

function toLocalInput(value?: string): string {
  if (!value) return "";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";

  const pad = (n: number) => String(n).padStart(2, "0");

  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
    date.getDate(),
  )} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(
    date.getSeconds(),
  )}`;
}

function toRFC3339(value: string): string {
  return new Date(value.replace(" ", "T")).toISOString();
}

function fillForm(data: ProductionPlan) {
  form.plannedQuantity = data.plannedQuantity;
  form.plannedStartAt = toLocalInput(data.plannedStartAt);
  form.plannedEndAt = toLocalInput(data.plannedEndAt);
  form.description = data.description ?? "";
}

async function loadPlan() {
  if (!Number.isInteger(planId.value) || planId.value <= 0) {
    loadError.value = "生产计划 ID 无效";
    return;
  }

  loading.value = true;
  loadError.value = "";

  try {
    const detail = await getProductionPlan(planId.value);

    plan.value = detail;
    fillForm(detail);
    editing.value = false;
  } catch (error: unknown) {
    plan.value = null;
    loadError.value =
      error instanceof Error ? error.message : "生产计划详情加载失败";
  } finally {
    loading.value = false;
  }
}

function startEdit() {
  if (!plan.value || !canEdit.value) return;
  fillForm(plan.value);
  editing.value = true;
}

function discardEdit() {
  if (plan.value) fillForm(plan.value);
  editing.value = false;
}

async function saveChanges() {
  if (!plan.value || saving.value) return;

  if (!Number.isInteger(form.plannedQuantity) || form.plannedQuantity <= 0) {
    ElMessage.warning("计划数量必须是大于 0 的整数");
    return;
  }

  const start = new Date(form.plannedStartAt.replace(" ", "T"));
  const end = new Date(form.plannedEndAt.replace(" ", "T"));

  if (
    Number.isNaN(start.getTime()) ||
    Number.isNaN(end.getTime()) ||
    end.getTime() <= start.getTime()
  ) {
    ElMessage.warning("计划结束时间必须晚于开始时间");
    return;
  }

  const payload: UpdateProductionPlanRequest = {
    plannedQuantity: form.plannedQuantity,
    plannedStartAt: toRFC3339(form.plannedStartAt),
    plannedEndAt: toRFC3339(form.plannedEndAt),
    description: form.description.trim(),
  };

  saving.value = true;

  try {
    await updateProductionPlan(plan.value.id, payload);
    ElMessage.success("生产计划已更新");
    await loadPlan();
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "更新生产计划失败",
    );
  } finally {
    saving.value = false;
  }
}

async function releasePlan() {
  if (!plan.value || actionLoading.value) return;

  try {
    await ElMessageBox.confirm(
      `确认发布生产计划「${plan.value.planNo}」吗？发布后将不能编辑计划内容。`,
      "发布生产计划",
      {
        type: "warning",
        confirmButtonText: "确认发布",
        cancelButtonText: "暂不发布",
      },
    );

    actionLoading.value = true;
    await releaseProductionPlan(plan.value.id);
    ElMessage.success("生产计划已发布");
    await loadPlan();
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;

    ElMessage.error(
      error instanceof Error ? error.message : "发布生产计划失败",
    );
  } finally {
    actionLoading.value = false;
  }
}

async function cancelPlan() {
  if (!plan.value || actionLoading.value) return;

  try {
    await ElMessageBox.confirm(
      `确认取消生产计划「${plan.value.planNo}」吗？请确认该计划不再需要执行。`,
      "取消生产计划",
      {
        type: "warning",
        confirmButtonText: "确认取消",
        cancelButtonText: "返回",
      },
    );

    actionLoading.value = true;
    await cancelProductionPlan(plan.value.id);
    ElMessage.success("生产计划已取消");
    await loadPlan();
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;

    ElMessage.error(
      error instanceof Error ? error.message : "取消生产计划失败",
    );
  } finally {
    actionLoading.value = false;
  }
}

onMounted(() => {
  void loadPlan();
});
</script>

<template>
  <div class="detail-page" v-loading="loading">
    <header class="page-header">
      <div>
        <button
          class="back-link"
          type="button"
          @click="router.push('/mes/production-plans')"
        >
          <el-icon><ArrowLeft /></el-icon>
          返回生产计划列表
        </button>

        <div class="eyebrow">MES / PRODUCTION PLAN DETAIL</div>

        <div class="title-row">
          <h1>{{ plan?.planNo ?? "生产计划详情" }}</h1>
          <el-tag v-if="plan" :type="statusType(plan.status)" effect="light">
            {{ statusLabel(plan.status) }}
          </el-tag>
        </div>

        <p class="description">查看生产计划的基础信息、生产数量及计划周期。</p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadPlan">
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
          @click="releasePlan"
        >
          发布计划
        </el-button>

        <el-button
          v-if="canCancel && !editing"
          type="danger"
          plain
          :icon="CircleClose"
          :loading="actionLoading"
          @click="cancelPlan"
        >
          取消计划
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

    <template v-if="plan">
      <section class="info-grid">
        <article class="info-card">
          <span>计划编号</span>
          <strong>{{ plan.planNo }}</strong>
        </article>
        <article class="info-card">
          <span>计划数量</span>
          <strong>{{ plan.plannedQuantity }}</strong>
        </article>
        <article class="info-card">
          <span>产品 ID</span>
          <strong>{{ plan.productId }}</strong>
        </article>
        <article class="info-card">
          <span>工厂 ID</span>
          <strong>{{ plan.factoryId }}</strong>
        </article>
      </section>

      <section class="panel">
        <div class="section-heading">
          <div>
            <h2>计划信息</h2>
            <p>生产计划基础信息与计划周期</p>
          </div>
          <span class="section-id">ID {{ plan.id }}</span>
        </div>

        <el-form label-position="top">
          <div class="form-grid">
            <el-form-item label="计划编号">
              <el-input :model-value="plan.planNo" disabled />
            </el-form-item>

            <el-form-item label="状态">
              <el-input :model-value="statusLabel(plan.status)" disabled />
            </el-form-item>

            <el-form-item label="产品 ID">
              <el-input :model-value="String(plan.productId)" disabled />
            </el-form-item>

            <el-form-item label="工厂 ID">
              <el-input :model-value="String(plan.factoryId)" disabled />
            </el-form-item>

            <el-form-item label="计划生产数量" required>
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
                :model-value="String(plan.plannedQuantity)"
                disabled
              />
            </el-form-item>

            <el-form-item label="计划开始时间" required>
              <el-date-picker
                v-if="editing"
                v-model="form.plannedStartAt"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
              <el-input
                v-else
                :model-value="formatDate(plan.plannedStartAt)"
                disabled
              />
            </el-form-item>

            <el-form-item label="计划结束时间" required>
              <el-date-picker
                v-if="editing"
                v-model="form.plannedEndAt"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
              <el-input
                v-else
                :model-value="formatDate(plan.plannedEndAt)"
                disabled
              />
            </el-form-item>
          </div>

          <el-form-item label="计划描述">
            <el-input
              v-if="editing"
              v-model="form.description"
              type="textarea"
              :rows="4"
              maxlength="500"
              show-word-limit
            />
            <el-input
              v-else
              :model-value="plan.description || '暂无描述'"
              type="textarea"
              :rows="3"
              disabled
            />
          </el-form-item>
        </el-form>

        <div v-if="editing" class="edit-actions">
          <el-button @click="discardEdit">放弃修改</el-button>
          <el-button type="primary" :loading="saving" @click="saveChanges">
            保存修改
          </el-button>
        </div>
      </section>

      <section class="panel metadata-panel">
        <div class="section-heading">
          <div>
            <h2>记录信息</h2>
            <p>系统记录的创建与更新时间</p>
          </div>
        </div>

        <div class="metadata-grid">
          <div>
            <span>创建时间</span>
            <strong>{{ formatDate(plan.createdAt) }}</strong>
          </div>
          <div>
            <span>最近更新</span>
            <strong>{{ formatDate(plan.updatedAt) }}</strong>
          </div>
        </div>
      </section>

      <el-alert
        title="BOM 与工艺路线不属于生产计划字段，将在后续工单创建阶段选择。"
        type="info"
        :closable="false"
        show-icon
      />
    </template>
  </div>
</template>

<style scoped>
.detail-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 34px 0 56px;
  color: #17283b;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 26px;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin: 0 0 22px;
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
  font-weight: 750;
}

.description {
  margin: 12px 0 0;
  color: #718096;
  font-size: 13px;
  line-height: 1.8;
}

.header-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-top: 24px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 22px;
}

.info-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
  border: 1px solid #e6ebf1;
  border-radius: 10px;
  background: #fff;
}

.info-card span {
  color: #718096;
  font-size: 12px;
}

.info-card strong {
  overflow-wrap: anywhere;
  color: #0e1f33;
  font-size: 20px;
  font-variant-numeric: tabular-nums;
}

.panel {
  margin-bottom: 22px;
  padding: 24px;
  border: 1px solid #e6ebf1;
  border-radius: 10px;
  background: #fff;
  box-shadow: 0 3px 12px rgb(14 31 51 / 3%);
}

.section-heading {
  display: flex;
  align-items: flex-start;
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

.section-id {
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
