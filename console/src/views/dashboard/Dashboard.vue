<script setup lang="ts">
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";

import { useAuthStore } from "../../store/auth";

const router = useRouter();
const authStore = useAuthStore();

async function handleLogout() {
  try {
    await authStore.logout();

    ElMessage.success("已退出登录");

    await router.replace("/login");
  } catch {
    authStore.clearSession();

    await router.replace("/login");
  }
}
</script>

<template>
  <div class="dashboard">
    <header class="header">
      <div class="brand">
        <div class="brand-mark">ID</div>

        <div>
          <div class="brand-title">设备数字身份平台</div>

          <div class="brand-subtitle">OpenIndustrial</div>
        </div>
      </div>

      <div class="header-right">
        <div v-if="authStore.user" class="user-info">
          <div class="user-name">
            {{ authStore.user.name || authStore.user.email }}
          </div>

          <div class="user-email">
            {{ authStore.user.email }}
          </div>
        </div>

        <button class="logout" type="button" @click="handleLogout">
          退出登录
        </button>
      </div>
    </header>

    <main class="content">
      <section class="welcome">
        <div class="eyebrow">CONTROL CENTER</div>

        <h1>欢迎回来</h1>

        <p>设备数字身份平台控制台</p>
      </section>

      <section class="cards">
        <div class="card">
          <div class="label">当前用户</div>

          <div class="value">
            {{ authStore.user?.name || "-" }}
          </div>

          <div class="detail">
            {{ authStore.user?.email || "-" }}
          </div>
        </div>

        <div class="card">
          <div class="label">Tenant</div>

          <div class="value">
            {{ authStore.user?.tenant_id || "-" }}
          </div>

          <div class="detail">当前工厂</div>
        </div>

        <div class="card">
          <div class="label">登录状态</div>

          <div class="value status">已登录</div>

          <div class="detail">Access Token 已建立</div>
        </div>
      </section>

      <section class="placeholder">
        <div class="placeholder-title">Industrial Cloud</div>

        <div class="placeholder-text">
          CRM · MES · WMS · IoT · Device Identity
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100vh;

  background: #f5f7fa;
}

.header {
  height: 68px;

  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 0 32px;

  background: #0e1f33;

  color: #ffffff;
}

.brand {
  display: flex;
  align-items: center;
}

.brand-mark {
  width: 34px;
  height: 34px;

  display: flex;
  align-items: center;
  justify-content: center;

  margin-right: 11px;

  background: #f0a030;

  color: #0e1f33;

  font-size: 12px;
  font-weight: 800;
}

.brand-title {
  font-size: 15px;
  font-weight: 600;
}

.brand-subtitle {
  margin-top: 2px;

  color: rgba(255, 255, 255, 0.42);

  font-size: 9px;

  letter-spacing: 0.5px;
}

.header-right {
  display: flex;
  align-items: center;

  gap: 22px;
}

.user-info {
  text-align: right;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
}

.user-email {
  margin-top: 2px;

  color: rgba(255, 255, 255, 0.5);

  font-size: 11px;
}

.logout {
  height: 34px;

  padding: 0 14px;

  border: 1px solid rgba(255, 255, 255, 0.2);

  border-radius: 4px;

  background: transparent;

  color: rgba(255, 255, 255, 0.85);

  font-size: 12px;
}

.logout:hover {
  border-color: #f0a030;

  color: #f0a030;
}

.content {
  max-width: 1200px;

  margin: 0 auto;

  padding: 48px 32px;
}

.welcome {
  margin-bottom: 36px;
}

.eyebrow {
  margin-bottom: 8px;

  color: #f0a030;

  font-size: 11px;
  font-weight: 600;

  letter-spacing: 1.5px;
}

.welcome h1 {
  margin: 0;

  color: #0e1f33;

  font-size: 30px;
  font-weight: 600;
}

.welcome p {
  margin-top: 8px;

  color: #667085;

  font-size: 14px;
}

.cards {
  display: grid;

  grid-template-columns: repeat(3, minmax(0, 1fr));

  gap: 20px;
}

.card {
  min-height: 140px;

  padding: 24px;

  background: #ffffff;

  border: 1px solid #eaecf0;

  border-radius: 6px;
}

.label {
  color: #667085;

  font-size: 12px;
}

.value {
  margin-top: 18px;

  color: #0e1f33;

  font-size: 20px;
  font-weight: 600;

  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.value.status {
  color: #16794a;
}

.detail {
  margin-top: 7px;

  color: #98a2b3;

  font-size: 12px;

  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.placeholder {
  margin-top: 24px;

  padding: 28px;

  background: #ffffff;

  border: 1px solid #eaecf0;

  border-radius: 6px;
}

.placeholder-title {
  color: #0e1f33;

  font-size: 16px;
  font-weight: 600;
}

.placeholder-text {
  margin-top: 8px;

  color: #98a2b3;

  font-size: 12px;

  letter-spacing: 0.5px;
}

@media (max-width: 800px) {
  .header {
    padding: 0 18px;
  }

  .content {
    padding: 32px 18px;
  }

  .cards {
    grid-template-columns: 1fr;
  }

  .user-info {
    display: none;
  }
}
</style>
