package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain/entity"
	"github.com/google/uuid"
)

type InvoiceRepository struct {
	DB *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{DB: db}
}

func (r *InvoiceRepository) CreateInvoice(ctx context.Context, invoice *entity.Invoice) error {
	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	now := time.Now()
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx for invoice: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO invoices (id, tenant_id, customer_name, total_amount, status, stripe_payment_intent_id, idempotency_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = tx.ExecContext(ctx, query,
		invoice.ID, invoice.TenantID, invoice.CustomerName, invoice.TotalAmount, invoice.Status,
		invoice.StripePaymentIntentID, invoice.IdempotencyKey, invoice.CreatedAt, invoice.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert invoice: %w", err)
	}

	for _, item := range invoice.Items {
		item.ID = uuid.New()
		item.TenantID = invoice.TenantID
		item.InvoiceID = invoice.ID

		itemQuery := `
			INSERT INTO invoice_items (id, tenant_id, invoice_id, item_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.ExecContext(ctx, itemQuery, item.ID, item.TenantID, item.InvoiceID, item.ItemID, item.Quantity, item.UnitPrice)
		if err != nil {
			return fmt.Errorf("failed to insert invoice item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *InvoiceRepository) GetInvoiceByID(ctx context.Context, tenantID, invoiceID uuid.UUID) (*entity.Invoice, error) {
	query := `
		SELECT id, tenant_id, customer_name, total_amount, status, COALESCE(stripe_payment_intent_id, ''), idempotency_key, created_at, updated_at
		FROM invoices
		WHERE tenant_id = $1 AND id = $2`

	var inv entity.Invoice
	err := r.DB.QueryRowContext(ctx, query, tenantID, invoiceID).Scan(
		&inv.ID, &inv.TenantID, &inv.CustomerName, &inv.TotalAmount, &inv.Status,
		&inv.StripePaymentIntentID, &inv.IdempotencyKey, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice: %w", err)
	}

	itemsQuery := `SELECT id, tenant_id, invoice_id, item_id, quantity, unit_price FROM invoice_items WHERE tenant_id = $1 AND invoice_id = $2`
	rows, err := r.DB.QueryContext(ctx, itemsQuery, tenantID, invoiceID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it entity.InvoiceItem
			if err := rows.Scan(&it.ID, &it.TenantID, &it.InvoiceID, &it.ItemID, &it.Quantity, &it.UnitPrice); err == nil {
				inv.Items = append(inv.Items, it)
			}
		}
	}

	return &inv, nil
}

func (r *InvoiceRepository) GetInvoiceByStripePaymentIntentID(ctx context.Context, paymentIntentID string) (*entity.Invoice, error) {
	query := `
		SELECT id, tenant_id, customer_name, total_amount, status, COALESCE(stripe_payment_intent_id, ''), idempotency_key, created_at, updated_at
		FROM invoices
		WHERE stripe_payment_intent_id = $1`

	var inv entity.Invoice
	err := r.DB.QueryRowContext(ctx, query, paymentIntentID).Scan(
		&inv.ID, &inv.TenantID, &inv.CustomerName, &inv.TotalAmount, &inv.Status,
		&inv.StripePaymentIntentID, &inv.IdempotencyKey, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice by stripe payment intent: %w", err)
	}

	itemsQuery := `SELECT id, tenant_id, invoice_id, item_id, quantity, unit_price FROM invoice_items WHERE tenant_id = $1 AND invoice_id = $2`
	rows, err := r.DB.QueryContext(ctx, itemsQuery, inv.TenantID, inv.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it entity.InvoiceItem
			if err := rows.Scan(&it.ID, &it.TenantID, &it.InvoiceID, &it.ItemID, &it.Quantity, &it.UnitPrice); err == nil {
				inv.Items = append(inv.Items, it)
			}
		}
	}

	return &inv, nil
}

func (r *InvoiceRepository) UpdateInvoiceStatus(ctx context.Context, tenantID, invoiceID uuid.UUID, status entity.InvoiceStatus) error {
	query := `UPDATE invoices SET status = $1, updated_at = NOW() WHERE tenant_id = $2 AND id = $3`
	_, err := r.DB.ExecContext(ctx, query, status, tenantID, invoiceID)
	if err != nil {
		fallbackQuery := `UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE tenant_id = $2 AND id = $3`
		_, err = r.DB.ExecContext(ctx, fallbackQuery, status, tenantID, invoiceID)
		if err != nil {
			return fmt.Errorf("failed to update invoice status: %w", err)
		}
	}
	return nil
}
