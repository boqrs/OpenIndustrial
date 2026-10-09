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
    path: "/",
    redirect: "/dashboard",
  },

  {
    path: "/(.)",
    redirect: "/dashboard",
  },
];
