import { useState } from "react";
import { AuthContext } from "./AuthContext";
import type { Role } from "./auth.types";
import { setAuth, clearAuth, getRole } from "../utils/tokenMemory";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [role, setRole] = useState<Role | null>(getRole() as Role | null);

  const loginSuccess = (token: string, role: Role) => {
    setAuth(token, role);
    setRole(role);
  };

  const logout = () => {
    clearAuth();
    setRole(null);
  };

  return (
    <AuthContext.Provider value={{ role, loginSuccess, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
