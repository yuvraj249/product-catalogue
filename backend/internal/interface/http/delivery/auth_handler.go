package delivery

import (
	"encoding/json"
	"net/http"

	"backend/internal/usecase"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type registerRequest struct {
	TenantName string `json:"tenant_name"`
	Subdomain  string `json:"subdomain"`
	AdminEmail string `json:"admin_email"`
	Password   string `json:"password"`
}

type loginRequest struct {
	Subdomain string `json:"subdomain"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json request payload"}`, http.StatusBadRequest)
		return
	}

	tenant, user, err := h.authUsecase.RegisterTenant(r.Context(), req.TenantName, req.Subdomain, req.AdminEmail, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "tenant registered successfully",
		"tenant":  tenant,
		"user":    user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json request payload"}`, http.StatusBadRequest)
		return
	}

	token, user, err := h.authUsecase.Login(r.Context(), req.Subdomain, req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
