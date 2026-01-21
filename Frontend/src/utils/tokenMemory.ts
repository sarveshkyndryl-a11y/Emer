let accessToken: string | null = null;
let role: string | null = null;

export function setAuth(token: string, userRole: string) {
  accessToken = token;
  role = userRole;
}

export function clearAuth() {
  accessToken = null;
  role = null;
}

export function getAccessToken() {
  return accessToken;
}

export function getRole() {
  return role;
}
