<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
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

function statusLabel(status: string): string {
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

function formatDate(value?: string | null): string {
  if (!value) return "—";

  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}

/**
 * 找到第一道尚未完成的工序。
 *
 * 前端这里只负责限制操作入口。
 * 工序顺序和状态转换仍必须由后端最终校验。
 */
const currentOperationId = computed<number | null>(() => {
  const sorted = [...operations.value].sort((a, b) => a.sequence - b.sequence);

  const current = sorted.find(
    (operation) =>
      operation.status !== "completed" &&
      operation.status !== "skipped" &&
      operation.status !== "failed",
  );

  return current?.id ?? null;
});

function canOperate(operation: ExecutionOperation): boolean {
  return (
    props.executionStatus === "in_progress" &&
    submittingId.value === null &&
    (operation.status === "pending" || operation.status === "in_progress")
  );
}

function canStart(operation: ExecutionOperation): boolean {
  return (
    canOperate(operation) &&
    operation.status === "pending" &&
    operation.id === currentOperationId.value
  );
}

function canSubmitResult(operation: ExecutionOperation): boolean {
  return canOperate(operation) && operation.status === "in_progress";
}

function formatResult(result?: Record<string, unknown>): string {
  if (!result || Object.keys(result).length === 0) {
    return "暂无结果数据";
  }

  try {
    return JSON.stringify(result, null, 2);
  } catch {
    return "结果数据无法显示";
  }
}

async function loadOperations(): Promise<void> {
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

async function handleStart(operation: ExecutionOperation): Promise<void> {
  if (!canStart(operation)) {
    ElMessage.warning("请按工序顺序操作");
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确认开始工序「${operation.name}」吗？`,
      "开始工序",
      {
        confirmButtonText: "开始",
        cancelButtonText: "取消",
        type: "warning",
      },
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
): Promise<void> {
  if (!canSubmitResult(operation)) {
    ElMessage.warning("当前工序状态不允许提交结果");
    return;
  }

  const isFail = action === "fail";
  const title = isFail ? "标记工序失败" : "完成工序";

  const initialValue = isFail
    ? JSON.stringify({ reason: "工序执行失败" }, null, 2)
    : "{}";

  let resultText: string;

  try {
    const response = await ElMessageBox.prompt(
      isFail
        ? "请输入失败原因或失败结果，内容必须是 JSON 对象。"
        : "请输入工序结果，内容必须是 JSON 对象。没有额外结果时填写 {}。",
      title,
      {
        inputType: "textarea",
        inputValue: initialValue,
        inputPlaceholder:
          '{\n  "serial_number": "SN-001",\n  "hardware_id": "HW-001"\n}',
        inputValidator: (value) => {
          if (!value?.trim()) {
            return "请输入 JSON 对象";
          }

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

    resultText = response.value;
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
    const payload = { result };

    if (isFail) {
      await failOperation(props.executionId, operation.id, payload);

      ElMessage.warning("工序已标记为失败");
    } else {
      await completeOperation(props.executionId, operation.id, payload);

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

watch(
  () => props.executionId,
  () => {
    operations.value = [];
    void loadOperations();
  },
);
</script>

<template>
  <section class="operation-section">
    <div class="section-header">
      <div>
        <h3>工序作业</h3>
        <p class="section-description">
          按工艺路线顺序执行工序，并记录每道工序的实际结果。
        </p>
      </div>

      <el-button :loading="loading" @click="loadOperations">
        刷新工序
      </el-button>
    </div>

    <el-alert
      v-if="executionStatus !== 'in_progress'"
      title="请先启动生产执行，工序操作才会开放。"
      type="info"
      :closable="false"
      show-icon
      class="operation-alert"
    />

    <el-alert
      v-else-if="currentOperationId === null && operations.length > 0"
      title="所有工序均已进入终态，请刷新执行单并检查生产结果。"
      type="success"
      :closable="false"
      show-icon
      class="operation-alert"
    />

    <el-table
      v-loading="loading"
      :data="[...operations].sort((a, b) => a.sequence - b.sequence)"
      row-key="id"
      empty-text="当前执行没有工序"
      border
    >
      <el-table-column prop="sequence" label="顺序" width="75" />

      <el-table-column prop="code" label="工序编码" min-width="125" />

      <el-table-column prop="name" label="工序名称" min-width="150" />

      <el-table-column prop="workstationId" label="工作站 ID" width="105">
        <template #default="{ row }">
          {{ row.workstationId || "—" }}
        </template>
      </el-table-column>

      <el-table-column label="工序状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="开始时间" min-width="165">
        <template #default="{ row }">
          {{ formatDate(row.startedAt) }}
        </template>
      </el-table-column>

      <el-table-column label="完成时间" min-width="165">
        <template #default="{ row }">
          {{ formatDate(row.completedAt) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" fixed="right" width="205">
        <template #default="{ row }">
          <template v-if="canOperate(row)">
            <el-button
              v-if="canStart(row)"
              link
              type="primary"
              :loading="submittingId === row.id"
              @click="handleStart(row)"
            >
              开始
            </el-button>

            <el-button
              v-if="canSubmitResult(row)"
              link
              type="success"
              :loading="submittingId === row.id"
              @click="submitOperationResult(row, 'complete')"
            >
              完成
            </el-button>

            <el-button
              v-if="canSubmitResult(row)"
              link
              type="danger"
              :loading="submittingId === row.id"
              @click="submitOperationResult(row, 'fail')"
            >
              标记失败
            </el-button>

            <span
              v-if="!canStart(row) && !canSubmitResult(row)"
              class="muted-text"
            >
              等待前序工序
            </span>
          </template>

          <span v-else class="muted-text"> 暂无可用操作 </span>
        </template>
      </el-table-column>

      <el-table-column label="工序结果" min-width="130">
        <template #default="{ row }">
          <el-popover
            v-if="row.result && Object.keys(row.result).length > 0"
            placement="left"
            title="工序结果"
            :width="360"
            trigger="click"
          >
            <template #reference>
              <el-button link type="primary"> 查看结果 </el-button>
            </template>

            <pre class="result-json">{{ formatResult(row.result) }}</pre>
          </el-popover>

          <span v-else class="muted-text"> 暂无结果 </span>
        </template>
      </el-table-column>
    </el-table>

    <p class="table-hint">
      工序操作规则由后端最终校验。若操作失败，请根据接口返回的错误信息检查工序状态及执行顺序。
    </p>
  </section>
</template>

<style scoped>
.operation-section {
  min-width: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 17px;
  font-weight: 600;
  color: #172b42;
}

.section-description {
  margin: 7px 0 0;
  color: #7b8794;
  font-size: 13px;
  line-height: 1.6;
}

.operation-alert {
  margin-bottom: 16px;
}

.muted-text {
  color: #909399;
  font-size: 12px;
}

.result-json {
  max-height: 320px;
  overflow: auto;
  margin: 0;
  padding: 12px;
  border-radius: 6px;
  background: #f5f7fa;
  color: #344054;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.table-hint {
  margin: 12px 0 0;
  color: #909399;
  font-size: 12px;
  line-height: 1.6;
}
</style>
