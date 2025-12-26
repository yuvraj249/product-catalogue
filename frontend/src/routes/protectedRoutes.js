import { isAuthenticated } from "../utils/auth";
import { Navigate } from "react-router-dom";

export const ProtectedRoute = ({ children })  => {
  return isAuthenticated() ? children : <Navigate to="/" replace />;
}