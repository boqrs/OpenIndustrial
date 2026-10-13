<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Plus,
  Refresh,
  Search,
  OfficeBuilding,
  Connection,
} from "@element-plus/icons-vue";

import {
  createFactory,
  deleteFactory,
  listFactories,
  updateFactory,
} from "../api";

import type {
  CreateFactoryRequest,
  Factory,
  UpdateFactoryRequest,
} from "../types";

const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const factories = ref<Factory[]>([]);
const keyword = ref("");

const dialogVisible = ref(false);
const editingId = ref<number | null>(null);

const emptyForm = (): CreateFactoryRequest => ({
  name: "",
  code: "",
  address: "",
  timezone: "Asia/Shanghai",
});

const form = ref<CreateFactoryRequest>(emptyForm());

const rules = {
  name: [{ required: true, message: "请输入工厂名称", trigger: "blur" }],
  code: [{ required: true, message: "请输入工厂编码", trigger: "blur" }],
  timezone: [{ required: true, message: "请输入时区", trigger: "blur" }],
};

const formRef = ref();

const filteredFactories = () => {
  const text = keyword.value.trim().toLowerCase();

  if (!text) return factories.value;

  return factories.value.filter((item) =>
    [item.name, item.code, item.address ?? "", item.timezone].some((value) =>
      value.toLowerCase().includes(text),
    ),
  );
};

async function loadFactories() {
  loading.value = true;

  try {
    factories.value = await listFactories();
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "获取工厂列表失败",
    );
  } finally {
    loading.value = false;
  }
}

function openCreateDialog() {
  editingId.value = null;
  form.value = emptyForm();
  dialogVisible.value = true;
}

function openEditDialog(row: Factory) {
  editingId.value = row.id;

  form.value = {
    name: row.name,
    code: row.code,
    address: row.address ?? "",
    timezone: row.timezone || "Asia/Shanghai",
  };

  dialogVisible.value = true;
}

async function submitForm() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return;

    saving.value = true;

    try {
      const payload: CreateFactoryRequest = {
        name: form.value.name.trim(),
        code: form.value.code.trim(),
        address: form.value.address?.trim() ?? "",
        timezone: form.value.timezone?.trim() || "Asia/Shanghai",
      };

      if (editingId.value !== null) {
        const updatePayload: UpdateFactoryRequest = payload;

        await updateFactory(editingId.value, updatePayload);
        ElMessage.success("工厂信息已更新");
      } else {
        await createFactory(payload);
        ElMessage.success("工厂创建成功");
      }

      dialogVisible.value = false;
      await loadFactories();
    } catch (error) {
      ElMessage.error(error instanceof Error ? error.message : "保存工厂失败");
    } finally {
      saving.value = false;
    }
  });
}

async function removeFactory(row: Factory) {
  try {
    await ElMessageBox.confirm(
      `确定删除工厂「${row.name}」吗？如果工厂下仍有生产线或工作站，后端会拒绝删除。`,
      "删除工厂",
      {
        type: "warning",
        confirmButtonText: "确认删除",
        cancelButtonText: "取消",
      },
    );

    await deleteFactory(row.id);
    ElMessage.success("工厂已删除");
    await loadFactories();
  } catch (error) {
    if (error === "cancel" || error === "close") return;

    ElMessage.error(error instanceof Error ? error.message : "删除工厂失败");
  }
}

function openTopology(row: Factory) {
  router.push({
    name: "ProductionTopology",
    query: { factory_id: String(row.id) },
  });
}

function statusText(status: string) {
  const map: Record<string, string> = {
    active: "启用",
    inactive: "停用",
    disabled: "停用",
  };

  return map[status.toLowerCase()] ?? status;
}

onMounted(loadFactories);
</script>

<template>
  <div class="factory-page">
    <section class="page-hero">
      <div class="hero-icon">
        <el-icon :size="26"><OfficeBuilding /></el-icon>
      </div>

      <div class="hero-copy">
        <div class="eyebrow">FACTORY MANAGEMENT</div>
        <h1>工厂管理</h1>
        <p>维护租户下的工厂信息，并管理工厂对应的生产组织结构。</p>
      </div>

      <div class="hero-actions">
        <el-button :icon="Refresh" @click="loadFactories"> 刷新 </el-button>
        <el-button type="primary" :icon="Plus" @click="openCreateDialog">
          新建工厂
        </el-button>
      </div>
    </section>

    <section class="content-card">
      <div class="toolbar">
        <div>
          <h2>工厂列表</h2>
          <p>当前租户的工厂基础信息</p>
        </div>

        <el-input
          v-model="keyword"
          clearable
          placeholder="搜索名称、编码、地址或时区"
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <el-table
        v-loading="loading"
        :data="filteredFactories()"
        row-key="id"
        stripe
        empty-text="暂无工厂，请先创建工厂"
      >
        <el-table-column prop="code" label="工厂编码" min-width="130" />

        <el-table-column label="工厂名称" min-width="180">
          <template #default="{ row }">
            <div class="factory-name">
              <span class="name-mark">
                <el-icon><OfficeBuilding /></el-icon>
              </span>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column
          prop="address"
          label="地址"
          min-width="200"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ row.address || "—" }}
          </template>
        </el-table-column>

        <el-table-column prop="timezone" label="时区" min-width="150" />

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'active' ? 'success' : 'info'"
              effect="light"
            >
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" min-width="180">
          <template #default="{ row }">
            {{
              row.created_at ? new Date(row.created_at).toLocaleString() : "—"
            }}
          </template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="250">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :icon="Connection"
              @click="openTopology(row)"
            >
              生产拓扑
            </el-button>
            <el-button link type="primary" @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button link type="danger" @click="removeFactory(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-footer">共 {{ filteredFactories().length }} 个工厂</div>
    </section>

    <el-dialog
      v-model="dialogVisible"
      :title="editingId === null ? '新建工厂' : '编辑工厂'"
      width="560px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="工厂名称" prop="name">
          <el-input
            v-model="form.name"
            maxlength="100"
            placeholder="例如：苏州工厂"
          />
        </el-form-item>

        <el-form-item label="工厂编码" prop="code">
          <el-input
            v-model="form.code"
            maxlength="64"
            placeholder="例如：SUZHOU-01"
          />
        </el-form-item>

        <el-form-item label="工厂地址">
          <el-input
            v-model="form.address"
            maxlength="255"
            placeholder="请输入工厂地址"
          />
        </el-form-item>

        <el-form-item label="时区" prop="timezone">
          <el-input v-model="form.timezone" placeholder="例如：Asia/Shanghai" />
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
.factory-page {
  min-height: 100%;
  padding: 26px 30px;
  background: #f5f7fa;
  color: #172b42;
}

.page-hero {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 26px 30px;
  border-radius: 14px;
  color: #fff;
  background: linear-gradient(115deg, #0e1f33, #203a55);
}

.hero-icon {
  display: grid;
  width: 56px;
  height: 56px;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid rgb(240 160 48 / 35%);
  border-radius: 14px;
  color: #f0a030;
  background: rgb(255 255 255 / 5%);
}

.hero-copy {
  flex: 1;
  min-width: 0;
}

.eyebrow {
  color: #a9bdcf;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1.8px;
}

.hero-copy h1 {
  margin: 9px 0;
  font-size: 25px;
}

.hero-copy p {
  margin: 0;
  color: #c4d0dd;
  font-size: 12px;
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.content-card {
  margin-top: 22px;
  padding: 22px;
  border: 1px solid #e6ebf1;
  border-radius: 12px;
  background: #fff;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 20px;
}

.toolbar h2 {
  margin: 0;
  font-size: 17px;
}

.toolbar p {
  margin: 7px 0 0;
  color: #8190a0;
  font-size: 12px;
}

.search-input {
  max-width: 340px;
}

.factory-name {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
}

.name-mark {
  display: grid;
  width: 30px;
  height: 30px;
  place-items: center;
  border-radius: 8px;
  color: #547ba1;
  background: #edf3f9;
}

.table-footer {
  padding-top: 16px;
  color: #8290a0;
  font-size: 12px;
}

@media (max-width: 900px) {
  .factory-page {
    padding: 18px;
  }

  .page-hero {
    align-items: flex-start;
    flex-wrap: wrap;
    padding: 22px;
  }

  .hero-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .toolbar {
    align-items: flex-start;
    flex-direction: column;
  }

  .search-input {
    width: 100%;
    max-width: none;
  }
}
</style>
