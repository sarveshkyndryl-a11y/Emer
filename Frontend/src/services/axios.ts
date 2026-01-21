import axios from "axios";
import { tokenMemory } from "../utils/tokenMemory";
import { refreshToken } from "./auth.api";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
  withCredentials: true,
});

/* Attach access token */
api.interceptors.request.use((config) => {
  const token = tokenMemory.get();
  if (token && config.headers) {
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
      !originalRequest?._retry
    ) {
      originalRequest._retry = true;

      try {
        const { accessToken } = await refreshToken();

        // ✅ only update token memory
        tokenMemory.set(accessToken);

        if (originalRequest.headers) {
          originalRequest.headers.Authorization =
            `Bearer ${accessToken}`;
        }

        return api(originalRequest);
      } catch (err) {
        // ✅ let AuthProvider + route guards decide
        tokenMemory.clear();
        return Promise.reject(err);
      }
    }

    return Promise.reject(error);
  }
);

export default api;
