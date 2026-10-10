<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Refresh,
  Check,
  Plus,
  Delete,
} from "@element-plus/icons-vue";

import {
  activateProductModel,
  archiveProductModel,
  deactivateProductModel,
  getAttributeDefinitions,
  getProductModel,
  updateAttributeDefinitions,
  updateProductModel,
} from "../api";
import type {
  AttributeDataType,
  AttributeDefinition,
  AttributeDefinitionInput,
  ProductDetail,
} from "../types";

const route = useRoute();
const router = useRouter();

const productId = computed(() => Number(route.params.id));
const loading = ref(false);
const saving = ref(false);
const loadError = ref("");
const editMode = ref(route.query.edit === "1");

const product = ref<ProductDetail | null>(null);
const definitions = ref<AttributeDefinition[]>([]);

const form = reactive({
  name: "",
  category: "",
  description: "",
});

interface AttributeRow {
  name: string;
  label: string;
  description: string;
  data_type: AttributeDataType;
  unit: string;
  required: boolean;
}

const attributeRows = ref<AttributeRow[]>([]);

const supportedTypes: { label: string; value: AttributeDataType }[] = [
  { label: "字符串", value: "string" },
  { label: "长文本", value: "text" },
  { label: "整数", value: "integer" },
  { label: "浮点数", value: "float" },
  { label: "布尔值", value: "boolean" },
  { label: "日期时间", value: "datetime" },
  { label: "JSON", value: "json" },
];

const canEditSchema = computed(() => product.value?.status === "pending");

const statusLabel = computed(() => {
  const labels: Record<string, string> = {
    pending: "待启用",
    active: "已启用",
    inactive: "已停用",
    archived: "已归档",
  };

  return labels[product.value?.status ?? ""] ?? product.value?.status ?? "未知";
});

const statusTagType = computed(() => {
  switch (product.value?.status) {
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
});

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

function copyProductToForm(data: ProductDetail) {
  form.name = data.name ?? "";
  form.category = data.category ?? "";
  form.description = data.description ?? "";
}

function copyDefinitionsToRows(items: AttributeDefinition[]) {
  attributeRows.value = items.map((item) => ({
    name: item.name,
    label: item.label || item.name,
    description: item.description || "",
    data_type: normalizeType(item.data_type),
    unit: item.unit || "",
    required: Boolean(item.required),
  }));
}

function normalizeType(value: string): AttributeDataType {
  const supported = supportedTypes.some((item) => item.value === value);
  return supported ? (value as AttributeDataType) : "string";
}

async function loadProduct() {
  if (!Number.isInteger(productId.value) || productId.value <= 0) {
    loadError.value = "产品型号 ID 无效";
    return;
  }

  loading.value = true;
  loadError.value = "";

  try {
    const [detail, attrs] = await Promise.all([
      getProductModel(productId.value),
      getAttributeDefinitions(productId.value),
    ]);

    product.value = detail;
    definitions.value = attrs;
    copyProductToForm(detail);
    copyDefinitionsToRows(attrs);

    if (route.query.edit === "1" && detail.status !== "pending") {
      editMode.value = false;
      ElMessage.info("只有待启用的产品型号可以编辑");
      void router.replace(`/mes/products/${productId.value}`);
    }
  } catch (error: unknown) {
    loadError.value =
      error instanceof Error ? error.message : "产品型号详情加载失败";
  } finally {
    loading.value = false;
  }
}

function addAttribute() {
  attributeRows.value.push({
    name: "",
    label: "",
    description: "",
    data_type: "string",
    unit: "",
    required: false,
  });
}

function removeAttribute(index: number) {
  attributeRows.value.splice(index, 1);
}

function buildAttributePayload():
  | Record<string, AttributeDefinitionInput>
  | string {
  const result: Record<string, AttributeDefinitionInput> = {};

  for (const row of attributeRows.value) {
    const name = row.name.trim();

    if (!name) return "请填写所有属性的标识名";

    if (!/^[a-zA-Z][a-zA-Z0-9_]*$/.test(name)) {
      return `属性标识“${name}”格式无效`;
    }

    if (Object.prototype.hasOwnProperty.call(result, name)) {
      return `属性标识“${name}”重复`;
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

async function saveProduct() {
  if (!product.value || saving.value) return;

  if (!form.name.trim()) {
    ElMessage.warning("产品名称不能为空");
    return;
  }

  if (!form.category.trim()) {
    ElMessage.warning("产品分类不能为空");
    return;
  }

  const attributePayload = buildAttributePayload();

  if (typeof attributePayload === "string") {
    ElMessage.warning(attributePayload);
    return;
  }

  saving.value = true;

  try {
    // 先更新基础信息，再替换属性定义。
    // 两个请求不是一个数据库事务；若第二步失败，需要重试属性保存。
    await updateProductModel(productId.value, {
      name: form.name.trim(),
      category: form.category.trim(),
      description: form.description.trim(),
    });

    if (canEditSchema.value) {
      await updateAttributeDefinitions(productId.value, attributePayload);
    }

    ElMessage.success("产品型号保存成功");
    editMode.value = false;
    await router.replace(`/mes/products/${productId.value}`);
    await loadProduct();
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : "保存产品型号失败",
    );
    await loadProduct();
  } finally {
    saving.value = false;
  }
}

function cancelEdit() {
  editMode.value = false;
  copyProductToForm(product.value!);
  copyDefinitionsToRows(definitions.value);
  void router.replace(`/mes/products/${productId.value}`);
}

async function runStatusAction(action: "activate" | "deactivate" | "archive") {
  if (!product.value) return;

  const actionText = {
    activate: "启用",
    deactivate: "停用",
    archive: "归档",
  }[action];

  try {
    await ElMessageBox.confirm(
      `确认${actionText}产品型号“${product.value.name}（${product.value.code}）”吗？`,
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
    if (action === "activate") await activateProductModel(productId.value);
    if (action === "deactivate") await deactivateProductModel(productId.value);
    if (action === "archive") await archiveProductModel(productId.value);

    ElMessage.success(`已提交${actionText}操作`);
    await loadProduct();
  } catch (error: unknown) {
    ElMessage.error(
      error instanceof Error ? error.message : `${actionText}失败`,
    );
  }
}

onMounted(() => {
  void loadProduct();
});
</script>

<template>
  <div class="product-detail-page" v-loading="loading">
    <header class="page-header">
      <div>
        <button class="back-link" @click="router.push('/mes/products')">
          <el-icon><ArrowLeft /></el-icon>
          返回产品型号列表
        </button>

        <div class="page-eyebrow">MES / PRODUCT MODEL / DETAIL</div>
        <h1>{{ product?.name || "产品型号详情" }}</h1>
        <p class="page-description">
          查看产品定义、维护静态属性结构，并管理产品型号生命周期。
        </p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadProduct">
          刷新
        </el-button>

        <template v-if="editMode">
          <el-button :disabled="saving" @click="cancelEdit">取消</el-button>
          <el-button
            type="primary"
            :icon="Check"
            :loading="saving"
            @click="saveProduct"
          >
            保存修改
          </el-button>
        </template>

        <el-button
          v-else-if="canEditSchema"
          type="primary"
          :icon="Check"
          @click="editMode = true"
        >
          编辑产品
        </el-button>
      </div>
    </header>

    <el-alert
      v-if="loadError"
      :title="loadError"
      type="error"
      show-icon
      :closable="false"
      class="error-alert"
    />

    <template v-if="product">
      <section class="hero-card">
        <div class="hero-product-mark">P</div>

        <div class="hero-main">
          <div class="hero-title-row">
            <h2>{{ product.name }}</h2>
            <el-tag :type="statusTagType" effect="light" round>
              {{ statusLabel }}
            </el-tag>
          </div>

          <div class="hero-identifiers">
            <span
              >产品编码：<strong>{{ product.code }}</strong></span
            >
            <span
              >版本：<strong>v{{ product.version }}</strong></span
            >
            <span
              >分类：<strong>{{ product.category }}</strong></span
            >
          </div>

          <p class="hero-description">
            {{ product.description || "暂无产品说明" }}
          </p>
        </div>

        <div class="hero-actions">
          <el-button
            v-if="product.status === 'pending' || product.status === 'inactive'"
            type="success"
            @click="runStatusAction('activate')"
          >
            启用
          </el-button>

          <el-button
            v-if="product.status === 'active'"
            type="warning"
            plain
            @click="runStatusAction('deactivate')"
          >
            停用
          </el-button>

          <el-button
            v-if="product.status === 'pending' || product.status === 'inactive'"
            type="danger"
            plain
            @click="runStatusAction('archive')"
          >
            归档
          </el-button>
        </div>
      </section>

      <section class="detail-grid">
        <article class="content-panel basic-panel">
          <div class="panel-heading">
            <div>
              <h2>基本信息</h2>
              <p>产品型号的身份信息和说明</p>
            </div>
          </div>

          <el-form v-if="editMode" label-position="top" class="basic-form">
            <el-form-item label="产品名称">
              <el-input v-model.trim="form.name" maxlength="255" />
            </el-form-item>

            <el-form-item label="产品编码">
              <el-input :model-value="product.code" disabled />
            </el-form-item>

            <el-form-item label="版本号">
              <el-input :model-value="product.version" disabled />
            </el-form-item>

            <el-form-item label="产品分类">
              <el-input v-model.trim="form.category" maxlength="100" />
            </el-form-item>

            <el-form-item label="产品说明">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="4"
                maxlength="2000"
                show-word-limit
              />
            </el-form-item>
          </el-form>

          <dl v-else class="info-list">
            <div>
              <dt>产品名称</dt>
              <dd>{{ product.name || "—" }}</dd>
            </div>
            <div>
              <dt>产品编码</dt>
              <dd>{{ product.code }}</dd>
            </div>
            <div>
              <dt>版本号</dt>
              <dd>{{ product.version }}</dd>
            </div>
            <div>
              <dt>产品分类</dt>
              <dd>{{ product.category }}</dd>
            </div>
            <div>
              <dt>资源 ID</dt>
              <dd>{{ product.resource_id }}</dd>
            </div>
            <div>
              <dt>创建时间</dt>
              <dd>{{ formatDate(product.created_at) }}</dd>
            </div>
            <div>
              <dt>更新时间</dt>
              <dd>{{ formatDate(product.updated_at) }}</dd>
            </div>
            <div class="info-description">
              <dt>产品说明</dt>
              <dd>{{ product.description || "暂无说明" }}</dd>
            </div>
          </dl>
        </article>

        <article class="content-panel attributes-panel">
          <div class="panel-heading attributes-heading">
            <div>
              <h2>属性定义</h2>
              <p>{{ definitions.length }} 项属性 · 描述产品的数据结构</p>
            </div>

            <el-button
              v-if="editMode && canEditSchema"
              :icon="Plus"
              plain
              @click="addAttribute"
            >
              添加属性
            </el-button>
          </div>

          <el-alert
            v-if="!canEditSchema"
            title="当前状态下属性结构不可修改"
            description="Cloud 仅允许待启用产品型号修改属性定义。"
            type="info"
            :closable="false"
            class="schema-notice"
          />

          <div v-if="editMode && canEditSchema">
            <div v-if="attributeRows.length === 0" class="attributes-empty">
              暂无属性定义，可点击“添加属性”创建。
            </div>

            <div
              v-for="(attribute, index) in attributeRows"
              :key="index"
              class="attribute-editor"
            >
              <div class="attribute-editor-header">
                <strong>属性 {{ index + 1 }}</strong>
                <el-button
                  link
                  type="danger"
                  :icon="Delete"
                  @click="removeAttribute(index)"
                >
                  删除
                </el-button>
              </div>

              <div class="attribute-fields">
                <el-form-item label="标识名">
                  <el-input
                    v-model.trim="attribute.name"
                    placeholder="motor_speed"
                  />
                </el-form-item>

                <el-form-item label="显示名称">
                  <el-input
                    v-model.trim="attribute.label"
                    placeholder="电机转速"
                  />
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
                  <el-input
                    v-model.trim="attribute.unit"
                    placeholder="RPM、℃…"
                  />
                </el-form-item>
              </div>

              <el-form-item label="属性说明">
                <el-input v-model.trim="attribute.description" />
              </el-form-item>

              <el-checkbox v-model="attribute.required">必填属性</el-checkbox>
            </div>
          </div>

          <div v-else-if="definitions.length" class="definition-list">
            <article
              v-for="item in definitions"
              :key="item.id || item.name"
              class="definition-item"
            >
              <div class="definition-main">
                <div class="definition-title">
                  <strong>{{ item.label || item.name }}</strong>
                  <el-tag v-if="item.required" size="small" type="warning">
                    必填
                  </el-tag>
                </div>
                <code>{{ item.name }}</code>
                <p v-if="item.description">
                  {{ item.description }}
                </p>
              </div>

              <div class="definition-meta">
                <el-tag effect="plain">{{ item.data_type }}</el-tag>
                <span v-if="item.unit">{{ item.unit }}</span>
              </div>
            </article>
          </div>

          <div v-else class="attributes-empty">暂无属性定义。</div>

          <div v-if="editMode" class="inline-save">
            <el-button :disabled="saving" @click="cancelEdit">取消</el-button>
            <el-button type="primary" :loading="saving" @click="saveProduct">
              保存产品及属性
            </el-button>
          </div>
        </article>
      </section>
    </template>
  </div>
</template>

<style scoped>
.product-detail-page {
  min-width: 0;
  padding: 28px 32px 36px;
  color: var(--color-text, #101828);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
  margin-bottom: 24px;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 20px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #667085;
  cursor: pointer;
  font-size: 12px;
}

.back-link:hover {
  color: #b77b1c;
}

.page-eyebrow {
  margin-bottom: 9px;
  color: #b77b1c;
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 1.6px;
}

.page-header h1 {
  margin: 0;
  color: #0e1f33;
  font-size: 28px;
  font-weight: 750;
}

.page-description {
  margin: 9px 0 0;
  color: #667085;
  font-size: 13px;
  line-height: 1.8;
}

.header-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 9px;
}

.error-alert {
  margin-bottom: 18px;
}

.hero-card {
  display: flex;
  align-items: flex-start;
  gap: 20px;
  margin-bottom: 22px;
  padding: 25px;
  border: 1px solid #263c53;
  border-radius: 11px;
  background: #0e1f33;
  color: #fff;
}

.hero-product-mark {
  display: grid;
  width: 54px;
  height: 54px;
  flex-shrink: 0;
  place-items: center;
  border: 1px solid #6d5a3b;
  border-radius: 11px;
  background: #16293f;
  color: #f0a030;
  font-size: 21px;
  font-weight: 800;
}

.hero-main {
  flex: 1;
  min-width: 0;
}

.hero-title-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.hero-title-row h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 750;
}

.hero-identifiers {
  display: flex;
  flex-wrap: wrap;
  gap: 18px;
  margin-top: 14px;
  color: #b9c4d0;
  font-size: 12px;
}

.hero-identifiers strong {
  color: #fff;
  font-weight: 650;
}

.hero-description {
  margin: 15px 0 0;
  color: #c2cbd5;
  font-size: 12px;
  line-height: 1.8;
  white-space: pre-wrap;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;
}

.detail-grid {
  display: grid;
  grid-template-columns: minmax(280px, 0.85fr) minmax(0, 1.4fr);
  align-items: start;
  gap: 22px;
}

.content-panel {
  min-width: 0;
  overflow: hidden;
  border: 1px solid #eaecf0;
  border-radius: 10px;
  background: #fff;
}

.panel-heading {
  padding: 21px 22px 18px;
  border-bottom: 1px solid #f0f1f3;
}

.panel-heading h2 {
  margin: 0;
  color: #0e1f33;
  font-size: 15px;
  font-weight: 750;
}

.panel-heading p {
  margin: 7px 0 0;
  color: #667085;
  font-size: 11px;
}

.info-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  margin: 0;
  padding: 7px 22px 20px;
}

.info-list > div {
  padding: 14px 8px 14px 0;
  border-bottom: 1px solid #f2f4f7;
}

.info-list dt {
  margin-bottom: 7px;
  color: #98a2b3;
  font-size: 11px;
}

.info-list dd {
  overflow-wrap: anywhere;
  margin: 0;
  color: #344054;
  font-size: 12px;
  line-height: 1.7;
}

.info-description {
  grid-column: 1 / -1;
}

.basic-form {
  padding: 20px 22px 8px;
}

.attributes-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.schema-notice {
  margin: 16px 18px 0;
  width: auto;
}

.definition-list {
  padding: 0 22px;
}

.definition-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 18px;
  padding: 18px 0;
  border-bottom: 1px solid #f0f1f3;
}

.definition-item:last-child {
  border-bottom: 0;
}

.definition-main {
  min-width: 0;
}

.definition-title {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 9px;
  color: #344054;
  font-size: 13px;
}

.definition-main code {
  display: inline-block;
  margin-top: 7px;
  color: #8a5b16;
  font-size: 11px;
}

.definition-main p {
  margin: 7px 0 0;
  color: #667085;
  font-size: 11px;
  line-height: 1.7;
}

.definition-meta {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 8px;
  color: #667085;
  font-size: 11px;
}

.attributes-empty {
  margin: 18px;
  padding: 25px 12px;
  border: 1px dashed #d0d5dd;
  border-radius: 8px;
  color: #98a2b3;
  font-size: 12px;
  text-align: center;
}

.attribute-editor {
  margin: 16px 18px;
  padding: 15px;
  border: 1px solid #eaecf0;
  border-radius: 8px;
}

.attribute-editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
  color: #344054;
  font-size: 12px;
}

.attribute-fields {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 14px;
}

.inline-save {
  display: flex;
  justify-content: flex-end;
  gap: 9px;
  padding: 18px;
  border-top: 1px solid #f0f1f3;
}

@media (max-width: 1050px) {
  .product-detail-page {
    padding: 22px 20px 30px;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .hero-card {
    flex-wrap: wrap;
  }

  .hero-actions {
    width: 100%;
    justify-content: flex-start;
  }
}

@media (max-width: 680px) {
  .product-detail-page {
    padding: 16px 12px 24px;
  }

  .page-header {
    flex-direction: column;
  }

  .hero-card {
    padding: 18px;
  }

  .hero-identifiers {
    flex-direction: column;
    gap: 7px;
  }

  .info-list {
    grid-template-columns: 1fr;
    padding-right: 18px;
    padding-left: 18px;
  }

  .info-description {
    grid-column: auto;
  }

  .attribute-fields {
    grid-template-columns: 1fr;
  }

  .definition-item {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
