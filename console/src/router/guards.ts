import type { Router } from "vue-router";

import { hasPermission } from "../authz/access";
import type { User } from "../api/auth";
import type { Permission } from "../authz/permissions";

function readStoredUser(): User | null {
  const value = localStorage.getItem("auth_user");

  if (!value) {
    return null;
  }

  try {
    return JSON.parse(value) as User;
  } catch {
    localStorage.removeItem("auth_user");
    return null;
  }
}

export function setupRouterGuards(router: Router) {
  router.beforeEach((to) => {
    const accessToken = localStorage.getItem("access_token");
    const user = readStoredUser();

    const requiresAuth = to.meta.requiresAuth === true;
    const guestOnly = to.meta.guestOnly === true;
    const permission = to.meta.permission as Permission | undefined;

    if (requiresAuth && !accessToken) {
      return {
        path: "/login",
        query: {
          redirect: to.fullPath,
        },
      };
    }

    if (guestOnly && accessToken) {
      return "/dashboard";
    }
    if (requiresAuth && permission && !hasPermission(user, permission)) {
      if (to.path !== "/dashboard") {
        return "/dashboard";
      }

      // Dashboard 本身负责展示无权访问或角色未配置的提示。
      return true;
    }

    return true;
  });
}
