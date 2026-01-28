import { Routes, Route} from "react-router-dom";
import { LoginPage } from "../Components/Login";
import { ProtectedRoute } from "./protectedRoutes";
import Categories from "../Components/Categories/index";
import AdminComponent from "../Components/AdminComp";
import Users from "../Components/Users/index";


const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/admin" element={<ProtectedRoute><AdminComponent /></ProtectedRoute>} >
        <Route path="categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
        <Route path="users" element={<ProtectedRoute><Users /></ProtectedRoute>} />

      </Route>
    </Routes>
  );
}

export default AppRoutes

