import { Routes, Route} from "react-router-dom";
import { LoginPage } from "../Components/Login/LoginPage.Component";
import { ProtectedRoute } from "./protectedRoutes";
import Categories from "../Components/Categories/Categories.Component";
import AdminComponent from "../Components/AdminComp/Admin.Component";


const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/admin" element={<ProtectedRoute><AdminComponent /></ProtectedRoute>} >
        <Route path="categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
      </Route>
    </Routes>
  );
}

export default AppRoutes

