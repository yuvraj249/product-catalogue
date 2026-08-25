package usecase

import (
	"context"
	"fmt"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
	"github.com/google/uuid"
)

type PaymentUsecase struct {
	invoiceRepo   repository.InvoiceRepository
	inventoryRepo repository.InventoryRepository
}

func NewPaymentUsecase(invoiceRepo repository.InvoiceRepository, inventoryRepo repository.InventoryRepository) *PaymentUsecase {
	return &PaymentUsecase{
		invoiceRepo:   invoiceRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (u *PaymentUsecase) HandlePaymentIntentSucceeded(ctx context.Context, paymentIntentID string) error {
	invoice, err := u.invoiceRepo.GetInvoiceByStripePaymentIntentID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("failed to locate invoice for Stripe PaymentIntent %s: %w", paymentIntentID, err)
	}

	if invoice.Status == entity.InvoiceStatusPaid {
		return nil // Already processed idempotently
	}

	// Update Invoice Status to Paid
	if err := u.invoiceRepo.UpdateInvoiceStatus(ctx, invoice.TenantID, invoice.ID, entity.InvoiceStatusPaid); err != nil {
		return fmt.Errorf("failed to update invoice status to Paid: %w", err)
	}

	// Locate Main Warehouse and Customer Delivery Location
	mainWhID, customerLocID, err := u.inventoryRepo.GetOrCreateDefaultLocations(ctx, invoice.TenantID)
	if err != nil {
		return fmt.Errorf("failed to resolve default warehouse locations: %w", err)
	}

	// Trigger Stock Transfer for every item to Customer Delivery Location
	systemUserID := uuid.Nil // System Automated User
	for _, item := range invoice.Items {
		movement := entity.StockMovement{
			ID:                    uuid.New(),
			TenantID:              invoice.TenantID,
			ItemID:                item.ItemID,
			SourceLocationID:      mainWhID,
			DestinationLocationID: customerLocID,
			Quantity:              item.Quantity,
			ReferenceDocument:     fmt.Sprintf("INV-%s", invoice.ID.String()[:8]),
			CreatedBy:             systemUserID,
		}

		if err := u.inventoryRepo.TransferStock(ctx, movement); err != nil {
			return fmt.Errorf("failed to transfer stock for invoice item %s: %w", item.ItemID, err)
		}
	}

	return nil
}
