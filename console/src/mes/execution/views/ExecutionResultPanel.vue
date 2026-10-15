<template>
  <section class="result-panel">
    <div class="section-header">
      <div>
        <h3>生产结果</h3>
        <p class="section-description">
          上传单件产品的生产结果，并在校验完成后确认入账。
        </p>
      </div>
      <el-button
        :loading="loading"
        :disabled="!executionCompleted"
        @click="loadResult"
      >
        刷新结果
      </el-button>
    </div>

    <el-alert
      v-if="!executionCompleted"
      title="生产执行完成后，才能上传生产结果。"
      type="info"
      :closable="false"
      show-icon
    />

    <el-skeleton v-else-if="loading" :rows="3" animated />

    <template v-else>
      <el-alert
        v-if="loadError"
        title="查询生产结果失败，请刷新重试。为避免重复创建，暂不开放上传。"
        type="error"
        :closable="false"
        show-icon
        class="result-alert"
      />

      <template v-else-if="result">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="结果 ID">
            {{ result.id }}
          </el-descriptions-item>
          <el-descriptions-item label="执行单 ID">
            {{ result.execution_id }}
          </el-descriptions-item>
          <el-descriptions-item label="工单 ID">
            {{ result.work_order_id }}
          </el-descriptions-item>
          <el-descriptions-item label="结果状态">
            <el-tag :type="resultStatusType(result.status)">
              {{ resultStatusLabel(result.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="生产数量">
            {{ result.produced_quantity }}
          </el-descriptions-item>
          <el-descriptions-item label="合格数量">
            {{ result.qualified_quantity }}
          </el-descriptions-item>
          <el-descriptions-item label="不合格数量">
            {{ result.rejected_quantity }}
          </el-descriptions-item>
          <el-descriptions-item label="确认时间">
            {{ formatDate(result.confirmed_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ formatDate(result.created_at) }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="result.status === 'draft'" class="result-actions">
          <el-button
            v-if="canCancelResult"
            type="danger"
            plain
            :loading="submitting"
            @click="handleCancel"
          >
            取消草稿
          </el-button>
          <el-button
            v-if="canConfirmResult"
            type="primary"
            :loading="submitting"
            @click="handleConfirm"
          >
            确认生产结果
          </el-button>
        </div>

        <el-alert
          v-if="
            result.status === 'draft' && !canConfirmResult && !canCancelResult
          "
          title="当前账号没有确认或取消生产结果草稿的权限。"
          type="info"
          :closable="false"
          show-icon
          class="result-alert"
        />

        <el-alert
          v-else-if="result.status === 'confirmed'"
          title="生产结果已确认。后端已负责完成相应的工单数量更新及设备创建逻辑。"
          type="success"
          :closable="false"
          show-icon
          class="result-alert"
        />

        <el-alert
          v-else-if="result.status === 'cancelled'"
          title="该生产结果草稿已取消，不能继续确认。"
          type="info"
          :closable="false"
          show-icon
          class="result-alert"
        />
      </template>

      <template v-else>
        <el-empty description="该执行单尚无生产结果" :image-size="80" />

        <el-form
          v-if="canCreateResult"
          :model="form"
          label-width="110px"
          class="result-form"
        >
          <el-form-item label="生产数量">
            <el-input :model-value="1" disabled />
          </el-form-item>

          <el-form-item label="检验结果">
            <el-radio-group v-model="form.qualified">
              <el-radio :value="true">合格</el-radio>
              <el-radio :value="false">不合格</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="合格数量">
            <el-input :model-value="form.qualified ? 1 : 0" disabled />
          </el-form-item>

          <el-form-item label="不合格数量">
            <el-input :model-value="form.qualified ? 0 : 1" disabled />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              :loading="submitting"
              @click="handleCreate"
            >
              上传生产结果
            </el-button>
          </el-form-item>
        </el-form>

        <el-alert
          v-else
          title="当前账号没有上传生产结果的权限。"
          type="info"
          :closable="false"
          show-icon
          class="result-alert"
        />
      </template>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  cancelExecutionResult,
  confirmExecutionResult,
  createExecutionResult,
  getExecutionResultByExecutionId,
} from "../../executionresult/api";
import type { ExecutionResult } from "../../executionresult/api";
import { hasPermission } from "../../../authz/access";
import { Permissions } from "../../../authz/permissions";
import { useAuthStore } from "../../../store/auth";

const props = defineProps<{
  executionId: number;
  executionStatus: string;
}>();

const emit = defineEmits<{
  changed: [];
}>();

const authStore = useAuthStore();

const canCreateResult = computed(() =>
  hasPermission(authStore.user, Permissions.EXECUTION_RESULT_CREATE),
);
const canConfirmResult = computed(() =>
  hasPermission(authStore.user, Permissions.EXECUTION_RESULT_CONFIRM),
);
const canCancelResult = computed(() =>
  hasPermission(authStore.user, Permissions.EXECUTION_RESULT_CANCEL),
);

const loading = ref(false);
const submitting = ref(false);
const loadError = ref(false);
const result = ref<ExecutionResult | null>(null);

const form = reactive({
  qualified: true,
});

const executionCompleted = computed(
  () => props.executionStatus === "completed",
);

async function loadResult(): Promise<void> {
  if (!executionCompleted.value) {
    result.value = null;
    loadError.value = false;
    return;
  }

  loading.value = true;
  loadError.value = false;

  try {
    result.value = await getExecutionResultByExecutionId(props.executionId);
  } catch (error) {
    // 查询失败不能当作“没有结果”，否则可能重复创建。
    console.error("查询生产结果失败：", error);
    result.value = null;
    loadError.value = true;
  } finally {
    loading.value = false;
  }
}

async function handleCreate(): Promise<void> {
  if (!hasPermission(authStore.user, Permissions.EXECUTION_RESULT_CREATE)) {
    ElMessage.warning("当前账号没有上传生产结果的权限");
    return;
  }

  if (!executionCompleted.value || loadError.value || result.value) {
    return;
  }

  try {
    await ElMessageBox.confirm(
      form.qualified
        ? "确认上传本次生产结果：1 件合格、0 件不合格？"
        : "确认上传本次生产结果：0 件合格、1 件不合格？",
      "上传生产结果",
      {
        confirmButtonText: "确认上传",
        cancelButtonText: "返回修改",
        type: "warning",
      },
    );
  } catch {
    return;
  }

  submitting.value = true;

  try {
    result.value = await createExecutionResult({
      execution_id: props.executionId,
      produced_quantity: 1,
      qualified_quantity: form.qualified ? 1 : 0,
      rejected_quantity: form.qualified ? 0 : 1,
    });

    ElMessage.success("生产结果已上传，当前状态为草稿");
    emit("changed");
  } catch (error) {
    console.error("上传生产结果失败：", error);
    ElMessage.error(
      error instanceof Error ? error.message : "上传生产结果失败",
    );
    await loadResult();
  } finally {
    submitting.value = false;
  }
}

async function handleConfirm(): Promise<void> {
  if (!hasPermission(authStore.user, Permissions.EXECUTION_RESULT_CONFIRM)) {
    ElMessage.warning("当前账号没有确认生产结果的权限");
    return;
  }

  if (!result.value || result.value.status !== "draft") {
    return;
  }

  try {
    await ElMessageBox.confirm(
      "确认后，生产结果将进入正式确认流程。请确认数量和检验结果无误。",
      "确认生产结果",
      {
        confirmButtonText: "确认结果",
        cancelButtonText: "取消",
        type: "warning",
      },
    );
  } catch {
    return;
  }

  submitting.value = true;

  try {
    await confirmExecutionResult(result.value.id);
    ElMessage.success("生产结果已确认");
    await loadResult();
    emit("changed");
  } catch (error) {
    console.error("确认生产结果失败：", error);
    ElMessage.error(
      error instanceof Error ? error.message : "确认生产结果失败",
    );
    await loadResult();
  } finally {
    submitting.value = false;
  }
}

async function handleCancel(): Promise<void> {
  if (!hasPermission(authStore.user, Permissions.EXECUTION_RESULT_CANCEL)) {
    ElMessage.warning("当前账号没有取消生产结果的权限");
    return;
  }

  if (!result.value || result.value.status !== "draft") {
    return;
  }

  try {
    await ElMessageBox.confirm(
      "确定取消该生产结果草稿吗？取消后不能继续确认。",
      "取消生产结果",
      {
        confirmButtonText: "确定取消",
        cancelButtonText: "返回",
        type: "warning",
      },
    );
  } catch {
    return;
  }

  submitting.value = true;

  try {
    await cancelExecutionResult(result.value.id);
    ElMessage.success("生产结果草稿已取消");
    await loadResult();
    emit("changed");
  } catch (error) {
    console.error("取消生产结果失败：", error);
    ElMessage.error(
      error instanceof Error ? error.message : "取消生产结果失败",
    );
    await loadResult();
  } finally {
    submitting.value = false;
  }
}

function resultStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    draft: "待确认",
    confirmed: "已确认",
    cancelled: "已取消",
  };
  return labels[status] ?? status;
}

function resultStatusType(
  status: string,
): "success" | "warning" | "danger" | "info" {
  const types: Record<string, "success" | "warning" | "danger" | "info"> = {
    draft: "warning",
    confirmed: "success",
    cancelled: "info",
  };
  return types[status] ?? "info";
}

function formatDate(value?: string | null): string {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}

onMounted(loadResult);

watch(
  () => [props.executionId, props.executionStatus],
  () => {
    result.value = null;
    loadError.value = false;
    void loadResult();
  },
);
</script>

<style scoped>
.result-panel {
  margin-top: 24px;
}

.section-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.section-header h3 {
  margin: 0;
  font-size: 17px;
  font-weight: 600;
}

.section-description {
  margin: 7px 0 0;
  color: #909399;
  font-size: 13px;
}

.result-alert {
  margin-top: 16px;
}

.result-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 18px;
}

.result-form {
  max-width: 560px;
  margin: 8px auto 0;
}
</style>
