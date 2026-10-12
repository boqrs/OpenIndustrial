<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Plus,
  Delete,
  Refresh,
  Check,
} from "@element-plus/icons-vue";

import { getProductModel } from "../../product/api";
import type { ProductModel } from "../../product/types";
import {
  getBOM,
  listMaterials,
  updateBOM,
  releaseBOM,
  obsoleteBOM,
} from "../api";
import type { BOM, BOMItemInput, Material } from "../types";

const route = useRoute();
const router = useRouter();

const bomId = computed(() => Number(route.params.id));
const loading = ref(false);
const saving = ref(false);
const publishing = ref(false);
const bom = ref<BOM | null>(null);
const product = ref<ProductModel | null>(null);
const materials = ref<Material[]>([]);
const loadError = ref("");

interface ItemRow {
  material_id: number | undefined;
  quantity: number;
  unit: string;
  sequence: number;
  operation_code: string;
  is_optional: boolean;
  description: string;
}

const form = reactive({
  description: "",
});

const rows = ref<ItemRow[]>([]);

const editable = computed(() => bom.value?.status === "draft");

const statusText = computed(() => {
  const labels: Record<string, string> = {
    draft: "草稿",
    released: "已发布",
    obsolete: "已作废",
  };
  return labels[bom.value?.status ?? ""] ?? bom.value?.status ?? "未知";
});

const statusType = computed(() => {
  if (bom.value?.status === "released") return "success";
  if (bom.value?.status === "obsolete") return "danger";
  return "info";
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

function materialLabel(id?: number): string {
  const item = materials.value.find((material) => material.id === id);
  return item ? `${item.code} — ${item.name}` : "";
}

function addItem() {
  rows.value.push({
    material_id: undefined,
    quantity: 1,
    unit: "",
    sequence: rows.value.length + 1,
    operation_code: "",
    is_optional: false,
    description: "",
  });
}

function removeItem(index: number) {
  rows.value.splice(index, 1);
  rows.value.forEach((item, position) => {
    item.sequence = position + 1;
  });
}

function onMaterialChange(row: ItemRow) {
  const material = materials.value.find((item) => item.id === row.material_id);
  if (material) row.unit = material.unit;
}

async function loadDetail() {
  if (!Number.isInteger(bomId.value) || bomId.value <= 0) {
    loadError.value = "BOM ID 无效";
    return;
  }

  loading.value = true;
  loadError.value = "";

  try {
    const [detail, materialList] = await Promise.all([
      getBOM(bomId.value),
      listMaterials(),
    ]);

    bom.value = detail;
    materials.value = materialList;
    form.description = detail.description ?? "";

    rows.value = (detail.items ?? []).map((item) => ({
      material_id: item.material_id,
      quantity: item.quantity,
      unit: item.unit,
      sequence: item.sequence,
      operation_code: item.operation_code ?? "",
      is_optional: item.is_optional,
      description: item.description ?? "",
    }));

    try {
      product.value = await getProductModel(detail.product_id);
    } catch {
      product.value = null;
    }
  } catch (error: unknown) {
    loadError.value =
      error instanceof Error ? error.message : "BOM 详情加载失败";
  } finally {
    loading.value = false;
  }
}

function buildItems(): BOMItemInput[] | string {
  const result: BOMItemInput[] = [];
  const seen = new Set<number>();

  for (let index = 0; index < rows.value.length; index += 1) {
    const row = rows.value[index];

    if (!row.material_id) {
      return `第 ${index + 1} 行请选择物料`;
    }

    if (seen.has(row.material_id)) {
      return `物料「${materialLabel(row.material_id)}」重复，请合并为一行`;
    }
    seen.add(row.material_id);

    if (!Number.isFinite(row.quantity) || row.quantity <= 0) {
      return `第 ${index + 1} 行用量必须大于 0`;
    }

    if (!row.unit.trim()) {
      return `第 ${index + 1} 行请填写单位`;
    }

    result.push({
      material_id: row.material_id,
      quantity: row.quantity,
      unit: row.unit.trim(),
      sequence: index + 1,
      operation_code: row.operation_code.trim(),
      is_optional: row.is_optional,
      description: row.description.trim(),
    });
  }

  return result;
}

async function saveDraft() {
  if (!bom.value || !editable.value || saving.value) return;

  const items = buildItems();
  if (typeof items === "string") {
    ElMessage.warning(items);
    return;
  }

  saving.value = true;
  try {
    await updateBOM(bom.value.id, {
      description: form.description.trim(),
      items,
    });
    ElMessage.success("BOM 草稿已保存");
    await loadDetail();
  } catch (error: unknown) {
    ElMessage.error(error instanceof Error ? error.message : "保存失败");
  } finally {
    saving.value = false;
  }
}

async function publishBOM() {
  if (!bom.value || !editable.value) return;

  if (rows.value.length === 0) {
    ElMessage.warning("请至少添加一条物料明细后再发布");
    return;
  }

  try {
    await ElMessageBox.confirm(
      "发布后该 BOM 将不可编辑。请确认物料及用量均已核对。",
      "发布 BOM",
      {
        type: "warning",
        confirmButtonText: "确认发布",
        cancelButtonText: "继续编辑",
      },
    );
  } catch {
    return;
  }

  // 发布前先保存明细，避免发布旧数据。
  const items = buildItems();
  if (typeof items === "string") {
    ElMessage.warning(items);
    return;
  }

  publishing.value = true;
  try {
    await updateBOM(bom.value.id, {
      description: form.description.trim(),
      items,
    });
    await releaseBOM(bom.value.id);
    ElMessage.success("BOM 已发布");
    await loadDetail();
  } catch (error: unknown) {
    ElMessage.error(error instanceof Error ? error.message : "发布失败");
  } finally {
    publishing.value = false;
  }
}

async function makeObsolete() {
  if (!bom.value || bom.value.status !== "released") return;

  try {
    await ElMessageBox.confirm("作废后该 BOM 不可恢复，请确认。", "作废 BOM", {
      type: "warning",
      confirmButtonText: "确认作废",
      cancelButtonText: "取消",
    });

    await obsoleteBOM(bom.value.id);
    ElMessage.success("BOM 已作废");
    await loadDetail();
  } catch (error: unknown) {
    if (error === "cancel" || error === "close") return;
    if (error instanceof Error) ElMessage.error(error.message);
  }
}

onMounted(() => {
  void loadDetail();
});
</script>

<template>
  <div class="bom-detail-page" v-loading="loading">
    <header class="page-header">
      <div>
        <el-button text :icon="ArrowLeft" @click="router.push('/mes/boms')">
          返回 BOM 列表
        </el-button>
        <div class="eyebrow">MES / BOM DETAIL</div>
        <h1>{{ bom?.bom_no || "BOM 详情" }}</h1>
        <p class="description">查看物料组成与工艺信息，并维护 BOM 版本。</p>
      </div>

      <div class="header-actions">
        <el-button :icon="Refresh" :loading="loading" @click="loadDetail">
          刷新
        </el-button>
        <el-button
          v-if="editable"
          :icon="Check"
          type="primary"
          :loading="saving"
          @click="saveDraft"
        >
          保存草稿
        </el-button>
        <el-button
          v-if="editable"
          type="success"
          :loading="publishing"
          @click="publishBOM"
        >
          发布 BOM
        </el-button>
        <el-button
          v-if="bom?.status === 'released'"
          type="danger"
          @click="makeObsolete"
        >
          作废
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

    <template v-if="bom">
      <section class="overview-grid">
        <article class="overview-card">
          <span>状态</span>
          <el-tag :type="statusType" effect="light">{{ statusText }}</el-tag>
        </article>
        <article class="overview-card">
          <span>BOM 版本</span>
          <strong>V{{ bom.version }}</strong>
        </article>
        <article class="overview-card">
          <span>物料明细</span>
          <strong>{{ rows.length }} <small>行</small></strong>
        </article>
        <article class="overview-card">
          <span>最后更新</span>
          <strong class="date-value">{{ formatDate(bom.updated_at) }}</strong>
        </article>
      </section>

      <section class="panel">
        <div class="section-heading">
          <div>
            <h2>基本信息</h2>
            <p>产品型号与 BOM 编号用于识别物料清单版本。</p>
          </div>
        </div>

        <div class="basic-grid">
          <div class="info-field">
            <label>产品型号</label>
            <strong>{{ product?.name || `产品 ID ${bom.product_id}` }}</strong>
            <small>{{ product?.code || "—" }}</small>
          </div>
          <div class="info-field">
            <label>BOM 编号</label>
            <strong>{{ bom.bom_no }}</strong>
          </div>
          <div class="info-field">
            <label>创建时间</label>
            <strong>{{ formatDate(bom.created_at) }}</strong>
          </div>
        </div>

        <el-form label-position="top" class="description-form">
          <el-form-item label="BOM 描述">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="2"
              maxlength="1000"
              show-word-limit
              :disabled="!editable"
              placeholder="填写 BOM 的用途、变更或工艺说明"
            />
          </el-form-item>
        </el-form>
      </section>

      <section class="panel items-panel">
        <div class="section-heading">
          <div>
            <h2>物料明细</h2>
            <p>每种物料单独占一行，数量必须大于 0。</p>
          </div>
          <el-button
            v-if="editable"
            type="primary"
            plain
            :icon="Plus"
            @click="addItem"
          >
            添加物料
          </el-button>
        </div>

        <el-table :data="rows" row-key="sequence" style="width: 100%">
          <el-table-column label="#" width="65" align="center">
            <template #default="{ $index }">
              {{ $index + 1 }}
            </template>
          </el-table-column>

          <el-table-column label="物料" min-width="250">
            <template #default="{ row }">
              <el-select
                v-model="row.material_id"
                filterable
                clearable
                :disabled="!editable"
                placeholder="选择物料"
                style="width: 100%"
                @change="onMaterialChange(row)"
              >
                <el-option
                  v-for="material in materials"
                  :key="material.id"
                  :label="`${material.code} — ${material.name}`"
                  :value="material.id"
                />
              </el-select>
            </template>
          </el-table-column>

          <el-table-column label="用量" width="140">
            <template #default="{ row }">
              <el-input-number
                v-model="row.quantity"
                :min="0.000001"
                :precision="6"
                :controls="false"
                :disabled="!editable"
                style="width: 100%"
              />
            </template>
          </el-table-column>

          <el-table-column label="单位" width="100">
            <template #default="{ row }">
              <el-input
                v-model="row.unit"
                :disabled="!editable"
                placeholder="单位"
              />
            </template>
          </el-table-column>

          <el-table-column label="工序编码" width="150">
            <template #default="{ row }">
              <el-input
                v-model="row.operation_code"
                :disabled="!editable"
                placeholder="可选"
              />
            </template>
          </el-table-column>

          <el-table-column label="可选物料" width="105" align="center">
            <template #default="{ row }">
              <el-checkbox v-model="row.is_optional" :disabled="!editable" />
            </template>
          </el-table-column>

          <el-table-column label="说明" min-width="170">
            <template #default="{ row }">
              <el-input
                v-model="row.description"
                :disabled="!editable"
                placeholder="可选"
              />
            </template>
          </el-table-column>

          <el-table-column
            v-if="editable"
            label="操作"
            width="80"
            fixed="right"
          >
            <template #default="{ $index }">
              <el-button
                link
                type="danger"
                :icon="Delete"
                @click="removeItem($index)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <div v-if="materials.length === 0" class="empty-note">
          当前没有可选择的物料，请先完成物料档案维护。
        </div>

        <div v-if="editable" class="bottom-actions">
          <el-button :icon="Plus" @click="addItem">添加一行</el-button>
          <span>草稿可随时修改；发布前必须至少包含一条有效物料明细。</span>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.bom-detail-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 26px 0 48px;
  color: #17283b;
}

.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.page-header > div:first-child {
  min-width: 0;
}

.eyebrow {
  margin-top: 12px;
  color: #9b6a22;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.5px;
}

h1 {
  margin: 8px 0 0;
  color: #0e1f33;
  font-size: 27px;
  overflow-wrap: anywhere;
}

.description {
  margin: 9px 0 0;
  color: #738195;
  font-size: 13px;
}

.header-actions {
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 20px;
}

.overview-card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
  min-height: 92px;
  padding: 18px 20px;
  border: 1px solid #e6eaf0;
  border-radius: 10px;
  background: #fff;
}

.overview-card > span {
  color: #738195;
  font-size: 12px;
}

.overview-card > strong {
  color: #0e1f33;
  font-size: 20px;
}

.overview-card small {
  font-size: 12px;
  font-weight: 400;
}

.overview-card .date-value {
  font-size: 12px;
}

.panel {
  margin-bottom: 20px;
  padding: 22px;
  border: 1px solid #e6eaf0;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 4px 16px rgb(14 31 51 / 3%);
}

.section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.section-heading h2 {
  margin: 0;
  color: #0e1f33;
  font-size: 16px;
}

.section-heading p {
  margin: 7px 0 0;
  color: #8995a4;
  font-size: 12px;
}

.basic-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
  margin-bottom: 22px;
}

.info-field {
  display: flex;
  flex-direction: column;
  gap: 7px;
  min-width: 0;
}

.info-field label {
  color: #8995a4;
  font-size: 11px;
}

.info-field strong {
  color: #17283b;
  font-size: 13px;
  overflow-wrap: anywhere;
}

.info-field small {
  color: #8c99a8;
  font-size: 11px;
}

.description-form {
  max-width: 850px;
}

.items-panel {
  overflow: hidden;
}

.empty-note {
  padding: 20px 0 4px;
  color: #8995a4;
  font-size: 12px;
  text-align: center;
}

.bottom-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
  margin-top: 18px;
}

.bottom-actions span {
  color: #8995a4;
  font-size: 12px;
}

@media (max-width: 900px) {
  .bom-detail-page {
    width: calc(100% - 32px);
  }

  .page-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .overview-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 560px) {
  .bom-detail-page {
    width: calc(100% - 24px);
    padding-top: 18px;
  }

  .overview-grid {
    gap: 10px;
  }

  .overview-card {
    padding: 14px;
  }

  .panel {
    padding: 14px;
  }

  .basic-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .header-actions {
    justify-content: flex-start;
  }
}
</style>
