<template>
  <div class="login-page">
    <!-- 左侧品牌区 -->
    <div class="brand-side">
      <div class="brand-content">
        <!-- Logo -->
        <div class="logo-lockup">
          <svg
            class="logo-mark"
            viewBox="0 0 60 60"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <!-- 外圈不闭合圆环 -->
            <path
              d="M30 4 A26 26 0 1 1 4 30"
              stroke="#F0A030"
              stroke-width="2.4"
              fill="none"
              stroke-linecap="round"
            />
            <path
              d="M4 30 A26 26 0 0 1 12 12"
              stroke="#F0A030"
              stroke-width="2.4"
              fill="none"
              stroke-linecap="round"
              opacity="0.35"
            />

            <!-- 中心设备轮廓 -->
            <rect
              x="20"
              y="20"
              width="20"
              height="16"
              rx="3"
              stroke="#E8EDF5"
              stroke-width="1.8"
              fill="none"
            />
            <circle cx="30" cy="28" r="2.4" fill="#F0A030" />
            <line
              x1="24"
              y1="36"
              x2="24"
              y2="40"
              stroke="#E8EDF5"
              stroke-width="1.6"
              stroke-linecap="round"
            />
            <line
              x1="36"
              y1="36"
              x2="36"
              y2="40"
              stroke="#E8EDF5"
              stroke-width="1.6"
              stroke-linecap="round"
            />

            <!-- 编码线 -->
            <line
              x1="22"
              y1="44"
              x2="38"
              y2="44"
              stroke="#F0A030"
              stroke-width="1.6"
              stroke-linecap="round"
              opacity="0.7"
            />
            <line
              x1="26"
              y1="47"
              x2="34"
              y2="47"
              stroke="#F0A030"
              stroke-width="1.6"
              stroke-linecap="round"
              opacity="0.4"
            />

            <!-- 三阶段节点 -->
            <circle cx="30" cy="4" r="3" fill="#F0A030" />
            <circle cx="56" cy="30" r="3" fill="#F0A030" opacity="0.55" />
            <circle cx="4" cy="30" r="3" fill="#F0A030" opacity="0.55" />
          </svg>

          <div>
            <div class="logo-text-main">设备数字身份平台</div>
            <div class="logo-text-sub">Device Identity Platform</div>
          </div>
        </div>

        <h1 class="brand-title">
          一物一码<br />
          从制造到使用<em>全程可溯</em>
        </h1>

        <p class="brand-desc">
          为每一台设备赋予唯一数字身份。<br />
          打通 CRM、MES、WMS 与
          IoT，记录设备从生产、运输到用户使用全过程的完整数据。
        </p>

        <!-- 生命周期三阶段 -->
        <div class="lifecycle">
          <div class="stage">
            <div class="stage-dot"></div>
            <div class="stage-label active">制造</div>
          </div>

          <div class="stage-line"></div>

          <div class="stage">
            <div class="stage-dot dim"></div>
            <div class="stage-label">运输</div>
          </div>

          <div class="stage-line"></div>

          <div class="stage">
            <div class="stage-dot dim"></div>
            <div class="stage-label">使用</div>
          </div>
        </div>

        <!-- 编码示例 -->
        <div class="code-sample">
          <span class="code-sample-label">Device ID</span>
          <span class="code-sample-value"> DV-2024-88A3-9F2C </span>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="form-side">
      <div class="form-box">
        <h2 class="form-title">欢迎回来</h2>

        <p class="form-subtitle">请登录设备数字身份管理平台</p>

        <form @submit.prevent="handleLogin">
          <!-- 租户编码 -->
          <div class="field">
            <label>租户编码</label>
            <input
              v-model.trim="form.tenant_code"
              type="text"
              autocomplete="organization"
              placeholder="请输入企业租户编码"
              :disabled="loading"
            />
          </div>

          <!-- 邮箱 -->
          <div class="field">
            <label>邮箱</label>
            <input
              v-model.trim="form.email"
              type="email"
              autocomplete="email"
              placeholder="请输入登录邮箱"
              :disabled="loading"
            />
          </div>

          <!-- 密码 -->
          <div class="field">
            <label>密码</label>
            <input
              v-model="form.password"
              type="password"
              autocomplete="current-password"
              placeholder="请输入登录密码"
              :disabled="loading"
            />
          </div>

          <div class="form-extra">
            <label class="checkbox-wrap">
              <input v-model="rememberLogin" type="checkbox" />
              <span>记住登录状态</span>
            </label>

            <a href="#" class="link" @click.prevent> 忘记密码？ </a>
          </div>

          <button type="submit" class="btn-login" :disabled="loading">
            {{ loading ? "登录中..." : "登 录" }}
          </button>
        </form>

        <p class="form-footer">
          首次使用？
          <RouterLink
            class="link"
            :to="{
              path: '/access-request',
              query: {
                tenant_code: form.tenant_code,
                email: form.email,
              },
            }"
          >
            联系管理员开通账户
          </RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRoute, useRouter, RouterLink } from "vue-router";
import { ElMessage } from "element-plus";

import { useAuthStore } from "../../store/auth";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const loading = ref(false);
const rememberLogin = ref(true);

const form = reactive({
  tenant_code: "",
  email: "",
  password: "",
});

async function handleLogin() {
  if (!form.tenant_code) {
    ElMessage.warning("请输入租户编码");
    return;
  }

  if (!form.email) {
    ElMessage.warning("请输入邮箱");
    return;
  }

  if (!form.password) {
    ElMessage.warning("请输入密码");
    return;
  }

  loading.value = true;

  try {
    await authStore.login({
      tenant_code: form.tenant_code,
      email: form.email,
      password: form.password,
    });

    ElMessage.success("登录成功");

    const redirect =
      typeof route.query.redirect === "string" &&
      route.query.redirect.startsWith("/")
        ? route.query.redirect
        : "/dashboard";

    await router.replace(redirect);
  } catch (error: unknown) {
    const message =
      error instanceof Error
        ? error.message
        : "登录失败，请检查租户编码、邮箱和密码";

    ElMessage.error(message);
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  --navy: #0e1f33;
  --navy-light: #16293f;
  --gold: #f0a030;
  --gold-dim: rgba(240, 160, 48, 0.12);
  --text-light: #e8edf5;
  --text-muted: #7a8a9e;
  --input-bg: #f5f7fa;

  width: 100%;
  height: 100vh;
  display: flex;
  overflow: hidden;
  font-family:
    -apple-system, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
}

/* ==================== 左侧品牌区 ==================== */

.brand-side {
  width: 55%;
  background: var(--navy);
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 8%;
  overflow: hidden;
}

.brand-side::before {
  content: "";
  position: absolute;
  inset: 0;
  background-image:
    repeating-linear-gradient(
      0deg,
      rgba(240, 160, 48, 0.035) 0px,
      rgba(240, 160, 48, 0.035) 1px,
      transparent 1px,
      transparent 64px
    ),
    repeating-linear-gradient(
      90deg,
      rgba(240, 160, 48, 0.035) 0px,
      rgba(240, 160, 48, 0.035) 1px,
      transparent 1px,
      transparent 64px
    );

  mask-image: radial-gradient(ellipse at 35% 45%, #000 15%, transparent 75%);

  -webkit-mask-image: radial-gradient(
    ellipse at 35% 45%,
    #000 15%,
    transparent 75%
  );
}

.brand-side::after {
  content: "";
  position: absolute;
  width: 700px;
  height: 700px;
  border-radius: 50%;
  background: radial-gradient(
    circle,
    rgba(240, 160, 48, 0.1) 0%,
    transparent 65%
  );
  top: 45%;
  left: 25%;
  transform: translate(-50%, -50%);
  pointer-events: none;
}

.brand-content {
  position: relative;
  z-index: 2;
}

/* ==================== Logo ==================== */

.logo-lockup {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 52px;
}

.logo-mark {
  width: 60px;
  height: 60px;
  flex-shrink: 0;
}

.logo-text-main {
  font-size: 25px;
  font-weight: 700;
  color: var(--text-light);
  letter-spacing: 1px;
  line-height: 1.2;
}

.logo-text-sub {
  font-size: 11px;
  color: var(--gold);
  letter-spacing: 2.8px;
  text-transform: uppercase;
  margin-top: 4px;
  opacity: 0.85;
}

.brand-title {
  font-size: 38px;
  font-weight: 700;
  color: var(--text-light);
  line-height: 1.4;
  letter-spacing: 0.5px;
  margin: 0 0 20px;
}

.brand-title em {
  font-style: normal;
  color: var(--gold);
}

.brand-desc {
  font-size: 15px;
  color: var(--text-muted);
  line-height: 1.9;
  max-width: 460px;
  margin: 0;
}

/* ==================== 生命周期 ==================== */

.lifecycle {
  display: flex;
  align-items: center;
  gap: 0;
  margin-top: 56px;
  position: relative;
  z-index: 2;
}

.stage {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.stage-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--gold);
  box-shadow: 0 0 0 4px rgba(240, 160, 48, 0.15);
}

.stage-dot.dim {
  background: var(--navy-light);
  border: 2px solid var(--gold);
  box-shadow: none;
  opacity: 0.6;
}

.stage-label {
  font-size: 13px;
  color: var(--text-muted);
  letter-spacing: 1px;
}

.stage-label.active {
  color: var(--gold);
  font-weight: 600;
}

.stage-line {
  width: 72px;
  height: 2px;
  background: linear-gradient(90deg, var(--gold), rgba(240, 160, 48, 0.25));
  margin: 0 8px;
  margin-bottom: 26px;
}

/* ==================== Device ID ==================== */

.code-sample {
  margin-top: 48px;
  position: relative;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 10px 18px;
  border: 1px solid rgba(240, 160, 48, 0.25);
  border-radius: 6px;
  background: rgba(240, 160, 48, 0.04);
  width: fit-content;
}

.code-sample-label {
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 1.5px;
  text-transform: uppercase;
}

.code-sample-value {
  font-family: "SF Mono", "Consolas", monospace;
  font-size: 14px;
  color: var(--gold);
  letter-spacing: 1.5px;
  font-weight: 600;
}

/* ==================== 右侧表单 ==================== */

.form-side {
  width: 45%;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.form-box {
  width: 100%;
  max-width: 380px;
}

.form-title {
  font-size: 26px;
  font-weight: 700;
  color: #1a2a3d;
  margin: 0 0 8px;
}

.form-subtitle {
  font-size: 14px;
  color: #8a94a6;
  margin: 0 0 40px;
}

.field {
  margin-bottom: 22px;
}

.field label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #3d4757;
  margin-bottom: 8px;
  letter-spacing: 0.3px;
}

.field input {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  font-size: 15px;
  color: #1a2a3d;
  background: var(--input-bg);
  border: 1.5px solid transparent;
  border-radius: 8px;
  outline: none;
  transition: all 0.2s;
  font-family: inherit;
}

.field input::placeholder {
  color: #a8b2c1;
}

.field input:focus {
  background: #ffffff;
  border-color: var(--gold);
  box-shadow: 0 0 0 4px var(--gold-dim);
}

.field input:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.form-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  font-size: 13px;
}

.checkbox-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #5a6474;
  cursor: pointer;
  user-select: none;
}

.checkbox-wrap input {
  width: 16px;
  height: 16px;
  accent-color: var(--gold);
  cursor: pointer;
}

.link {
  color: #d08a20;
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}

.btn-login {
  width: 100%;
  height: 50px;
  background: var(--gold);
  color: #0e1f33;
  font-size: 16px;
  font-weight: 700;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  letter-spacing: 1px;
  transition: all 0.2s;
  font-family: inherit;
}

.btn-login:hover:not(:disabled) {
  background: #e09428;
  box-shadow: 0 6px 20px rgba(240, 160, 48, 0.35);
  transform: translateY(-1px);
}

.btn-login:active:not(:disabled) {
  transform: translateY(0);
}

.btn-login:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}

.form-footer {
  margin: 32px 0 0;
  text-align: center;
  font-size: 13px;
  color: #8a94a6;
}

/* ==================== 响应式 ==================== */

@media (max-width: 960px) {
  .login-page {
    flex-direction: column;
    height: auto;
    min-height: 100vh;
    overflow-y: auto;
  }

  .brand-side {
    width: 100%;
    padding: 48px 32px;
    min-height: auto;
  }

  .brand-title {
    font-size: 28px;
  }

  .lifecycle {
    margin-top: 36px;
  }

  .stage-line {
    width: 40px;
  }

  .form-side {
    width: 100%;
    padding: 48px 24px;
  }
}
</style>
