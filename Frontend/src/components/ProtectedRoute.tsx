import { Navigate } from "react-router-dom";
import { useAuth } from "../auth/useAuth";
import type { ReactNode } from "react";
import type { Role } from "../auth/auth.types";

export default function ProtectedRoute({
  role,
  children,
}: {
  role: Role;
  children: ReactNode;
}) {
  const { role: userRole, loading } = useAuth();

  // ⛔ CRITICAL: wait for auth bootstrap
  if (loading) {
    return <div>Loading...</div>;
  }

  if (!userRole) {
    return <Navigate to="/login" replace />;
  }

  if (userRole !== role) {
    return <Navigate to="/404" replace />;
  }

  return <>{children}</>;
}
