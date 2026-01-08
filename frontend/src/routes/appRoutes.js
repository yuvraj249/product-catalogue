import { Routes, Route} from "react-router-dom";
import { LoginPage } from "../Components/Login";
import { ProtectedRoute } from "./protectedRoutes";
import Categories from "../Components/Categories/index";
import AdminComponent from "../Components/AdminComp";
import Suppliers from "../Components/Suppliers";


const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/admin" element={<ProtectedRoute><AdminComponent /></ProtectedRoute>} >
        <Route path="categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
        <Route path="suppliers" element={<ProtectedRoute><Suppliers /></ProtectedRoute>} />
      </Route>
    </Routes>
  );
}

export default AppRoutes

