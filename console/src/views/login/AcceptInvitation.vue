<template>
  <div class="page">
    <div class="card">
      <template v-if="!submitted">
        <div class="header">
          <div class="brand-mark">ID</div>

          <h1>激活账号</h1>

          <p>
            您已收到设备数字身份平台的账号邀请。 请设置登录密码完成账号激活。
          </p>
        </div>

        <form class="form" @submit.prevent="handleSubmit">
          <div class="field">
            <label>设置密码</label>

            <input
              v-model="password"
              type="password"
              autocomplete="new-password"
              placeholder="请输入登录密码"
              :disabled="loading"
            />
          </div>

          <div class="field">
            <label>确认密码</label>

            <input
              v-model="confirmPassword"
              type="password"
              autocomplete="new-password"
              placeholder="请再次输入密码"
              :disabled="loading"
            />
          </div>

          <button type="submit" :disabled="loading">
            {{ loading ? "激活中..." : "激活账号" }}
          </button>
        </form>

        <div class="footer">
          <button type="button" class="back" @click="backToLogin">
            返回登录
          </button>
        </div>
      </template>

      <template v-else>
        <div class="success">
          <div class="success-icon">✓</div>

          <h1>账号激活成功</h1>

          <p>您的账号已经激活，可以使用设置的密码登录系统。</p>

          <button type="button" @click="backToLogin">返回登录</button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";

import { acceptInvitation } from "../../api/auth";

const router = useRouter();
const route = useRoute();

const password = ref("");
const confirmPassword = ref("");

const loading = ref(false);
const submitted = ref(false);

const token = String(route.query.token || "");

async function handleSubmit() {
  if (!token) {
    ElMessage.error("邀请链接无效或缺少邀请凭证");
    return;
  }

  if (!password.value) {
    ElMessage.warning("请输入登录密码");
    return;
  }

  if (password.value.length < 8) {
    ElMessage.warning("密码长度不能少于 8 位");
    return;
  }

  if (password.value !== confirmPassword.value) {
    ElMessage.warning("两次输入的密码不一致");
    return;
  }

  loading.value = true;

  try {
    await acceptInvitation({
      token,
      password: password.value,
    });

    submitted.value = true;

    ElMessage.success("账号激活成功");
  } catch (error: unknown) {
    ElMessage.error(error instanceof Error ? error.message : "账号激活失败");
  } finally {
    loading.value = false;
  }
}

function backToLogin() {
  router.replace("/login");
}
</script>

<style scoped>
.page {
  min-height: 100vh;

  display: flex;
  align-items: center;
  justify-content: center;

  padding: 40px 24px;

  background: #0e1f33;
}

.card {
  width: 100%;
  max-width: 460px;

  padding: 42px;

  background: #ffffff;

  border-radius: 8px;

  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
}

.header {
  margin-bottom: 32px;
}

.brand-mark {
  width: 42px;
  height: 42px;

  margin-bottom: 22px;

  display: flex;
  align-items: center;
  justify-content: center;

  background: #0e1f33;

  color: #f0a030;

  font-size: 14px;
  font-weight: 700;

  letter-spacing: 1px;
}

.header h1,
.success h1 {
  margin: 0 0 12px;

  color: #0e1f33;

  font-size: 26px;
  font-weight: 600;
}

.header p,
.success p {
  color: #667085;

  font-size: 14px;

  line-height: 1.7;
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

  padding: 0 14px;

  border: 1px solid #d0d5dd;

  border-radius: 4px;

  outline: none;

  font-size: 14px;
}

.field input:focus {
  border-color: #0e1f33;

  box-shadow: 0 0 0 2px rgba(14, 31, 51, 0.08);
}

.form > button,
.success > button {
  height: 44px;

  border: 0;
  border-radius: 4px;

  background: #0e1f33;

  color: #ffffff;

  font-size: 14px;
  font-weight: 500;

  cursor: pointer;
}

.form > button:hover,
.success > button:hover {
  background: #16293f;
}

.form > button:disabled {
  opacity: 0.6;

  cursor: not-allowed;
}

.footer {
  margin-top: 24px;

  text-align: center;
}

.back {
  padding: 0;

  border: 0;

  background: transparent;

  color: #667085;

  font-size: 14px;

  cursor: pointer;
}

.back:hover {
  color: #0e1f33;
}

.success {
  text-align: center;
}

.success-icon {
  width: 56px;
  height: 56px;

  margin: 0 auto 24px;

  display: flex;
  align-items: center;
  justify-content: center;

  border-radius: 50%;

  background: #edf7ed;

  color: #2e7d32;

  font-size: 28px;
  font-weight: 600;
}

.success > button {
  width: 100%;

  margin-top: 28px;
}
</style>
