import { computed, ref } from "vue";
import { defineStore } from "pinia";

import type { LoginResponse, User } from "../api/auth";

const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const AUTH_USER_KEY = "auth_user";

export const useAuthStore = defineStore("auth", () => {
  const accessToken = ref<string>(localStorage.getItem(ACCESS_TOKEN_KEY) || "");

  const refreshToken = ref<string>(
    localStorage.getItem(REFRESH_TOKEN_KEY) || "",
  );

  const user = ref<User | null>(readUser());

  const isAuthenticated = computed(() => Boolean(accessToken.value));

  function setSession(response: LoginResponse, remember = true) {
    accessToken.value = response.access_token;
    refreshToken.value = response.refresh_token;
    user.value = response.user;

    if (remember) {
      localStorage.setItem(ACCESS_TOKEN_KEY, response.access_token);

      localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);

      localStorage.setItem(AUTH_USER_KEY, JSON.stringify(response.user));
    }
  }

  function clearSession() {
    accessToken.value = "";
    refreshToken.value = "";
    user.value = null;

    localStorage.removeItem(ACCESS_TOKEN_KEY);

    localStorage.removeItem(REFRESH_TOKEN_KEY);

    localStorage.removeItem(AUTH_USER_KEY);
  }

  return {
    accessToken,
    refreshToken,
    user,
    isAuthenticated,
    setSession,
    clearSession,
  };
});

function readUser(): User | null {
  const value = localStorage.getItem(AUTH_USER_KEY);

  if (!value) {
    return null;
  }

  try {
    return JSON.parse(value) as User;
  } catch {
    localStorage.removeItem(AUTH_USER_KEY);

    return null;
  }
}
