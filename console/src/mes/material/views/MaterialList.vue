<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  createMaterial,
  deleteMaterial,
  listMaterials,
  updateMaterial,
} from "../api";
import type { CreateMaterialRequest, Material, MaterialType } from "../types";

const loading = ref(false);
const saving = ref(false);
const materials = ref<Material[]>([]);
const keyword = ref("");
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);

const emptyForm = (): CreateMaterialRequest => ({
  code: "",
  name: "",
  material_type: "raw",
  unit: "",
  description: "",
});

const form = reactive<CreateMaterialRequest>(emptyForm());

const rules = {
  code: [{ required: true, message: "请输入物料编码", trigger: "blur" }],
  name: [{ required: true, message: "请输入物料名称", trigger: "blur" }],
  material_type: [
    { required: true, message: "请选择物料类型", trigger: "change" },
  ],
  unit: [{ required: true, message: "请输入计量单位", trigger: "blur" }],
};

const formRef = ref();

const filteredMaterials = computed(() => {
  const value = keyword.value.trim().toLowerCase();

  if (!value) return materials.value;

  return materials.value.filter((item) =>
    [item.code, item.name, item.unit, item.description].some((field) =>
      field?.toLowerCase().includes(value),
    ),
  );
});

function materialTypeLabel(type: string) {
  const labels: Record<string, string> = {
    raw: "原材料",
    component: "零部件",
    consumable: "耗材",
    packaging: "包装材料",
  };

  return labels[type] ?? type;
}

async function loadMaterials() {
  loading.value = true;

  try {
    materials.value = await listMaterials({ offset: 0, limit: 100 });
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "加载物料列表失败",
    );
  } finally {
    loading.value = false;
  }
}

function openCreateDialog() {
  editingId.value = null;
  Object.assign(form, emptyForm());
  dialogVisible.value = true;
}

function openEditDialog(item: Material) {
  editingId.value = item.id;

  Object.assign(form, {
    code: item.code,
    name: item.name,
    material_type: item.material_type as MaterialType,
    unit: item.unit,
    description: item.description ?? "",
  });

  dialogVisible.value = true;
}

async function submitForm() {
  if (!formRef.value) return;

  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;

  try {
    const payload: CreateMaterialRequest = {
      code: form.code.trim(),
      name: form.name.trim(),
      material_type: form.material_type,
      unit: form.unit.trim(),
      description: form.description?.trim() ?? "",
    };

    if (editingId.value === null) {
      await createMaterial(payload);
      ElMessage.success("物料创建成功");
    } else {
      await updateMaterial(editingId.value, payload);
      ElMessage.success("物料更新成功");
    }

    dialogVisible.value = false;
    await loadMaterials();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "保存物料失败");
  } finally {
    saving.value = false;
  }
}

async function handleDelete(item: Material) {
  try {
    await ElMessageBox.confirm(
      `确定删除物料“${item.code} - ${item.name}”吗？`,
      "删除物料",
      {
        type: "warning",
        confirmButtonText: "确认删除",
        cancelButtonText: "取消",
      },
    );

    await deleteMaterial(item.id);
    ElMessage.success("物料已删除");
    await loadMaterials();
  } catch (error) {
    if (error === "cancel" || error === "close") return;

    ElMessage.error(error instanceof Error ? error.message : "删除物料失败");
  }
}

onMounted(loadMaterials);
</script>

<template>
  <div class="material-page">
    <div class="page-header">
      <div>
        <div class="eyebrow">MES · MASTER DATA</div>
        <h1>物料管理</h1>
        <p>维护制造业务使用的物料档案，为 BOM 物料清单提供基础数据。</p>
      </div>

      <el-button type="primary" @click="openCreateDialog"> 新建物料 </el-button>
    </div>

    <el-card shadow="never" class="list-card">
      <div class="toolbar">
        <el-input
          v-model="keyword"
          clearable
          placeholder="搜索物料编码、名称、单位或描述"
          style="max-width: 360px"
        />

        <el-button :loading="loading" @click="loadMaterials"> 刷新 </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="filteredMaterials"
        row-key="id"
        stripe
        empty-text="暂无物料，请先新建物料"
      >
        <el-table-column prop="code" label="物料编码" min-width="150" />
        <el-table-column prop="name" label="物料名称" min-width="180" />

        <el-table-column label="物料类型" width="130">
          <template #default="{ row }">
            {{ materialTypeLabel(row.material_type) }}
          </template>
        </el-table-column>

        <el-table-column prop="unit" label="基本单位" width="110" />

        <el-table-column
          prop="description"
          label="描述"
          min-width="200"
          show-overflow-tooltip
        />

        <el-table-column label="创建时间" min-width="180">
          <template #default="{ row }">
            {{
              row.created_at ? new Date(row.created_at).toLocaleString() : "-"
            }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="list-footer">
        当前显示 {{ filteredMaterials.length }} 条物料
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="editingId === null ? '新建物料' : '编辑物料'"
      width="560px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="物料编码" prop="code">
          <el-input
            v-model="form.code"
            placeholder="请输入唯一的物料编码"
            maxlength="100"
          />
        </el-form-item>

        <el-form-item label="物料名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入物料名称"
            maxlength="255"
          />
        </el-form-item>

        <el-form-item label="物料类型" prop="material_type">
          <el-select v-model="form.material_type" style="width: 100%">
            <el-option label="原材料" value="raw" />
            <el-option label="零部件" value="component" />
            <el-option label="耗材" value="consumable" />
            <el-option label="包装材料" value="packaging" />
          </el-select>
        </el-form-item>

        <el-form-item label="基本单位" prop="unit">
          <el-input
            v-model="form.unit"
            placeholder="例如：pcs、kg、m"
            maxlength="32"
          />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="可选"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitForm">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.material-page {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding: 30px 0 48px;
  color: #17283b;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.eyebrow {
  margin-bottom: 10px;
  color: #9a743b;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1.8px;
}

.page-header h1 {
  margin: 0;
  color: #20354a;
  font-size: 28px;
  font-weight: 750;
}

.page-header p {
  margin-top: 10px;
  color: #728296;
  font-size: 13px;
  line-height: 1.8;
}

.list-card {
  border-color: #e5eaf0;
  border-radius: 12px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.list-footer {
  padding-top: 16px;
  color: #728296;
  font-size: 12px;
}

@media (max-width: 760px) {
  .material-page {
    width: calc(100% - 32px);
    padding-top: 20px;
  }

  .page-header {
    flex-direction: column;
  }

  .toolbar {
    flex-direction: column;
  }
}
</style>
