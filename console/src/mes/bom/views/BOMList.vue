<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, Refresh, Search, View } from "@element-plus/icons-vue";

import { listProductModels } from "../../product/api";
import type { ProductModel } from "../../product/types";
import { createBOM, listBOMsByProduct } from "../api";
import type { BOM, CreateBOMRequest } from "../types";

const router = useRouter();

const loading = ref(false);
const productsLoading = ref(false);
const submitting = ref(false);

const products = ref<ProductModel[]>([]);
const boms = ref<BOM[]>([]);
const selectedProductId = ref<number | undefined>();
const keyword = ref("");
const statusFilter = ref("");

const dialogVisible = ref(false);
const formError = ref("");

const form = ref({
  bom_no: "",
  version: 1,
  description: "",
});

const filteredBOMs = computed(() => {
  const query = keyword.value.trim().toLowerCase();

  return boms.value.filter((bom) => {
    const matchesKeyword =
      !query ||
      bom.bom_no.toLowerCase().includes(query) ||
      bom.description?.toLowerCase().includes(query);

    const matchesStatus =
      !statusFilter.value || bom.status === statusFilter.value;

    return matchesKeyword && matchesStatus;
  });
});

const selectedProduct = computed(() =>
  products.value.find((item) => item.id === selectedProductId.value),
);

const summary = computed(() => ({
  total: boms.value.length,
  draft: boms.value.filter((item) => item.status === "draft").length,
  released: boms.value.filter((item) => item.status === "released").length,
  obsolete: boms.value.filter((item) => item.status === "obsolete").length,
}));

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    draft: "草稿",
    released: "已发布",
    obsolete: "已作废",
  };
  return labels[status] ?? status;
}

function statusType(status: string): "info" | "success" | "danger" | "warning" {
  if (status === "released") return "success";
  if (status === "obsolete") return "danger";
  return "info";
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

async function loadProducts() {
  productsLoading.value = true;
  try {
    const response = await listProductModels({
      current_page: 1,
      page_size: 100,
    });

    products.value = response.items ?? [];

    if (selectedProductId.value === undefined && products.value.length > 0) {
      selectedProductId.value = products.value[0].id;
    }
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "产品列表加载失败",
    );
  } finally {
    productsLoading.value = false;
  }
}

async function loadBOMs() {
  if (!selectedProductId.value) {
    boms.value = [];
    return;
  }

  loading.value = true;
  try {
    boms.value = await listBOMsByProduct(selectedProductId.value);
  } catch (error: unknown) {
    boms.value = [];
    ElMessage.error(
      error instanceof Error ? error.message : "BOM 列表加载失败",
    );
  } finally {
    loading.value = false;
  }
}

watch(selectedProductId, () => {
  keyword.value = "";
  statusFilter.value = "";
  void loadBOMs();
});

function openCreate() {
  if (!selectedProductId.value) {
    ElMessage.warning("请先选择产品型号");
    return;
  }

  const existingVersions = boms.value
    .filter((item) => item.status !== "obsolete")
    .map((item) => item.version);

  form.value = {
    bom_no: `${selectedProduct.value?.code ?? "BOM"}-BOM`,
    version: Math.max(0, ...existingVersions) + 1,
    description: "",
  };
  formError.value = "";
  dialogVisible.value = true;
}

async function submitCreate() {
  if (submitting.value || !selectedProductId.value) return;

  if (!form.value.bom_no.trim()) {
    formError.value = "请输入 BOM 编号";
    return;
  }

  if (!Number.isInteger(form.value.version) || form.value.version <= 0) {
    formError.value = "版本必须是大于 0 的整数";
    return;
  }

  const payload: CreateBOMRequest = {
    product_id: selectedProductId.value,
    bom_no: form.value.bom_no.trim(),
    version: form.value.version,
    description: form.value.description.trim(),
    items: [],
  };

  submitting.value = true;
  formError.value = "";

  try {
    const created = await createBOM(payload);
    ElMessage.success("BOM 草稿创建成功");
    dialogVisible.value = false;

    if (created?.id) {
      await router.push(`/mes/boms/${created.id}`);
    } else {
      await loadBOMs();
    }
  } catch (error: unknown) {
    formError.value = error instanceof Error ? error.message : "创建 BOM 失败";
  } finally {
    submitting.value = false;
  }
}

async function confirmRefresh() {
  await loadBOMs();
}

async function confirmObsolete(row: BOM) {
  try {
    await ElMessageBox.confirm(
      `确认作废 BOM「${row.bom_no} / V${row.version}」吗？作废后不可恢复。`,
      "作废 BOM",
      {
        type: "warning",
        confirmButtonText: "确认作废",
        cancelButtonText: "取消",
      },
    );

    const { obsoleteBOM } = await import("../api");
    await obsoleteBOM(row.id);
    ElMessage.success("BOM 已作废");
    await loadBOMs();
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;
    if (error instanceof Error) ElMessage.error(error.message);
  }
}

onMounted(() => {
  void loadProducts();
});
</script>

<template>
  <div class="bom-page">
    <header class="page-header">
      <div>
        <div class="eyebrow">MES / BILL OF MATERIALS</div>
        <h1>BOM 物料清单</h1>
        <p class="description">
          管理产品的物料组成、用量及版本，发布后的 BOM 可供后续生产流程使用。
        </p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="confirmRefresh">
          刷新
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">
          新建 BOM
        </el-button>
      </div>
    </header>

    <section class="summary-grid">
      <article class="summary-card">
        <span>当前产品 BOM 总数</span>
        <strong>{{ summary.total }}</strong>
      </article>
      <article class="summary-card">
        <span>草稿</span>
        <strong>{{ summary.draft }}</strong>
      </article>
      <article class="summary-card">
        <span>已发布</span>
        <strong>{{ summary.released }}</strong>
      </article>
      <article class="summary-card">
        <span>已作废</span>
        <strong>{{ summary.obsolete }}</strong>
      </article>
    </section>

    <section class="panel">
      <div class="filter-row">
        <div class="product-select">
          <label>产品型号</label>
          <el-select
            v-model="selectedProductId"
            placeholder="选择产品型号"
            filterable
            :loading="productsLoading"
            style="width: 280px"
          >
            <el-option
              v-for="product in products"
              :key="product.id"
              :label="`${product.name}（${product.code}）`"
              :value="product.id"
            />
          </el-select>
        </div>

        <el-input
          v-model="keyword"
          clearable
          :prefix-icon="Search"
          placeholder="搜索 BOM 编号或描述"
          style="width: 250px"
        />

        <el-select
          v-model="statusFilter"
          clearable
          placeholder="全部状态"
          style="width: 140px"
        >
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="released" />
          <el-option label="已作废" value="obsolete" />
        </el-select>

        <span class="result-count"> 共 {{ filteredBOMs.length }} 条记录 </span>
      </div>

      <div v-if="selectedProduct" class="selected-product">
        <span class="product-dot"></span>
        <span
          >当前产品：<strong>{{ selectedProduct.name }}</strong></span
        >
        <span class="product-code">{{ selectedProduct.code }}</span>
        <span class="product-version"
          >产品版本 {{ selectedProduct.version }}</span
        >
      </div>

      <el-table
        v-loading="loading"
        :data="filteredBOMs"
        row-key="id"
        empty-text="暂无 BOM，请选择产品并新建 BOM 草稿"
        style="width: 100%"
      >
        <el-table-column label="BOM 编号" min-width="190">
          <template #default="{ row }">
            <div class="bom-number">{{ row.bom_no }}</div>
            <div class="secondary-text">ID {{ row.id }}</div>
          </template>
        </el-table-column>

        <el-table-column label="版本" width="100">
          <template #default="{ row }">
            <span class="version-pill">V{{ row.version }}</span>
          </template>
        </el-table-column>

        <el-table-column label="物料行数" width="110" align="center">
          <template #default="{ row }">
            {{ row.items?.length ?? 0 }}
          </template>
        </el-table-column>

        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="描述" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.description || "—" }}
          </template>
        </el-table-column>

        <el-table-column label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="190" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :icon="View"
              @click="router.push(`/mes/boms/${row.id}`)"
            >
              查看详情
            </el-button>
            <el-button
              v-if="row.status === 'released'"
              link
              type="danger"
              @click="confirmObsolete(row)"
            >
              作废
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!productsLoading && products.length === 0" class="empty-hint">
        暂无产品型号。请先在产品模块创建产品型号，再创建 BOM。
      </div>
    </section>

    <el-dialog
      v-model="dialogVisible"
      title="新建 BOM 草稿"
      width="560px"
      destroy-on-close
    >
      <div class="dialog-context">
        为产品
        <strong>{{ selectedProduct?.name }}</strong>
        创建新的 BOM 版本。创建后可继续添加物料明细。
      </div>

      <el-form label-position="top">
        <el-form-item label="BOM 编号" required>
          <el-input
            v-model="form.bom_no"
            maxlength="100"
            placeholder="例如：PRODUCT-BOM"
          />
        </el-form-item>

        <el-form-item label="版本号" required>
          <el-input-number
            v-model="form.version"
            :min="1"
            :precision="0"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            maxlength="1000"
            show-word-limit
            placeholder="填写该 BOM 版本的用途或变更说明"
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
.bom-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 32px 0 48px;
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
  color: #9b6a22;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.5px;
}

h1 {
  margin: 0;
  color: #0e1f33;
  font-size: 28px;
  font-weight: 700;
}

.description {
  margin: 10px 0 0;
  color: #738195;
  font-size: 13px;
  line-height: 1.7;
}

.header-actions {
  display: flex;
  gap: 10px;
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
  gap: 13px;
  padding: 20px 22px;
  border: 1px solid #e6eaf0;
  border-radius: 10px;
  background: #fff;
}

.summary-card span {
  color: #738195;
  font-size: 12px;
}

.summary-card strong {
  color: #0e1f33;
  font-size: 27px;
  font-weight: 700;
}

.panel {
  padding: 22px;
  border: 1px solid #e6eaf0;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 4px 16px rgb(14 31 51 / 3%);
}

.filter-row {
  display: flex;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 14px;
  margin-bottom: 20px;
}

.product-select {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.product-select label {
  color: #56667a;
  font-size: 12px;
  font-weight: 600;
}

.result-count {
  margin-left: auto;
  padding-bottom: 9px;
  color: #7b8796;
  font-size: 12px;
}

.selected-product {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 9px;
  margin-bottom: 18px;
  padding: 12px 14px;
  border: 1px solid #e7edf4;
  border-radius: 8px;
  color: #5e6f83;
  background: #f7f9fc;
  font-size: 12px;
}

.selected-product strong {
  color: #0e1f33;
}

.product-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #f0a030;
}

.product-code,
.product-version,
.version-pill {
  padding: 4px 8px;
  border-radius: 5px;
  color: #52677d;
  background: #eaf0f6;
  font-size: 11px;
}

.bom-number {
  color: #17283b;
  font-size: 13px;
  font-weight: 700;
}

.secondary-text {
  margin-top: 5px;
  color: #9aa5b2;
  font-size: 11px;
}

.empty-hint {
  padding: 24px 8px 4px;
  color: #8995a4;
  font-size: 12px;
  text-align: center;
}

.dialog-context {
  margin-bottom: 22px;
  padding: 12px 14px;
  border-radius: 8px;
  color: #627286;
  background: #f4f7fa;
  font-size: 12px;
  line-height: 1.7;
}

.dialog-context strong {
  color: #0e1f33;
}

@media (max-width: 900px) {
  .bom-page {
    width: calc(100% - 32px);
  }

  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .page-header {
    flex-direction: column;
  }
}

@media (max-width: 560px) {
  .bom-page {
    width: calc(100% - 24px);
    padding-top: 22px;
  }

  .summary-grid {
    gap: 10px;
  }

  .summary-card {
    padding: 15px;
  }

  .panel {
    padding: 14px;
  }

  .filter-row > * {
    width: 100% !important;
  }

  .result-count {
    margin-left: 0;
  }
}
</style>
