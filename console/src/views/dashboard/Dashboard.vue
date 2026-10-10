<template>
  <div class="dashboard">
    <!-- 顶部平台导航 -->
    <header class="topbar">
      <div class="topbar-inner">
        <div class="brand">
          <div class="brand-mark">
            <svg viewBox="0 0 64 64" aria-hidden="true">
              <path
                d="M8 52V12H17L32 31L47 12H56V52H46V28L32 46L18 28V52Z"
                fill="#F0A030"
              />
              <path
                d="M24 52H40"
                stroke="#F0A030"
                stroke-width="3"
                stroke-linecap="round"
              />
            </svg>
          </div>

          <div class="brand-copy">
            <div class="brand-title">OpenIndustrial</div>
            <div class="brand-subtitle">设备数字身份平台</div>
          </div>
        </div>

        <div class="topbar-right">
          <div class="user-summary">
            <div class="user-avatar">{{ avatarText }}</div>
            <div class="user-info">
              <span class="user-name">{{ displayName }}</span>
              <span class="user-role">{{ roleLabel }}</span>
            </div>
          </div>

          <el-button
            class="logout-button"
            :icon="SwitchButton"
            text
            @click="handleLogout"
          >
            退出登录
          </el-button>
        </div>
      </div>
    </header>

    <main class="dashboard-main">
      <!-- 欢迎区域 -->
      <section class="welcome-section">
        <div class="welcome-content">
          <div class="eyebrow">
            <span class="eyebrow-line"></span>
            INDUSTRIAL CLOUD PLATFORM
          </div>

          <h1>
            欢迎回来，<span class="welcome-name">{{ displayName }}</span>
          </h1>

          <p class="welcome-description">
            连接工厂、生产、仓储与设备，让工业数据与设备数字身份贯穿产品全生命周期。
          </p>

          <div class="welcome-meta">
            <span class="status-dot"></span>
            <span>平台运行中</span>
            <span class="meta-divider"></span>
            <span>{{ roleLabel }}</span>
          </div>
        </div>

        <div class="hero-visual" aria-hidden="true">
          <div class="hero-orbit orbit-one"></div>
          <div class="hero-orbit orbit-two"></div>
          <div class="hero-orbit orbit-three"></div>

          <svg class="robot-illustration" viewBox="0 0 300 260">
            <defs>
              <linearGradient id="robotGold" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stop-color="#FFD18A" />
                <stop offset="100%" stop-color="#F0A030" />
              </linearGradient>
              <linearGradient id="robotBlue" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stop-color="#344D67" />
                <stop offset="100%" stop-color="#172C43" />
              </linearGradient>
            </defs>

            <g fill="none" stroke="#7C91A6" stroke-width="1" opacity=".35">
              <path d="M35 210H265" />
              <path d="M50 225H250" />
              <path d="M70 240H230" />
              <path d="M75 190L150 150L225 190" />
              <path d="M75 190V230M150 150V230M225 190V230" />
            </g>

            <g class="robot-arm">
              <path
                d="M115 205L115 165L143 135L188 107"
                stroke="#233B54"
                stroke-width="24"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M115 205L115 165L143 135L188 107"
                stroke="url(#robotGold)"
                stroke-width="15"
                stroke-linecap="round"
                stroke-linejoin="round"
              />

              <path
                d="M188 107L214 121L235 101"
                stroke="#233B54"
                stroke-width="18"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M188 107L214 121L235 101"
                stroke="url(#robotGold)"
                stroke-width="10"
                stroke-linecap="round"
                stroke-linejoin="round"
              />

              <circle
                cx="115"
                cy="165"
                r="19"
                fill="url(#robotBlue)"
                stroke="#F0A030"
                stroke-width="3"
              />
              <circle cx="115" cy="165" r="7" fill="#F0A030" />

              <circle
                cx="143"
                cy="135"
                r="17"
                fill="url(#robotBlue)"
                stroke="#F0A030"
                stroke-width="3"
              />
              <circle cx="143" cy="135" r="6" fill="#F0A030" />

              <circle
                cx="188"
                cy="107"
                r="16"
                fill="url(#robotBlue)"
                stroke="#F0A030"
                stroke-width="3"
              />
              <circle cx="188" cy="107" r="6" fill="#F0A030" />

              <path
                d="M228 95L241 82M237 105L254 105M226 113L238 126"
                stroke="#F0A030"
                stroke-width="4"
                stroke-linecap="round"
              />
            </g>

            <g>
              <path
                d="M65 205H164L153 226H76Z"
                fill="#172C43"
                stroke="#425B74"
                stroke-width="2"
              />
              <rect
                x="86"
                y="188"
                width="58"
                height="19"
                rx="4"
                fill="#243D56"
              />
              <rect
                x="92"
                y="193"
                width="46"
                height="4"
                rx="2"
                fill="#F0A030"
                opacity=".9"
              />
            </g>

            <g fill="#F0A030">
              <circle cx="61" cy="70" r="3" />
              <circle cx="245" cy="61" r="4" />
              <circle cx="254" cy="163" r="3" />
              <circle cx="51" cy="151" r="2" />
            </g>

            <g fill="none" stroke="#F0A030" stroke-width="1.5" opacity=".75">
              <path d="M43 83H72V97" />
              <path d="M220 47H246V57" />
              <path d="M244 178H265V191" />
            </g>
          </svg>

          <div class="scan-line"></div>
          <div class="visual-caption">
            <span class="caption-dot"></span>
            DIGITAL IDENTITY · CONNECTED
          </div>
        </div>

        <div class="hero-decoration"></div>
      </section>

      <!-- 管理员账户统计 -->
      <section v-if="isAdmin" class="stats-section">
        <div class="section-heading">
          <div>
            <div class="section-kicker">OVERVIEW</div>
            <h2>账户概览</h2>
            <p class="section-description">
              查看当前工厂的用户账户状态与分布。
            </p>
          </div>

          <el-button
            class="refresh-button"
            :icon="Refresh"
            :loading="loading"
            @click="loadOverview"
          >
            刷新数据
          </el-button>
        </div>

        <div v-if="loadError" class="error-banner">
          <el-icon><WarningFilled /></el-icon>
          <span>{{ loadError }}</span>
          <el-button link @click="loadOverview">重试</el-button>
        </div>

        <div class="stats-grid">
          <button
            v-for="card in statCards"
            :key="card.key"
            type="button"
            class="stat-card"
            :class="[
              `stat-card--${card.key}`,
              { 'stat-card--loading': loading },
            ]"
            @click="openStatCard(card.key)"
          >
            <div class="stat-card-top">
              <span class="stat-label">{{ card.label }}</span>
              <span
                class="stat-indicator"
                :class="`indicator--${card.key}`"
              ></span>
            </div>

            <div v-if="loading" class="stat-skeleton"></div>
            <div v-else class="stat-value">
              {{ card.value ?? "—" }}
            </div>

            <div class="stat-card-bottom">
              <span>{{ card.description }}</span>
              <span class="stat-arrow">↗</span>
            </div>
          </button>
        </div>

        <div class="distribution-panel">
          <div class="distribution-header">
            <div>
              <h3>账户状态分布</h3>
              <p>各状态账户数量及占比</p>
            </div>
            <div class="distribution-total">
              <span>总账户</span>
              <strong>{{ loading ? "—" : (stats?.total ?? 0) }}</strong>
            </div>
          </div>

          <div v-if="loading" class="distribution-loading">
            <div
              v-for="index in 4"
              :key="index"
              class="distribution-skeleton"
            ></div>
          </div>

          <template v-else>
            <div class="distribution-bar" aria-label="账户状态比例">
              <div
                v-for="item in statusRows"
                :key="item.key"
                class="distribution-segment"
                :class="`segment--${item.key}`"
                :style="{ width: `${item.percentage}%` }"
                :title="`${item.label}: ${item.value}`"
              ></div>
            </div>

            <div class="distribution-legend">
              <div
                v-for="item in statusRows"
                :key="item.key"
                class="legend-item"
              >
                <span class="legend-dot" :class="`dot--${item.key}`"></span>
                <span class="legend-label">{{ item.label }}</span>
                <strong>{{ item.value }}</strong>
                <span class="legend-percentage">{{ item.percentage }}%</span>
              </div>
            </div>
          </template>
        </div>
      </section>

      <!-- 第一层：四大业务模块 -->
      <section class="modules-section">
        <div class="section-heading">
          <div>
            <div class="section-kicker">BUSINESS MODULES</div>
            <h2>业务中心</h2>
            <p class="section-description">
              覆盖客户与订单、生产制造、仓储交付和设备运行的完整工业业务流程。
            </p>
          </div>
        </div>

        <div class="module-grid business-module-grid">
          <!-- CRM 客户关系与订单管理 -->
          <div class="module-card module-card--crm">
            <div class="module-card-heading">
              <div class="module-symbol crm-symbol">
                <svg viewBox="0 0 48 48" aria-hidden="true">
                  <circle
                    cx="18"
                    cy="15"
                    r="7"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                  <path
                    d="M5 38V34C5 28.5 10 25 18 25C26 25 31 28.5 31 34V38"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                  <path
                    d="M31 12H42V29H35L29 34V29H27"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linejoin="round"
                  />
                  <path
                    d="M33 18H38M33 23H38"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                </svg>
              </div>
              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">CUSTOMER RELATIONSHIP MANAGEMENT</div>
            <h3>客户与订单 CRM</h3>
            <p class="module-description">
              统一管理客户资料与业务订单，为工厂生产计划和后续交付提供业务来源。
            </p>

            <div class="module-card-footer">
              <span class="module-tag">客户与订单</span>
              <span class="module-status">建设中</span>
            </div>
          </div>

          <!-- MES 制造执行 -->
          <button
            v-if="canAccessProducts"
            type="button"
            class="module-card module-card--mes"
            @click="openProducts"
          >
            <div class="module-card-heading">
              <div class="mes-logo" aria-hidden="true">
                <svg viewBox="0 0 64 64">
                  <path
                    d="M8 52V12H17L32 31L47 12H56V52H46V28L32 46L18 28V52Z"
                    fill="#F0A030"
                  />
                  <path
                    d="M24 52H40"
                    stroke="#F0A030"
                    stroke-width="3"
                    stroke-linecap="round"
                  />
                </svg>
              </div>

              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">MANUFACTURING EXECUTION</div>
            <h3>制造执行 MES</h3>
            <p class="module-description">
              统一管理产品、物料清单与生产制造基础数据，为生产执行和设备身份生成提供支撑。
            </p>

            <div class="module-card-footer">
              <span class="mes-tag">生产制造</span>
              <span class="module-enter">进入模块 <span>→</span></span>
            </div>
          </button>

          <!-- WMS 仓储管理 -->
          <div class="module-card module-card--wms">
            <div class="module-card-heading">
              <div class="module-symbol wms-symbol">
                <svg viewBox="0 0 48 48" aria-hidden="true">
                  <path
                    d="M5 18L24 7L43 18V41H5Z"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linejoin="round"
                  />
                  <path
                    d="M5 18H43M15 18V41M33 18V41M15 29H33"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                </svg>
              </div>
              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">WAREHOUSE MANAGEMENT</div>
            <h3>仓储管理 WMS</h3>
            <p class="module-description">
              管理成品入库、库存、预留与发运流程，让生产结果与设备交付过程保持一致。
            </p>

            <div class="module-card-footer">
              <span class="module-tag">仓储物流</span>
              <span class="module-status">业务模块</span>
            </div>
          </div>

          <!-- IoT 设备连接 -->
          <div class="module-card module-card--iot">
            <div class="module-card-heading">
              <div class="module-symbol iot-symbol">
                <svg viewBox="0 0 48 48" aria-hidden="true">
                  <rect
                    x="17"
                    y="17"
                    width="14"
                    height="14"
                    rx="3"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                  <path
                    d="M24 5V12M24 36V43M5 24H12M36 24H43M10 10L15 15M33 33L38 38M38 10L33 15M15 33L10 38"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                  <circle cx="24" cy="24" r="3" fill="#F0A030" />
                </svg>
              </div>
              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">INTERNET OF THINGS</div>
            <h3>物联网 IoT</h3>
            <p class="module-description">
              通过统一的设备接入与协议适配能力，连接现场设备，逐步建立设备状态和运行数据闭环。
            </p>

            <div class="module-card-footer">
              <span class="module-tag">设备连接</span>
              <span class="module-status">业务模块</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 第二层：共享平台能力 -->
      <section class="shared-capabilities-section">
        <div class="section-heading">
          <div>
            <div class="section-kicker">SHARED PLATFORM CAPABILITIES</div>
            <h2>共享平台能力</h2>
            <p class="section-description">
              为 CRM、MES、WMS 和 IoT 提供统一的身份、权限与安全基础设施。
            </p>
          </div>
        </div>

        <div class="module-grid shared-capabilities-grid">
          <!-- 设备数字身份 -->
          <div class="module-card module-card--identity">
            <div class="module-card-heading">
              <div class="module-symbol identity-symbol">
                <svg viewBox="0 0 48 48" aria-hidden="true">
                  <rect
                    x="9"
                    y="5"
                    width="30"
                    height="38"
                    rx="5"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                  <path
                    d="M17 15H31M17 21H31M17 27H25"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                  <circle cx="31" cy="32" r="5" fill="#F0A030" />
                  <path
                    d="M29 32L31 34L34 30"
                    fill="none"
                    stroke="#0E1F33"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </div>
              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">DIGITAL IDENTITY</div>
            <h3>设备数字身份</h3>
            <p class="module-description">
              以唯一设备身份连接生产履历、设备证书与运行数据，建立可信的设备生命周期档案。
            </p>

            <div class="module-card-footer">
              <span class="module-tag">设备身份</span>
              <span class="module-status">平台核心能力</span>
            </div>
          </div>

          <!-- 用户与权限 -->
          <div class="module-card module-card--users">
            <div class="module-card-heading">
              <div class="module-symbol users-symbol">
                <svg viewBox="0 0 48 48" aria-hidden="true">
                  <circle
                    cx="18"
                    cy="15"
                    r="7"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                  <path
                    d="M5 39V35C5 29 10 25 18 25C23 25 27 27 29 31"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                  <path
                    d="M32 25L40 28V34C40 39 36 42 32 44C28 42 24 39 24 34V28Z"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linejoin="round"
                  />
                  <path
                    d="M29 34L31 36L35 32"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </div>
              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">USERS &amp; PERMISSIONS</div>
            <h3>用户与权限</h3>
            <p class="module-description">
              统一管理工厂用户、角色和业务访问权限，确保员工根据职责访问对应功能。
            </p>

            <div class="module-card-footer">
              <span class="module-tag">账号与角色</span>
              <span class="module-status">平台核心能力</span>
            </div>
          </div>

          <!-- 安全与认证 -->
          <div class="module-card module-card--security">
            <div class="module-card-heading">
              <div class="module-symbol security-symbol">
                <svg viewBox="0 0 48 48" aria-hidden="true">
                  <path
                    d="M24 5L39 11V22C39 32 32 39 24 43C16 39 9 32 9 22V11Z"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linejoin="round"
                  />
                  <rect
                    x="17"
                    y="21"
                    width="14"
                    height="11"
                    rx="2"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                  <path
                    d="M20 21V17C20 12 28 12 28 17V21"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                  <circle cx="24" cy="26" r="1.5" fill="#F0A030" />
                </svg>
              </div>
              <span class="module-arrow">↗</span>
            </div>

            <div class="module-category">SECURITY &amp; AUTHENTICATION</div>
            <h3>安全与认证</h3>
            <p class="module-description">
              通过统一身份认证、访问控制与设备证书机制，为平台用户和工业设备提供可信安全保障。
            </p>

            <div class="module-card-footer">
              <span class="module-tag">身份认证</span>
              <span class="module-status">平台核心能力</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 非管理员工作区 -->
      <section
        v-if="!isAdmin && isKnownRole(authStore.user)"
        class="workspace-section"
      >
        <div class="workspace-panel">
          <div class="workspace-icon">
            <svg viewBox="0 0 48 48" aria-hidden="true">
              <rect
                x="6"
                y="8"
                width="36"
                height="32"
                rx="5"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              />
              <path
                d="M6 17H42M16 8V17M32 8V17"
                stroke="currentColor"
                stroke-width="2"
              />
              <path
                d="M15 26H21M27 26H33M15 32H21M27 32H33"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              />
            </svg>
          </div>
          <div class="workspace-content">
            <div class="section-kicker">MY WORKSPACE</div>
            <h2>我的工作台</h2>
            <p>
              当前以{{
                roleLabel
              }}身份登录。可访问的业务模块将根据账户权限开放。
            </p>
          </div>
        </div>
      </section>

      <div class="page-bottom-space"></div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { Refresh, SwitchButton, WarningFilled } from "@element-plus/icons-vue";

import { getDashboardOverview } from "../../api/dashboard";
import type { UserStats } from "../../api/dashboard";
import { useAuthStore } from "../../store/auth";
import { hasPermission, isKnownRole } from "../../authz/access";
import { Permissions } from "../../authz/permissions";
import { getRoleLabel } from "../../authz/roles";

type StatKey = "total" | "init" | "invited" | "active" | "disabled";

interface StatCard {
  key: StatKey;
  label: string;
  value: number | undefined;
  description: string;
}

const router = useRouter();
const authStore = useAuthStore();

const stats = ref<UserStats | null>(null);
const loading = ref(false);
const loadError = ref("");

const isAdmin = computed(() =>
  hasPermission(authStore.user, Permissions.USER_LIST),
);

const canAccessProducts = computed(() =>
  hasPermission(authStore.user, Permissions.PRODUCT_LIST),
);

const displayName = computed(() => {
  const user = authStore.user as Record<string, unknown> | null;

  if (!user) return "用户";

  const candidates = [
    user.displayName,
    user.display_name,
    user.name,
    user.email,
  ];

  const value = candidates.find(
    (item) => typeof item === "string" && item.trim().length > 0,
  );

  if (typeof value !== "string") return "用户";

  if (value.includes("@")) {
    return value.split("@")[0] || "用户";
  }

  return value;
});

const avatarText = computed(() => {
  const name = displayName.value.trim();
  return name ? Array.from(name)[0].toUpperCase() : "U";
});

const roleLabel = computed(() =>
  getRoleLabel(authStore.user?.user_type ?? null),
);

const statCards = computed<StatCard[]>(() => [
  {
    key: "total",
    label: "全部账户",
    value: stats.value?.total,
    description: "当前工厂用户总数",
  },
  {
    key: "init",
    label: "待处理",
    value: stats.value?.init,
    description: "尚未完成邀请流程",
  },
  {
    key: "invited",
    label: "待激活",
    value: stats.value?.invited,
    description: "已邀请，等待激活",
  },
  {
    key: "active",
    label: "正常账户",
    value: stats.value?.active,
    description: "可正常登录使用",
  },
  {
    key: "disabled",
    label: "已停用",
    value: stats.value?.disabled,
    description: "当前不可登录",
  },
]);

const statusRows = computed(() => {
  const total = stats.value?.total ?? 0;

  return [
    {
      key: "init",
      label: "待处理",
      value: stats.value?.init ?? 0,
      percentage: getPercentage(stats.value?.init ?? 0, total),
    },
    {
      key: "invited",
      label: "待激活",
      value: stats.value?.invited ?? 0,
      percentage: getPercentage(stats.value?.invited ?? 0, total),
    },
    {
      key: "active",
      label: "正常账户",
      value: stats.value?.active ?? 0,
      percentage: getPercentage(stats.value?.active ?? 0, total),
    },
    {
      key: "disabled",
      label: "已停用",
      value: stats.value?.disabled ?? 0,
      percentage: getPercentage(stats.value?.disabled ?? 0, total),
    },
  ];
});

function getPercentage(value: number, total: number): number {
  if (total <= 0) return 0;
  return Math.round((value / total) * 1000) / 10;
}

function openStatCard(key: StatKey) {
  if (!isAdmin.value) return;

  if (key === "total") {
    void router.push("/users");
    return;
  }

  void router.push({
    path: "/users",
    query: { status: key },
  });
}

function openProducts() {
  void router.push("/mes");
}

async function loadOverview() {
  if (!isAdmin.value) return;

  loading.value = true;
  loadError.value = "";

  try {
    const overview = await getDashboardOverview();
    stats.value = overview.users;
  } catch (error: unknown) {
    console.error("Failed to load dashboard overview:", error);
    loadError.value = "账户统计加载失败，请检查网络或稍后重试。";
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
  if (isAdmin.value) {
    void loadOverview();
  }
});
</script>

<style scoped>
.dashboard {
  --navy: #0e1f33;
  --navy-light: #16293f;
  --navy-soft: #223b55;
  --gold: #f0a030;
  --gold-light: #ffd18a;
  --text-primary: #17283b;
  --text-secondary: #637488;
  --border: #e5eaf0;
  --surface: #ffffff;
  --page-bg: #f4f6f9;

  min-height: 100vh;
  color: var(--text-primary);
  background: var(--page-bg);
}

/* 顶部导航 */
.topbar {
  position: sticky;
  top: 0;
  z-index: 20;
  height: 76px;
  background: var(--navy);
  border-bottom: 1px solid rgb(255 255 255 / 8%);
}

.topbar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: min(1440px, calc(100% - 64px));
  height: 100%;
  margin: 0 auto;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.brand-mark {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  flex: 0 0 42px;
  overflow: hidden;
  border: 1px solid rgb(240 160 48 / 28%);
  border-radius: 10px;
  background: rgb(255 255 255 / 4%);
}

.brand-mark svg {
  width: 30px;
  height: 30px;
}

.brand-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.brand-title {
  color: #fff;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.brand-subtitle {
  color: #a9b7c6;
  font-size: 12px;
  letter-spacing: 1px;
}

.topbar-right,
.user-summary {
  display: flex;
  align-items: center;
}

.topbar-right {
  gap: 24px;
}

.user-summary {
  gap: 10px;
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  flex: 0 0 36px;
  border: 1px solid rgb(240 160 48 / 50%);
  border-radius: 50%;
  color: var(--gold-light);
  background: rgb(240 160 48 / 12%);
  font-size: 14px;
  font-weight: 700;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  max-width: 180px;
  overflow: hidden;
  color: #f5f7fa;
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-role {
  color: #9babbc;
  font-size: 11px;
}

.logout-button {
  color: #c6d0db !important;
  font-size: 12px;
}

.logout-button:hover {
  color: var(--gold-light) !important;
}

/* 页面主体 */
.dashboard-main {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
}

.welcome-section {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 300px;
  margin-top: 28px;
  padding: 34px 52px;
  overflow: hidden;
  border: 1px solid rgb(255 255 255 / 8%);
  border-radius: 18px;
  color: #fff;
  background:
    radial-gradient(ellipse at 80% 20%, rgb(48 76 103 / 52%), transparent 42%),
    linear-gradient(115deg, #0e1f33 0%, #16293f 62%, #203a53 100%);
  box-shadow: 0 12px 32px rgb(14 31 51 / 10%);
}

.welcome-content {
  position: relative;
  z-index: 2;
  width: 57%;
}

.eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  color: #a9b9c9;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 2.1px;
}

.eyebrow-line {
  width: 25px;
  height: 2px;
  background: var(--gold);
}

.welcome-section h1 {
  margin: 0;
  color: #f8fafc;
  font-size: clamp(25px, 3vw, 36px);
  font-weight: 600;
  letter-spacing: -0.7px;
  line-height: 1.35;
}

.welcome-name {
  color: var(--gold-light);
}

.welcome-description {
  max-width: 540px;
  margin: 17px 0 24px;
  color: #bac6d2;
  font-size: 13px;
  line-height: 1.9;
}

.welcome-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #b8c5d2;
  font-size: 11px;
}

.status-dot,
.caption-dot {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 50%;
  background: #65c6a0;
  box-shadow: 0 0 0 4px rgb(101 198 160 / 12%);
}

.meta-divider {
  width: 1px;
  height: 12px;
  margin: 0 3px;
  background: #516579;
}

.hero-visual {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40%;
  max-width: 430px;
  height: 250px;
}

.robot-illustration {
  position: relative;
  z-index: 2;
  width: min(100%, 330px);
  height: auto;
  overflow: visible;
}

.hero-orbit {
  position: absolute;
  top: 50%;
  left: 50%;
  border: 1px solid rgb(151 176 198 / 14%);
  border-radius: 50%;
  transform: translate(-50%, -50%);
}

.orbit-one {
  width: 150px;
  height: 150px;
}

.orbit-two {
  width: 215px;
  height: 215px;
}

.orbit-three {
  width: 280px;
  height: 280px;
}

.scan-line {
  position: absolute;
  z-index: 3;
  top: 20%;
  left: 15%;
  width: 70%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent,
    rgb(240 160 48 / 75%),
    transparent
  );
  opacity: 0.55;
  animation: scan 5s ease-in-out infinite alternate;
}

.visual-caption {
  position: absolute;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  gap: 9px;
  color: #9eafbf;
  font-size: 9px;
  letter-spacing: 1.1px;
}

.caption-dot {
  width: 5px;
  height: 5px;
  flex-basis: 5px;
}

.hero-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 230px;
  height: 100%;
  pointer-events: none;
  opacity: 0.14;
  background-image: radial-gradient(#a6bbce 0.7px, transparent 0.7px);
  background-size: 13px 13px;
  mask-image: linear-gradient(90deg, transparent, #000);
}

@keyframes scan {
  from {
    transform: translateY(0);
  }

  to {
    transform: translateY(145px);
  }
}

/* 统一标题 */
.stats-section,
.modules-section,
.shared-capabilities-section {
  margin-top: 42px;
}

.section-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 20px;
}

.section-kicker {
  margin-bottom: 7px;
  color: #a36b24;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1.8px;
}

.section-heading h2,
.workspace-content h2 {
  margin: 0;
  color: var(--text-primary);
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.4px;
}

.section-description {
  margin: 8px 0 0;
  color: var(--text-secondary);
  font-size: 12px;
  line-height: 1.7;
}

.refresh-button {
  height: 36px;
  border-color: #dce3eb;
  color: #40566c;
  background: #fff;
}

.refresh-button:hover {
  border-color: var(--gold);
  color: #9b671f;
  background: #fffaf2;
}

/* 账户统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 155px;
  padding: 21px 20px 17px;
  border: 1px solid var(--border);
  border-radius: 12px;
  text-align: left;
  background: #fff;
  box-shadow: 0 3px 12px rgb(26 43 61 / 2%);
  cursor: pointer;
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    box-shadow 180ms ease;
}

.stat-card:hover {
  border-color: #d5b27f;
  box-shadow: 0 9px 24px rgb(26 43 61 / 7%);
  transform: translateY(-3px);
}

.stat-card-top,
.stat-card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.stat-label {
  color: #66778a;
  font-size: 12px;
}

.stat-indicator {
  width: 8px;
  height: 8px;
  flex: 0 0 8px;
  border-radius: 50%;
  background: #9ca9b7;
}

.indicator--total {
  background: #354f6a;
}

.indicator--init {
  background: #9aa7b4;
}

.indicator--invited {
  background: #e3ad4c;
}

.indicator--active {
  background: #49a987;
}

.indicator--disabled {
  background: #d47b72;
}

.stat-value {
  margin: 17px 0 12px;
  color: #172b40;
  font-size: 32px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
  line-height: 1;
}

.stat-card-bottom {
  margin-top: auto;
  color: #91a0af;
  font-size: 10px;
}

.stat-arrow {
  color: #b0bac4;
  font-size: 16px;
  transition: color 180ms ease;
}

.stat-card:hover .stat-arrow {
  color: #b47b2e;
}

.stat-skeleton {
  width: 60px;
  height: 30px;
  margin: 17px 0 12px;
  border-radius: 5px;
  background: linear-gradient(90deg, #eef1f5 25%, #f7f8fa 50%, #eef1f5 75%);
  background-size: 200% 100%;
  animation: skeleton 1.5s infinite;
}

@keyframes skeleton {
  to {
    background-position: -200% 0;
  }
}

/* 状态分布 */
.distribution-panel {
  margin-top: 18px;
  padding: 24px 26px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: #fff;
}

.distribution-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 23px;
}

.distribution-header h3 {
  margin: 0;
  color: #24384d;
  font-size: 14px;
  font-weight: 650;
}

.distribution-header p {
  margin: 6px 0 0;
  color: #8a99a8;
  font-size: 11px;
}

.distribution-total {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #8796a5;
  font-size: 11px;
}

.distribution-total strong {
  color: #253a50;
  font-size: 23px;
  font-weight: 650;
}

.distribution-bar {
  display: flex;
  width: 100%;
  height: 8px;
  overflow: hidden;
  border-radius: 8px;
  background: #edf0f4;
}

.distribution-segment {
  min-width: 0;
  transition: width 300ms ease;
}

.segment--init {
  background: #9aa7b4;
}

.segment--invited {
  background: #e3ad4c;
}

.segment--active {
  background: #49a987;
}

.segment--disabled {
  background: #d47b72;
}

.distribution-legend {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 22px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.legend-dot {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 50%;
}

.dot--init {
  background: #9aa7b4;
}

.dot--invited {
  background: #e3ad4c;
}

.dot--active {
  background: #49a987;
}

.dot--disabled {
  background: #d47b72;
}

.legend-label {
  overflow: hidden;
  color: #728296;
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.legend-item strong {
  margin-left: auto;
  color: #283d52;
  font-size: 12px;
  font-weight: 650;
  font-variant-numeric: tabular-nums;
}

.legend-percentage {
  color: #9aa6b3;
  font-size: 10px;
  font-variant-numeric: tabular-nums;
}

.distribution-loading {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.distribution-skeleton {
  height: 12px;
  border-radius: 4px;
  background: #eef1f5;
  animation: skeleton 1.5s infinite;
}

/* 错误提示 */
.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding: 12px 16px;
  border: 1px solid #f0d4d0;
  border-radius: 8px;
  color: #a94e44;
  background: #fff8f7;
  font-size: 12px;
}

.error-banner .el-button {
  margin-left: auto;
  color: #a94e44;
}

/* 业务模块与共享平台能力 */
.module-grid {
  display: grid;
  gap: 18px;
}

/* 四大业务模块 */
.business-module-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

/* 下方共享平台能力 */
.shared-capabilities-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.module-card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 245px;
  padding: 25px 26px 21px;
  overflow: hidden;
  border: 1px solid var(--border);
  border-radius: 13px;
  text-align: left;
  color: var(--text-primary);
  background: #fff;
  box-shadow: 0 3px 12px rgb(26 43 61 / 2%);
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    box-shadow 180ms ease;
}

button.module-card {
  font: inherit;
  cursor: pointer;
}

.module-card:hover {
  border-color: #d8c3a3;
  box-shadow: 0 10px 26px rgb(26 43 61 / 6%);
  transform: translateY(-3px);
}

.module-card-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 21px;
}

.module-arrow {
  color: #9eabb8;
  font-size: 19px;
  transition:
    color 180ms ease,
    transform 180ms ease;
}

.module-card:hover .module-arrow {
  color: #b17b32;
  transform: translate(2px, -2px);
}

/* MES 专属 Logo：深海军蓝 + 金橙色 */
.mes-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 58px;
  height: 58px;
  flex: 0 0 58px;
  border: 1px solid rgb(240 160 48 / 25%);
  border-radius: 13px;
  background: #0e1f33;
  box-shadow: 0 5px 13px rgb(14 31 51 / 12%);
  transition:
    transform 180ms ease,
    box-shadow 180ms ease;
}

.mes-logo svg {
  display: block;
  width: 43px;
  height: 43px;
}

.module-card--mes:hover .mes-logo {
  box-shadow: 0 7px 18px rgb(14 31 51 / 20%);
  transform: translateY(-2px);
}

.module-category {
  margin-bottom: 8px;
  color: #8a99a8;
  font-size: 9px;
  font-weight: 650;
  letter-spacing: 1.5px;
}

.module-card h3 {
  margin: 0;
  color: #20354a;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.2px;
}

.module-description {
  max-width: 560px;
  margin: 12px 0 22px;
  color: #728296;
  font-size: 12px;
  line-height: 1.9;
}

.module-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid #edf0f4;
}

.mes-tag {
  padding: 5px 9px;
  border: 1px solid rgb(240 160 48 / 25%);
  border-radius: 5px;
  color: #9d6b27;
  background: #fff8eb;
  font-size: 10px;
}

.module-enter {
  color: #324b63;
  font-size: 11px;
  font-weight: 600;
}

.module-enter span {
  display: inline-block;
  margin-left: 5px;
  transition: transform 180ms ease;
}

.module-card--mes:hover .module-enter span {
  transform: translateX(3px);
}

.module-symbol {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 58px;
  height: 58px;
  flex: 0 0 58px;
  border-radius: 13px;
}

.module-symbol svg {
  width: 39px;
  height: 39px;
}

.crm-symbol {
  color: #9a743b;
  background: #fbf3e7;
}

.identity-symbol {
  color: #315572;
  background: #edf3f8;
}

.wms-symbol {
  color: #538a77;
  background: #edf6f1;
}

.iot-symbol {
  color: #637aa1;
  background: #eff2fa;
}

.users-symbol {
  color: #527c83;
  background: #edf5f5;
}

.security-symbol {
  color: #6d6b9a;
  background: #f1f0fa;
}

.module-tag {
  padding: 5px 9px;
  border-radius: 5px;
  color: #526b80;
  background: #f0f4f7;
  font-size: 10px;
}

.module-status {
  color: #98a4b0;
  font-size: 10px;
}

/* 普通员工工作区 */
.workspace-section {
  margin-top: 42px;
}

.workspace-panel {
  display: flex;
  align-items: center;
  gap: 22px;
  padding: 28px;
  border: 1px solid var(--border);
  border-radius: 13px;
  background: #fff;
}

.workspace-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  flex: 0 0 64px;
  border-radius: 13px;
  color: #35546e;
  background: #edf3f8;
}

.workspace-icon svg {
  width: 40px;
  height: 40px;
}

.workspace-content h2 {
  margin-bottom: 8px;
  font-size: 19px;
}

.workspace-content p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 12px;
  line-height: 1.8;
}

.page-bottom-space {
  height: 48px;
}

/* 响应式布局 */
@media (max-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .welcome-section {
    padding: 32px;
  }

  .hero-visual {
    width: 38%;
  }

  .distribution-legend {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .business-module-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .shared-capabilities-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .topbar {
    height: 68px;
  }

  .topbar-inner,
  .dashboard-main {
    width: calc(100% - 32px);
  }

  .topbar-right {
    gap: 10px;
  }

  .brand-mark {
    width: 36px;
    height: 36px;
    flex-basis: 36px;
  }

  .brand-mark svg {
    width: 26px;
    height: 26px;
  }

  .brand-title {
    font-size: 14px;
  }

  .brand-subtitle {
    font-size: 10px;
  }

  .user-info {
    display: none;
  }

  .logout-button {
    padding: 7px !important;
  }

  .welcome-section {
    min-height: 0;
    margin-top: 16px;
    padding: 28px 24px;
  }

  .welcome-content {
    width: 100%;
  }

  .welcome-section h1 {
    font-size: 25px;
  }

  .welcome-description {
    max-width: 100%;
    font-size: 12px;
  }

  .hero-visual {
    display: none;
  }

  .hero-decoration {
    width: 45%;
  }

  .stats-section,
  .modules-section,
  .shared-capabilities-section {
    margin-top: 30px;
  }

  .section-heading h2 {
    font-size: 20px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
  }

  .stat-card {
    min-height: 140px;
    padding: 17px 15px;
  }

  .stat-value {
    font-size: 28px;
  }

  .business-module-grid,
  .shared-capabilities-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .module-card {
    min-height: 230px;
    padding: 22px;
  }

  .distribution-panel {
    padding: 20px 16px;
  }

  .distribution-legend {
    gap: 16px 10px;
  }

  .legend-item {
    gap: 6px;
  }

  .legend-percentage {
    display: none;
  }
}

@media (max-width: 420px) {
  .brand {
    gap: 8px;
  }

  .brand-title {
    font-size: 13px;
  }

  .brand-subtitle {
    letter-spacing: 0;
  }

  .topbar-right {
    gap: 6px;
  }

  .welcome-section {
    padding: 24px 20px;
  }

  .welcome-section h1 {
    font-size: 22px;
  }

  .section-heading {
    align-items: flex-start;
  }

  .refresh-button {
    padding: 8px 10px;
    font-size: 11px;
  }

  .distribution-total {
    gap: 7px;
  }

  .distribution-total strong {
    font-size: 20px;
  }

  .workspace-panel {
    align-items: flex-start;
    gap: 15px;
    padding: 20px;
  }

  .workspace-icon {
    width: 48px;
    height: 48px;
    flex-basis: 48px;
  }

  .workspace-icon svg {
    width: 30px;
    height: 30px;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    scroll-behavior: auto !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
