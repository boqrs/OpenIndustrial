```vue
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { createExecution, listExecutions } from "../api";
import type { ProductionExecution } from "../types";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const creating = ref(false);
const executions = ref<ProductionExecution[]>([]);
const statusFilter = ref("");
const workOrderFilter = ref<number | undefined>();
const createDialogVisible = ref(false);
const newWorkOrderId = ref<number | undefined>();

const statusOptions = [
  { label: "待执行", value: "pending" },
  { label: "执行中", value: "in_progress" },
  { label: "已完成", value: "completed" },
  { label: "失败", value: "failed" },
  { label: "已取消", value: "cancelled" },
];

function statusLabel(status: string) {
  return statusOptions.find((item) => item.value === status)?.label ?? status;
}

function statusType(
  status: string,
): "info" | "primary" | "success" | "warning" | "danger" {
  switch (status) {
    case "pending":
      return "info";
    case "in_progress":
      return "primary";
    case "completed":
      return "success";
    case "failed":
      return "danger";
    case "cancelled":
      return "warning";
    default:
      return "info";
  }
}

function formatDate(value?: string | null) {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}

function readWorkOrderId(value: unknown): number | undefined {
  const raw = Array.isArray(value) ? value[0] : value;
  if (typeof raw !== "string" && typeof raw !== "number") {
    return undefined;
  }

  const id = Number(raw);
  return Number.isSafeInteger(id) && id > 0 ? id : undefined;
}

async function loadExecutions() {
  loading.value = true;
  try {
    executions.value = await listExecutions({
      ...(statusFilter.value ? { status: statusFilter.value } : {}),
      ...(workOrderFilter.value ? { workOrderId: workOrderFilter.value } : {}),
    });
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "获取生产执行列表失败",
    );
  } finally {
    loading.value = false;
  }
}

function openDetail(row: ProductionExecution) {
  void router.push(`/mes/executions/${row.id}`);
}

function openCreateDialog() {
  newWorkOrderId.value = workOrderFilter.value;
  createDialogVisible.value = true;
}

async function submitCreate() {
  if (
    !newWorkOrderId.value ||
    !Number.isSafeInteger(newWorkOrderId.value) ||
    newWorkOrderId.value <= 0
  ) {
    ElMessage.warning("请输入有效的工单 ID");
    return;
  }

  creating.value = true;
  try {
    const created = await createExecution(newWorkOrderId.value);
    ElMessage.success("生产执行创建成功");
    createDialogVisible.value = false;

    if (created?.id) {
      await router.push(`/mes/executions/${created.id}`);
      return;
    }

    await loadExecutions();
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "创建生产执行失败",
    );
  } finally {
    creating.value = false;
  }
}

async function confirmCreate() {
  try {
    await ElMessageBox.confirm(
      `确认根据工单 ${newWorkOrderId.value ?? ""} 创建生产执行吗？`,
      "创建生产执行",
      {
        type: "warning",
        confirmButtonText: "确认创建",
        cancelButtonText: "取消",
      },
    );
    await submitCreate();
  } catch {
    // 用户取消确认时不执行操作。
  }
}

async function initializePage() {
  const queryWorkOrderId = readWorkOrderId(
    route.query.work_order_id ?? route.query.workOrderId,
  );

  if (queryWorkOrderId) {
    workOrderFilter.value = queryWorkOrderId;
    newWorkOrderId.value = queryWorkOrderId;
  }

  await loadExecutions();

  // 只有从工单详情明确发起创建时，才自动打开创建对话框。
  if (queryWorkOrderId && route.query.create === "1") {
    createDialogVisible.value = true;
  }
}

onMounted(() => {
  void initializePage();
});
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h2>生产执行</h2>
        <p>从生产工单创建执行记录，并跟踪工序操作和生产结果。</p>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        创建生产执行
      </el-button>
    </div>

    <el-card shadow="never" class="filter-card">
      <el-form inline @submit.prevent="loadExecutions">
        <el-form-item label="工单 ID">
          <el-input-number
            v-model="workOrderFilter"
            :min="1"
            :precision="0"
            clearable
            placeholder="输入工单 ID"
          />
        </el-form-item>

        <el-form-item label="执行状态">
          <el-select
            v-model="statusFilter"
            clearable
            placeholder="全部状态"
            style="width: 160px"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="loadExecutions"> 查询 </el-button>
          <el-button
            @click="
              statusFilter = '';
              workOrderFilter = undefined;
              loadExecutions();
            "
          >
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <el-table
        v-loading="loading"
        :data="executions"
        row-key="id"
        empty-text="暂无生产执行记录"
        @row-click="openDetail"
      >
        <el-table-column prop="id" label="执行 ID" width="100" />
        <el-table-column prop="workOrderId" label="工单 ID" width="110" />
        <el-table-column prop="productId" label="产品 ID" width="110" />
        <el-table-column prop="routingId" label="工艺路线 ID" width="130" />
        <el-table-column prop="routingVersion" label="工艺版本" width="100" />

        <el-table-column label="执行状态" width="130">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="开始时间" min-width="175">
          <template #default="{ row }">
            {{ formatDate(row.startedAt) }}
          </template>
        </el-table-column>

        <el-table-column label="创建时间" min-width="175">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="110">
          <template #default="{ row }">
            <el-button link type="primary" @click.stop="openDetail(row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="createDialogVisible"
      title="创建生产执行"
      width="440px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="工单 ID" required>
          <el-input-number
            v-model="newWorkOrderId"
            :min="1"
            :precision="0"
            :step="1"
            placeholder="输入已发布工单的 ID"
            style="width: 100%"
          />
          <div class="field-hint">
            只能为已发布或执行中的工单创建执行记录，最终以服务端校验为准。
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="createDialogVisible = false"> 取消 </el-button>
        <el-button type="primary" :loading="creating" @click="confirmCreate">
          确认创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 8px;
  font-size: 22px;
  color: #172b42;
}

.page-header p {
  margin: 0;
  color: #718096;
  font-size: 13px;
}

.filter-card {
  margin-bottom: 16px;
}

.field-hint {
  margin-top: 6px;
  color: #909399;
  font-size: 12px;
  line-height: 1.6;
}
</style>
```
