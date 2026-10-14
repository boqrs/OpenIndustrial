import type { RouteRecordRaw } from "vue-router";

import { Permissions } from "../authz/permissions";

export const routes: RouteRecordRaw[] = [
  // ============================================================
  // 公开路由：未登录用户可以访问
  // ============================================================
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/login/Login.vue"),
    meta: {
      guestOnly: true,
    },
  },
  {
    path: "/access-request",
    name: "AccessRequest",
    component: () => import("../views/login/AccessRequest.vue"),
    meta: {
      guestOnly: true,
    },
  },
  {
    path: "/invitation/accept",
    name: "AcceptInvitation",
    component: () => import("../views/login/AcceptInvitation.vue"),
  },

  // ============================================================
  // 登录后的主布局
  // 所有内部业务路由都必须放在 children 中
  // ============================================================
  {
    path: "/",
    component: () => import("../layouts/AuthenticatedLayout.vue"),
    meta: {
      requiresAuth: true,
    },
    children: [
      // 默认进入 Dashboard
      {
        path: "",
        redirect: "/dashboard",
      },

      // --------------------------------------------------------
      // Dashboard
      // --------------------------------------------------------
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("../views/dashboard/Dashboard.vue"),
      },

      // --------------------------------------------------------
      // 用户管理
      // --------------------------------------------------------
      {
        path: "users",
        name: "Users",
        component: () => import("../views/users/Users.vue"),
        meta: {
          permission: Permissions.USER_LIST,
        },
      },

      // --------------------------------------------------------
      // MES 首页
      // --------------------------------------------------------
      {
        path: "mes",
        name: "MESHome",
        component: () => import("../views/mes/MESHome.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 产品管理
      // --------------------------------------------------------
      {
        path: "mes/products",
        name: "ProductList",
        component: () => import("../mes/product/views/ProductList.vue"),
        meta: {
          permission: Permissions.PRODUCT_LIST,
        },
      },
      {
        path: "mes/products/:id",
        name: "ProductDetail",
        component: () => import("../mes/product/views/ProductDetail.vue"),
        props: true,
        meta: {
          permission: Permissions.PRODUCT_LIST,
        },
      },

      // --------------------------------------------------------
      // MES - BOM 管理
      // 保留现有 BOM 路由，避免重复注册
      // --------------------------------------------------------
      {
        path: "mes/boms",
        name: "BOMList",
        component: () => import("../mes/bom/views/BOMList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
      {
        path: "mes/boms/:id",
        name: "BOMDetail",
        component: () => import("../mes/bom/views/BOMDetail.vue"),
        props: true,
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 物料管理
      // 必须放在 children 内，因此这里使用相对路径
      // --------------------------------------------------------
      {
        path: "mes/materials",
        name: "MaterialList",
        component: () => import("../mes/material/views/MaterialList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 生产计划
      // --------------------------------------------------------
      {
        path: "mes/production-plans",
        name: "ProductionPlanList",
        component: () =>
          import("../mes/production/views/ProductionPlanList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
      {
        path: "mes/production-plans/:id",
        name: "ProductionPlanDetail",
        component: () =>
          import("../mes/production/views/ProductionPlanDetail.vue"),
        props: true,
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 工厂管理
      // --------------------------------------------------------
      {
        path: "mes/factories",
        name: "FactoryList",
        component: () => import("../mes/factory/views/FactoryList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
      {
        path: "mes/production-topology",
        name: "ProductionTopology",
        component: () => import("../mes/factory/views/ProductionTopology.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 工艺路线
      // --------------------------------------------------------
      {
        path: "mes/routings",
        name: "RoutingList",
        component: () => import("../mes/routing/views/RoutingList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
      {
        path: "mes/routings/:id",
        name: "RoutingDetail",
        component: () => import("../mes/routing/views/RoutingDetail.vue"),
        props: true,
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 工单管理
      // --------------------------------------------------------
      {
        path: "mes/work-orders",
        name: "WorkOrderList",
        component: () => import("../mes/production/views/WorkOrderList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
      {
        path: "mes/work-orders/:id",
        name: "WorkOrderDetail",
        component: () => import("../mes/production/views/WorkOrderDetail.vue"),
        props: true,
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },

      // --------------------------------------------------------
      // MES - 生产执行
      // --------------------------------------------------------
      {
        path: "mes/executions",
        name: "ExecutionList",
        component: () => import("../mes/execution/views/ExecutionList.vue"),
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
      {
        path: "mes/executions/:id",
        name: "ExecutionDetail",
        component: () => import("../mes/execution/views/ExecutionDetail.vue"),
        props: true,
        meta: {
          permission: Permissions.MES_VIEW,
        },
      },
    ],
  },
  // ============================================================
  // 兜底路由
  // ============================================================
  {
    path: "/:pathMatch(.*)*",
    redirect: "/dashboard",
  },
];
