import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../auth/useAuth";

export default function ProtectedLayout() {
  const { role, loading } = useAuth();

  // ⛔ DO NOTHING until auth is restored
  if (loading) {
    return <div>Loading...</div>; // or spinner
  }

  if (!role) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
