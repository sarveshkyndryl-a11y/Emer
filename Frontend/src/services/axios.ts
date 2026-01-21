import axios from "axios";
import { getAccessToken, setAuth, clearAuth } from "../utils/tokenMemory";
import { refreshToken } from "./auth.api";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
  withCredentials: true, // 🔑 sends refresh_token cookie
});

/* Attach access token */
api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

/* Handle 401 → refresh token */
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const originalRequest = error.config;

    if (
      error.response?.status === 401 &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true;

      try {
        const data = await refreshToken();
        setAuth(data.accessToken, data.role);

        originalRequest.headers.Authorization =
          `Bearer ${data.accessToken}`;

        return api(originalRequest);
      } catch {
        clearAuth();
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  }
);

export default api;
