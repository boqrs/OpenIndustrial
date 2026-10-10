import type { RouteRecordRaw } from "vue-router";

import { Permissions } from "../authz/permissions";

export const routes: RouteRecordRaw[] = [
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

  {
    path: "/",
    component: () => import("../layouts/AuthenticatedLayout.vue"),
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        redirect: "/dashboard",
      },

      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("../views/dashboard/Dashboard.vue"),
      },

      {
        path: "users",
        name: "Users",
        component: () => import("../views/users/Users.vue"),
        meta: {
          permission: Permissions.USER_LIST,
        },
      },

      {
        path: "mes",
        name: "MESHome",
        component: () => import("../views/mes/MESHome.vue"),
      },

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
    ],
  },

  {
    path: "/:pathMatch(.*)*",
    redirect: "/dashboard",
  },
];
