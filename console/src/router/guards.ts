import type { Router } from "vue-router";

export function setupRouterGuards(router: Router) {
  router.beforeEach((to) => {
    const accessToken = localStorage.getItem("access_token");

    const requiresAuth = to.meta.requiresAuth === true;

    const isGuestOnly = to.meta.guestOnly === true;

    if (requiresAuth && !accessToken) {
      return {
        path: "/login",
        query: {
          redirect: to.fullPath,
        },
      };
    }

    if (isGuestOnly && accessToken) {
      return "/dashboard";
    }

    return true;
  });
}
