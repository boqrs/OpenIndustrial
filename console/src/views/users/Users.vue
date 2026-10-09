<template>
  <div class="users-page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div>
        <div class="page-eyebrow">IDENTITY MANAGEMENT</div>
        <h1>用户管理</h1>
        <p class="page-description">
          管理当前工厂的用户账户、激活状态及账户邀请。
        </p>
      </div>

      <el-button
        type="primary"
        :icon="Refresh"
        :loading="loading"
        @click="loadUsers"
      >
        刷新列表
      </el-button>
    </header>

    <!-- 用户统计 -->
    <section class="summary-grid">
      <button
        v-for="item in summaryCards"
        :key="item.key"
        type="button"
        class="summary-card"
        :class="{ 'summary-card-active': selectedStatus === item.key }"
        @click="selectStatus(item.key)"
      >
        <span class="summary-label">{{ item.label }}</span>
        <strong class="summary-value">
          {{ item.value }}
        </strong>
        <span class="summary-description">{{ item.description }}</span>
      </button>
    </section>

    <!-- 用户列表 -->
    <section class="content-panel">
      <div class="panel-header">
        <div>
          <h2>账户列表</h2>
          <p>查看用户信息，并为待审核账户发送邀请。</p>
        </div>

        <el-button type="primary" :icon="Plus" @click="openInviteDialog">
          邀请用户
        </el-button>
      </div>

      <!-- 筛选条件 -->
      <div class="filter-bar">
        <el-input
          v-model="keyword"
          clearable
          placeholder="搜索姓名或邮箱"
          class="keyword-input"
          :prefix-icon="Search"
          @keyup.enter="searchUsers"
          @clear="searchUsers"
        />

        <el-select
          v-model="selectedStatus"
          class="status-select"
          placeholder="全部状态"
          @change="handleStatusChange"
        >
          <el-option label="全部状态" value="" />
          <el-option label="待审核" value="init" />
          <el-option label="待激活" value="invited" />
          <el-option label="正常" value="active" />
          <el-option label="已停用" value="disabled" />
        </el-select>

        <el-button type="primary" @click="searchUsers"> 搜索 </el-button>

        <el-button @click="resetFilters"> 重置 </el-button>
      </div>

      <!-- 加载错误 -->
      <el-alert
        v-if="loadError"
        :title="loadError"
        type="error"
        show-icon
        :closable="false"
        class="error-alert"
      />

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="users"
        row-key="id"
        class="users-table"
        empty-text="暂无符合条件的用户"
      >
        <el-table-column label="用户" min-width="230">
          <template #default="{ row }">
            <div class="user-cell">
              <div class="user-avatar">
                {{ getInitial(row.name) }}
              </div>

              <div class="user-info">
                <span class="user-name">{{ row.name }}</span>
                <span class="user-email">{{ row.email }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="用户类型" width="130">
          <template #default="{ row }">
            <span class="user-type">
              {{ formatUserType(row.user_type) }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="账户状态" width="130">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" effect="light" round>
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" min-width="180">
          <template #default="{ row }">
            <span class="date-text">
              {{ formatDate(row.created_at) }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'init'"
              type="primary"
              link
              :loading="invitingUserId === String(row.id)"
              @click="openInviteDialog(row)"
            >
              发送邀请
            </el-button>

            <span v-else class="no-action">—</span>
          </template>
        </el-table-column>
      </el-table>

      <footer class="table-footer">
        <span class="result-count"> 共 {{ total }} 条记录 </span>

        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="sizes, prev, pager, next"
          background
          @current-change="loadUsers"
          @size-change="handlePageSizeChange"
        />
      </footer>
    </section>

    <!-- 邀请用户弹窗 -->
    <el-dialog
      v-model="inviteDialogVisible"
      title="邀请用户加入"
      width="520px"
      :close-on-click-modal="!submittingInvite"
      :close-on-press-escape="!submittingInvite"
      @closed="resetInviteForm"
    >
      <div class="invite-intro">
        <div class="invite-icon">
          <el-icon :size="22">
            <Message />
          </el-icon>
        </div>

        <div>
          <h3>发送账户激活邀请</h3>
          <p>提交后，系统将通过邮件发送邀请，用户完成激活后即可登录。</p>
        </div>
      </div>

      <el-alert
        v-if="inviteError"
        :title="inviteError"
        type="error"
        show-icon
        :closable="false"
        class="invite-error"
      />

      <el-form
        ref="inviteFormRef"
        :model="inviteForm"
        :rules="inviteRules"
        label-position="top"
        @submit.prevent="submitInvitation"
      >
        <el-form-item label="用户姓名" prop="name">
          <el-input
            v-model.trim="inviteForm.name"
            placeholder="请输入用户姓名"
            maxlength="100"
            clearable
          />
        </el-form-item>

        <el-form-item label="电子邮箱" prop="email">
          <el-input
            v-model.trim="inviteForm.email"
            placeholder="请输入用户邮箱"
            maxlength="254"
            clearable
          />
        </el-form-item>

        <el-form-item label="账户角色" prop="role_id">
          <el-select
            v-model="inviteForm.role_id"
            placeholder="请选择角色"
            style="width: 100%"
            :loading="rolesLoading"
            :disabled="rolesLoading || roles.length === 0"
            no-data-text="暂无可用角色"
          >
            <el-option
              v-for="role in roles"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            >
              <div class="role-option">
                <span>{{ role.name }}</span>
                <span v-if="role.description" class="role-description">
                  {{ role.description }}
                </span>
              </div>
            </el-option>
          </el-select>

          <div v-if="!rolesLoading && roles.length === 0" class="field-tip">
            暂未获取到可用角色，请检查角色接口及当前账户权限。
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button
            :disabled="submittingInvite"
            @click="inviteDialogVisible = false"
          >
            取消
          </el-button>

          <el-button
            type="primary"
            :loading="submittingInvite"
            :disabled="roles.length === 0 || rolesLoading"
            @click="submitInvitation"
          >
            发送邀请
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { Message, Plus, Refresh, Search } from "@element-plus/icons-vue";

import {
  inviteUser,
  listRoles,
  listUsers,
  type IdentityRole,
  type ManagedUser,
  type UserStatus,
} from "../../api/users";

type StatusFilter = "" | UserStatus;

const route = useRoute();
const router = useRouter();

const users = ref<ManagedUser[]>([]);
const total = ref(0);

const loading = ref(false);
const loadError = ref("");

const keyword = ref("");
const selectedStatus = ref<StatusFilter>("");
const currentPage = ref(1);
const pageSize = ref(20);

const inviteDialogVisible = ref(false);
const submittingInvite = ref(false);
const inviteError = ref("");
const invitingUserId = ref<string | null>(null);

const roles = ref<IdentityRole[]>([]);
const rolesLoading = ref(false);

const inviteFormRef = ref<FormInstance>();

const inviteForm = reactive({
  name: "",
  email: "",
  role_id: undefined as number | undefined,
});

const inviteRules: FormRules = {
  name: [
    {
      required: true,
      message: "请输入用户姓名",
      trigger: "blur",
    },
    {
      min: 2,
      max: 100,
      message: "姓名长度应为 2 到 100 个字符",
      trigger: "blur",
    },
  ],
  email: [
    {
      required: true,
      message: "请输入电子邮箱",
      trigger: "blur",
    },
    {
      type: "email",
      message: "请输入有效的电子邮箱地址",
      trigger: ["blur", "change"],
    },
  ],
  role_id: [
    {
      required: true,
      message: "请选择账户角色",
      trigger: "change",
    },
  ],
};

const summaryCards = computed(() => {
  const count = (status: UserStatus | "") => {
    if (status === "") {
      return total.value;
    }

    return users.value.filter((user) => user.status === status).length;
  };

  return [
    {
      key: "" as StatusFilter,
      label: "用户总数",
      value: count(""),
      description: "当前查询结果总数",
    },
    {
      key: "init" as StatusFilter,
      label: "待审核",
      value: count("init"),
      description: "等待管理员邀请",
    },
    {
      key: "invited" as StatusFilter,
      label: "待激活",
      value: count("invited"),
      description: "等待用户完成激活",
    },
    {
      key: "active" as StatusFilter,
      label: "正常用户",
      value: count("active"),
      description: "已激活的账户",
    },
    {
      key: "disabled" as StatusFilter,
      label: "已停用",
      value: count("disabled"),
      description: "当前已停用的账户",
    },
  ];
});

function getErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  if (
    typeof error === "object" &&
    error !== null &&
    "message" in error &&
    typeof error.message === "string"
  ) {
    return error.message;
  }

  return fallback;
}

function getInitial(name: string): string {
  const value = name.trim();

  return value ? value.slice(0, 1).toUpperCase() : "U";
}

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    init: "待审核",
    invited: "待激活",
    active: "正常",
    disabled: "已停用",
  };

  return labels[status] ?? "未知状态";
}

function statusTagType(
  status: string,
): "success" | "warning" | "danger" | "info" {
  const types: Record<string, "success" | "warning" | "danger" | "info"> = {
    init: "warning",
    invited: "info",
    active: "success",
    disabled: "danger",
  };

  return types[status] ?? "info";
}

function formatUserType(userType: string): string {
  const labels: Record<string, string> = {
    admin: "管理员",
    employee: "员工",
    operator: "操作员",
    viewer: "只读用户",
  };

  return labels[userType] ?? userType ?? "普通用户";
}

function formatDate(value?: string): string {
  if (!value) {
    return "—";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(date);
}

async function loadUsers(): Promise<void> {
  loading.value = true;
  loadError.value = "";

  try {
    const result = await listUsers({
      status: selectedStatus.value || undefined,
      keyword: keyword.value.trim() || undefined,
      limit: pageSize.value,
      offset: (currentPage.value - 1) * pageSize.value,
    });

    users.value = result.users;
    total.value = result.total;
  } catch (error: unknown) {
    loadError.value = getErrorMessage(
      error,
      "获取用户列表失败，请检查网络或登录状态。",
    );

    users.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function searchUsers(): void {
  currentPage.value = 1;
  void loadUsers();
}

function handleStatusChange(): void {
  currentPage.value = 1;
  void syncStatusToRoute();
  void loadUsers();
}

function selectStatus(status: StatusFilter): void {
  selectedStatus.value = status;
  currentPage.value = 1;
  void syncStatusToRoute();
  void loadUsers();
}

function resetFilters(): void {
  keyword.value = "";
  selectedStatus.value = "";
  currentPage.value = 1;
  void syncStatusToRoute();
  void loadUsers();
}

function handlePageSizeChange(): void {
  currentPage.value = 1;
  void loadUsers();
}

async function syncStatusToRoute(): Promise<void> {
  const query = { ...route.query };

  if (selectedStatus.value) {
    query.status = selectedStatus.value;
  } else {
    delete query.status;
  }

  if (query.status === route.query.status) {
    return;
  }

  await router.replace({ query });
}

function readStatusFromRoute(): StatusFilter {
  const value = route.query.status;
  const status = Array.isArray(value) ? value[0] : value;

  if (
    status === "init" ||
    status === "invited" ||
    status === "active" ||
    status === "disabled"
  ) {
    return status;
  }

  return "";
}

watch(
  () => route.query.status,
  (value) => {
    const status = Array.isArray(value) ? value[0] : value;
    const nextStatus: StatusFilter =
      status === "init" ||
      status === "invited" ||
      status === "active" ||
      status === "disabled"
        ? status
        : "";

    if (selectedStatus.value !== nextStatus) {
      selectedStatus.value = nextStatus;
      currentPage.value = 1;
      void loadUsers();
    }
  },
);

async function loadRoles(): Promise<void> {
  rolesLoading.value = true;

  try {
    roles.value = await listRoles();
  } catch (error: unknown) {
    roles.value = [];

    ElMessage.error(getErrorMessage(error, "获取角色列表失败，请稍后重试。"));
  } finally {
    rolesLoading.value = false;
  }
}

async function openInviteDialog(user?: ManagedUser): Promise<void> {
  inviteError.value = "";
  inviteDialogVisible.value = true;

  inviteForm.name = user?.name && user.name !== "未设置姓名" ? user.name : "";
  inviteForm.email = user?.email ?? "";
  inviteForm.role_id = undefined;

  if (roles.value.length === 0 && !rolesLoading.value) {
    await loadRoles();
  }
}

function resetInviteForm(): void {
  inviteError.value = "";
  inviteForm.name = "";
  inviteForm.email = "";
  inviteForm.role_id = undefined;
  inviteFormRef.value?.clearValidate();
}

async function submitInvitation(): Promise<void> {
  if (!inviteFormRef.value || submittingInvite.value) {
    return;
  }

  inviteError.value = "";

  try {
    await inviteFormRef.value.validate();
  } catch {
    return;
  }

  if (inviteForm.role_id === undefined) {
    inviteError.value = "请选择账户角色。";
    return;
  }

  submittingInvite.value = true;

  try {
    await inviteUser({
      name: inviteForm.name.trim(),
      email: inviteForm.email.trim(),
      role_id: inviteForm.role_id,
    });

    ElMessage.success("邀请已提交。");

    inviteDialogVisible.value = false;

    await loadUsers();
    await loadRoles();
  } catch (error: unknown) {
    inviteError.value = getErrorMessage(
      error,
      "发送邀请失败，请检查输入信息后重试。",
    );
  } finally {
    submittingInvite.value = false;
  }
}

onMounted(() => {
  selectedStatus.value = readStatusFromRoute();
  void loadUsers();
});
</script>

<style scoped>
.users-page {
  min-height: 100%;
  padding: 28px;
  color: #1d2b3d;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 24px;
}

.page-eyebrow {
  margin-bottom: 8px;
  color: #9b762d;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1.8px;
}

.page-header h1 {
  margin: 0;
  color: #0e1f33;
  font-size: 28px;
  font-weight: 700;
}

.page-description {
  margin: 10px 0 0;
  color: #7c8999;
  font-size: 13px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card {
  display: flex;
  min-width: 0;
  flex-direction: column;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid #e5eaf0;
  border-radius: 10px;
  background: #fff;
  color: inherit;
  text-align: left;
  cursor: pointer;
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease;
}

.summary-card:hover {
  transform: translateY(-2px);
  border-color: #d3b16d;
  box-shadow: 0 6px 20px rgb(14 31 51 / 6%);
}

.summary-card-active {
  border-color: #c49a45;
  box-shadow: inset 0 0 0 1px #c49a45;
}

.summary-card:focus-visible {
  outline: 2px solid #c49a45;
  outline-offset: 3px;
}

.summary-label {
  color: #748195;
  font-size: 13px;
}

.summary-value {
  margin: 14px 0 8px;
  color: #0e1f33;
  font-size: 30px;
  font-weight: 700;
  line-height: 1.2;
}

.summary-description {
  color: #98a2b1;
  font-size: 12px;
}

.content-panel {
  padding: 24px;
  border: 1px solid #e7ebf0;
  border-radius: 10px;
  background: #fff;
  box-shadow: 0 3px 12px rgb(14 31 51 / 2%);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 22px;
}

.panel-header h2 {
  margin: 0;
  color: #17283d;
  font-size: 18px;
  font-weight: 700;
}

.panel-header p {
  margin: 8px 0 0;
  color: #8792a1;
  font-size: 12px;
}

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 20px;
}

.keyword-input {
  width: 280px;
}

.status-select {
  width: 160px;
}

.error-alert {
  margin-bottom: 16px;
}

.users-table {
  width: 100%;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.user-avatar {
  display: flex;
  width: 38px;
  height: 38px;
  flex: 0 0 38px;
  align-items: center;
  justify-content: center;
  border-radius: 9px;
  background: #eaf0f7;
  color: #24415e;
  font-size: 15px;
  font-weight: 700;
}

.user-info {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 5px;
}

.user-name {
  overflow: hidden;
  color: #243449;
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-email {
  overflow: hidden;
  color: #8994a3;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-type {
  color: #526176;
  font-size: 13px;
}

.date-text {
  color: #7c8999;
  font-size: 12px;
}

.no-action {
  color: #c1c8d1;
}

.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-top: 20px;
}

.result-count {
  color: #8994a3;
  font-size: 12px;
}

.invite-intro {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 24px;
}

.invite-icon {
  display: flex;
  width: 44px;
  height: 44px;
  flex: 0 0 44px;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: #f7f0df;
  color: #9b762d;
}

.invite-intro h3 {
  margin: 2px 0 8px;
  color: #1d2b3d;
  font-size: 15px;
}

.invite-intro p {
  margin: 0;
  color: #8490a0;
  font-size: 12px;
  line-height: 1.7;
}

.invite-error {
  margin-bottom: 18px;
}

.role-option {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 4px 0;
}

.role-description {
  color: #8b96a5;
  font-size: 11px;
}

.field-tip {
  margin-top: 7px;
  color: #a66c26;
  font-size: 12px;
  line-height: 1.6;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 1200px) {
  .summary-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .users-page {
    padding: 16px;
  }

  .page-header {
    align-items: flex-start;
  }

  .page-header h1 {
    font-size: 23px;
  }

  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
  }

  .summary-card {
    padding: 15px;
  }

  .summary-value {
    font-size: 26px;
  }

  .content-panel {
    padding: 16px;
  }

  .panel-header {
    align-items: flex-start;
  }

  .filter-bar {
    flex-direction: column;
  }

  .keyword-input,
  .status-select {
    width: 100%;
  }

  .table-footer {
    align-items: flex-start;
    flex-direction: column;
  }

  .table-footer .el-pagination {
    max-width: 100%;
    flex-wrap: wrap;
  }
}
</style>
