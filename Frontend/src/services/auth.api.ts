import baseApi from "./axios.base";
import type { Role } from "../auth/auth.types";

/* =========================
   TYPES
========================= */

// What backend ACTUALLY returns
type AuthApiResponse = {
  access_token: string;
  role: Role;
};

// What frontend SHOULD use
export type AuthResponse = {
  accessToken: string;
  role: Role;
};

/* =========================
   LOGIN
========================= */

export async function login(
  email: string,
  password: string
): Promise<AuthResponse> {
  const res = await baseApi.post<AuthApiResponse>("/auth/login", {
    email,
    password,
  });

  return {
    accessToken: res.data.access_token,
    role: res.data.role,
  };
}

/* =========================
   REFRESH TOKEN
========================= */

export async function refreshToken(): Promise<AuthResponse> {
  const res = await baseApi.post<AuthApiResponse>("/auth/refresh");

  return {
    accessToken: res.data.access_token,
    role: res.data.role,
  };
}

/* =========================
   RESTORE SESSION (APP LOAD)
========================= */

export async function restoreSession(): Promise<AuthResponse> {
  const res = await baseApi.post<AuthApiResponse>("/auth/refresh");

  return {
    accessToken: res.data.access_token,
    role: res.data.role,
  };
}
