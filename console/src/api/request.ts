import axios from "axios";

const request = axios.create({
  baseURL: "/api/v1/external",
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

request.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

request.interceptors.response.use(
  (response) => {
    const body = response.data;

    // 后端 ginx.Success 通常返回业务响应对象。
    // 如果响应本身已经是业务数据，则直接返回。
    if (body && typeof body === "object" && "data" in body) {
      return body.data;
    }

    return body;
  },
  (error) => {
    const status = error?.response?.status;

    if (status === 401) {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("auth_user");
    }

    const message =
      error?.response?.data?.message ||
      error?.response?.data?.error ||
      error?.message ||
      "请求失败";

    return Promise.reject(new Error(message));
  },
);

export default request;
