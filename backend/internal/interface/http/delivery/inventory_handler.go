package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domain/entity"
	"backend/internal/interface/http/middleware"
	"backend/internal/usecase"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	inventoryUsecase *usecase.InventoryUsecase
}

func NewInventoryHandler(inventoryUsecase *usecase.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{inventoryUsecase: inventoryUsecase}
}

func (h *InventoryHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.CustomJWTClaims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var item entity.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, `{"error":"invalid json request body"}`, http.StatusBadRequest)
		return
	}
	item.TenantID = claims.TenantID

	if err := h.inventoryUsecase.CreateItem(r.Context(), &item); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.CustomJWTClaims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	q := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	items, total, err := h.inventoryUsecase.ListItems(r.Context(), claims.TenantID, q, limit, offset)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"items": items,
		"total": total,
	})
}

func (h *InventoryHandler) GetItemByBarcode(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.CustomJWTClaims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	barcode := r.URL.Query().Get("barcode")
	if barcode == "" {
		http.Error(w, `{"error":"barcode parameter required"}`, http.StatusBadRequest)
		return
	}

	item, err := h.inventoryUsecase.GetItemByBarcode(r.Context(), claims.TenantID, barcode)
	if err != nil {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

type transferRequest struct {
	ItemID                uuid.UUID `json:"item_id"`
	SourceLocationID      uuid.UUID `json:"source_location_id"`
	DestinationLocationID uuid.UUID `json:"destination_location_id"`
	Quantity              int       `json:"quantity"`
	ReferenceDocument     string    `json:"reference_document"`
}

func (h *InventoryHandler) TransferStock(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.CustomJWTClaims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
		return
	}

	movement := entity.StockMovement{
		ID:                    uuid.New(),
		TenantID:              claims.TenantID,
		ItemID:                req.ItemID,
		SourceLocationID:      req.SourceLocationID,
		DestinationLocationID: req.DestinationLocationID,
		Quantity:              req.Quantity,
		ReferenceDocument:     req.ReferenceDocument,
		CreatedBy:             claims.UserID,
	}

	if err := h.inventoryUsecase.TransferStock(r.Context(), movement); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "stock movement completed successfully"})
}

type scrapRequest struct {
	ItemID            uuid.UUID `json:"item_id"`
	SourceLocationID  uuid.UUID `json:"source_location_id"`
	Quantity          int       `json:"quantity"`
	ReferenceDocument string    `json:"reference_document"`
}

func (h *InventoryHandler) ScrapStock(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.CustomJWTClaims)
	if !ok || claims == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req scrapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
		return
	}

	err := h.inventoryUsecase.ScrapStock(r.Context(), claims.TenantID, req.ItemID, req.SourceLocationID, claims.UserID, req.Quantity, req.ReferenceDocument)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "items moved to scrap location successfully"})
}
