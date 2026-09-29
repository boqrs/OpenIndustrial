```vue
<template>
  <div class="login-page">
    <!-- 左侧品牌区 -->
    <section class="brand-side">
      <div class="brand-content">
        <!-- Logo -->
        <div class="logo-lockup">
          <svg
            class="logo-mark"
            viewBox="0 0 60 60"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            aria-hidden="true"
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

          <div class="logo-text">
            <div class="logo-text-main">设备数字身份平台</div>

            <div class="logo-text-sub">Device Identity Platform</div>
          </div>
        </div>

        <!-- 主标题 -->
        <h1 class="brand-title">
          一物一码
          <br />
          从制造到使用<em>全程可溯</em>
        </h1>

        <!-- 产品介绍 -->
        <p class="brand-desc">
          为每一台设备赋予唯一数字身份。<br />
          打通 CRM、MES、WMS 与
          IoT，记录设备从生产、运输到用户使用全过程的完整数据。
        </p>

        <!-- 生命周期 -->
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

        <!-- Device ID 示例 -->
        <div class="code-sample">
          <span class="code-sample-label"> Device ID </span>

          <span class="code-sample-value"> DV-2024-88A3-9F2C </span>
        </div>
      </div>
    </section>

    <!-- 右侧登录区 -->
    <section class="form-side">
      <div class="form-box">
        <h2 class="form-title">欢迎回来</h2>

        <p class="form-subtitle">请登录设备数字身份管理平台</p>

        <form @submit.prevent="handleLogin">
          <!-- 账户 -->
          <div class="field">
            <label for="account"> 账户名 </label>

            <input
              id="account"
              v-model="account"
              type="text"
              autocomplete="username"
              placeholder="请输入企业账户或邮箱"
              :disabled="loading"
            />
          </div>

          <!-- 密码 -->
          <div class="field">
            <label for="password"> 密码 </label>

            <input
              id="password"
              v-model="password"
              type="password"
              autocomplete="current-password"
              placeholder="请输入登录密码"
              :disabled="loading"
              @keyup.enter="handleLogin"
            />
          </div>

          <!-- 记住登录 / 忘记密码 -->
          <div class="form-extra">
            <label class="checkbox-wrap">
              <input v-model="remember" type="checkbox" :disabled="loading" />

              <span> 记住登录状态 </span>
            </label>

            <a class="link" href="#" @click.prevent> 忘记密码？ </a>
          </div>

          <!-- 登录 -->
          <button class="btn-login" type="submit" :disabled="loading">
            <span v-if="loading"> 登录中... </span>

            <span v-else> 登 录 </span>
          </button>
        </form>

        <!-- 错误提示 -->
        <p v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </p>

        <!-- 底部 -->
        <p class="form-footer">
          首次使用？
          <a class="link" href="#" @click.prevent> 联系管理员开通账户 </a>
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const account = ref("");
const password = ref("");
const remember = ref(false);

const loading = ref(false);
const errorMessage = ref("");

const handleLogin = async () => {
  errorMessage.value = "";

  if (!account.value.trim()) {
    errorMessage.value = "请输入账户名";
    return;
  }

  if (!password.value) {
    errorMessage.value = "请输入登录密码";
    return;
  }

  loading.value = true;

  try {
    /*
     * TODO:
     * 后续在这里接入真实登录 API。
     *
     * 例如：
     *
     * const response = await login({
     *   email: account.value,
     *   password: password.value,
     * })
     *
     * 登录成功后：
     *
     * router.replace('/dashboard')
     */

    // 当前阶段仅用于验证页面交互。
    await new Promise((resolve) => {
      window.setTimeout(resolve, 500);
    });
  } finally {
    loading.value = false;
  }
};
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
  min-height: 100vh;

  display: flex;

  overflow: hidden;

  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC",
    "Microsoft YaHei", sans-serif;
}

/* =========================================================
   左侧品牌区
   ========================================================= */

.brand-side {
  width: 55%;

  position: relative;

  display: flex;
  flex-direction: column;
  justify-content: center;

  padding: 0 8%;

  overflow: hidden;

  background: var(--navy);
}

/* 数字编码纹理 */

.brand-side::before {
  content: "";

  position: absolute;
  inset: 0;

  background-image:
    repeating-linear-gradient(
      0deg,
      rgba(240, 160, 48, 0.035) 0,
      rgba(240, 160, 48, 0.035) 1px,
      transparent 1px,
      transparent 64px
    ),
    repeating-linear-gradient(
      90deg,
      rgba(240, 160, 48, 0.035) 0,
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

  pointer-events: none;
}

/* 光晕 */

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

/* =========================================================
   Logo
   ========================================================= */

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
  margin-top: 4px;

  font-size: 11px;

  color: var(--gold);

  letter-spacing: 2.8px;

  text-transform: uppercase;

  opacity: 0.85;
}

/* =========================================================
   品牌标题
   ========================================================= */

.brand-title {
  margin: 0 0 20px;

  font-size: 38px;

  font-weight: 700;

  color: var(--text-light);

  line-height: 1.4;

  letter-spacing: 0.5px;
}

.brand-title em {
  font-style: normal;

  color: var(--gold);
}

/* =========================================================
   品牌描述
   ========================================================= */

.brand-desc {
  max-width: 460px;

  margin: 0;

  font-size: 15px;

  color: var(--text-muted);

  line-height: 1.9;
}

/* =========================================================
   生命周期
   ========================================================= */

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

  margin: 0 8px 26px;

  background: linear-gradient(90deg, var(--gold), rgba(240, 160, 48, 0.25));
}

/* =========================================================
   Device ID
   ========================================================= */

.code-sample {
  width: fit-content;

  display: inline-flex;
  align-items: center;

  gap: 12px;

  margin-top: 48px;

  padding: 10px 18px;

  border: 1px solid rgba(240, 160, 48, 0.25);

  border-radius: 6px;

  background: rgba(240, 160, 48, 0.04);

  position: relative;

  z-index: 2;
}

.code-sample-label {
  font-size: 11px;

  color: var(--text-muted);

  letter-spacing: 1.5px;

  text-transform: uppercase;
}

.code-sample-value {
  font-family: "SF Mono", "Consolas", "Liberation Mono", monospace;

  font-size: 14px;

  color: var(--gold);

  letter-spacing: 1.5px;

  font-weight: 600;
}

/* =========================================================
   右侧登录区
   ========================================================= */

.form-side {
  width: 45%;

  display: flex;
  align-items: center;
  justify-content: center;

  padding: 40px;

  background: #ffffff;
}

.form-box {
  width: 100%;
  max-width: 380px;
}

/* =========================================================
   标题
   ========================================================= */

.form-title {
  margin: 0 0 8px;

  font-size: 26px;

  font-weight: 700;

  color: #1a2a3d;
}

.form-subtitle {
  margin: 0 0 40px;

  font-size: 14px;

  color: #8a94a6;
}

/* =========================================================
   输入框
   ========================================================= */

.field {
  margin-bottom: 22px;
}

.field label {
  display: block;

  margin-bottom: 8px;

  font-size: 13px;

  font-weight: 600;

  color: #3d4757;

  letter-spacing: 0.3px;
}

.field input {
  width: 100%;
  height: 48px;

  padding: 0 16px;

  border: 1.5px solid transparent;

  border-radius: 8px;

  outline: none;

  background: var(--input-bg);

  color: #1a2a3d;

  font-family: inherit;

  font-size: 15px;

  transition:
    border-color 0.2s,
    background 0.2s,
    box-shadow 0.2s;
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

/* =========================================================
   登录辅助区域
   ========================================================= */

.form-extra {
  display: flex;
  align-items: center;
  justify-content: space-between;

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

.checkbox-wrap input:disabled {
  cursor: not-allowed;
}

.link {
  color: #d08a20;

  font-weight: 500;

  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}

/* =========================================================
   登录按钮
   ========================================================= */

.btn-login {
  width: 100%;
  height: 50px;

  border: none;

  border-radius: 8px;

  background: var(--gold);

  color: var(--navy);

  font-family: inherit;

  font-size: 16px;

  font-weight: 700;

  letter-spacing: 1px;

  cursor: pointer;

  transition:
    background 0.2s,
    box-shadow 0.2s,
    transform 0.2s;
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
  opacity: 0.7;

  cursor: wait;
}

/* =========================================================
   错误提示
   ========================================================= */

.error-message {
  margin: 18px 0 0;

  padding: 10px 12px;

  border-radius: 6px;

  background: #fff4f1;

  color: #c6533b;

  font-size: 13px;

  line-height: 1.5;

  text-align: center;
}

/* =========================================================
   Footer
   ========================================================= */

.form-footer {
  margin: 32px 0 0;

  color: #8a94a6;

  font-size: 13px;

  text-align: center;
}

/* =========================================================
   Responsive
   ========================================================= */

@media (max-width: 960px) {
  .login-page {
    min-height: 100vh;

    flex-direction: column;

    overflow-y: auto;
  }

  .brand-side {
    width: 100%;

    min-height: auto;

    padding: 48px 32px;
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

@media (max-width: 520px) {
  .brand-side {
    padding: 36px 24px;
  }

  .logo-lockup {
    margin-bottom: 36px;
  }

  .logo-mark {
    width: 52px;
    height: 52px;
  }

  .logo-text-main {
    font-size: 20px;
  }

  .brand-title {
    font-size: 26px;
  }

  .brand-desc {
    font-size: 14px;
  }

  .lifecycle {
    margin-top: 32px;
  }

  .stage-line {
    width: 24px;
    margin-left: 5px;
    margin-right: 5px;
  }

  .code-sample {
    margin-top: 36px;

    gap: 8px;

    padding: 9px 12px;
  }

  .code-sample-value {
    font-size: 12px;
  }

  .form-side {
    padding: 40px 20px;
  }
}
</style>
```
