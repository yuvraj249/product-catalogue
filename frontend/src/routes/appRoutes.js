import { Routes, Route } from "react-router-dom";
import { LoginPage } from "../Components/Login/LoginPage.Component";
import Dashboard from "../Components/Dashboard/Dashboard.Component";
import { ProtectedRoute } from "./protectedRoutes";
import Products from "../Components/Products/Products.Component";
import Suppliers from "../Components/Suppliers/Suppliers.Component";

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
      <Route path="/products" element={<ProtectedRoute><Products /></ProtectedRoute>} />
      <Route path="/suppliers" element={<ProtectedRoute><Suppliers /></ProtectedRoute> } />
    </Routes>
  );
}
