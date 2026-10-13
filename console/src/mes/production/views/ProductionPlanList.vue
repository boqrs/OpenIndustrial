<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Plus, Refresh, Search, View } from "@element-plus/icons-vue";

import { listProductModels } from "../../product/api";
import type { ProductModel } from "../../product/types";
import { createProductionPlan, listProductionPlans } from "../api";
import type {
  CreateProductionPlanRequest,
  ProductionPlan,
  ProductionPlanStatus,
} from "../types";

const router = useRouter();

const loading = ref(false);
const productsLoading = ref(false);
const submitting = ref(false);

const plans = ref<ProductionPlan[]>([]);
const products = ref<ProductModel[]>([]);
const keyword = ref("");
const statusFilter = ref<ProductionPlanStatus | "">("");

const dialogVisible = ref(false);
const formError = ref("");

const form = reactive({
  planNo: "",
  productID: undefined as number | undefined,
  factoryID: undefined as number | undefined,
  plannedQuantity: 1,
  plannedStartAt: "",
  plannedEndAt: "",
  description: "",
});

const filteredPlans = computed(() => {
  const query = keyword.value.trim().toLowerCase();

  return plans.value.filter((plan) => {
    const matchesKeyword =
      !query ||
      plan.planNo.toLowerCase().includes(query) ||
      String(plan.productId).includes(query) ||
      String(plan.factoryId).includes(query) ||
      plan.description?.toLowerCase().includes(query);

    const matchesStatus =
      !statusFilter.value || plan.status === statusFilter.value;

    return matchesKeyword && matchesStatus;
  });
});

const summary = computed(() => ({
  total: plans.value.length,
  draft: plans.value.filter((item) => item.status === "draft").length,
  released: plans.value.filter((item) => item.status === "released").length,
  inProgress: plans.value.filter((item) => item.status === "in_progress")
    .length,
  completed: plans.value.filter((item) => item.status === "completed").length,
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

function toRFC3339(value: string): string {
  // 日期控件使用本地时间展示，提交时转换为后端 time.Time 可解析的 RFC3339。
  return new Date(value.replace(" ", "T")).toISOString();
}

async function loadPlans() {
  loading.value = true;

  try {
    plans.value = await listProductionPlans(statusFilter.value || undefined);
  } catch (error: unknown) {
    plans.value = [];
    ElMessage.error(
      error instanceof Error ? error.message : "生产计划加载失败",
    );
  } finally {
    loading.value = false;
  }
}

async function loadProducts() {
  productsLoading.value = true;

  try {
    const response = await listProductModels({
      current_page: 1,
      page_size: 100,
    });

    products.value = response.items ?? [];
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "产品型号加载失败",
    );
  } finally {
    productsLoading.value = false;
  }
}

function openCreate() {
  const now = new Date();
  const start = new Date(now.getTime() + 60 * 60 * 1000);
  const end = new Date(start.getTime() + 8 * 60 * 60 * 1000);

  const formatLocalInput = (date: Date) => {
    const pad = (n: number) => String(n).padStart(2, "0");

    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
      date.getDate(),
    )} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(
      date.getSeconds(),
    )}`;
  };

  Object.assign(form, {
    planNo: "",
    productID: products.value[0]?.id,
    factoryID: undefined,
    plannedQuantity: 1,
    plannedStartAt: formatLocalInput(start),
    plannedEndAt: formatLocalInput(end),
    description: "",
  });

  formError.value = "";
  dialogVisible.value = true;
}

async function submitCreate() {
  if (submitting.value) return;

  if (!form.planNo.trim()) {
    formError.value = "请输入生产计划编号";
    return;
  }

  if (!form.productID) {
    formError.value = "请选择产品型号";
    return;
  }

  if (!form.factoryID || !Number.isInteger(form.factoryID)) {
    formError.value = "请输入有效的工厂 ID";
    return;
  }

  if (!Number.isInteger(form.plannedQuantity) || form.plannedQuantity <= 0) {
    formError.value = "计划数量必须是大于 0 的整数";
    return;
  }

  const start = new Date(form.plannedStartAt.replace(" ", "T"));
  const end = new Date(form.plannedEndAt.replace(" ", "T"));

  if (
    Number.isNaN(start.getTime()) ||
    Number.isNaN(end.getTime()) ||
    end.getTime() <= start.getTime()
  ) {
    formError.value = "计划结束时间必须晚于开始时间";
    return;
  }

  const payload: CreateProductionPlanRequest = {
    planNo: form.planNo.trim(),
    productID: form.productID,
    factoryID: form.factoryID,
    plannedQuantity: form.plannedQuantity,
    plannedStartAt: toRFC3339(form.plannedStartAt),
    plannedEndAt: toRFC3339(form.plannedEndAt),
    description: form.description.trim(),
  };

  submitting.value = true;
  formError.value = "";

  try {
    const created = await createProductionPlan(payload);

    ElMessage.success("生产计划草稿创建成功");
    dialogVisible.value = false;

    await loadPlans();

    if (created?.id) {
      await router.push(`/mes/production-plans/${created.id}`);
    }
  } catch (error: unknown) {
    formError.value =
      error instanceof Error ? error.message : "创建生产计划失败";
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  void loadProducts();
  void loadPlans();
});
</script>

<template>
  <div class="production-page">
    <header class="page-header">
      <div>
        <div class="eyebrow">MES / PRODUCTION PLANNING</div>
        <h1>生产计划</h1>
        <p class="description">
          管理工厂生产计划，跟踪计划状态与生产周期，为后续工单创建提供计划依据。
        </p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadPlans">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">
          新建计划
        </el-button>
      </div>
    </header>

    <section class="summary-grid">
      <article class="summary-card">
        <span>计划总数</span>
        <strong>{{ summary.total }}</strong>
        <small>当前查询结果</small>
      </article>
      <article class="summary-card">
        <span>草稿</span>
        <strong>{{ summary.draft }}</strong>
        <small>尚未发布</small>
      </article>
      <article class="summary-card">
        <span>已发布</span>
        <strong>{{ summary.released }}</strong>
        <small>等待后续执行</small>
      </article>
      <article class="summary-card">
        <span>执行中</span>
        <strong>{{ summary.inProgress }}</strong>
        <small>生产流程推进中</small>
      </article>
      <article class="summary-card">
        <span>已完成</span>
        <strong>{{ summary.completed }}</strong>
        <small>已结束的生产计划</small>
      </article>
    </section>

    <section class="panel">
      <div class="filter-row">
        <el-input
          v-model="keyword"
          clearable
          :prefix-icon="Search"
          placeholder="搜索计划编号、产品 ID、工厂 ID"
          style="width: 300px"
        />

        <el-select
          v-model="statusFilter"
          clearable
          placeholder="全部状态"
          style="width: 150px"
          @change="loadPlans"
        >
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="released" />
          <el-option label="执行中" value="in_progress" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>

        <span class="result-count"> 共 {{ filteredPlans.length }} 条记录 </span>
      </div>

      <el-table
        v-loading="loading"
        :data="filteredPlans"
        row-key="id"
        empty-text="暂无生产计划，请新建生产计划草稿"
        style="width: 100%"
      >
        <el-table-column label="计划编号" min-width="190">
          <template #default="{ row }">
            <button
              type="button"
              class="plan-link"
              @click="router.push(`/mes/production-plans/${row.id}`)"
            >
              {{ row.planNo }}
            </button>
            <div class="secondary-text">ID {{ row.id }}</div>
          </template>
        </el-table-column>

        <el-table-column label="产品 ID" prop="productId" width="110" />
        <el-table-column label="工厂 ID" prop="factoryId" width="110" />

        <el-table-column label="计划数量" width="120" align="right">
          <template #default="{ row }">
            <strong class="quantity">{{ row.plannedQuantity }}</strong>
          </template>
        </el-table-column>

        <el-table-column label="计划开始" min-width="165">
          <template #default="{ row }">
            {{ formatDate(row.plannedStartAt) }}
          </template>
        </el-table-column>

        <el-table-column label="计划结束" min-width="165">
          <template #default="{ row }">
            {{ formatDate(row.plannedEndAt) }}
          </template>
        </el-table-column>

        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :icon="View"
              @click="router.push(`/mes/production-plans/${row.id}`)"
            >
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>

    <el-dialog
      v-model="dialogVisible"
      title="新建生产计划"
      width="640px"
      destroy-on-close
    >
      <el-alert
        title="新建的计划将以草稿状态保存。BOM 和工艺路线在后续工单创建阶段选择。"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <el-form label-position="top" @submit.prevent="submitCreate">
        <el-form-item label="生产计划编号" required>
          <el-input
            v-model="form.planNo"
            maxlength="64"
            placeholder="例如：PLAN-2026-001"
          />
        </el-form-item>

        <div class="form-grid">
          <el-form-item label="产品型号" required>
            <el-select
              v-model="form.productID"
              placeholder="选择产品型号"
              filterable
              :loading="productsLoading"
              style="width: 100%"
            >
              <el-option
                v-for="product in products"
                :key="product.id"
                :label="`${product.name}（${product.code}）`"
                :value="product.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="工厂 ID" required>
            <el-input-number
              v-model="form.factoryID"
              :min="1"
              :precision="0"
              controls-position="right"
              placeholder="输入工厂 ID"
              style="width: 100%"
            />
          </el-form-item>
        </div>

        <el-form-item label="计划生产数量" required>
          <el-input-number
            v-model="form.plannedQuantity"
            :min="1"
            :precision="0"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>

        <div class="form-grid">
          <el-form-item label="计划开始时间" required>
            <el-date-picker
              v-model="form.plannedStartAt"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择开始时间"
              style="width: 100%"
            />
          </el-form-item>

          <el-form-item label="计划结束时间" required>
            <el-date-picker
              v-model="form.plannedEndAt"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择结束时间"
              style="width: 100%"
            />
          </el-form-item>
        </div>

        <el-form-item label="计划描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
            placeholder="填写本次生产计划的说明（选填）"
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
          创建草稿
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.production-page {
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

.eyebrow {
  margin-bottom: 10px;
  color: #a16a19;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.5px;
}

.page-header h1 {
  margin: 0;
  color: #0e1f33;
  font-size: 28px;
  font-weight: 750;
  letter-spacing: -0.5px;
}

.description {
  max-width: 720px;
  margin: 12px 0 0;
  color: #718096;
  font-size: 13px;
  line-height: 1.8;
}

.header-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
  padding-top: 24px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 22px;
}

.summary-card {
  display: flex;
  flex-direction: column;
  min-height: 116px;
  padding: 20px;
  border: 1px solid #e6ebf1;
  border-radius: 10px;
  background: #fff;
  box-shadow: 0 3px 12px rgb(14 31 51 / 3%);
}

.summary-card span {
  color: #718096;
  font-size: 12px;
}

.summary-card strong {
  margin-top: 12px;
  color: #0e1f33;
  font-size: 27px;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.summary-card small {
  margin-top: 10px;
  color: #9aa6b5;
  font-size: 11px;
}

.panel {
  padding: 22px;
  border: 1px solid #e6ebf1;
  border-radius: 10px;
  background: #fff;
  box-shadow: 0 3px 12px rgb(14 31 51 / 3%);
}

.filter-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 20px;
}

.result-count {
  margin-left: auto;
  color: #8a97a8;
  font-size: 12px;
}

.plan-link {
  padding: 0;
  border: 0;
  color: #245b91;
  background: transparent;
  font: inherit;
  font-weight: 650;
  cursor: pointer;
}

.plan-link:hover {
  color: #d28b25;
}

.secondary-text {
  margin-top: 5px;
  color: #9aa6b5;
  font-size: 11px;
}

.quantity {
  color: #0e1f33;
  font-variant-numeric: tabular-nums;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

@media (max-width: 1000px) {
  .summary-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .production-page {
    width: calc(100% - 32px);
  }
}

@media (max-width: 640px) {
  .page-header {
    flex-direction: column;
  }

  .header-actions {
    padding-top: 0;
  }

  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .panel {
    padding: 14px;
  }

  .form-grid {
    grid-template-columns: 1fr;
    gap: 0;
  }
}
</style>
