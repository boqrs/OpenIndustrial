<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import {
  ArrowRight,
  Box,
  Calendar,
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

/**
 * 权限控制
 *
 * 产品管理：需要 PRODUCT_LIST 权限
 * 原料管理、BOM 管理：需要 MES_VIEW 权限
 *
 * 其他模块暂时保留入口展示，后续接入对应页面权限时再细化。
 */
const canAccessProducts = computed(() =>
  hasPermission(authStore.user, Permissions.PRODUCT_LIST),
);

const canAccessMaterials = computed(() =>
  hasPermission(authStore.user, Permissions.MES_VIEW),
);

const canAccessBOMs = computed(() =>
  hasPermission(authStore.user, Permissions.MES_VIEW),
);

interface ResourceModule {
  id: string;
  title: string;
  description: string;
  icon: typeof Box;
  path?: string;
  status: string;
  permission?: boolean;
  color: string;
}

const resourceModules = computed<ResourceModule[]>(() => [
  {
    id: "products",
    title: "产品管理",
    description: "维护产品定义、产品编码及产品基础信息。",
    icon: Box,
    path: "/mes/products",
    status: "产品定义",
    permission: canAccessProducts.value,
    color: "blue",
  },
  {
    id: "materials",
    title: "原料管理",
    description: "维护生产所需原料、物料编码及基础属性。",
    icon: Tools,
    path: "/mes/materials",
    status: "物料基础数据",
    permission: canAccessMaterials.value,
    color: "orange",
  },
  {
    id: "bom",
    title: "BOM 管理",
    description: "维护产品结构，定义产品与生产原料之间的组成关系。",
    icon: Operation,
    path: "/mes/boms",
    status: "产品结构",
    permission: canAccessBOMs.value,
    color: "purple",
  },
  {
    id: "factory",
    title: "工厂管理",
    description: "维护工厂基础信息，并配置生产线、工作站等生产组织结构。",
    icon: OfficeBuilding,
    path: "/mes/factories",
    status: "工厂与生产组织",
    permission: canAccessMaterials.value,
    color: "slate",
  },
  {
    id: "planning",
    title: "生产计划",
    description: "规划生产任务，协调产品需求与生产资源。",
    icon: Calendar,
    path: "/mes/production-plans",
    status: "计划管理",
    permission: canAccessMaterials.value,
    color: "green",
  },
  {
    id: "execution",
    title: "生产执行",
    description: "跟踪工单执行、工序操作和生产过程记录。",
    icon: Setting,
    status: "过程执行",
    color: "cyan",
  },
]);

function openModule(module: ResourceModule) {
  console.log("[MESHome] 点击模块", {
    id: module.id,
    title: module.title,
    path: module.path,
    permission: module.permission,
    canOpen: canOpenModule(module),
  });

  if (!module.path || module.permission === false) {
    return;
  }

  router.push(module.path);
}

function canOpenModule(module: ResourceModule) {
  return Boolean(module.path) && module.permission !== false;
}
</script>

<template>
  <div class="mes-home">
    <!-- 页面顶部 -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-eyebrow">
          <span class="eyebrow-line" />
          MANUFACTURING EXECUTION SYSTEM
        </div>

        <h1>制造执行系统</h1>

        <p class="hero-description">
          从产品定义、物料管理与产品结构配置，到生产计划、生产执行及生产结果，
          建立完整的制造业务流程。
        </p>

        <div class="hero-tags">
          <span class="hero-tag">
            <span class="tag-dot tag-dot--blue" />
            产品与物料
          </span>
          <span class="hero-tag">
            <span class="tag-dot tag-dot--orange" />
            生产组织
          </span>
          <span class="hero-tag">
            <span class="tag-dot tag-dot--green" />
            制造执行
          </span>
        </div>
      </div>

      <div class="hero-decoration" aria-hidden="true">
        <div class="decoration-ring decoration-ring--outer" />
        <div class="decoration-ring decoration-ring--middle" />
        <div class="decoration-ring decoration-ring--inner" />

        <div class="decoration-center">
          <el-icon :size="38">
            <DataAnalysis />
          </el-icon>
        </div>

        <div class="decoration-node decoration-node--one" />
        <div class="decoration-node decoration-node--two" />
        <div class="decoration-node decoration-node--three" />
        <div class="decoration-node decoration-node--four" />
      </div>

      <div class="hero-index" aria-hidden="true">
        <span>MES</span>
        <span class="hero-index-number">01</span>
      </div>
    </section>

    <!-- 制造业务模块 -->
    <section class="section-block">
      <div class="section-heading">
        <div>
          <div class="section-kicker">BUSINESS MODULES</div>
          <h2>制造业务</h2>
          <p>配置制造业务的基础数据与生产流程，逐步建立完整的制造执行体系。</p>
        </div>

        <div class="section-heading-mark">
          <span class="heading-mark-line" />
          <span>MES / CORE</span>
        </div>
      </div>

      <div class="module-grid">
        <article
          v-for="module in resourceModules"
          :key="module.id"
          class="module-card"
          :class="[
            `module-card--${module.color}`,
            { 'module-card--disabled': !canOpenModule(module) },
          ]"
          :role="canOpenModule(module) ? 'button' : undefined"
          :tabindex="canOpenModule(module) ? 0 : -1"
          @click="openModule(module)"
          @keydown.enter="openModule(module)"
          @keydown.space.prevent="openModule(module)"
        >
          <div class="module-card-top">
            <div class="module-icon">
              <el-icon :size="23">
                <component :is="module.icon" />
              </el-icon>
            </div>

            <span class="module-status">
              <span class="module-status-dot" />
              {{ module.status }}
            </span>
          </div>

          <div class="module-card-content">
            <h3>{{ module.title }}</h3>
            <p>{{ module.description }}</p>
          </div>

          <div class="module-card-bottom">
            <span>
              {{
                canOpenModule(module)
                  ? "进入模块"
                  : module.permission === false
                    ? "暂无访问权限"
                    : "即将接入"
              }}
            </span>

            <el-icon v-if="canOpenModule(module)" :size="16">
              <ArrowRight />
            </el-icon>
            <span v-else class="module-card-placeholder">—</span>
          </div>
        </article>
      </div>
    </section>

    <!-- 制造业务主流程：独立 HTML 动画 -->
    <section class="section-block flow-section">
      <div class="section-heading">
        <div>
          <div class="section-kicker">PRODUCTION FLOW</div>
          <h2>制造业务主流程</h2>
          <p>
            通过动态演示了解从产品定义、原料与 BOM 配置，到生产计划、
            生产执行和生产结果的完整业务链路。
          </p>
        </div>

        <div class="flow-heading-badge">
          <span class="flow-heading-badge-dot" />
          <span>PROCESS OVERVIEW</span>
        </div>
      </div>

      <div class="production-animation">
        <iframe
          src="/manufacturing-flow.html"
          title="制造业务主流程动态演示"
          class="production-animation__frame"
          loading="lazy"
        />
      </div>

      <div class="flow-caption">
        <div class="flow-caption-icon">
          <el-icon :size="18">
            <Document />
          </el-icon>
        </div>

        <div class="flow-caption-content">
          <strong>业务流程说明</strong>
          <span>
            产品定义是制造业务的起点；原料与 BOM 为生产提供物料和结构依据；
            后续通过生产计划、生产工单与生产执行形成生产结果。
          </span>
        </div>
      </div>
    </section>

    <!-- 当前工作区 -->
    <section class="section-block workspace-section">
      <div class="section-heading">
        <div>
          <div class="section-kicker">WORKSPACE</div>
          <h2>工作区</h2>
          <p>当前制造业务模块的接入情况。</p>
        </div>
      </div>

      <div class="workspace-card">
        <div class="workspace-icon">
          <el-icon :size="22">
            <DataAnalysis />
          </el-icon>
        </div>

        <div class="workspace-main">
          <strong>制造业务工作区</strong>
          <span>
            产品管理、原料管理和 BOM 管理已提供页面入口；
            生产计划、生产执行及工厂管理可在后续开发中逐步接入。
          </span>
        </div>

        <div class="workspace-status">
          <span class="workspace-status-dot" />
          <span>持续建设中</span>
        </div>
      </div>
    </section>

    <footer class="mes-footer">
      <div class="footer-brand">
        <span class="footer-brand-mark" />
        OPENINDUSTRIAL
      </div>

      <span class="footer-divider" />

      <span>Manufacturing Execution System</span>

      <span class="footer-version">MES / WORKSPACE</span>
    </footer>
  </div>
</template>

<style scoped>
.mes-home {
  --mes-navy: #0e1f33;
  --mes-navy-light: #16293f;
  --mes-gold: #f0a030;
  --mes-text: #172b42;
  --mes-muted: #718096;
  --mes-border: #e6ebf1;

  width: 100%;
  min-width: 0;
  padding: 26px 30px 20px;
  color: var(--mes-text);
  background: #f5f7fa;
  box-sizing: border-box;
}

.hero-section {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 238px;
  padding: 36px 42px;
  overflow: hidden;
  border-radius: 16px;
  color: #fff;
  background:
    radial-gradient(circle at 78% 48%, rgb(57 91 125 / 30%) 0, transparent 32%),
    linear-gradient(115deg, #0e1f33 0%, #16293f 58%, #203a55 100%);
  box-sizing: border-box;
}

.hero-content {
  position: relative;
  z-index: 2;
  max-width: 680px;
}

.hero-eyebrow {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #c4d0dd;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 2px;
}

.eyebrow-line {
  display: inline-block;
  width: 22px;
  height: 2px;
  background: var(--mes-gold);
}

.hero-content h1 {
  margin: 21px 0 12px;
  color: #fff;
  font-size: 32px;
  font-weight: 650;
  letter-spacing: 1px;
  line-height: 1.25;
}

.hero-description {
  max-width: 590px;
  margin: 0;
  color: #c0cbd7;
  font-size: 13px;
  line-height: 1.9;
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-top: 25px;
}

.hero-tag {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  color: #e2e8ef;
  font-size: 11px;
}

.tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.tag-dot--blue {
  background: #77b7f4;
}

.tag-dot--orange {
  background: #f0a030;
}

.tag-dot--green {
  background: #64c9a0;
}

.hero-decoration {
  position: absolute;
  top: 50%;
  right: 12%;
  width: 230px;
  height: 230px;
  transform: translateY(-50%);
  opacity: 0.88;
}

.decoration-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  border: 1px solid rgb(173 199 224 / 17%);
  border-radius: 50%;
  transform: translate(-50%, -50%);
}

.decoration-ring--outer {
  width: 220px;
  height: 220px;
}

.decoration-ring--middle {
  width: 160px;
  height: 160px;
  border-style: dashed;
}

.decoration-ring--inner {
  width: 100px;
  height: 100px;
  border-color: rgb(240 160 48 / 40%);
}

.decoration-center {
  position: absolute;
  top: 50%;
  left: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 62px;
  height: 62px;
  color: var(--mes-gold);
  border: 1px solid rgb(240 160 48 / 40%);
  border-radius: 18px;
  background: rgb(14 31 51 / 70%);
  transform: translate(-50%, -50%);
}

.decoration-node {
  position: absolute;
  width: 9px;
  height: 9px;
  border: 2px solid #d6e2ed;
  border-radius: 50%;
  background: #284762;
  box-shadow: 0 0 12px rgb(140 188 225 / 35%);
}

.decoration-node--one {
  top: 22px;
  left: 108px;
}

.decoration-node--two {
  top: 108px;
  right: 3px;
}

.decoration-node--three {
  bottom: 16px;
  left: 69px;
}

.decoration-node--four {
  top: 68px;
  left: 15px;
  border-color: var(--mes-gold);
}

.hero-index {
  position: absolute;
  right: 28px;
  bottom: 20px;
  display: flex;
  align-items: baseline;
  gap: 8px;
  color: rgb(255 255 255 / 38%);
  font-size: 9px;
  letter-spacing: 2px;
}

.hero-index-number {
  color: rgb(255 255 255 / 70%);
  font-size: 18px;
  font-weight: 300;
}

.section-block {
  margin-top: 34px;
}

.section-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 19px;
}

.section-kicker {
  margin-bottom: 7px;
  color: #8a98a9;
  font-size: 9px;
  font-weight: 750;
  letter-spacing: 2px;
}

.section-heading h2 {
  margin: 0;
  color: #182c43;
  font-size: 20px;
  font-weight: 650;
  letter-spacing: 0.2px;
}

.section-heading p {
  margin: 8px 0 0;
  color: var(--mes-muted);
  font-size: 12px;
  line-height: 1.7;
}

.section-heading-mark {
  display: flex;
  align-items: center;
  gap: 9px;
  padding-bottom: 3px;
  color: #97a3b1;
  font-size: 9px;
  letter-spacing: 1.3px;
  white-space: nowrap;
}

.heading-mark-line {
  width: 22px;
  height: 1px;
  background: var(--mes-gold);
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.module-card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 194px;
  padding: 21px 22px 17px;
  overflow: hidden;
  border: 1px solid var(--mes-border);
  border-radius: 12px;
  background: #fff;
  box-sizing: border-box;
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    box-shadow 180ms ease;
}

.module-card[role="button"] {
  cursor: pointer;
}

.module-card[role="button"]:hover {
  z-index: 1;
  border-color: #cbd7e3;
  box-shadow: 0 8px 22px rgb(23 43 66 / 7%);
  transform: translateY(-3px);
}

.module-card[role="button"]:focus-visible {
  outline: 2px solid #4a83b8;
  outline-offset: 3px;
}

.module-card--disabled {
  background: #fbfcfd;
}

.module-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.module-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 43px;
  height: 43px;
  flex: 0 0 auto;
  border-radius: 11px;
}

.module-card--blue .module-icon {
  color: #3c78b4;
  background: #eaf3fc;
}

.module-card--orange .module-icon {
  color: #c7832d;
  background: #fff3e2;
}

.module-card--purple .module-icon {
  color: #7861bd;
  background: #f0edfb;
}

.module-card--green .module-icon {
  color: #329572;
  background: #e8f7f0;
}

.module-card--cyan .module-icon {
  color: #258a9d;
  background: #e5f5f8;
}

.module-card--slate .module-icon {
  color: #64788d;
  background: #edf1f5;
}

.module-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #8a96a4;
  font-size: 10px;
  white-space: nowrap;
}

.module-status-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #a9b4c0;
}

.module-card--blue .module-status-dot {
  background: #6195cb;
}

.module-card--orange .module-status-dot {
  background: #e3a14d;
}

.module-card--purple .module-status-dot {
  background: #927bce;
}

.module-card-content {
  flex: 1;
  padding-top: 18px;
}

.module-card-content h3 {
  margin: 0;
  color: #20344b;
  font-size: 15px;
  font-weight: 650;
}

.module-card-content p {
  margin: 8px 0 0;
  color: #8290a0;
  font-size: 11px;
  line-height: 1.8;
}

.module-card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 15px;
  padding-top: 12px;
  border-top: 1px solid #f0f2f5;
  color: #7e8d9e;
  font-size: 10px;
}

.module-card[role="button"] .module-card-bottom {
  color: #4b7197;
}

.module-card-placeholder {
  color: #b6c0cb;
}

.flow-heading-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 11px;
  border: 1px solid #e4eaf0;
  border-radius: 6px;
  color: #718296;
  background: #fff;
  font-size: 9px;
  letter-spacing: 1px;
  white-space: nowrap;
}

.flow-heading-badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #56b995;
  box-shadow: 0 0 0 3px #e8f6f0;
}

.production-animation {
  width: 100%;
  overflow: hidden;
  border: 1px solid #e5eaf0;
  border-radius: 13px;
  background: #f5f6f8;
  box-shadow: 0 3px 12px rgb(26 45 66 / 3%);
}

.production-animation__frame {
  display: block;
  width: 100%;
  height: 440px;
  border: 0;
  background: transparent;
}

.flow-caption {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-top: 13px;
  padding: 15px 17px;
  border: 1px solid #e8edf2;
  border-radius: 9px;
  background: #fff;
}

.flow-caption-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  flex: 0 0 auto;
  border-radius: 8px;
  color: #6e87a2;
  background: #eff3f7;
}

.flow-caption-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding-top: 1px;
}

.flow-caption-content strong {
  color: #30455b;
  font-size: 12px;
  font-weight: 650;
}

.flow-caption-content span {
  color: #8290a0;
  font-size: 11px;
  line-height: 1.8;
}

.workspace-card {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 19px 21px;
  border: 1px solid var(--mes-border);
  border-radius: 11px;
  background: #fff;
}

.workspace-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 43px;
  height: 43px;
  flex: 0 0 auto;
  border-radius: 11px;
  color: #547ba1;
  background: #edf3f9;
}

.workspace-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.workspace-main strong {
  color: #263c53;
  font-size: 13px;
  font-weight: 650;
}

.workspace-main span {
  color: #8290a0;
  font-size: 11px;
  line-height: 1.7;
}

.workspace-status {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 7px 10px;
  border: 1px solid #e5eee8;
  border-radius: 5px;
  color: #568b72;
  background: #f5faf7;
  font-size: 10px;
  white-space: nowrap;
}

.workspace-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #5cb58a;
}

.mes-footer {
  display: flex;
  align-items: center;
  gap: 13px;
  margin-top: 35px;
  padding: 17px 0 4px;
  border-top: 1px solid #e5eaf0;
  color: #a0aab6;
  font-size: 9px;
}

.footer-brand {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  color: #718196;
  font-size: 9px;
  font-weight: 750;
  letter-spacing: 1.3px;
}

.footer-brand-mark {
  width: 7px;
  height: 7px;
  border-radius: 2px;
  background: var(--mes-gold);
  transform: rotate(45deg);
}

.footer-divider {
  width: 1px;
  height: 12px;
  background: #dce2e9;
}

.footer-version {
  margin-left: auto;
  letter-spacing: 1px;
}

@media (max-width: 1100px) {
  .mes-home {
    padding-right: 22px;
    padding-left: 22px;
  }

  .hero-decoration {
    right: 6%;
    opacity: 0.55;
  }

  .module-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .mes-home {
    padding: 16px 13px 18px;
  }

  .hero-section {
    min-height: 245px;
    padding: 27px 23px;
  }

  .hero-content {
    max-width: 100%;
  }

  .hero-content h1 {
    font-size: 26px;
  }

  .hero-description {
    max-width: 100%;
    font-size: 12px;
  }

  .hero-decoration {
    top: auto;
    right: -75px;
    bottom: -95px;
    transform: none;
    opacity: 0.24;
  }

  .hero-index {
    right: 17px;
    bottom: 13px;
  }

  .section-block {
    margin-top: 27px;
  }

  .section-heading {
    align-items: flex-start;
  }

  .section-heading h2 {
    font-size: 18px;
  }

  .section-heading p {
    max-width: 100%;
  }

  .section-heading-mark,
  .flow-heading-badge {
    display: none;
  }

  .module-grid {
    grid-template-columns: minmax(0, 1fr);
    gap: 12px;
  }

  .module-card {
    min-height: 175px;
    padding: 17px 18px 14px;
  }

  .production-animation__frame {
    height: 540px;
  }

  .workspace-card {
    align-items: flex-start;
    flex-wrap: wrap;
    padding: 16px;
  }

  .workspace-main {
    flex-basis: calc(100% - 60px);
  }

  .workspace-status {
    margin-left: 58px;
  }

  .mes-footer {
    flex-wrap: wrap;
    gap: 9px;
    line-height: 1.6;
  }

  .footer-version {
    width: 100%;
    margin-left: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .module-card {
    transition: none;
  }

  .module-card[role="button"]:hover {
    transform: none;
  }
}
</style>
