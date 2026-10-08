import type { RouteRecordRaw } from "vue-router";

export const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/Login.vue"),
    meta: {
      guestOnly: true,
    },
  },

  {
    path: "/access-request",
    name: "AccessRequest",
    component: () => import("../views/AccessRequest.vue"),
    meta: {
      guestOnly: true,
    },
  },

  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("../views/Dashboard.vue"),
    meta: {
      requiresAuth: true,
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
