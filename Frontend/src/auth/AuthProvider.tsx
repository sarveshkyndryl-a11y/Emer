import { useEffect, useState } from "react";
import { AuthContext } from "./AuthContext";
import type { Role } from "./auth.types";
import { tokenMemory } from "../utils/tokenMemory";
import { restoreSession } from "../services/auth.api";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [role, setRole] = useState<Role | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function bootstrapAuth() {
      try {
        const { accessToken, role } = await restoreSession();
        tokenMemory.set(accessToken);
        setRole(role);
      } catch {
        tokenMemory.clear();
        setRole(null);
      } finally {
        setLoading(false);
      }
    }

    bootstrapAuth();
  }, []);

  const loginSuccess = (token: string, role: Role) => {
    tokenMemory.set(token);
    setRole(role);
  };

  const logout = () => {
    tokenMemory.clear();
    setRole(null);
  };

  return (
    <AuthContext.Provider value={{ role, loading, loginSuccess, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
