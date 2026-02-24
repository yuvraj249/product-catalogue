import { isAuthenticated } from "../utils/auth";
import { Navigate, useLocation } from "react-router-dom";

export const ProtectedRoute = ({ children }) => {
  const location = useLocation();

  if (!isAuthenticated()) {
    return <Navigate to="/" replace />;
  }

  if (location.pathname.startsWith("/admin")) {
    localStorage.setItem("lastAdminRoute", location.pathname);
  }

  return children;
};
