import { jwtDecode } from "jwt-decode";
import axios from "axios";

const TOKEN_KEY = "auth_token";

export const setToken = (token) => {
  localStorage.setItem(TOKEN_KEY, token);
};

export const getToken = () => {
  return localStorage.getItem(TOKEN_KEY);
};

export const clearToken = () => {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem("lastAdminRoute");
};

export const getUserInfo = () => {
  const token = getToken();
  if (!token) return null;

  try {
    return jwtDecode(token);
  } catch (e) {
    console.error("Invalid token", e);
    clearToken();
    return null;
  }
};

export const isAuthenticated = () => {
  const token = getToken();
  if (!token) return false;

  try {
    const decoded = jwtDecode(token);
    if (decoded.exp * 1000 < Date.now()) {
      // Don't clear immediately if we have a refresh cookie — axios interceptor will handle refresh
      return true;
    }

    return true;
  } catch (err) {
    clearToken();
    return false;
  }
};

export const handleLogout = async () => {
  try {
    const token = getToken();
    await axios.post(
      `${process.env.REACT_APP_API_URL || "http://localhost:8080"}/auth/logout`,
      {},
      {
        withCredentials: true,
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      }
    );
  } catch (err) {
    console.error("Logout error", err);
  } finally {
    clearToken();
    window.location.href = "/";
  }
};
