<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  completeOperation,
  failOperation,
  listExecutionOperations,
  startOperation,
} from "../api";
import type { ExecutionOperation } from "../types";

const props = defineProps<{
  executionId: number;
  executionStatus: string;
}>();

const emit = defineEmits<{
  changed: [];
}>();

const loading = ref(false);
const submittingId = ref<number | null>(null);
const operations = ref<ExecutionOperation[]>([]);

const statusOptions: Record<string, string> = {
  pending: "待执行",
  in_progress: "执行中",
  completed: "已完成",
  skipped: "已跳过",
  failed: "失败",
};

function statusLabel(status: string) {
  return statusOptions[status] ?? status;
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
    case "skipped":
      return "warning";
    case "failed":
      return "danger";
    default:
      return "info";
  }
}

function formatDate(value?: string | null) {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}

function canOperate(operation: ExecutionOperation) {
  return (
    props.executionStatus === "in_progress" &&
    operation.status !== "completed" &&
    operation.status !== "skipped" &&
    operation.status !== "failed"
  );
}

async function loadOperations() {
  loading.value = true;
  try {
    operations.value = await listExecutionOperations(props.executionId);
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "获取工序列表失败",
    );
  } finally {
    loading.value = false;
  }
}

async function handleStart(operation: ExecutionOperation) {
  try {
    await ElMessageBox.confirm(
      `确认开始工序「${operation.name}」吗？`,
      "开始工序",
      { type: "warning" },
    );
  } catch {
    return;
  }

  submittingId.value = operation.id;
  try {
    await startOperation(props.executionId, operation.id);
    ElMessage.success("工序已开始");
    await loadOperations();
    emit("changed");
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "开始工序失败");
  } finally {
    submittingId.value = null;
  }
}

async function submitOperationResult(
  operation: ExecutionOperation,
  action: "complete" | "fail",
) {
  const isFail = action === "fail";
  const title = isFail ? "标记工序失败" : "完成工序";

  let resultText = isFail
    ? JSON.stringify({ reason: "工序执行失败" }, null, 2)
    : "{}";

  try {
    const { value } = await ElMessageBox.prompt(
      isFail
        ? "请输入失败原因或失败结果 JSON。"
        : "请输入工序结果 JSON；没有额外结果时填写 {}。",
      title,
      {
        inputType: "textarea",
        inputValue: resultText,
        inputPlaceholder: '{\n  "serial_number": "SN-..." \n}',
        inputValidator: (value) => {
          if (!value?.trim()) return "请输入合法 JSON";
          try {
            const parsed: unknown = JSON.parse(value);
            if (
              parsed === null ||
              Array.isArray(parsed) ||
              typeof parsed !== "object"
            ) {
              return "结果必须是 JSON 对象";
            }
            return true;
          } catch {
            return "JSON 格式不正确";
          }
        },
        confirmButtonText: isFail ? "确认失败" : "确认完成",
        cancelButtonText: "取消",
        type: isFail ? "warning" : "info",
      },
    );

    resultText = value;
  } catch {
    return;
  }

  let result: Record<string, unknown>;
  try {
    const parsed: unknown = JSON.parse(resultText);
    if (
      parsed === null ||
      Array.isArray(parsed) ||
      typeof parsed !== "object"
    ) {
      ElMessage.error("工序结果必须是 JSON 对象");
      return;
    }
    result = parsed as Record<string, unknown>;
  } catch {
    ElMessage.error("工序结果 JSON 格式不正确");
    return;
  }

  submittingId.value = operation.id;
  try {
    const data = { result };
    if (isFail) {
      await failOperation(props.executionId, operation.id, data);
      ElMessage.warning("工序已标记为失败");
    } else {
      await completeOperation(props.executionId, operation.id, data);
      ElMessage.success("工序已完成");
    }

    await loadOperations();
    emit("changed");
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : `${title}失败`);
  } finally {
    submittingId.value = null;
  }
}

onMounted(loadOperations);
watch(() => props.executionId, loadOperations);
</script>

<template>
  <div>
    <el-alert
      v-if="executionStatus !== 'in_progress'"
      title="请先启动生产执行，工序操作才会开放。"
      type="info"
      :closable="false"
      show-icon
      class="operation-alert"
    />

    <el-table
      v-loading="loading"
      :data="operations"
      row-key="id"
      empty-text="当前执行没有工序"
    >
      <el-table-column prop="sequence" label="顺序" width="80" />
      <el-table-column prop="code" label="工序编码" width="150" />
      <el-table-column prop="name" label="工序名称" min-width="150" />
      <el-table-column prop="workstationId" label="工作站 ID" width="120">
        <template #default="{ row }">
          {{ row.workstationId || "—" }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="开始时间" min-width="170">
        <template #default="{ row }">
          {{ formatDate(row.startedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="完成时间" min-width="170">
        <template #default="{ row }">
          {{ formatDate(row.completedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="230">
        <template #default="{ row }">
          <template v-if="canOperate(row)">
            <el-button
              v-if="row.status === 'pending'"
              link
              type="primary"
              :loading="submittingId === row.id"
              @click="handleStart(row)"
            >
              开始
            </el-button>
            <el-button
              v-if="row.status === 'in_progress'"
              link
              type="success"
              :loading="submittingId === row.id"
              @click="submitOperationResult(row, 'complete')"
            >
              完成
            </el-button>
            <el-button
              v-if="row.status === 'in_progress'"
              link
              type="danger"
              :loading="submittingId === row.id"
              @click="submitOperationResult(row, 'fail')"
            >
              标记失败
            </el-button>
          </template>
          <span v-else class="muted-text">暂无可用操作</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.operation-alert {
  margin-bottom: 16px;
}

.muted-text {
  color: #909399;
  font-size: 12px;
}
</style>
