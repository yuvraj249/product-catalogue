import {jwtDecode} from "jwt-decode"

const TOKEN_KEY = "auth_token";

export const setToken = (token) => {
  localStorage.setItem(TOKEN_KEY, token)
};

export const getToken = () => {
  return localStorage.getItem(TOKEN_KEY)
};

export const clearToken = () => {
  localStorage.removeItem(TOKEN_KEY)
};

export const getUserInfo = () => {
    const token = getToken()
    if (!token) return null

    try{
    return jwtDecode(token)
    }catch(e){
    console.error("Invalid token", e)
    clearToken()
    return null
    }
}


export const isAuthenticated = () => {
  const token = getToken();
  if (!token) return false;

  try {
    const decoded = jwtDecode(token);
    if (decoded.exp * 1000 < Date.now()) {
      clearToken();
      return false;
    }

    return true;
  } catch (err) {
    clearToken();
    return false;
  }
};

export const handleLogout = () => {
  clearToken();
  window.location.href = "/";
};
