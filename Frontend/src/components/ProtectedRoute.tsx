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
  const auth = useAuth();

  if (!auth.role) return <Navigate to="/login" replace />;
  if (auth.role !== role) return <Navigate to="/404" replace />;

  return <>{children}</>;
}
