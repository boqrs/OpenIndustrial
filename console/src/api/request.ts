import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

const request = axios.create({
  baseURL: "/api/v1/external",
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

let refreshing = false;

let refreshPromise: Promise<string | null> | null = null;

function getAccessToken(): string {
  return localStorage.getItem("access_token") || "";
}

function getRefreshToken(): string {
  return localStorage.getItem("refresh_token") || "";
}

function clearAuth() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("auth_user");
}

function isPublicAuthRequest(url: string | undefined): boolean {
  if (!url) {
    return false;
  }

  return (
    url.includes("/login") ||
    url.includes("/refresh") ||
    url.includes("/identity/access-requests") ||
    url.includes("/identity/invitations/accept")
  );
}

async function refreshAccessToken(): Promise<string | null> {
  const token = getRefreshToken();

  if (!token) {
    return null;
  }

  if (refreshing && refreshPromise) {
    return refreshPromise;
  }

  refreshing = true;

  refreshPromise = axios
    .post(
      "/api/v1/external/refresh",
      {
        refresh_token: token,
      },
      {
        headers: {
          "Content-Type": "application/json",
        },
      },
    )
    .then((response) => {
      const body = response.data;

      const data =
        body && typeof body === "object" && "data" in body ? body.data : body;

      const accessToken = data?.access_token;
      const refreshToken = data?.refresh_token;

      if (!accessToken || !refreshToken) {
        clearAuth();
        return null;
      }

      localStorage.setItem("access_token", accessToken);

      localStorage.setItem("refresh_token", refreshToken);

      return accessToken;
    })
    .catch(() => {
      clearAuth();
      return null;
    })
    .finally(() => {
      refreshing = false;
      refreshPromise = null;
    });

  return refreshPromise;
}

request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getAccessToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

request.interceptors.response.use(
  (response) => {
    const body = response.data;

    if (body && typeof body === "object" && "data" in body) {
      return body.data;
    }

    return body;
  },

  async (error: AxiosError) => {
    const status = error.response?.status;

    const originalRequest = error.config as
      | (InternalAxiosRequestConfig & {
          _retry?: boolean;
        })
      | undefined;

    const isPublicRequest = isPublicAuthRequest(originalRequest?.url);

    if (
      status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !isPublicRequest
    ) {
      originalRequest._retry = true;

      const accessToken = await refreshAccessToken();

      if (accessToken) {
        originalRequest.headers.Authorization = `Bearer ${accessToken}`;

        return request(originalRequest);
      }

      clearAuth();

      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }

    const responseData = error.response?.data as
      | {
          message?: string;
          error?: string;
        }
      | undefined;

    const message =
      responseData?.message ||
      responseData?.error ||
      error.message ||
      "请求失败";

    return Promise.reject(new Error(message));
  },
);

export default request;
