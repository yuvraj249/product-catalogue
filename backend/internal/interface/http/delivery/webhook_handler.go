package delivery

import (
	"encoding/json"
	"io"
	"net/http"

	"backend/internal/usecase"
)

type WebhookHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewWebhookHandler(paymentUsecase *usecase.PaymentUsecase) *WebhookHandler {
	return &WebhookHandler{paymentUsecase: paymentUsecase}
}

type stripeEvent struct {
	Type string `json:"type"`
	Data struct {
		Object struct {
			ID string `json:"id"`
		} `json:"object"`
	} `json:"data"`
}

func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read webhook body"}`, http.StatusBadRequest)
		return
	}

	var event stripeEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, `{"error":"invalid json webhook event"}`, http.StatusBadRequest)
		return
	}

	if event.Type == "payment_intent.succeeded" {
		paymentIntentID := event.Data.Object.ID
		if paymentIntentID != "" {
			err := h.paymentUsecase.HandlePaymentIntentSucceeded(r.Context(), paymentIntentID)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))
}
