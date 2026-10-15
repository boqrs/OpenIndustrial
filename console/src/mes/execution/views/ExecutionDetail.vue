```vue
<template>
  <div class="execution-detail-page">
    <el-card shadow="never">
      <template #header>
        <div class="page-header">
          <div>
            <h2>生产执行详情</h2>
            <p class="subtitle">查看生产执行状态及工序作业进度</p>
          </div>
          <el-button @click="goBack">返回列表</el-button>
        </div>
      </template>

      <el-skeleton v-if="loading" :rows="5" animated />

      <el-empty v-else-if="!execution" description="未找到生产执行记录" />

      <template v-else>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="执行单 ID">
            {{ execution.id }}
          </el-descriptions-item>

          <el-descriptions-item label="工单 ID">
            {{ execution.workOrderId ?? "-" }}
          </el-descriptions-item>

          <el-descriptions-item label="产品 ID">
            {{ execution.productId ?? "-" }}
          </el-descriptions-item>

          <el-descriptions-item label="工艺路线 ID">
            {{ execution.routingId ?? "-" }}
          </el-descriptions-item>

          <el-descriptions-item label="工艺路线版本">
            {{ execution.routingVersion ?? "-" }}
          </el-descriptions-item>

          <el-descriptions-item label="执行状态">
            <el-tag :type="statusTagType(execution.status)">
              {{ statusLabel(execution.status) }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="开始时间">
            {{ formatDate(execution.startedAt) }}
          </el-descriptions-item>

          <el-descriptions-item label="完成时间">
            {{ formatDate(execution.completedAt) }}
          </el-descriptions-item>

          <el-descriptions-item label="创建时间" :span="2">
            {{ formatDate(execution.createdAt) }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="action-bar">
          <el-button
            v-if="execution.status === 'pending'"
            type="primary"
            :loading="starting"
            @click="handleStart"
          >
            开始生产执行
          </el-button>

          <el-button :loading="refreshing" @click="loadExecution">
            刷新
          </el-button>
        </div>

        <el-divider />

        <ExecutionOperation
          :execution-id="execution.id"
          :execution-status="execution.status"
          @changed="loadExecution"
        />
        <ExecutionResultPanel
          :execution-id="execution.id"
          :execution-status="execution.status"
          @changed="loadExecution"
        />
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import ExecutionResultPanel from "./ExecutionResultPanel.vue";
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { getExecution, startExecution } from "../api";
import type { ProductionExecution } from "../types";
import ExecutionOperation from "./ExecutionOperation.vue";

const route = useRoute();
const router = useRouter();

const execution = ref<ProductionExecution | null>(null);
const loading = ref(false);
const refreshing = ref(false);
const starting = ref(false);

function getExecutionId(): number {
  const id = Number(route.params.id);

  if (!Number.isInteger(id) || id <= 0) {
    return 0;
  }

  return id;
}

async function loadExecution(): Promise<void> {
  const id = getExecutionId();

  if (!id) {
    execution.value = null;
    ElMessage.error("无效的生产执行 ID");
    return;
  }

  loading.value = true;
  refreshing.value = true;

  try {
    execution.value = await getExecution(id);
  } catch (error) {
    console.error("加载生产执行详情失败：", error);
    ElMessage.error("加载生产执行详情失败");
    execution.value = null;
  } finally {
    loading.value = false;
    refreshing.value = false;
  }
}

async function handleStart(): Promise<void> {
  if (!execution.value) {
    return;
  }

  try {
    await ElMessageBox.confirm("确定开始该生产执行吗？", "开始生产执行", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });
  } catch {
    return;
  }

  starting.value = true;

  try {
    await startExecution(execution.value.id);
    ElMessage.success("生产执行已启动");
    await loadExecution();
  } catch (error) {
    console.error("启动生产执行失败：", error);
    ElMessage.error("启动生产执行失败");
  } finally {
    starting.value = false;
  }
}

function goBack(): void {
  router.push("/mes/executions");
}

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: "待执行",
    in_progress: "执行中",
    completed: "已完成",
    failed: "失败",
    cancelled: "已取消",
  };

  return labels[status] ?? status;
}

function statusTagType(
  status: string,
): "success" | "warning" | "danger" | "info" | "primary" {
  const types: Record<
    string,
    "success" | "warning" | "danger" | "info" | "primary"
  > = {
    pending: "info",
    in_progress: "primary",
    completed: "success",
    failed: "danger",
    cancelled: "info",
  };

  return types[status] ?? "info";
}

function formatDate(value?: string | null): string {
  if (!value) {
    return "-";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
}

onMounted(loadExecution);
</script>

<style scoped>
.execution-detail-page {
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  margin: 8px 0 0;
  color: #909399;
  font-size: 13px;
}

.action-bar {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}
</style>
```
