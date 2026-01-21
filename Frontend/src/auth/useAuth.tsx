import { useContext } from "react";
import { AuthContext } from "./AuthContext";
import type { AuthContextType } from "./auth.types";

export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return ctx;
}
