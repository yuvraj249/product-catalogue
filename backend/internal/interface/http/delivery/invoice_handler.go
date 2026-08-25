package delivery

import (
	"encoding/json"
	"net/http"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
	"backend/internal/interface/http/middleware"
)

type InvoiceHandler struct {
	invoiceRepo repository.InvoiceRepository
}

func NewInvoiceHandler(invoiceRepo repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{invoiceRepo: invoiceRepo}
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.CustomJWTClaims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var inv entity.Invoice
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, `{"error":"invalid json request payload"}`, http.StatusBadRequest)
		return
	}
	inv.TenantID = claims.TenantID

	if err := h.invoiceRepo.CreateInvoice(r.Context(), &inv); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(inv)
}
