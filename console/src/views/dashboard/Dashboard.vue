<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  Connection,
  DataAnalysis,
  Refresh,
  SwitchButton,
  WarningFilled,
} from "@element-plus/icons-vue";

import { getDashboardOverview } from "../../api/dashboard";
import type { UserStats } from "../../api/dashboard";
import { useAuthStore } from "../../store/auth";

const router = useRouter();
const authStore = useAuthStore();

const stats = ref<UserStats | null>(null);
const loading = ref(false);
const loadError = ref("");

const displayName = computed(
  () => authStore.user?.name || authStore.user?.email || "管理员",
);

const statCards = computed(() => [
  {
    key: "total",
    label: "用户总数",
    value: stats.value?.total,
    description: "当前工厂的全部用户",
    icon: "U",
    tone: "navy",
  },
  {
    key: "init",
    label: "待审核申请",
    value: stats.value?.init,
    description: "等待管理员邀请",
    icon: "A",
    tone: "gold",
  },
  {
    key: "invited",
    label: "待激活用户",
    value: stats.value?.invited,
    description: "已邀请，等待完成激活",
    icon: "I",
    tone: "blue",
  },
  {
    key: "active",
    label: "正常用户",
    value: stats.value?.active,
    description: "已激活的用户账号",
    icon: "✓",
    tone: "green",
  },
  {
    key: "disabled",
    label: "已停用用户",
    value: stats.value?.disabled,
    description: "当前不可正常登录",
    icon: "—",
    tone: "gray",
  },
]);

const statusRows = computed(() => {
  const data = stats.value;

  if (!data) {
    return [];
  }

  const total = data.total;

  return [
    {
      key: "active",
      label: "正常用户",
      value: data.active,
      percent: total > 0 ? (data.active / total) * 100 : 0,
      color: "var(--color-success, #16a34a)",
    },
    {
      key: "init",
      label: "待审核申请",
      value: data.init,
      percent: total > 0 ? (data.init / total) * 100 : 0,
      color: "#f0a030",
    },
    {
      key: "invited",
      label: "待激活用户",
      value: data.invited,
      percent: total > 0 ? (data.invited / total) * 100 : 0,
      color: "#5288c7",
    },
    {
      key: "disabled",
      label: "已停用用户",
      value: data.disabled,
      percent: total > 0 ? (data.disabled / total) * 100 : 0,
      color: "#98a2b3",
    },
  ];
});

async function loadOverview() {
  loading.value = true;
  loadError.value = "";

  try {
    const overview = await getDashboardOverview();
    stats.value = overview.users;
  } catch (error: unknown) {
    loadError.value = "暂时无法获取平台统计数据，请检查网络或稍后重试。";
    console.error("Failed to load dashboard overview:", error);
  } finally {
    loading.value = false;
  }
}

async function handleLogout() {
  try {
    await authStore.logout();
  } finally {
    await router.replace("/login");
  }
}

onMounted(() => {
  void loadOverview();
});
</script>

<template>
  <div class="dashboard-page">
    <header class="topbar">
      <div class="brand">
        <div class="brand-mark">
          <span class="brand-mark-inner">OI</span>
        </div>

        <div class="brand-copy">
          <div class="brand-name">OpenIndustrial</div>
          <div class="brand-subtitle">设备数字身份平台</div>
        </div>
      </div>

      <div class="topbar-right">
        <div class="system-status">
          <span class="status-dot"></span>
          <span>工业云平台</span>
        </div>

        <div class="topbar-divider"></div>

        <div class="account-info">
          <div class="account-avatar">
            {{ displayName.slice(0, 1).toUpperCase() }}
          </div>

          <div class="account-copy">
            <span class="account-name">{{ displayName }}</span>
            <span class="account-role">平台用户</span>
          </div>
        </div>

        <el-button class="logout-button" text @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          <span>退出登录</span>
        </el-button>
      </div>
    </header>

    <main class="main-content">
      <section class="welcome-section">
        <div>
          <div class="eyebrow">
            <span class="eyebrow-line"></span>
            INDUSTRIAL CLOUD
          </div>

          <h1>工作台</h1>

          <p class="welcome-description">
            欢迎回来，{{ displayName }}。这里是您的工业云管理入口。
          </p>
        </div>

        <div class="welcome-decoration" aria-hidden="true">
          <div class="decoration-ring ring-one"></div>
          <div class="decoration-ring ring-two"></div>
          <div class="decoration-center">OI</div>
          <span class="decoration-node node-one"></span>
          <span class="decoration-node node-two"></span>
          <span class="decoration-node node-three"></span>
        </div>
      </section>

      <section class="stats-section">
        <div class="section-heading">
          <div>
            <h2>账号概览</h2>
            <p>当前工厂的用户账号状态</p>
          </div>

          <el-button
            class="refresh-button"
            :loading="loading"
            @click="loadOverview"
          >
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
        </div>

        <div v-if="loadError" class="error-panel" role="alert">
          <div class="error-content">
            <el-icon class="error-icon"><WarningFilled /></el-icon>
            <div>
              <strong>数据加载失败</strong>
              <p>{{ loadError }}</p>
            </div>
          </div>

          <el-button type="primary" :loading="loading" @click="loadOverview">
            重试
          </el-button>
        </div>

        <div class="stats-grid" :aria-busy="loading">
          <article v-for="card in statCards" :key="card.key" class="stat-card">
            <div class="stat-card-top">
              <span class="stat-label">{{ card.label }}</span>

              <span
                class="stat-icon"
                :class="`tone-${card.tone}`"
                aria-hidden="true"
              >
                {{ card.icon }}
              </span>
            </div>

            <div class="stat-value">
              <span v-if="loading && !stats" class="skeleton-value"></span>
              <span v-else>{{ card.value ?? "—" }}</span>
            </div>

            <div class="stat-description">{{ card.description }}</div>
          </article>
        </div>
      </section>

      <section class="details-grid">
        <article class="panel status-panel">
          <div class="panel-heading">
            <div>
              <h2>账号状态分布</h2>
              <p>根据后端实时统计数据计算</p>
            </div>

            <span class="panel-symbol">
              <el-icon><DataAnalysis /></el-icon>
            </span>
          </div>

          <div v-if="loading && !stats" class="panel-loading">
            正在加载统计数据…
          </div>

          <div v-else-if="loadError && !stats" class="panel-empty">
            暂无可用数据
          </div>

          <div v-else-if="stats" class="status-list">
            <div v-for="row in statusRows" :key="row.key" class="status-row">
              <div class="status-row-heading">
                <span class="status-row-label">
                  <span
                    class="status-row-dot"
                    :style="{ backgroundColor: row.color }"
                  ></span>
                  {{ row.label }}
                </span>

                <span class="status-row-value">{{ row.value }}</span>
              </div>

              <div class="progress-track">
                <div
                  class="progress-fill"
                  :style="{
                    width: `${row.percent}%`,
                    backgroundColor: row.color,
                  }"
                ></div>
              </div>
            </div>
          </div>
        </article>

        <article class="panel lifecycle-panel">
          <div class="panel-heading">
            <div>
              <h2>平台业务流程</h2>
              <p>从生产制造到设备运行</p>
            </div>

            <span class="panel-symbol">
              <el-icon><Connection /></el-icon>
            </span>
          </div>

          <div class="lifecycle-list">
            <div class="lifecycle-item">
              <div class="lifecycle-number">01</div>
              <div class="lifecycle-copy">
                <strong>生产制造</strong>
                <span>MES · 生产过程中建立设备身份</span>
              </div>
              <span class="lifecycle-marker"></span>
            </div>

            <div class="lifecycle-connector"></div>

            <div class="lifecycle-item">
              <div class="lifecycle-number">02</div>
              <div class="lifecycle-copy">
                <strong>仓储与交付</strong>
                <span>WMS · 记录设备库存和交付状态</span>
              </div>
              <span class="lifecycle-marker"></span>
            </div>

            <div class="lifecycle-connector"></div>

            <div class="lifecycle-item">
              <div class="lifecycle-number">03</div>
              <div class="lifecycle-copy">
                <strong>设备连接</strong>
                <span>IoT · 设备激活并接入工业云</span>
              </div>
              <span class="lifecycle-marker"></span>
            </div>

            <div class="lifecycle-connector"></div>

            <div class="lifecycle-item">
              <div class="lifecycle-number">04</div>
              <div class="lifecycle-copy">
                <strong>数字身份</strong>
                <span>统一追踪设备身份与运行数据</span>
              </div>
              <span class="lifecycle-marker"></span>
            </div>
          </div>
        </article>
      </section>

      <footer class="page-footer">
        <span>OpenIndustrial Industrial Cloud</span>
        <span>设备数字身份 · 全生命周期管理</span>
      </footer>
    </main>
  </div>
</template>

<style scoped>
.dashboard-page {
  min-height: 100vh;
  background: #f5f7fa;
  color: #101828;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 76px;
  padding: 0 40px;
  background: #0e1f33;
  color: #fff;
  box-shadow: 0 2px 12px rgb(14 31 51 / 12%);
}

.brand,
.topbar-right,
.account-info {
  display: flex;
  align-items: center;
}

.brand {
  gap: 13px;
}

.brand-mark {
  display: grid;
  width: 40px;
  height: 40px;
  place-items: center;
  border: 1px solid rgb(240 160 48 / 70%);
  border-radius: 9px;
  background: rgb(240 160 48 / 10%);
}

.brand-mark-inner {
  color: #f0a030;
  font-size: 15px;
  font-weight: 800;
  letter-spacing: -1px;
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.brand-subtitle {
  margin-top: 3px;
  color: #aab8c8;
  font-size: 12px;
}

.topbar-right {
  gap: 22px;
}

.system-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #c8d4e1;
  font-size: 12px;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #4bc58a;
  box-shadow: 0 0 0 4px rgb(75 197 138 / 12%);
}

.topbar-divider {
  width: 1px;
  height: 30px;
  background: rgb(255 255 255 / 15%);
}

.account-info {
  gap: 10px;
}

.account-avatar {
  display: grid;
  width: 35px;
  height: 35px;
  place-items: center;
  border: 1px solid rgb(240 160 48 / 45%);
  border-radius: 50%;
  background: #20364d;
  color: #f0a030;
  font-size: 14px;
  font-weight: 700;
}

.account-copy {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.account-name {
  max-width: 180px;
  overflow: hidden;
  color: #f5f7fa;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.account-role {
  color: #aab8c8;
  font-size: 11px;
}

.logout-button {
  color: #d1d9e2;
}

.logout-button:hover {
  color: #f0a030;
}

.main-content {
  width: min(1440px, 100%);
  margin: 0 auto;
  padding: 34px 40px 24px;
}

.welcome-section {
  position: relative;
  display: flex;
  min-height: 175px;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  padding: 30px 36px;
  border: 1px solid #e4e9f0;
  border-radius: 12px;
  background: linear-gradient(110deg, #fff 0%, #fff 58%, #f0f4f8 100%);
}

.eyebrow {
  display: flex;
  align-items: center;
  gap: 9px;
  color: #8a6a36;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 2px;
}

.eyebrow-line {
  width: 22px;
  height: 2px;
  background: #f0a030;
}

.welcome-section h1 {
  margin: 13px 0 9px;
  color: #0e1f33;
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.welcome-description {
  margin: 0;
  color: #667085;
  font-size: 13px;
  line-height: 1.7;
}

.welcome-decoration {
  position: relative;
  width: 190px;
  height: 145px;
  margin-right: 35px;
  flex: 0 0 auto;
}

.decoration-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  border: 1px solid rgb(14 31 51 / 13%);
  border-radius: 50%;
  transform: translate(-50%, -50%);
}

.ring-one {
  width: 100px;
  height: 100px;
}

.ring-two {
  width: 145px;
  height: 145px;
  border-style: dashed;
}

.decoration-center {
  position: absolute;
  top: 50%;
  left: 50%;
  display: grid;
  width: 54px;
  height: 54px;
  place-items: center;
  border: 1px solid rgb(240 160 48 / 55%);
  border-radius: 13px;
  background: #0e1f33;
  color: #f0a030;
  font-size: 18px;
  font-weight: 800;
  transform: translate(-50%, -50%);
}

.decoration-node {
  position: absolute;
  width: 9px;
  height: 9px;
  border: 2px solid #fff;
  border-radius: 50%;
  background: #f0a030;
  box-shadow: 0 0 0 3px rgb(240 160 48 / 16%);
}

.node-one {
  top: 14px;
  left: 94px;
}

.node-two {
  right: 16px;
  bottom: 37px;
  background: #426b92;
  box-shadow: 0 0 0 3px rgb(66 107 146 / 15%);
}

.node-three {
  bottom: 27px;
  left: 29px;
}

.stats-section {
  margin-top: 32px;
}

.section-heading,
.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-heading {
  margin-bottom: 17px;
}

.section-heading h2,
.panel-heading h2 {
  margin: 0;
  color: #182b40;
  font-size: 16px;
  font-weight: 700;
}

.section-heading p,
.panel-heading p {
  margin: 6px 0 0;
  color: #8a94a3;
  font-size: 12px;
}

.refresh-button {
  border-color: #d8e0e9;
  color: #344b63;
  background: #fff;
}

.refresh-button:hover {
  border-color: #f0a030;
  color: #9b641a;
  background: #fffaf2;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 16px;
}

.stat-card {
  min-width: 0;
  padding: 20px;
  border: 1px solid #e5eaf0;
  border-radius: 10px;
  background: #fff;
  transition:
    border-color 160ms ease,
    box-shadow 160ms ease,
    transform 160ms ease;
}

.stat-card:hover {
  border-color: #d5dfea;
  box-shadow: 0 7px 22px rgb(14 31 51 / 5%);
  transform: translateY(-2px);
}

.stat-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.stat-label {
  overflow: hidden;
  color: #667085;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stat-icon {
  display: grid;
  width: 32px;
  height: 32px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 700;
}

.tone-navy {
  background: #eaf0f6;
  color: #0e1f33;
}

.tone-gold {
  background: #fff4df;
  color: #ad721b;
}

.tone-blue {
  background: #eaf3ff;
  color: #376da6;
}

.tone-green {
  background: #e8f7ef;
  color: #208052;
}

.tone-gray {
  background: #f0f2f5;
  color: #737e8c;
}

.stat-value {
  min-height: 46px;
  margin-top: 15px;
  color: #0e1f33;
  font-size: 32px;
  font-weight: 750;
  line-height: 1.4;
  font-variant-numeric: tabular-nums;
}

.stat-description {
  color: #98a2b3;
  font-size: 11px;
  line-height: 1.6;
}

.skeleton-value {
  display: inline-block;
  width: 65px;
  height: 27px;
  border-radius: 5px;
  background: linear-gradient(90deg, #edf0f4 25%, #f7f8fa 50%, #edf0f4 75%);
  background-size: 200% 100%;
  animation: skeleton 1.4s ease infinite;
}

@keyframes skeleton {
  to {
    background-position: -200% 0;
  }
}

.error-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 16px;
  padding: 16px 18px;
  border: 1px solid #f3d3cd;
  border-radius: 9px;
  background: #fff8f6;
}

.error-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.error-icon {
  margin-top: 2px;
  color: #c24132;
  font-size: 20px;
}

.error-content strong {
  color: #923c32;
  font-size: 13px;
}

.error-content p {
  margin: 5px 0 0;
  color: #9d5b53;
  font-size: 12px;
  line-height: 1.6;
}

.details-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 18px;
  margin-top: 24px;
}

.panel {
  min-width: 0;
  padding: 24px;
  border: 1px solid #e5eaf0;
  border-radius: 10px;
  background: #fff;
}

.panel-heading {
  margin-bottom: 25px;
}

.panel-symbol {
  display: grid;
  width: 36px;
  height: 36px;
  place-items: center;
  border-radius: 9px;
  background: #edf2f7;
  color: #29445f;
  font-size: 17px;
}

.status-list {
  display: flex;
  flex-direction: column;
  gap: 22px;
}

.status-row-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 9px;
}

.status-row-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #475467;
  font-size: 12px;
}

.status-row-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.status-row-value {
  color: #182b40;
  font-size: 12px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}

.progress-track {
  height: 6px;
  overflow: hidden;
  border-radius: 6px;
  background: #eef1f5;
}

.progress-fill {
  height: 100%;
  min-width: 0;
  border-radius: inherit;
  transition: width 300ms ease;
}

.panel-loading,
.panel-empty {
  display: grid;
  min-height: 180px;
  place-items: center;
  color: #98a2b3;
  font-size: 13px;
}

.lifecycle-list {
  display: flex;
  flex-direction: column;
}

.lifecycle-item {
  display: flex;
  align-items: center;
  gap: 13px;
  min-height: 51px;
}

.lifecycle-number {
  display: grid;
  width: 34px;
  height: 34px;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid #e2e8ef;
  border-radius: 9px;
  background: #f7f9fb;
  color: #49627b;
  font-size: 11px;
  font-weight: 700;
}

.lifecycle-copy {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 5px;
}

.lifecycle-copy strong {
  color: #263c52;
  font-size: 12px;
  font-weight: 650;
}

.lifecycle-copy span {
  color: #8a94a3;
  font-size: 11px;
  line-height: 1.5;
}

.lifecycle-marker {
  width: 7px;
  height: 7px;
  flex: 0 0 auto;
  border-radius: 50%;
  background: #f0a030;
}

.lifecycle-connector {
  width: 1px;
  height: 13px;
  margin-left: 16px;
  background: #dce4ec;
}

.page-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 27px;
  padding: 18px 2px 0;
  border-top: 1px solid #e5eaf0;
  color: #98a2b3;
  font-size: 11px;
}

@media (max-width: 1200px) {
  .main-content {
    padding-right: 28px;
    padding-left: 28px;
  }

  .topbar {
    padding-right: 28px;
    padding-left: 28px;
  }

  .stats-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .topbar {
    min-height: 68px;
    padding: 0 18px;
  }

  .brand-mark {
    width: 36px;
    height: 36px;
  }

  .brand-name {
    font-size: 14px;
  }

  .brand-subtitle {
    font-size: 11px;
  }

  .topbar-right {
    gap: 12px;
  }

  .system-status,
  .topbar-divider,
  .account-copy {
    display: none;
  }

  .account-info {
    gap: 0;
  }

  .logout-button {
    padding: 7px;
  }

  .logout-button span {
    display: none;
  }

  .main-content {
    padding: 22px 16px;
  }

  .welcome-section {
    min-height: 150px;
    padding: 25px 22px;
  }

  .welcome-section h1 {
    font-size: 24px;
  }

  .welcome-decoration {
    width: 100px;
    height: 100px;
    margin-right: -35px;
    opacity: 0.6;
  }

  .ring-two {
    width: 95px;
    height: 95px;
  }

  .ring-one {
    width: 65px;
    height: 65px;
  }

  .decoration-center {
    width: 42px;
    height: 42px;
    font-size: 14px;
  }

  .node-one {
    top: 4px;
    left: 49px;
  }

  .node-two {
    right: 1px;
    bottom: 20px;
  }

  .node-three {
    bottom: 15px;
    left: 7px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 11px;
  }

  .stat-card {
    padding: 16px;
  }

  .stat-value {
    font-size: 28px;
  }

  .details-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .panel {
    padding: 20px;
  }

  .page-footer {
    flex-direction: column;
    align-items: flex-start;
  }
}

@media (max-width: 380px) {
  .brand {
    gap: 8px;
  }

  .brand-name {
    font-size: 12px;
  }

  .brand-subtitle {
    font-size: 10px;
  }

  .welcome-decoration {
    display: none;
  }

  .stat-card {
    padding: 13px;
  }

  .stat-label {
    font-size: 11px;
  }

  .stat-description {
    font-size: 10px;
  }
}
</style>
