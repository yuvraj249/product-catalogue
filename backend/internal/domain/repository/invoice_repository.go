package repository

import (
	"context"

	"backend/internal/domain/entity"
	"github.com/google/uuid"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *entity.Invoice) error
	GetInvoiceByID(ctx context.Context, tenantID, invoiceID uuid.UUID) (*entity.Invoice, error)
	GetInvoiceByStripePaymentIntentID(ctx context.Context, paymentIntentID string) (*entity.Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, tenantID, invoiceID uuid.UUID, status entity.InvoiceStatus) error
}
