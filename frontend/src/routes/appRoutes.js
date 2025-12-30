import { Routes, Route} from "react-router-dom";
import { LoginPage } from "../Components/Login/LoginPage.Component";
import Dashboard from "../Components/Dashboard/Dashboard.Component";
import { ProtectedRoute } from "./protectedRoutes";
import Suppliers from "../Components/Suppliers/Suppliers.Component";
import Categories from "../Components/Categories/Categories.Component";
import AdminComponent from "../Components/AdminComp/Admin.Component";


const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/admin" element={<AdminComponent />} >
        <Route path="dashboard" element={<ProtectedRoute><Dashboard/></ProtectedRoute>} />
        <Route path="suppliers" element={<ProtectedRoute><Suppliers /></ProtectedRoute> } />
        <Route path="categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
      </Route>
    </Routes>
  );
}

export default AppRoutes

