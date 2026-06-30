import { Navigate, Outlet } from "react-router-dom";

import { useAuth } from "../features/auth/AuthContext";

export function ProtectedRoute() {
  const auth = useAuth();
  return auth.token === null ? <Navigate to="/login" replace /> : <Outlet />;
}
