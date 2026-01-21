import api from "./axios";

export interface AuthResponse {
  accessToken: string;
  role: "admin" | "support" | "customer";
}

/* LOGIN */
export async function login(
  email: string,
  password: string
): Promise<AuthResponse> {
  const res = await api.post<AuthResponse>("/auth/login", {
    email,
    password,
  });
  return res.data;
}

/* REFRESH TOKEN */
export async function refreshToken(): Promise<AuthResponse> {
  const res = await api.post<AuthResponse>("/auth/refresh");
  return res.data;
}
