import type { RouteRecordRaw } from "vue-router";

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
    meta: {
      guestOnly: true,
    },
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
