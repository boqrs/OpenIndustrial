<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import {
  ElMessage,
  ElMessageBox,
  type FormInstance,
  type FormRules,
} from "element-plus";
import {
  Plus,
  Refresh,
  Search,
  View,
  EditPen,
  CircleCheck,
  CircleClose,
  Folder,
} from "@element-plus/icons-vue";

import {
  activateProductModel,
  archiveProductModel,
  createProductModel,
  deactivateProductModel,
  listProductModels,
} from "../api";
import type {
  AttributeDataType,
  AttributeDefinitionInput,
  CreateProductModelRequest,
  ProductModel,
} from "../types";

const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const loadError = ref("");
const products = ref<ProductModel[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);

const keyword = ref("");
const statusFilter = ref("");
const categoryFilter = ref("");

const createDialogVisible = ref(false);
const createFormRef = ref<FormInstance>();
const formError = ref("");

const createForm = reactive({
  name: "",
  code: "",
  version: "1.0.0",
  category: "",
  description: "",
});

interface AttributeFormRow {
  name: string;
  label: string;
  description: string;
  data_type: AttributeDataType;
  unit: string;
  required: boolean;
}

const attributeRows = ref<AttributeFormRow[]>([]);

const supportedTypes: { label: string; value: AttributeDataType }[] = [
  { label: "字符串", value: "string" },
  { label: "长文本", value: "text" },
  { label: "整数", value: "integer" },
  { label: "浮点数", value: "float" },
  { label: "布尔值", value: "boolean" },
  { label: "日期时间", value: "datetime" },
  { label: "JSON", value: "json" },
];

const formRules: FormRules = {
  name: [{ required: true, message: "请输入产品名称", trigger: "blur" }],
  code: [{ required: true, message: "请输入产品编码", trigger: "blur" }],
  version: [{ required: true, message: "请输入产品版本", trigger: "blur" }],
  category: [{ required: true, message: "请输入产品分类", trigger: "blur" }],
};

const categories = computed(() => {
  const values = products.value
    .map((item) => item.category)
    .filter((value): value is string => Boolean(value?.trim()));

  return [...new Set(values)].sort((a, b) => a.localeCompare(b));
});

const summaryCards = computed(() => [
  {
    key: "all",
    label: "产品型号总数",
    value: total.value,
    description: "当前筛选条件下的记录",
  },
  {
    key: "pending",
    label: "待启用",
    value: products.value.filter((item) => item.status === "pending").length,
    description: "等待完成定义后启用",
  },
  {
    key: "active",
    label: "已启用",
    value: products.value.filter((item) => item.status === "active").length,
    description: "可用于后续业务流程",
  },
  {
    key: "inactive",
    label: "已停用",
    value: products.value.filter((item) => item.status === "inactive").length,
    description: "暂不用于新业务",
  },
]);

function statusLabel(status?: string): string {
  const labels: Record<string, string> = {
    pending: "待启用",
    active: "已启用",
    inactive: "已停用",
    archived: "已归档",
  };

  return labels[status ?? ""] ?? status ?? "未知";
}

function statusTagType(
  status?: string,
): "success" | "warning" | "info" | "danger" {
  switch (status) {
    case "active":
      return "success";
    case "pending":
      return "warning";
    case "inactive":
      return "info";
    case "archived":
      return "danger";
    default:
      return "info";
  }
}

/** 为各个状态提供独立的视觉样式。 */
function statusClass(status?: string): string {
  switch (status) {
    case "active":
      return "status-tag--active";
    case "pending":
      return "status-tag--pending";
    case "inactive":
      return "status-tag--inactive";
    case "archived":
      return "status-tag--archived";
    default:
      return "status-tag--unknown";
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

function resetCreateForm() {
  createForm.name = "";
  createForm.code = "";
  createForm.version = "1.0.0";
  createForm.category = "";
  createForm.description = "";
  attributeRows.value = [];
  formError.value = "";
  createFormRef.value?.clearValidate();
}

function openCreateDialog() {
  resetCreateForm();
  createDialogVisible.value = true;
}

function addAttributeRow() {
  attributeRows.value.push({
    name: "",
    label: "",
    description: "",
    data_type: "string",
    unit: "",
    required: false,
  });
}

function removeAttributeRow(index: number) {
  attributeRows.value.splice(index, 1);
}

function buildAttributes(): Record<string, AttributeDefinitionInput> | string {
  const result: Record<string, AttributeDefinitionInput> = {};

  for (const row of attributeRows.value) {
    const name = row.name.trim();

    if (!name) {
      return "请填写所有属性的英文标识名";
    }

    if (!/^[a-zA-Z][a-zA-Z0-9_]*$/.test(name)) {
      return `属性标识“${name}”必须以英文字母开头，且只能包含字母、数字和下划线`;
    }

    if (Object.prototype.hasOwnProperty.call(result, name)) {
      return `属性标识“${name}”重复，请修改后再提交`;
    }

    if (!row.label.trim()) {
      return `请填写属性“${name}”的显示名称`;
    }

    result[name] = {
      label: row.label.trim(),
      description: row.description.trim() || undefined,
      data_type: row.data_type,
      unit: row.unit.trim() || undefined,
      required: row.required,
    };
  }

  return result;
}

async function loadProducts() {
  loading.value = true;
  loadError.value = "";

  try {
    const response = await listProductModels({
      code: keyword.value.trim() || undefined,
      category: categoryFilter.value || undefined,
      status: statusFilter.value || undefined,
      current_page: currentPage.value,
      page_size: pageSize.value,
    });

    products.value = response.items ?? [];
    total.value = response.total ?? 0;
  } catch (error: unknown) {
    loadError.value =
      error instanceof Error
        ? error.message
        : "产品型号列表加载失败，请稍后重试";
  } finally {
    loading.value = false;
  }
}

function searchProducts() {
  currentPage.value = 1;
  void loadProducts();
}

function resetFilters() {
  keyword.value = "";
  statusFilter.value = "";
  categoryFilter.value = "";
  currentPage.value = 1;
  void loadProducts();
}

function handlePageSizeChange() {
  currentPage.value = 1;
  void loadProducts();
}

function openDetail(row: ProductModel) {
  void router.push(`/mes/products/${row.id}`);
}

function openEdit(row: ProductModel) {
  void router.push(`/mes/products/${row.id}?edit=1`);
}

async function submitCreate() {
  if (submitting.value) return;

  const valid = await createFormRef.value
    ?.validate()
    .then(() => true)
    .catch(() => false);

  if (!valid) return;

  const attributes = buildAttributes();
  if (typeof attributes === "string") {
    formError.value = attributes;
    return;
  }

  const payload: CreateProductModelRequest = {
    name: createForm.name.trim(),
    code: createForm.code.trim(),
    version: createForm.version.trim(),
    category: createForm.category.trim(),
    description: createForm.description.trim() || undefined,
    attributes,
  };

  submitting.value = true;
  formError.value = "";

  try {
    const created = await createProductModel(payload);
    ElMessage.success("产品型号创建成功");
    createDialogVisible.value = false;
    currentPage.value = 1;
    await loadProducts();

    if (created?.id) {
      void router.push(`/mes/products/${created.id}`);
    }
  } catch (error: unknown) {
    formError.value =
      error instanceof Error ? error.message : "产品型号创建失败";
  } finally {
    submitting.value = false;
  }
}

async function changeStatus(
  row: ProductModel,
  action: "activate" | "deactivate" | "archive",
) {
  const actionText = {
    activate: "启用",
    deactivate: "停用",
    archive: "归档",
  }[action];

  try {
    await ElMessageBox.confirm(
      `确认${actionText}产品型号“${row.name}（${row.code}）”吗？`,
      `确认${actionText}`,
      {
        type: action === "archive" ? "warning" : "info",
        confirmButtonText: `确认${actionText}`,
        cancelButtonText: "取消",
      },
    );
  } catch {
    return;
  }

  try {
    if (action === "activate") await activateProductModel(row.id);
    if (action === "deactivate") await deactivateProductModel(row.id);
    if (action === "archive") await archiveProductModel(row.id);

    ElMessage.success(`已提交${actionText}操作`);
    await loadProducts();
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : `${actionText}操作失败`,
    );
  }
}

onMounted(() => {
  void loadProducts();
});
</script>

<template>
  <div class="product-page">
    <header class="page-header">
      <div>
        <div class="page-eyebrow">MES / PRODUCT MODEL</div>
        <h1>产品型号</h1>
        <p class="page-description">
          定义产品型号、版本与静态属性，为生产制造及后续设备身份创建提供统一依据。
        </p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadProducts">
          刷新列表
        </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreateDialog">
          新建产品型号
        </el-button>
      </div>
    </header>

    <section class="summary-grid">
      <article
        v-for="(card, index) in summaryCards"
        :key="card.key"
        class="summary-card"
        :class="`summary-card-${index}`"
      >
        <span class="summary-label">{{ card.label }}</span>
        <strong class="summary-value">{{ card.value }}</strong>
        <span class="summary-description">{{ card.description }}</span>
      </article>
    </section>

    <section class="content-panel">
      <div class="panel-header">
        <div>
          <h2>产品型号列表</h2>
          <p>产品型号描述产品本身，不代表已生产的单台设备。</p>
        </div>
      </div>

      <div class="filter-bar">
        <el-input
          v-model="keyword"
          clearable
          :prefix-icon="Search"
          placeholder="按产品编码搜索"
          class="keyword-input"
          @keyup.enter="searchProducts"
          @clear="searchProducts"
        />

        <el-select
          v-model="categoryFilter"
          clearable
          placeholder="全部分类"
          class="filter-select"
          @change="searchProducts"
        >
          <el-option
            v-for="category in categories"
            :key="category"
            :label="category"
            :value="category"
          />
        </el-select>

        <el-select
          v-model="statusFilter"
          clearable
          placeholder="全部状态"
          class="filter-select"
          @change="searchProducts"
        >
          <el-option label="待启用" value="pending" />
          <el-option label="已启用" value="active" />
          <el-option label="已停用" value="inactive" />
          <el-option label="已归档" value="archived" />
        </el-select>

        <el-button type="primary" @click="searchProducts">搜索</el-button>
        <el-button @click="resetFilters">重置</el-button>
      </div>

      <el-alert
        v-if="loadError"
        :title="loadError"
        type="error"
        show-icon
        :closable="false"
        class="error-alert"
      />

      <el-table
        v-loading="loading"
        :data="products"
        row-key="id"
        class="product-table"
        empty-text="暂无产品型号，请创建第一个产品型号"
      >
        <el-table-column label="产品型号" min-width="240">
          <template #default="{ row }">
            <div class="product-name-cell">
              <div class="product-mark">P</div>
              <div class="product-name-info">
                <button class="product-name-link" @click="openDetail(row)">
                  {{ row.name || "未命名产品" }}
                </button>
                <span class="product-code">{{ row.code }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="版本" width="110">
          <template #default="{ row }">
            <span class="version-label">v{{ row.version || "—" }}</span>
          </template>
        </el-table-column>

        <el-table-column label="产品分类" min-width="140">
          <template #default="{ row }">
            <el-tag effect="plain" type="info">
              {{ row.category || "未分类" }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              :class="['product-status-tag', statusClass(row.status)]"
              effect="light"
              round
            >
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="属性数" width="100" align="center">
          <template #default="{ row }">
            <span class="attribute-count">{{
              row.attributes?.length ?? "—"
            }}</span>
          </template>
        </el-table-column>

        <el-table-column label="更新时间" min-width="175">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.updated_at) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" min-width="250" fixed="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button
                link
                type="primary"
                :icon="View"
                @click="openDetail(row)"
              >
                详情
              </el-button>

              <el-button
                v-if="row.status === 'pending'"
                link
                type="primary"
                :icon="EditPen"
                @click="openEdit(row)"
              >
                编辑
              </el-button>

              <el-button
                v-if="row.status === 'pending' || row.status === 'inactive'"
                link
                type="success"
                :icon="CircleCheck"
                @click="changeStatus(row, 'activate')"
              >
                启用
              </el-button>

              <el-button
                v-if="row.status === 'active'"
                link
                type="warning"
                :icon="CircleClose"
                @click="changeStatus(row, 'deactivate')"
              >
                停用
              </el-button>

              <el-button
                v-if="row.status === 'inactive' || row.status === 'pending'"
                link
                type="danger"
                :icon="Folder"
                @click="changeStatus(row, 'archive')"
              >
                归档
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <footer class="table-footer">
        <span class="result-count">共 {{ total }} 条记录</span>

        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="sizes, prev, pager, next, jumper"
          background
          @current-change="loadProducts"
          @size-change="handlePageSizeChange"
        />
      </footer>
    </section>

    <el-dialog
      v-model="createDialogVisible"
      title="新建产品型号"
      width="760px"
      top="5vh"
      destroy-on-close
      :close-on-click-modal="!submitting"
      @closed="resetCreateForm"
    >
      <div class="dialog-intro">
        <div class="dialog-intro-mark">P</div>
        <div>
          <h3>建立产品定义</h3>
          <p>编码与版本用于识别产品型号；属性定义描述该型号的静态数据结构。</p>
        </div>
      </div>

      <el-alert
        v-if="formError"
        :title="formError"
        type="error"
        show-icon
        :closable="false"
        class="error-alert"
      />

      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="formRules"
        label-position="top"
        @submit.prevent="submitCreate"
      >
        <div class="form-grid">
          <el-form-item label="产品名称" prop="name">
            <el-input
              v-model.trim="createForm.name"
              maxlength="255"
              placeholder="例如：工业网关"
            />
          </el-form-item>

          <el-form-item label="产品编码" prop="code">
            <el-input
              v-model.trim="createForm.code"
              maxlength="100"
              placeholder="例如：GW-100"
            />
          </el-form-item>

          <el-form-item label="版本号" prop="version">
            <el-input
              v-model.trim="createForm.version"
              maxlength="50"
              placeholder="例如：1.0.0"
            />
          </el-form-item>

          <el-form-item label="产品分类" prop="category">
            <el-input
              v-model.trim="createForm.category"
              maxlength="100"
              placeholder="例如：gateway、plc、robot"
            />
          </el-form-item>
        </div>

        <el-form-item label="产品说明">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            maxlength="2000"
            show-word-limit
            placeholder="说明产品用途、适用场景等"
          />
        </el-form-item>

        <div class="attributes-heading">
          <div>
            <h3>属性定义</h3>
            <p>可选。激活产品型号后，属性结构不可修改。</p>
          </div>

          <el-button :icon="Plus" plain @click="addAttributeRow">
            添加属性
          </el-button>
        </div>

        <div v-if="attributeRows.length === 0" class="attributes-empty">
          暂无属性定义。你可以先创建产品型号，再在详情页维护属性结构。
        </div>

        <div
          v-for="(attribute, index) in attributeRows"
          :key="index"
          class="attribute-form-row"
        >
          <div class="attribute-row-header">
            <strong>属性 {{ index + 1 }}</strong>
            <el-button link type="danger" @click="removeAttributeRow(index)">
              删除
            </el-button>
          </div>

          <div class="form-grid attribute-grid">
            <el-form-item label="标识名">
              <el-input
                v-model.trim="attribute.name"
                placeholder="motor_speed"
              />
            </el-form-item>

            <el-form-item label="显示名称">
              <el-input v-model.trim="attribute.label" placeholder="电机转速" />
            </el-form-item>

            <el-form-item label="数据类型">
              <el-select v-model="attribute.data_type" style="width: 100%">
                <el-option
                  v-for="type in supportedTypes"
                  :key="type.value"
                  :label="type.label"
                  :value="type.value"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="单位">
              <el-input v-model.trim="attribute.unit" placeholder="RPM、℃…" />
            </el-form-item>
          </div>

          <el-form-item label="属性说明">
            <el-input
              v-model.trim="attribute.description"
              placeholder="说明属性的含义"
            />
          </el-form-item>

          <el-checkbox v-model="attribute.required">必填属性</el-checkbox>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button
            :disabled="submitting"
            @click="createDialogVisible = false"
          >
            取消
          </el-button>
          <el-button type="primary" :loading="submitting" @click="submitCreate">
            创建产品型号
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.product-page {
  min-width: 0;
  padding: 28px 32px 36px;
  color: var(--color-text, #101828);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
  margin-bottom: 26px;
}

.page-eyebrow {
  margin-bottom: 10px;
  color: #b77b1c;
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 1.8px;
}

.page-header h1 {
  margin: 0;
  color: var(--color-primary, #0e1f33);
  font-size: 28px;
  font-weight: 750;
  letter-spacing: -0.5px;
}

.page-description {
  max-width: 720px;
  margin: 10px 0 0;
  color: var(--color-text-secondary, #667085);
  font-size: 13px;
  line-height: 1.8;
}

.header-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card {
  display: flex;
  min-height: 124px;
  flex-direction: column;
  padding: 20px;
  border: 1px solid var(--color-border, #eaecf0);
  border-radius: 10px;
  background: #fff;
  text-align: left;
}

.summary-label {
  color: #667085;
  font-size: 12px;
}

.summary-value {
  margin-top: 10px;
  color: var(--color-primary, #0e1f33);
  font-size: 28px;
  font-weight: 750;
  line-height: 1.2;
}

.summary-description {
  margin-top: 8px;
  color: #98a2b3;
  font-size: 11px;
}

.summary-card-2 .summary-value {
  color: #15803d;
}

.summary-card-3 .summary-value {
  color: #64748b;
}

.content-panel {
  overflow: hidden;
  border: 1px solid var(--color-border, #eaecf0);
  border-radius: 10px;
  background: #fff;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 22px 24px 18px;
}

.panel-header h2 {
  margin: 0;
  color: var(--color-primary, #0e1f33);
  font-size: 16px;
  font-weight: 750;
}

.panel-header p {
  margin: 7px 0 0;
  color: #667085;
  font-size: 12px;
}

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 0 24px 20px;
}

.keyword-input {
  width: 260px;
}

.filter-select {
  width: 150px;
}

.error-alert {
  margin: 0 24px 16px;
  width: auto;
}

.product-table {
  width: 100%;
  border-top: 1px solid #f0f1f3;
}

/* 产品状态：列表和详情页采用一致的配色规则。 */
.product-status-tag {
  min-width: 74px;
  justify-content: center;
  font-weight: 650;
}

.product-status-tag.status-tag--active {
  --el-tag-bg-color: #ecfdf3;
  --el-tag-border-color: #abefc6;
  --el-tag-text-color: #067647;
}

.product-status-tag.status-tag--pending {
  --el-tag-bg-color: #fffaeb;
  --el-tag-border-color: #fedf89;
  --el-tag-text-color: #b54708;
}

.product-status-tag.status-tag--inactive {
  --el-tag-bg-color: #f2f4f7;
  --el-tag-border-color: #d0d5dd;
  --el-tag-text-color: #475467;
}

.product-status-tag.status-tag--archived {
  --el-tag-bg-color: #fff1f0;
  --el-tag-border-color: #fecdca;
  --el-tag-text-color: #b42318;
  font-weight: 750;
  box-shadow: inset 3px 0 0 #d92d20;
}

.product-status-tag.status-tag--unknown {
  --el-tag-bg-color: #f2f4f7;
  --el-tag-border-color: #d0d5dd;
  --el-tag-text-color: #475467;
}

.product-name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-mark,
.dialog-intro-mark {
  display: grid;
  width: 38px;
  height: 38px;
  flex-shrink: 0;
  place-items: center;
  border-radius: 9px;
  background: #0e1f33;
  color: #f0a030;
  font-size: 15px;
  font-weight: 800;
}

.product-name-info {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 5px;
}

.product-name-link {
  overflow: hidden;
  max-width: 100%;
  padding: 0;
  border: 0;
  background: transparent;
  color: #0e1f33;
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  font-weight: 700;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-name-link:hover {
  color: #b77b1c;
}

.product-code,
.date-text {
  color: #98a2b3;
  font-size: 11px;
}

.version-label {
  color: #344054;
  font-size: 12px;
  font-weight: 650;
}

.attribute-count {
  color: #344054;
  font-size: 12px;
  font-weight: 650;
}

.row-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 2px;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  padding: 18px 24px;
  border-top: 1px solid #f0f1f3;
}

.result-count {
  color: #667085;
  font-size: 12px;
}

.dialog-intro {
  display: flex;
  align-items: center;
  gap: 13px;
  margin-bottom: 22px;
  padding: 16px;
  border: 1px solid #e8ebef;
  border-radius: 9px;
  background: #f8fafc;
}

.dialog-intro h3 {
  margin: 0;
  color: #0e1f33;
  font-size: 14px;
}

.dialog-intro p {
  margin: 6px 0 0;
  color: #667085;
  font-size: 12px;
  line-height: 1.7;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 18px;
}

.attributes-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin: 10px 0 14px;
  padding-top: 18px;
  border-top: 1px solid #eaecf0;
}

.attributes-heading h3 {
  margin: 0;
  color: #0e1f33;
  font-size: 14px;
}

.attributes-heading p {
  margin: 6px 0 0;
  color: #667085;
  font-size: 11px;
}

.attributes-empty {
  padding: 20px 14px;
  border: 1px dashed #d0d5dd;
  border-radius: 8px;
  color: #98a2b3;
  font-size: 12px;
  text-align: center;
}

.attribute-form-row {
  margin-bottom: 14px;
  padding: 15px;
  border: 1px solid #eaecf0;
  border-radius: 8px;
  background: #fff;
}

.attribute-row-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
  color: #344054;
  font-size: 12px;
}

.attribute-grid {
  column-gap: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media (max-width: 1050px) {
  .product-page {
    padding: 22px 20px 30px;
  }

  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .page-header {
    flex-direction: column;
  }
}

@media (max-width: 680px) {
  .product-page {
    padding: 16px 12px 24px;
  }

  .summary-grid {
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }

  .summary-card {
    padding: 14px;
  }

  .filter-bar {
    padding: 0 14px 16px;
  }

  .keyword-input,
  .filter-select {
    width: 100%;
  }

  .panel-header {
    padding: 18px 14px;
  }

  .table-footer {
    flex-direction: column;
    align-items: flex-start;
    padding: 16px 14px;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .header-actions {
    width: 100%;
  }
}
</style>
