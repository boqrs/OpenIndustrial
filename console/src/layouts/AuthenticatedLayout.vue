<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { SwitchButton } from "@element-plus/icons-vue";

import { useAuthStore } from "../store/auth";
import { getRoleLabel } from "../authz/roles";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const displayName = computed(() => {
  const user = authStore.user;

  if (!user) return "用户";

  const value = user.name?.trim() || user.email?.trim() || "用户";

  return value.includes("@") ? value.split("@")[0] || "用户" : value;
});

const avatarText = computed(() => {
  return Array.from(displayName.value)[0]?.toUpperCase() || "U";
});

const roleLabel = computed(() =>
  getRoleLabel(authStore.user?.user_type ?? null),
);

const currentSection = computed(() => {
  if (route.path.startsWith("/mes")) return "制造执行 MES";
  if (route.path.startsWith("/users")) return "用户管理";
  if (route.path.startsWith("/wms")) return "仓储管理 WMS";
  if (route.path.startsWith("/iot")) return "IoT 设备连接";
  if (route.path.startsWith("/devices")) return "设备数字身份";
  return "平台工作台";
});

async function handleLogout() {
  try {
    await authStore.logout();
  } finally {
    await router.replace("/login");
  }
}
</script>

<template>
  <div class="authenticated-layout">
    <header class="global-topbar">
      <div class="topbar-inner">
        <button
          type="button"
          class="brand"
          aria-label="返回平台工作台"
          @click="router.push('/dashboard')"
        >
          <span class="brand-mark" aria-hidden="true">
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
          </span>

          <span class="brand-copy">
            <strong>OpenIndustrial</strong>
            <small>设备数字身份平台</small>
          </span>
        </button>

        <div class="topbar-context">
          <span class="context-indicator"></span>
          <span>{{ currentSection }}</span>
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

    <main class="layout-content">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.authenticated-layout {
  min-height: 100vh;
  color: #17283b;
  background: #f4f6f9;
}

.global-topbar {
  position: sticky;
  top: 0;
  z-index: 1000;
  height: 76px;
  border-bottom: 1px solid rgb(255 255 255 / 8%);
  color: #fff;
  background: #0e1f33;
  box-shadow: 0 4px 18px rgb(14 31 51 / 8%);
}

.topbar-inner {
  display: flex;
  align-items: center;
  width: min(1440px, calc(100% - 64px));
  height: 100%;
  margin: 0 auto;
  gap: 32px;
}

.brand {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 12px;
  padding: 0;
  border: 0;
  color: inherit;
  background: transparent;
  text-align: left;
  cursor: pointer;
}

.brand-mark {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border: 1px solid rgb(240 160 48 / 25%);
  border-radius: 10px;
  background: #16293f;
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

.brand-copy strong {
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.brand-copy small {
  color: #aab9c9;
  font-size: 11px;
  letter-spacing: 0.5px;
}

.topbar-context {
  display: flex;
  align-items: center;
  gap: 9px;
  min-width: 0;
  color: #d0dae4;
  font-size: 12px;
}

.context-indicator {
  width: 6px;
  height: 6px;
  flex: 0 0 6px;
  border-radius: 50%;
  background: #f0a030;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-left: auto;
}

.user-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  flex: 0 0 36px;
  border: 1px solid rgb(240 160 48 / 35%);
  border-radius: 50%;
  color: #ffd18a;
  background: #223b55;
  font-size: 13px;
  font-weight: 700;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.user-name {
  max-width: 180px;
  overflow: hidden;
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-role {
  color: #aab9c9;
  font-size: 10px;
}

.logout-button {
  color: #e1e8ef;
  font-size: 12px;
}

.logout-button:hover {
  color: #ffd18a;
  background: rgb(255 255 255 / 7%);
}

/* Dashboard 已有顶部栏由此布局统一接管，避免重复显示。 */
.layout-content :deep(.dashboard > .topbar) {
  display: none !important;
}

@media (max-width: 760px) {
  .global-topbar {
    height: 68px;
  }

  .topbar-inner {
    width: calc(100% - 32px);
    gap: 12px;
  }

  .brand {
    gap: 8px;
  }

  .brand-mark {
    width: 36px;
    height: 36px;
  }

  .brand-mark svg {
    width: 26px;
    height: 26px;
  }

  .brand-copy strong {
    font-size: 13px;
  }

  .brand-copy small {
    font-size: 9px;
  }

  .topbar-context {
    display: none;
  }

  .topbar-right {
    gap: 8px;
  }

  .user-info {
    display: none;
  }

  .logout-button {
    padding: 7px !important;
  }
}

@media (max-width: 420px) {
  .topbar-inner {
    width: calc(100% - 20px);
  }

  .brand-copy strong {
    font-size: 12px;
  }

  .brand-copy small {
    letter-spacing: 0;
  }
}
</style>
