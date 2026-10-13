<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, Refresh, Search, View } from "@element-plus/icons-vue";

import { listProductModels } from "../../product/api";
import type { ProductModel } from "../../product/types";
import {
  activateRouting,
  createRouting,
  deactivateRouting,
  listRoutings,
} from "../api";
import type { CreateRoutingRequest, Routing } from "../types";

const router = useRouter();

const loading = ref(false);
const productsLoading = ref(false);
const submitting = ref(false);
const actionId = ref<number | null>(null);

const products = ref<ProductModel[]>([]);
const routings = ref<Routing[]>([]);
const productId = ref<number | undefined>();
const statusFilter = ref("");
const keyword = ref("");

const dialogVisible = ref(false);
const formError = ref("");
const form = ref({ name: "", description: "" });

const selectedProduct = computed(
  () => products.value.find((item) => item.id === productId.value) ?? null,
);

const filteredRoutings = computed(() => {
  const q = keyword.value.trim().toLowerCase();

  return routings.value.filter((item) => {
    const matchesKeyword =
      !q ||
      item.name.toLowerCase().includes(q) ||
      item.code.toLowerCase().includes(q) ||
      item.description.toLowerCase().includes(q);

    return (
      matchesKeyword &&
      (!statusFilter.value || item.status === statusFilter.value)
    );
  });
});

const summary = computed(() => ({
  total: routings.value.length,
  active: routings.value.filter((item) => item.status === "active").length,
  draft: routings.value.filter((item) => item.status === "draft").length,
  inactive: routings.value.filter((item) => item.status === "inactive").length,
}));

function statusLabel(status: string) {
  const labels: Record<string, string> = {
    draft: "草稿",
    active: "已激活",
    inactive: "已停用",
    archived: "已归档",
  };
  return labels[status] ?? status;
}

function statusType(status: string): "info" | "success" | "warning" | "danger" {
  if (status === "active") return "success";
  if (status === "inactive") return "warning";
  if (status === "archived") return "danger";
  return "info";
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

async function loadProducts() {
  productsLoading.value = true;
  try {
    const response = await listProductModels({
      current_page: 1,
      page_size: 100,
    });
    products.value = response.items ?? [];
    if (productId.value === undefined && products.value.length > 0) {
      productId.value = products.value[0].id;
    }
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "产品列表加载失败",
    );
  } finally {
    productsLoading.value = false;
  }
}

async function loadRoutings() {
  loading.value = true;
  try {
    routings.value = await listRoutings({
      productId: productId.value,
    });
  } catch (error: unknown) {
    routings.value = [];
    ElMessage.error(
      error instanceof Error ? error.message : "工艺路线加载失败",
    );
  } finally {
    loading.value = false;
  }
}

watch(productId, () => {
  statusFilter.value = "";
  keyword.value = "";
  void loadRoutings();
});

function openCreate() {
  if (!productId.value) {
    ElMessage.warning("请先创建并选择产品型号");
    return;
  }

  form.value = {
    name: `${selectedProduct.value?.name ?? "产品"}标准工艺`,
    description: "",
  };
  formError.value = "";
  dialogVisible.value = true;
}

async function submitCreate() {
  if (submitting.value || !productId.value) return;

  const name = form.value.name.trim();
  if (!name) {
    formError.value = "请输入工艺路线名称";
    return;
  }

  const payload: CreateRoutingRequest = {
    productId: productId.value,
    name,
    description: form.value.description.trim(),
  };

  submitting.value = true;
  formError.value = "";

  try {
    const created = await createRouting(payload);
    ElMessage.success("工艺路线创建成功");
    dialogVisible.value = false;
    await loadRoutings();

    if (created?.id) {
      await router.push(`/mes/routings/${created.id}`);
    }
  } catch (error: unknown) {
    formError.value =
      error instanceof Error ? error.message : "创建工艺路线失败";
  } finally {
    submitting.value = false;
  }
}

function openDetail(row: Routing) {
  void router.push(`/mes/routings/${row.id}`);
}

async function changeStatus(row: Routing, action: "activate" | "deactivate") {
  const activating = action === "activate";
  const actionText = activating ? "激活" : "停用";

  try {
    await ElMessageBox.confirm(
      activating
        ? `确定激活工艺路线「${row.name}」吗？激活后，该产品的其他有效路线可能被切换。`
        : `确定停用工艺路线「${row.name}」吗？如果该路线存在未完成工单，后端将拒绝操作。`,
      `${actionText}工艺路线`,
      {
        type: "warning",
        confirmButtonText: `确认${actionText}`,
        cancelButtonText: "取消",
      },
    );

    actionId.value = row.id;
    if (activating) {
      await activateRouting(row.id);
    } else {
      await deactivateRouting(row.id);
    }

    ElMessage.success(`工艺路线${actionText}成功`);
    await loadRoutings();
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;
    ElMessage.error(
      error instanceof Error ? error.message : `${actionText}失败`,
    );
  } finally {
    actionId.value = null;
  }
}

onMounted(() => {
  void loadProducts();
});
</script>

<template>
  <div class="routing-page">
    <header class="page-header">
      <div>
        <div class="eyebrow">MES / PROCESS ROUTING</div>
        <h1>工艺路线</h1>
        <p class="description">
          定义产品的制造工序、执行顺序、工作站与标准工时，为后续生产计划和工单执行提供依据。
        </p>
      </div>
      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadRoutings">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">
          新建工艺路线
        </el-button>
      </div>
    </header>

    ```
    <section class="summary-grid">
      <article class="summary-card">
        <span class="summary-label">全部路线</span>
        <strong>{{ summary.total }}</strong>
        <span class="summary-note">当前产品的路线配置</span>
      </article>
      <article class="summary-card summary-card--green">
        <span class="summary-label">已激活</span>
        <strong>{{ summary.active }}</strong>
        <span class="summary-note">可供生产流程使用</span>
      </article>
      <article class="summary-card summary-card--blue">
        <span class="summary-label">草稿</span>
        <strong>{{ summary.draft }}</strong>
        <span class="summary-note">仍可配置工序</span>
      </article>
      <article class="summary-card summary-card--orange">
        <span class="summary-label">已停用</span>
        <strong>{{ summary.inactive }}</strong>
        <span class="summary-note">可按规则重新激活</span>
      </article>
    </section>

    <section class="content-card">
      <div class="filter-bar">
        <el-select
          v-model="productId"
          placeholder="选择产品型号"
          filterable
          :loading="productsLoading"
          class="product-select"
        >
          <el-option
            v-for="item in products"
            :key="item.id"
            :label="`${item.code} · ${item.name}`"
            :value="item.id"
          />
        </el-select>

        <el-input
          v-model="keyword"
          clearable
          :prefix-icon="Search"
          placeholder="搜索路线名称、编码或描述"
          class="keyword-input"
        />

        <el-select
          v-model="statusFilter"
          clearable
          placeholder="全部状态"
          class="status-select"
        >
          <el-option label="草稿" value="draft" />
          <el-option label="已激活" value="active" />
          <el-option label="已停用" value="inactive" />
          <el-option label="已归档" value="archived" />
        </el-select>

        <el-button :icon="Refresh" @click="loadRoutings">查询</el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="filteredRoutings"
        row-key="id"
        stripe
        class="routing-table"
        @row-dblclick="openDetail"
      >
        <el-table-column label="路线名称" min-width="210">
          <template #default="{ row }">
            <button class="route-name" type="button" @click="openDetail(row)">
              {{ row.name }}
            </button>
            <div class="route-code">{{ row.code }}</div>
          </template>
        </el-table-column>

        <el-table-column label="版本" width="90" align="center">
          <template #default="{ row }">V{{ row.version }}</template>
        </el-table-column>

        <el-table-column label="产品型号" min-width="180">
          <template #default>
            <span>{{ selectedProduct?.code ?? "—" }}</span>
            <div class="secondary-text">{{ selectedProduct?.name ?? "—" }}</div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="说明" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.description || "暂无描述" }}
          </template>
        </el-table-column>

        <el-table-column label="更新时间" min-width="165">
          <template #default="{ row }">{{
            formatDate(row.updatedAt)
          }}</template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="210" align="center">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :icon="View"
              @click="openDetail(row)"
            >
              详情
            </el-button>
            <el-button
              v-if="row.status !== 'active' && row.status !== 'archived'"
              link
              type="success"
              :loading="actionId === row.id"
              @click="changeStatus(row, 'activate')"
            >
              激活
            </el-button>
            <el-button
              v-if="row.status === 'active'"
              link
              type="warning"
              :loading="actionId === row.id"
              @click="changeStatus(row, 'deactivate')"
            >
              停用
            </el-button>
          </template>
        </el-table-column>

        <template #empty>
          <el-empty
            :description="
              productId ? '当前产品暂无工艺路线' : '请先选择一个产品型号'
            "
          >
            <el-button v-if="productId" type="primary" @click="openCreate">
              创建第一条工艺路线
            </el-button>
          </el-empty>
        </template>
      </el-table>

      <div class="table-footer">
        <span>共 {{ filteredRoutings.length }} 条路线</span>
        <span>工艺路线用于定义制造步骤，不代替 BOM 的物料结构。</span>
      </div>
    </section>

    <el-dialog
      v-model="dialogVisible"
      title="新建工艺路线"
      width="560px"
      destroy-on-close
    >
      <el-alert
        v-if="formError"
        :title="formError"
        type="error"
        :closable="false"
        show-icon
        class="form-alert"
      />

      <el-form label-position="top" @submit.prevent="submitCreate">
        <el-form-item label="关联产品型号">
          <el-input
            :model-value="
              selectedProduct
                ? `${selectedProduct.code} · ${selectedProduct.name}`
                : ''
            "
            disabled
          />
        </el-form-item>
        <el-form-item label="路线名称" required>
          <el-input
            v-model="form.name"
            maxlength="255"
            show-word-limit
            placeholder="例如：标准装配工艺"
            @keyup.enter="submitCreate"
          />
        </el-form-item>
        <el-form-item label="工艺说明">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            maxlength="1000"
            show-word-limit
            placeholder="说明适用范围、工艺要求或注意事项"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate">
          创建并配置工序
        </el-button>
      </template>
    </el-dialog>
    ```
  </div>
</template>

<style scoped>
.routing-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 34px 0 48px;
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
  color: #718197;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.6px;
}
h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 750;
  letter-spacing: -0.5px;
}
.description {
  max-width: 760px;
  margin: 10px 0 0;
  color: #718197;
  font-size: 13px;
  line-height: 1.8;
}
.header-actions {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}
.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 22px;
}
.summary-card {
  display: flex;
  flex-direction: column;
  min-height: 130px;
  padding: 20px;
  border: 1px solid #e5eaf0;
  border-radius: 12px;
  background: #fff;
}
.summary-label {
  color: #687a90;
  font-size: 12px;
}
.summary-card strong {
  margin-top: 8px;
  color: #17283b;
  font-size: 30px;
  line-height: 1.2;
}
.summary-note {
  margin-top: 8px;
  color: #9aa7b6;
  font-size: 11px;
}
.summary-card--green strong {
  color: #258354;
}
.summary-card--blue strong {
  color: #326bb3;
}
.summary-card--orange strong {
  color: #bd7b25;
}
.content-card {
  overflow: hidden;
  border: 1px solid #e5eaf0;
  border-radius: 12px;
  background: #fff;
}
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding: 18px 20px;
  border-bottom: 1px solid #edf0f4;
}
.product-select {
  width: 270px;
}
.keyword-input {
  width: 300px;
}
.status-select {
  width: 140px;
}
.routing-table {
  width: 100%;
}
.route-name {
  padding: 0;
  border: 0;
  color: #245fa8;
  background: transparent;
  font: inherit;
  font-size: 13px;
  font-weight: 700;
  text-align: left;
  cursor: pointer;
}
.route-name:hover {
  text-decoration: underline;
}
.route-code,
.secondary-text {
  margin-top: 5px;
  color: #8b9aac;
  font-size: 11px;
}
.table-footer {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 20px;
  color: #8b9aac;
  font-size: 11px;
}
.form-alert {
  margin-bottom: 16px;
}
@media (max-width: 900px) {
  .routing-page {
    width: calc(100% - 32px);
    padding-top: 24px;
  }
  .page-header {
    flex-direction: column;
  }
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .product-select,
  .keyword-input {
    width: 100%;
  }
  .table-footer {
    flex-direction: column;
  }
}
</style>
