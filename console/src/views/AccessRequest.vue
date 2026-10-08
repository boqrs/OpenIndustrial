<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";

import { requestAccess as requestAccessApi } from "../api/auth";

const router = useRouter();
const route = useRoute();

const tenantCode = ref(String(route.query.tenant_code || ""));

const name = ref("");

const email = ref(String(route.query.email || ""));

const loading = ref(false);
const submitted = ref(false);

const handleSubmit = async () => {
  if (!tenantCode.value.trim()) {
    ElMessage.error("请输入工厂代码");
    return;
  }

  if (!name.value.trim()) {
    ElMessage.error("请输入姓名");
    return;
  }

  if (!email.value.trim()) {
    ElMessage.error("请输入邮箱");
    return;
  }

  loading.value = true;

  try {
    await requestAccessApi({
      tenant_code: tenantCode.value.trim(),

      name: name.value.trim(),

      email: email.value.trim(),
    });

    submitted.value = true;

    ElMessage.success("申请已提交");
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "申请提交失败");
  } finally {
    loading.value = false;
  }
};

const backToLogin = () => {
  router.replace({
    path: "/login",
    query: {
      tenant_code: tenantCode.value.trim(),
      email: email.value.trim(),
    },
  });
};
</script>

<template>
  <div class="access-request-page">
    <div class="access-request-card">
      <template v-if="!submitted">
        <div class="header">
          <div class="brand-mark">ID</div>

          <h1>申请开通账号</h1>

          <p>
            如果您还没有系统账号，请填写以下信息，
            我们将通知工厂管理员为您开通。
          </p>
        </div>

        <form class="form" @submit.prevent="handleSubmit">
          <div class="field">
            <label> 工厂代码 </label>

            <input
              v-model="tenantCode"
              type="text"
              placeholder="请输入工厂代码"
              autocomplete="organization"
              :disabled="loading"
            />
          </div>

          <div class="field">
            <label> 姓名 </label>

            <input
              v-model="name"
              type="text"
              placeholder="请输入您的姓名"
              autocomplete="name"
              :disabled="loading"
            />
          </div>

          <div class="field">
            <label> 邮箱 </label>

            <input
              v-model="email"
              type="email"
              placeholder="请输入您的工作邮箱"
              autocomplete="email"
              :disabled="loading"
            />
          </div>

          <button type="submit" :disabled="loading">
            {{ loading ? "提交中..." : "提交申请" }}
          </button>
        </form>

        <div class="footer">
          <button type="button" class="back-button" @click="backToLogin">
            返回登录
          </button>
        </div>
      </template>

      <template v-else>
        <div class="success">
          <div class="success-icon">✓</div>

          <h1>申请已提交</h1>

          <p>您的账号开通申请已经提交给工厂管理员。</p>

          <p>管理员审核并发送邀请后， 您可以通过邀请链接完成账号激活。</p>

          <button type="button" @click="backToLogin">返回登录</button>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.access-request-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  background: #0e1f33;
}

.access-request-card {
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

h1 {
  margin: 0 0 12px;

  color: #0e1f33;

  font-size: 26px;
  font-weight: 600;
}

.header p,
.success p {
  margin: 0;

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
  box-sizing: border-box;

  padding: 12px 14px;

  border: 1px solid #d0d5dd;
  border-radius: 4px;

  outline: none;

  font-size: 14px;

  transition:
    border-color 0.2s,
    box-shadow 0.2s;
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

.back-button {
  border: 0;
  padding: 0;

  background: transparent;
  color: #667085;

  font-size: 14px;

  cursor: pointer;
}

.back-button:hover {
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

.success h1 {
  margin-bottom: 16px;
}

.success p + p {
  margin-top: 8px;
}

.success > button {
  width: 100%;
  margin-top: 28px;
}
</style>
