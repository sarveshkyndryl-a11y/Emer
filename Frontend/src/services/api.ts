import axios from "axios";
import { tokenMemory } from "../utils/tokenMemory";
import { refreshToken } from "./auth.api";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true, // refresh cookie
});

let isRefreshing = false;
let queue: any[] = [];

api.interceptors.request.use((config) => {
  const token = tokenMemory.get();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  res => res,
  async err => {
    const original = err.config;

    if (err.response?.status === 401 && !original._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          queue.push({ resolve, reject });
        }).then(token => {
          original.headers.Authorization = `Bearer ${token}`;
          return api(original);
        });
      }

      original._retry = true;
      isRefreshing = true;

      try {
        const { accessToken } = await refreshToken();
        tokenMemory.set(accessToken);

        queue.forEach(p => p.resolve(accessToken));
        queue = [];

        return api(original);
      } catch {
        tokenMemory.clear();
        window.location.href = "/login";
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(err);
  }
);
