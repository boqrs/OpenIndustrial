<script setup lang="ts">
import { ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";

import { useAuthStore } from "../../store/auth";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const tenantCode = ref(String(route.query.tenant_code || ""));

const email = ref(String(route.query.email || ""));

const password = ref("");

const loading = ref(false);

async function handleLogin() {
  if (!tenantCode.value.trim()) {
    ElMessage.error("请输入工厂代码");
    return;
  }

  if (!email.value.trim()) {
    ElMessage.error("请输入邮箱");
    return;
  }

  if (!password.value) {
    ElMessage.error("请输入密码");
    return;
  }

  loading.value = true;

  try {
    await authStore.login({
      tenant_code: tenantCode.value.trim(),

      email: email.value.trim(),

      password: password.value,
    });

    ElMessage.success("登录成功");

    const redirect =
      typeof route.query.redirect === "string" &&
      route.query.redirect.startsWith("/")
        ? route.query.redirect
        : "/dashboard";

    await router.replace(redirect);
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "登录失败，请检查账号和密码",
    );
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="login-page">
    <div class="background">
      <div class="grid"></div>
    </div>

    <main class="login-container">
      <div class="brand">
        <div class="brand-mark">ID</div>

        <div>
          <div class="brand-title">设备数字身份平台</div>

          <div class="brand-subtitle">INDUSTRIAL DEVICE IDENTITY PLATFORM</div>
        </div>
      </div>

      <section class="login-card">
        <div class="header">
          <h1>登录</h1>

          <p>使用您的工厂账号登录平台</p>
        </div>

        <form class="form" @submit.prevent="handleLogin">
          <div class="field">
            <label for="tenant-code"> 工厂代码 </label>

            <input
              id="tenant-code"
              v-model="tenantCode"
              type="text"
              placeholder="请输入工厂代码"
              autocomplete="organization"
              :disabled="loading"
            />
          </div>

          <div class="field">
            <label for="email"> 邮箱 </label>

            <input
              id="email"
              v-model="email"
              type="email"
              placeholder="请输入工作邮箱"
              autocomplete="email"
              :disabled="loading"
            />
          </div>

          <div class="field">
            <label for="password"> 密码 </label>

            <input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入密码"
              autocomplete="current-password"
              :disabled="loading"
            />
          </div>

          <button class="login-button" type="submit" :disabled="loading">
            {{ loading ? "登录中..." : "登录" }}
          </button>
        </form>

        <div class="footer">
          <span> 还没有平台账号？ </span>

          <RouterLink
            :to="{
              path: '/access-request',
              query: {
                tenant_code: tenantCode,
                email: email,
              },
            }"
          >
            联系工厂管理员开通账户
          </RouterLink>
        </div>
      </section>

      <div class="copyright">OpenIndustrial</div>
    </main>
  </div>
</template>

<style scoped>
.login-page {
  position: relative;

  min-height: 100vh;

  display: flex;
  align-items: center;
  justify-content: center;

  overflow: hidden;

  background: #0e1f33;
}

.background {
  position: absolute;
  inset: 0;

  pointer-events: none;
}

.grid {
  position: absolute;

  width: 160%;
  height: 70%;

  left: -30%;
  bottom: -10%;

  background-image:
    linear-gradient(rgba(255, 255, 255, 0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.025) 1px, transparent 1px);

  background-size: 48px 48px;

  transform: perspective(500px) rotateX(55deg);
}

.login-container {
  position: relative;
  z-index: 1;

  width: 100%;
  max-width: 420px;

  padding: 32px 24px;
}

.brand {
  display: flex;
  align-items: center;

  margin-bottom: 30px;
}

.brand-mark {
  width: 46px;
  height: 46px;

  display: flex;
  align-items: center;
  justify-content: center;

  margin-right: 14px;

  background: #f0a030;

  color: #0e1f33;

  font-size: 15px;
  font-weight: 800;

  letter-spacing: 1px;
}

.brand-title {
  color: #ffffff;

  font-size: 20px;
  font-weight: 600;
}

.brand-subtitle {
  margin-top: 4px;

  color: rgba(255, 255, 255, 0.42);

  font-size: 9px;

  letter-spacing: 0.5px;
}

.login-card {
  padding: 38px;

  background: #ffffff;

  border-radius: 6px;

  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.28);
}

.header {
  margin-bottom: 30px;
}

.header h1 {
  margin: 0 0 8px;

  color: #0e1f33;

  font-size: 28px;
  font-weight: 600;
}

.header p {
  color: #667085;

  font-size: 14px;
}

.form {
  display: flex;
  flex-direction: column;

  gap: 20px;
}

.field {
  display: flex;
  flex-direction: column;

  gap: 8px;
}

.field label {
  color: #344054;

  font-size: 14px;
  font-weight: 500;
}

.field input {
  width: 100%;
  height: 44px;

  box-sizing: border-box;

  padding: 0 13px;

  border: 1px solid #d0d5dd;

  border-radius: 4px;

  outline: none;

  color: #101828;

  font-size: 14px;
}

.field input::placeholder {
  color: #98a2b3;
}

.field input:focus {
  border-color: #0e1f33;

  box-shadow: 0 0 0 2px rgba(14, 31, 51, 0.08);
}

.field input:disabled {
  background: #f9fafb;
}

.login-button {
  width: 100%;
  height: 44px;

  margin-top: 4px;

  border: 0;
  border-radius: 4px;

  background: #0e1f33;

  color: #ffffff;

  font-size: 14px;
  font-weight: 500;

  transition:
    background 0.2s,
    transform 0.1s;
}

.login-button:hover:not(:disabled) {
  background: #16293f;
}

.login-button:active:not(:disabled) {
  transform: translateY(1px);
}

.login-button:disabled {
  opacity: 0.6;

  cursor: not-allowed;
}

.footer {
  display: flex;
  justify-content: center;

  margin-top: 26px;

  gap: 5px;

  color: #667085;

  font-size: 13px;
}

.footer a {
  color: #0e1f33;

  font-weight: 500;

  text-decoration: none;
}

.footer a:hover {
  color: #f0a030;

  text-decoration: underline;
}

.copyright {
  margin-top: 24px;

  text-align: center;

  color: rgba(255, 255, 255, 0.35);

  font-size: 11px;

  letter-spacing: 0.5px;
}

@media (max-width: 520px) {
  .login-container {
    padding: 24px 16px;
  }

  .login-card {
    padding: 28px 22px;
  }

  .brand-title {
    font-size: 18px;
  }
}
</style>
