import { Routes, Route, Navigate} from "react-router-dom";
import { LoginPage } from "../Components/Login";
import { ProtectedRoute } from "./protectedRoutes";
import Categories from "../Components/Categories/index";
import AdminComponent from "../Components/AdminComp";
import Suppliers from "../Components/Suppliers";
import StockMovements from "../Components/Stock_Movements";
import Products from "../Components/Products";
import Dashboard from "../Components/Dashboard";
import Users from "../Components/Users/index";


const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/admin" element={<ProtectedRoute><AdminComponent /></ProtectedRoute>} >
        <Route
          index
          element={
            <Navigate
              to={localStorage.getItem("lastAdminRoute") || "/admin/dashboard"}
              replace
            />
          }
        />
        <Route path="categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
        <Route path="suppliers" element={<ProtectedRoute><Suppliers /></ProtectedRoute>} />
        <Route path="stock-movements" element={<ProtectedRoute><StockMovements /></ProtectedRoute>} />
        <Route path="products" element={<ProtectedRoute><Products /></ProtectedRoute>} />
        <Route path="dashboard" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
        <Route path="users" element={<ProtectedRoute><Users /></ProtectedRoute>} />

      </Route>
    </Routes>
  );
}

export default AppRoutes

