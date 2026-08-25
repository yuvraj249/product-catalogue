package repository

import (
	"context"

	"backend/internal/domain/entity"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	TransferStock(ctx context.Context, movement entity.StockMovement) error
	CreateItem(ctx context.Context, item *entity.InventoryItem) error
	GetItemByID(ctx context.Context, tenantID, itemID uuid.UUID) (*entity.InventoryItem, error)
	GetItemByBarcode(ctx context.Context, tenantID uuid.UUID, barcode string) (*entity.InventoryItem, error)
	ListItems(ctx context.Context, tenantID uuid.UUID, search string, limit, offset int) ([]entity.InventoryItem, int, error)
	GetLocationByID(ctx context.Context, tenantID, locationID uuid.UUID) (*entity.Location, error)
	GetOrCreateScrapLocation(ctx context.Context, tenantID uuid.UUID) (*entity.Location, error)
	GetOrCreateDefaultLocations(ctx context.Context, tenantID uuid.UUID) (mainWhID, customerLocID uuid.UUID, err error)
}
