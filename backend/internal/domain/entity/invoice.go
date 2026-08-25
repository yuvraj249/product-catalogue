package entity

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "Draft"
	InvoiceStatusUnpaid    InvoiceStatus = "Unpaid"
	InvoiceStatusPaid      InvoiceStatus = "Paid"
	InvoiceStatusCancelled InvoiceStatus = "Cancelled"
)

type InvoiceItem struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	InvoiceID uuid.UUID `json:"invoice_id"`
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

type Invoice struct {
	ID                     uuid.UUID     `json:"id"`
	TenantID               uuid.UUID     `json:"tenant_id"`
	CustomerName           string        `json:"customer_name"`
	TotalAmount            float64       `json:"total_amount"`
	Status                 InvoiceStatus `json:"status"`
	StripePaymentIntentID string        `json:"stripe_payment_intent_id,omitempty"`
	IdempotencyKey         *uuid.UUID    `json:"idempotency_key,omitempty"`
	Items                  []InvoiceItem `json:"items,omitempty"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
}
