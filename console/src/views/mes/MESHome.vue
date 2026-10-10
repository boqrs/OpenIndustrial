<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import {
  ArrowRight,
  Box,
  Calendar,
  Connection,
  DataAnalysis,
  Document,
  OfficeBuilding,
  Operation,
  Setting,
  Tools,
} from "@element-plus/icons-vue";

import { useAuthStore } from "../../store/auth";
import { hasPermission } from "../../authz/access";
import { Permissions } from "../../authz/permissions";

const router = useRouter();
const authStore = useAuthStore();

const canAccessProducts = computed(() =>
  hasPermission(authStore.user, Permissions.PRODUCT_LIST),
);

const resourceModules = [
  {
    key: "factory",
    title: "工厂与生产资源",
    subtitle: "FACTORY RESOURCES",
    description: "逐步管理工厂、生产线、工位及其层级关系。",
    icon: OfficeBuilding,
    status: "建设中",
    statusType: "planned",
  },
  {
    key: "product",
    title: "产品与工艺",
    subtitle: "PRODUCT & PROCESS",
    description: "维护产品型号、物料清单与工艺路线，为生产执行提供基础数据。",
    icon: Box,
    status: "已开放",
    statusType: "available",
  },
  {
    key: "planning",
    title: "生产计划",
    subtitle: "PRODUCTION PLANNING",
    description: "规划生产任务，管理生产计划和生产工单。",
    icon: Calendar,
    status: "建设中",
    statusType: "planned",
  },
  {
    key: "execution",
    title: "生产执行",
    subtitle: "PRODUCTION EXECUTION",
    description: "组织工序操作、执行记录、生产结果和设备身份生成。",
    icon: Operation,
    status: "建设中",
    statusType: "planned",
  },
];

const productionStages = [
  { title: "产品", description: "定义产品型号", icon: Box },
  { title: "BOM", description: "定义物料结构", icon: Document },
  { title: "工艺路线", description: "定义生产流程", icon: Tools },
  { title: "生产计划", description: "安排生产任务", icon: Calendar },
  { title: "生产工单", description: "形成执行任务", icon: Setting },
  { title: "生产执行", description: "记录工序操作", icon: Operation },
  { title: "生产结果", description: "确认生产结果", icon: DataAnalysis },
];

function openProducts() {
  if (!canAccessProducts.value) return;
  void router.push("/mes/products");
}

function handleStageClick(title: string) {
  if (title === "产品") {
    openProducts();
  }
}
</script>

<template>
  <div class="mes-home">
    <section class="mes-hero">
      <div class="hero-copy">
        <div class="eyebrow">
          <span class="eyebrow-line"></span>
          MANUFACTURING EXECUTION SYSTEM
        </div>

        <h1>制造执行 <span>MES</span></h1>

        <p>
          从工厂资源与产品定义，到生产计划、工序执行和生产结果，
          在一个工作空间中逐步建立完整的制造业务流程。
        </p>

        <div class="hero-meta">
          <span class="meta-dot"></span>
          <span>制造业务工作台</span>
          <span class="meta-divider"></span>
          <span>功能持续建设中</span>
        </div>
      </div>

      <div class="hero-visual" aria-hidden="true">
        <div class="visual-ring ring-one"></div>
        <div class="visual-ring ring-two"></div>
        <div class="visual-center">
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
        <div class="visual-node node-one"></div>
        <div class="visual-node node-two"></div>
        <div class="visual-node node-three"></div>
      </div>
    </section>

    <section class="section-block">
      <div class="section-heading">
        <div>
          <div class="section-kicker">MANUFACTURING MODULES</div>
          <h2>制造业务</h2>
          <p>按制造业务的实际顺序组织功能，已开放的模块可以直接进入。</p>
        </div>
      </div>

      <div class="module-grid">
        <button
          v-for="item in resourceModules"
          :key="item.key"
          type="button"
          class="module-card"
          :class="[
            `module-card--${item.key}`,
            {
              'module-card--disabled':
                item.statusType !== 'available' ||
                (item.key === 'product' && !canAccessProducts),
            },
          ]"
          :disabled="
            item.statusType !== 'available' ||
            (item.key === 'product' && !canAccessProducts)
          "
          @click="item.key === 'product' && openProducts()"
        >
          <div class="module-card-top">
            <div class="module-icon">
              <el-icon :size="26">
                <component :is="item.icon" />
              </el-icon>
            </div>

            <el-tag
              :type="item.statusType === 'available' ? 'success' : 'info'"
              effect="plain"
              size="small"
            >
              {{
                item.key === "product" && !canAccessProducts
                  ? "无访问权限"
                  : item.status
              }}
            </el-tag>
          </div>

          <div class="module-subtitle">{{ item.subtitle }}</div>
          <h3>{{ item.title }}</h3>
          <p class="module-description">{{ item.description }}</p>

          <div class="module-footer">
            <span>
              {{
                item.key === "product" && !canAccessProducts
                  ? "请联系管理员申请权限"
                  : item.statusType === "available"
                    ? "进入模块"
                    : "后续逐步开放"
              }}
            </span>
            <el-icon v-if="item.statusType === 'available'">
              <ArrowRight />
            </el-icon>
          </div>
        </button>
      </div>
    </section>

    <section class="section-block">
      <div class="section-heading">
        <div>
          <div class="section-kicker">PRODUCTION FLOW</div>
          <h2>制造业务主流程</h2>
          <p>以现有 MES 领域模型为基础，逐步贯通产品定义到生产结果。</p>
        </div>
      </div>

      <div class="flow-panel">
        <div class="flow-list">
          <template
            v-for="(stage, index) in productionStages"
            :key="stage.title"
          >
            <article
              class="flow-stage"
              :class="{
                'flow-stage--clickable':
                  stage.title === '产品' && canAccessProducts,
              }"
              :role="
                stage.title === '产品' && canAccessProducts
                  ? 'button'
                  : undefined
              "
              :tabindex="
                stage.title === '产品' && canAccessProducts ? 0 : undefined
              "
              :aria-label="
                stage.title === '产品' && canAccessProducts
                  ? '进入产品列表'
                  : undefined
              "
              @click="handleStageClick(stage.title)"
              @keydown.enter.prevent="handleStageClick(stage.title)"
              @keydown.space.prevent="handleStageClick(stage.title)"
            >
              <div class="flow-icon">
                <el-icon :size="20">
                  <component :is="stage.icon" />
                </el-icon>
              </div>
              <strong>{{ stage.title }}</strong>
              <span>{{ stage.description }}</span>
            </article>

            <el-icon
              v-if="index < productionStages.length - 1"
              class="flow-arrow"
            >
              <ArrowRight />
            </el-icon>
          </template>
        </div>

        <div class="flow-note">
          <el-icon><Connection /></el-icon>
          <p>
            生产结果将作为后续设备实例生成和仓储入库的业务依据。
            设备实例不在本工作台中随意创建，而应遵循既定的制造流程。
          </p>
        </div>
      </div>
    </section>

    <section class="section-block">
      <div class="section-heading">
        <div>
          <div class="section-kicker">WORKSPACE STATUS</div>
          <h2>当前建设情况</h2>
          <p>只展示已经确定的模块状态，不使用模拟数据代替真实业务统计。</p>
        </div>
      </div>

      <div class="status-panel">
        <div class="status-icon">
          <el-icon :size="23"><DataAnalysis /></el-icon>
        </div>
        <div class="status-copy">
          <strong>制造业务工作台已建立</strong>
          <p>
            当前产品型号模块已经接入现有 Cloud API。
            生产资源、生产计划和生产执行模块将在后续阶段逐步实现。
          </p>
        </div>
        <el-tag type="info" effect="plain">持续建设中</el-tag>
      </div>
    </section>

    <div class="bottom-space"></div>
  </div>
</template>

<style scoped>
.mes-home {
  width: min(1440px, calc(100% - 64px));
  margin: 0 auto;
  padding-top: 28px;
  color: #17283b;
}

.mes-hero {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 260px;
  padding: 36px 42px;
  overflow: hidden;
  border: 1px solid #243b53;
  border-radius: 15px;
  color: #fff;
  background:
    radial-gradient(circle at 85% 50%, rgb(240 160 48 / 12%), transparent 30%),
    linear-gradient(120deg, #0e1f33 0%, #16293f 65%, #203951 100%);
}

.hero-copy {
  position: relative;
  z-index: 1;
  max-width: 720px;
}

.eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 18px;
  color: #b9c7d7;
  font-size: 10px;
  font-weight: 650;
  letter-spacing: 2px;
}

.eyebrow-line {
  width: 24px;
  height: 2px;
  background: #f0a030;
}

.mes-hero h1 {
  margin: 0;
  color: #fff;
  font-size: 32px;
  font-weight: 750;
  letter-spacing: -0.8px;
}

.mes-hero h1 span {
  margin-left: 8px;
  color: #f0a030;
  font-size: 21px;
  letter-spacing: 1px;
}

.mes-hero p {
  max-width: 650px;
  margin-top: 16px;
  color: #c3cfdb;
  font-size: 13px;
  line-height: 1.9;
}

.hero-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 23px;
  color: #b9c7d7;
  font-size: 11px;
}

.meta-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #f0a030;
}

.meta-divider {
  width: 1px;
  height: 12px;
  margin: 0 3px;
  background: #52667a;
}

.hero-visual {
  position: absolute;
  top: 50%;
  right: 8%;
  width: 210px;
  height: 210px;
  transform: translateY(-50%);
}

.visual-ring {
  position: absolute;
  inset: 8px;
  border: 1px solid rgb(240 160 48 / 20%);
  border-radius: 50%;
}

.ring-two {
  inset: 32px;
  border-color: rgb(188 205 222 / 18%);
}

.visual-center {
  position: absolute;
  top: 50%;
  left: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 84px;
  height: 84px;
  border: 1px solid rgb(240 160 48 / 35%);
  border-radius: 20px;
  background: #203951;
  box-shadow: 0 10px 32px rgb(0 0 0 / 12%);
  transform: translate(-50%, -50%);
}

.visual-center svg {
  width: 56px;
  height: 56px;
}

.visual-node {
  position: absolute;
  width: 9px;
  height: 9px;
  border: 2px solid #f0a030;
  border-radius: 50%;
  background: #16293f;
}

.node-one {
  top: 28px;
  right: 40px;
}

.node-two {
  bottom: 33px;
  left: 22px;
}

.node-three {
  right: 13px;
  bottom: 61px;
}

.section-block {
  margin-top: 34px;
}

.section-heading {
  margin-bottom: 18px;
}

.section-kicker {
  margin-bottom: 8px;
  color: #9a743b;
  font-size: 9px;
  font-weight: 700;
  letter-spacing: 1.8px;
}

.section-heading h2 {
  margin: 0;
  color: #20354a;
  font-size: 22px;
  font-weight: 750;
  letter-spacing: -0.4px;
}

.section-heading p {
  margin-top: 8px;
  color: #728296;
  font-size: 12px;
  line-height: 1.8;
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
}

.module-card {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 252px;
  padding: 23px 22px 19px;
  border: 1px solid #e5eaf0;
  border-radius: 13px;
  color: #17283b;
  background: #fff;
  box-shadow: 0 3px 12px rgb(26 43 61 / 2%);
  text-align: left;
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    box-shadow 180ms ease;
}

button.module-card:not(:disabled) {
  cursor: pointer;
}

button.module-card:not(:disabled):hover {
  border-color: #d8c3a3;
  box-shadow: 0 10px 26px rgb(26 43 61 / 6%);
  transform: translateY(-3px);
}

.module-card--disabled {
  background: #fcfcfd;
}

.module-card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 21px;
}

.module-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 54px;
  height: 54px;
  flex: 0 0 54px;
  border-radius: 12px;
  color: #8c682f;
  background: #fbf3e7;
}

.module-card--product .module-icon {
  color: #315572;
  background: #edf3f8;
}

.module-card--planning .module-icon {
  color: #538a77;
  background: #edf6f1;
}

.module-card--execution .module-icon {
  color: #637aa1;
  background: #eff2fa;
}

.module-subtitle {
  margin-bottom: 7px;
  color: #8a99a8;
  font-size: 9px;
  font-weight: 650;
  letter-spacing: 1.3px;
}

.module-card h3 {
  margin: 0;
  color: #20354a;
  font-size: 17px;
  font-weight: 700;
}

.module-description {
  margin: 12px 0 20px;
  color: #728296;
  font-size: 12px;
  line-height: 1.85;
}

.module-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: auto;
  padding-top: 15px;
  border-top: 1px solid #edf0f4;
  color: #536b80;
  font-size: 11px;
}

.module-card--disabled .module-footer {
  color: #98a4b0;
}

.flow-panel {
  padding: 25px;
  border: 1px solid #e5eaf0;
  border-radius: 13px;
  background: #fff;
}

.flow-list {
  display: flex;
  align-items: stretch;
  gap: 10px;
}

.flow-stage {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
  align-items: center;
  min-width: 0;
  text-align: center;
}

.flow-stage--clickable {
  cursor: pointer;
  border-radius: 10px;
  transition:
    background-color 180ms ease,
    transform 180ms ease;
}

.flow-stage--clickable:hover {
  background-color: #f5f7fa;
  transform: translateY(-2px);
}

.flow-stage--clickable:focus-visible {
  outline: 2px solid #f0a030;
  outline-offset: 3px;
}

.flow-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 45px;
  height: 45px;
  margin-bottom: 12px;
  border: 1px solid #e5eaf0;
  border-radius: 12px;
  color: #536f8a;
  background: #f5f7fa;
}

.flow-stage strong {
  color: #31485e;
  font-size: 12px;
}

.flow-stage > span {
  margin-top: 7px;
  color: #8a99a8;
  font-size: 10px;
  line-height: 1.6;
}

.flow-arrow {
  flex: 0 0 auto;
  margin-top: 14px;
  color: #c2cbd4;
}

.flow-note {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-top: 25px;
  padding: 15px 17px;
  border: 1px solid #f0e3ce;
  border-radius: 8px;
  color: #9a743b;
  background: #fffbf4;
}

.flow-note p {
  color: #756b5b;
  font-size: 11px;
  line-height: 1.8;
}

.status-panel {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 22px;
  border: 1px solid #e5eaf0;
  border-radius: 12px;
  background: #fff;
}

.status-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  flex: 0 0 48px;
  border-radius: 11px;
  color: #315572;
  background: #edf3f8;
}

.status-copy {
  flex: 1;
  min-width: 0;
}

.status-copy strong {
  color: #31485e;
  font-size: 13px;
}

.status-copy p {
  margin-top: 7px;
  color: #728296;
  font-size: 11px;
  line-height: 1.8;
}

.bottom-space {
  height: 48px;
}

@media (max-width: 1200px) {
  .module-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .hero-visual {
    right: 4%;
    opacity: 0.65;
  }

  .hero-copy {
    max-width: 70%;
  }
}

@media (max-width: 760px) {
  .mes-home {
    width: calc(100% - 32px);
    padding-top: 18px;
  }

  .mes-hero {
    min-height: 0;
    padding: 28px 23px;
  }

  .mes-hero h1 {
    font-size: 26px;
  }

  .hero-copy {
    max-width: 100%;
  }

  .hero-visual {
    display: none;
  }

  .section-block {
    margin-top: 28px;
  }

  .module-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .module-card {
    min-height: 220px;
  }

  .flow-panel {
    padding: 18px;
  }

  .flow-list {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 20px 12px;
  }

  .flow-stage {
    padding: 12px 5px;
  }

  .flow-arrow {
    display: none;
  }

  .status-panel {
    align-items: flex-start;
    flex-wrap: wrap;
    padding: 18px;
  }

  .status-copy {
    flex-basis: calc(100% - 65px);
  }

  .status-panel > .el-tag {
    margin-left: 64px;
  }
}
</style>
