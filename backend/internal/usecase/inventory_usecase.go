package usecase

import (
	"context"
	"fmt"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
	"github.com/google/uuid"
)

type InventoryUsecase struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryUsecase(inventoryRepo repository.InventoryRepository) *InventoryUsecase {
	return &InventoryUsecase{inventoryRepo: inventoryRepo}
}

func (u *InventoryUsecase) CreateItem(ctx context.Context, item *entity.InventoryItem) error {
	if item.SKU == "" {
		return fmt.Errorf("SKU is required")
	}
	if item.CostPrice < 0 || item.SalePrice < 0 {
		return fmt.Errorf("prices cannot be negative")
	}
	return u.inventoryRepo.CreateItem(ctx, item)
}

func (u *InventoryUsecase) TransferStock(ctx context.Context, movement entity.StockMovement) error {
	if movement.Quantity <= 0 {
		return entity.ErrNegativeAdjustment
	}
	return u.inventoryRepo.TransferStock(ctx, movement)
}

func (u *InventoryUsecase) ScrapStock(ctx context.Context, tenantID, itemID, sourceLocID, userID uuid.UUID, qty int, ref string) error {
	if qty <= 0 {
		return entity.ErrNegativeAdjustment
	}

	scrapLoc, err := u.inventoryRepo.GetOrCreateScrapLocation(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to locate scrap target: %w", err)
	}

	movement := entity.StockMovement{
		ID:                    uuid.New(),
		TenantID:              tenantID,
		ItemID:                itemID,
		SourceLocationID:      sourceLocID,
		DestinationLocationID: scrapLoc.ID,
		Quantity:              qty,
		ReferenceDocument:     ref,
		CreatedBy:             userID,
	}

	return u.inventoryRepo.TransferStock(ctx, movement)
}

func (u *InventoryUsecase) ListItems(ctx context.Context, tenantID uuid.UUID, search string, limit, offset int) ([]entity.InventoryItem, int, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	return u.inventoryRepo.ListItems(ctx, tenantID, search, limit, offset)
}

func (u *InventoryUsecase) GetItemByBarcode(ctx context.Context, tenantID uuid.UUID, barcode string) (*entity.InventoryItem, error) {
	return u.inventoryRepo.GetItemByBarcode(ctx, tenantID, barcode)
}
