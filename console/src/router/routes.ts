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
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("../views/dashboard/Dashboard.vue"),
    meta: {
      requiresAuth: true,
    },
  },

  {
    path: "/users",
    name: "Users",
    component: () => import("../views/users/Users.vue"),
    meta: {
      requiresAuth: true,
      permission: Permissions.USER_LIST,
    },
  },

  {
    path: "/mes/products",
    name: "ProductList",
    component: () => import("../mes/product/views/ProductList.vue"),
    meta: {
      requiresAuth: true,
      permission: Permissions.PRODUCT_LIST,
    },
  },

  {
    path: "/mes/products/:id",
    name: "ProductDetail",
    component: () => import("../mes/product/views/ProductDetail.vue"),
    props: true,
    meta: {
      requiresAuth: true,
      permission: Permissions.PRODUCT_LIST,
    },
  },

  {
    path: "/",
    redirect: "/dashboard",
  },

  {
    path: "/:pathMatch(.*)*",
    redirect: "/dashboard",
  },
];
