package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientStock     = errors.New("insufficient inventory available to complete operation")
	ErrNegativeAdjustment    = errors.New("inventory adjustment quantity cannot be negative")
	ErrIdenticalLocations    = errors.New("source and destination locations must be distinct")
	ErrUnauthorizedOperation = errors.New("user does not have clearance for this tenant or action")
)

type InventoryItem struct {
	ID                uuid.UUID `json:"id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	SKU               string    `json:"sku"`
	Name              string    `json:"name"`
	Barcode           string    `json:"barcode"`
	CostPrice         float64   `json:"cost_price"`
	SalePrice         float64   `json:"sale_price"`
	MinStockThreshold int       `json:"min_stock_threshold"`
	MaxStockThreshold int       `json:"max_stock_threshold"`
	CurrentStock      int       `json:"current_stock,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StockMovement struct {
	ID                    uuid.UUID `json:"id"`
	TenantID              uuid.UUID `json:"tenant_id"`
	ItemID                uuid.UUID `json:"item_id"`
	SourceLocationID      uuid.UUID `json:"source_location_id"`
	DestinationLocationID uuid.UUID `json:"destination_location_id"`
	Quantity              int       `json:"quantity"`
	ReferenceDocument     string    `json:"reference_document"`
	CreatedBy             uuid.UUID `json:"created_by"`
	CreatedAt             time.Time `json:"created_at"`
}

type StockQuantity struct {
	TenantID   uuid.UUID `json:"tenant_id"`
	LocationID uuid.UUID `json:"location_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Quantity   int       `json:"quantity"`
	UpdatedAt  time.Time `json:"updated_at"`
}
