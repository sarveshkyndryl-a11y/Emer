import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../auth/useAuth";

export default function ProtectedLayout() {
  const auth = useAuth();

  // Not logged in → send to login
  if (!auth.role) {
    return <Navigate to="/login" replace />;
  }

  // Logged in → allow nested routes
  return <Outlet />;
}
